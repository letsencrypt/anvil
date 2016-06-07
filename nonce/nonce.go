package nonce

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"
	"sync/atomic"
)

// MaxUsed defines the maximum number of Nonces we're willing to hold in
// memory.
const MaxUsed = 65536
const nonceLen = 32

var errInvalidNonceLength = errors.New("invalid nonce length")

// NonceService generates, cancels, and tracks Nonces.
type NonceService struct {
	mu       sync.RWMutex
	latest   int64
	earliest int64
	used     map[int64]bool
	gcm      cipher.AEAD
	maxUsed  int
}

// NewNonceService constructs a NonceService with defaults
func NewNonceService() (*NonceService, error) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		panic("Failure in NewCipher: " + err.Error())
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic("Failure in NewGCM: " + err.Error())
	}

	return &NonceService{
		earliest: 0,
		latest:   0,
		used:     make(map[int64]bool, MaxUsed),
		gcm:      gcm,
		maxUsed:  MaxUsed,
	}, nil
}

func (ns *NonceService) encrypt(counter int64) (string, error) {
	// Generate a nonce with upper 4 bytes zero
	nonce := make([]byte, 12)
	for i := 0; i < 4; i++ {
		nonce[i] = 0
	}
	if _, err := rand.Read(nonce[4:]); err != nil {
		return "", err
	}

	// Encode counter to plaintext
	pt := make([]byte, 8)
	binary.BigEndian.PutUint64(pt, uint64(counter))

	// Encrypt
	ret := make([]byte, nonceLen)
	ct := ns.gcm.Seal(nil, nonce, pt, nil)
	copy(ret, nonce[4:])
	copy(ret[8:], ct)
	return base64.RawURLEncoding.EncodeToString(ret), nil
}

func (ns *NonceService) decrypt(nonce string) (int64, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(nonce)
	if err != nil {
		return 0, err
	}
	if len(decoded) != nonceLen {
		return 0, errInvalidNonceLength
	}

	n := make([]byte, 12)
	for i := 0; i < 4; i++ {
		n[i] = 0
	}
	copy(n[4:], decoded[:8])

	pt, err := ns.gcm.Open(nil, n, decoded[8:], nil)
	if err != nil {
		return 0, err
	}

	return int64(binary.BigEndian.Uint64(pt)), nil
}

// Nonce provides a new Nonce.
func (ns *NonceService) Nonce() (string, error) {
	return ns.encrypt(atomic.AddInt64(&ns.latest, 1))
}

// minUsed returns the lowest key in the used map. Requires that a lock be held
// by caller.
func (ns *NonceService) minUsed() int64 {
	min := atomic.LoadInt64(&ns.latest)
	for t := range ns.used {
		if t < min {
			min = t
		}
	}
	return min
}

// Valid determines whether the provided Nonce string is valid, returning
// true if so.
func (ns *NonceService) Valid(nonce string) bool {
	c, err := ns.decrypt(nonce)
	if err != nil {
		return false
	}

	if c > atomic.LoadInt64(&ns.latest) {
		return false
	}

	if c <= atomic.LoadInt64(&ns.earliest) {
		return false
	}

	ns.mu.RLock()
	if ns.used[c] {
		ns.mu.RUnlock()
		return false
	}
	ns.mu.RUnlock()

	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.used[c] = true
	if len(ns.used) > ns.maxUsed {
		earliest := ns.minUsed()
		atomic.StoreInt64(&ns.earliest, earliest)
		delete(ns.used, earliest)
	}

	return true
}
