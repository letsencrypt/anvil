package ca

import (
	"encoding/hex"
	"testing"
	"time"

	blog "github.com/letsencrypt/boulder/log"
	"github.com/letsencrypt/boulder/metrics"
	"github.com/letsencrypt/boulder/test"
	"golang.org/x/crypto/ocsp"
)

// Set up an ocspLogQueue with a very long period and a large maxLen,
// to ensure any buffered entries get flushed on `.stop()`.
func TestOcspLogFlushOnExit(t *testing.T) {
	t.Parallel()
	log := blog.NewMock()
	stats := metrics.NoopRegisterer
	queue := newOCSPLogQueue(4000, 10000*time.Millisecond, stats, log)
	serial, err := hex.DecodeString("aabbccddeeffaabbccddeeff000102030405")
	if err != nil {
		t.Fatal(err)
	}
	go queue.loop()
	queue.enqueue(serial, time.Now(), ocsp.ResponseStatus(ocsp.Good))
	queue.stop()

	expected := []string{
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,",
	}
	test.AssertDeepEquals(t, log.GetAll(), expected)
}

// Ensure log lines are sent when they exceed maxLen.
func TestOcspFlushOnLength(t *testing.T) {
	t.Parallel()
	log := blog.NewMock()
	stats := metrics.NoopRegisterer
	queue := newOCSPLogQueue(100, 100*time.Millisecond, stats, log)
	serial, err := hex.DecodeString("aabbccddeeffaabbccddeeff000102030405")
	if err != nil {
		t.Fatal(err)
	}
	go queue.loop()
	for i := 0; i < 5; i++ {
		queue.enqueue(serial, time.Now(), ocsp.ResponseStatus(ocsp.Good))
	}
	queue.stop()

	expected := []string{
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,aabbccddeeffaabbccddeeff000102030405:0,",
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,aabbccddeeffaabbccddeeff000102030405:0,",
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,",
	}
	test.AssertDeepEquals(t, log.GetAll(), expected)
}

// Ensure log lines are sent after a timeout.
func TestOcspFlushOnTimeout(t *testing.T) {
	t.Parallel()
	log := blog.NewMock()
	stats := metrics.NoopRegisterer
	queue := newOCSPLogQueue(90000, 10*time.Millisecond, stats, log)
	serial, err := hex.DecodeString("aabbccddeeffaabbccddeeff000102030405")
	if err != nil {
		t.Fatal(err)
	}

	go queue.loop()
	for i := 0; i < 3; i++ {
		queue.enqueue(serial, time.Now(), ocsp.ResponseStatus(ocsp.Good))
		queue.enqueue(serial, time.Now(), ocsp.ResponseStatus(ocsp.Good))
		// This gets a little tricky: Each iteration of the `select` in
		// queue.loop() is a race between the channel that receives log
		// events and the `<-clk.After(n)` timer. Even if we used
		// a fake clock, our loop here would often win that race, producing
		// inconsistent logging results. For instance, it would be entirely
		// possible for all of these `enqueues` to win the race, putting
		// all log entries on one line.
		// To avoid that, sleep using the wall clock for 50ms.
		time.Sleep(50 * time.Millisecond)
	}
	queue.stop()

	expected := []string{
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,aabbccddeeffaabbccddeeff000102030405:0,",
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,aabbccddeeffaabbccddeeff000102030405:0,",
		"INFO: [AUDIT] OCSP signed: aabbccddeeffaabbccddeeff000102030405:0,aabbccddeeffaabbccddeeff000102030405:0,",
	}
	test.AssertDeepEquals(t, log.GetAll(), expected)
}
