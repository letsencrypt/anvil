package ca

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/ocsp"

	blog "github.com/letsencrypt/boulder/log"
)

// ocspLogQueue accumulates OCSP logging events and writes several of them
// in a single log line. This reduces the number of log lines and bytes,
// which would otherwise be quite high. As of Jan 2021 we do approximately
// 550 rps of OCSP generation events. We can turn that into about 5.5 rps
// of log lines if we accumulate 100 entries per line, which amounts to about
// 3900 bytes per log line.
// Summary of log line usage:
// serial in hex: 36 bytes, separator characters: 2 bytes, status: 1 byte
type ocspLogQueue struct {
	// Maximum length, in bytes, of a single log line.
	maxLogLen int
	// Maximum amount of time between OCSP logging events.
	period time.Duration
	queue  chan ocspLog
	// This allows the stop() function to block until we've drained the queue.
	wg     sync.WaitGroup
	depth  prometheus.Gauge
	logger blog.Logger
}

type ocspLog struct {
	serial []byte
	time   time.Time
	status ocsp.ResponseStatus
}

func newOCSPLogQueue(
	maxLogLen int,
	period time.Duration,
	stats prometheus.Registerer,
	logger blog.Logger,
) *ocspLogQueue {
	depth := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ocsp_log_queue_depth",
			Help: "Number of OCSP generation log entries waiting to be written",
		})
	stats.MustRegister(depth)
	return &ocspLogQueue{
		maxLogLen: maxLogLen,
		period:    period,
		queue:     make(chan ocspLog, 1000),
		wg:        sync.WaitGroup{},
		depth:     depth,
		logger:    logger,
	}
}

func (olq *ocspLogQueue) enqueue(serial []byte, time time.Time, status ocsp.ResponseStatus) {
	olq.queue <- ocspLog{
		serial: serial,
		time:   time,
		status: ocsp.ResponseStatus(status),
	}
}

// loop consumes events from the queue channel, batches them up, and
// logs them in batches of 100, or every 500 milliseconds, whichever comes first.
func (olq *ocspLogQueue) loop() {
	olq.wg.Add(1)
	defer olq.wg.Done()
	done := false
	for !done {
		var builder strings.Builder
		deadline := time.After(500 * time.Millisecond)
	inner:
		// To ensure we don't go over the max log line length,
		// use a safety margin greater than the expected length of
		// each entry.
		for builder.Len() < olq.maxLogLen-50 {
			olq.depth.Set(float64(len(olq.queue)))
			select {
			case ol, ok := <-olq.queue:
				if !ok {
					// Channel was closed, finish.
					done = true
					break inner
				} else {
					fmt.Fprintf(&builder, "%x:%d,", ol.serial, ol.status)
				}
			case <-deadline:
				break inner
			}
		}
		if builder.Len() > 0 {
			olq.logger.AuditInfof("OCSP signed: %s", builder.String())
		}
	}
}

// stop the loop, and wait for it to finish. This must be called only after
// it's guaranteed that nothing will call enqueue again (for instance, after
// the OCSPGenerator and CertificateAuthority services are shut down with
// no RPCs in flight). Otherwise, enqueue will panic.
func (olq *ocspLogQueue) stop() {
	close(olq.queue)
	olq.wg.Wait()
}
