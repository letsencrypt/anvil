package ca

import (
	"io/ioutil"
	"sync"

	"github.com/letsencrypt/boulder/log"
	"github.com/letsencrypt/boulder/reloader"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

// ECDSAAllowList acts as a container for a map of Registration IDs, a
// mutex, and a file reloader. This allows the map of IDs to be updated
// safely if changes to the allow list are detected.
type ECDSAAllowList struct {
	sync.RWMutex
	regIDsMap   map[int64]bool
	reloader    *reloader.Reloader
	logger      log.Logger
	statusGauge *prometheus.GaugeVec
}

// Update is an exported method (typically specified as a callback to a
// file reloader) that replaces the inner `regIDsMap` with the contents
// of a YAML list (as bytes)
func (e *ECDSAAllowList) Update(contents []byte) error {
	newRegIDsMap, err := unmarshalAllowList(contents)
	if err != nil {
		return err
	}
	e.Lock()
	defer e.Unlock()
	e.regIDsMap = newRegIDsMap
	// nil check for testing purposes
	if e.statusGauge != nil {
		e.statusGauge.WithLabelValues("succeeded").Set(float64(len(e.regIDsMap)))
	}
	return nil
}

// UpdateErr is an exported method,(typically specified as a callback to
// a file reloader) that records failed allow list file reload attempts.
func (e *ECDSAAllowList) UpdateErr(err error) {
	e.logger.Errf("error reloading ECDSA allowed list: %s", err)
	e.RLock()
	defer e.RUnlock()
	// nil check for testing purposes
	if e.statusGauge != nil {
		e.statusGauge.WithLabelValues("failed").Set(float64(len(e.regIDsMap)))
	}
}

// permitted checks if ECDSA issuance is permitted for the specified
// Registration ID.
func (e *ECDSAAllowList) permitted(regID int64) bool {
	e.RLock()
	defer e.RUnlock()
	return e.regIDsMap[regID]
}

// Stop stops an active allow list reloader. Typically called during
// boulder-ca shutdown.
func (e *ECDSAAllowList) Stop() {
	e.Lock()
	defer e.Unlock()
	if e.reloader != nil {
		e.reloader.Stop()
	}
}

func unmarshalAllowList(contents []byte) (map[int64]bool, error) {
	var regIDs []int64
	err := yaml.Unmarshal(contents, &regIDs)
	if err != nil {
		return nil, err
	}
	return makeRegIDsMap(regIDs), nil
}

func makeRegIDsMap(regIDs []int64) map[int64]bool {
	regIDsMap := make(map[int64]bool)
	for _, regID := range regIDs {
		regIDsMap[regID] = true
	}
	return regIDsMap
}

// NewECDSAAllowListFromFile is exported to allow `boulder-ca` to
// construct a new `ECDSAAllowList` object and set the initial allow
// list using the contents of a YAML file. An initial entry count is
// returned to `boulder-ca` for logging purposes.
func NewECDSAAllowListFromFile(filename string, reloader *reloader.Reloader, logger log.Logger, metric *prometheus.GaugeVec) (*ECDSAAllowList, int, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, 0, err
	}
	regIDsMap, err := unmarshalAllowList(contents)
	if err != nil {
		return nil, 0, err
	}
	return &ECDSAAllowList{regIDsMap: regIDsMap, reloader: reloader, logger: logger, statusGauge: metric}, len(regIDsMap), nil
}

// NewECDSAAllowListFromConfig is exported to allow `boulder-ca` to
// construct a new `ECDSAAllowList` object and set the inner `regIDsMap`
// from a list of registration IDs received in the CA config JSON.
//
// TODO(#5394): This is deprecated and exists to support deployability
// until `ECDSAAllowedAccounts` is replaced by `ECDSAAllowListFilename`
// in all staging and production configs.
func NewECDSAAllowListFromConfig(regIDs []int64) (*ECDSAAllowList, error) {
	regIDsMap := makeRegIDsMap(regIDs)
	return &ECDSAAllowList{regIDsMap: regIDsMap, reloader: nil, logger: nil, statusGauge: nil}, nil
}
