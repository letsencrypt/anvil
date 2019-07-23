package main

import (
	"errors"

	"github.com/letsencrypt/boulder/cmd"
)

// CleanupConfig describes common configuration parameters shared by all cleanup
// jobs.
type CleanupConfig struct {
	// Enabled controls whether the janitor will run this cleanup job.
	Enabled bool
	// GracePeriod controls when a resource is old enough to be cleaned up.
	GracePeriod cmd.ConfigDuration
	// BatchSize controls how many rows of the resource will be read from the DB
	// per-query.
	BatchSize int64
	// Parallelism controls how many independent go routines will run Delete
	// statements for old resources being cleaned up.
	Parallelism int
	// MaxDPS controls the maximum number of deletes which will be performed
	// per second in total for the resource's table across all of the parallel go
	// routines for this resource. This can be used to reduce the replication lag
	// caused by creating a very large numbers of delete statements.
	MaxDPS int
}

var (
	errInvalidGracePeriod   = errors.New("grace period must be > 0")
	errInvalidParallelism   = errors.New("parallelism must be > 0")
	errInvalidNegativeValue = errors.New("no numeric configuration values should be negative")
)

// Valid checks the cleanup config is valid or returns an error.
func (c CleanupConfig) Valid() error {
	if c.GracePeriod.Duration <= 0 {
		return errInvalidGracePeriod
	}
	if c.Parallelism <= 0 {
		return errInvalidParallelism
	}
	if c.BatchSize < 0 || c.MaxDPS < 0 {
		return errInvalidNegativeValue
	}
	return nil
}

// Config describes the overall Janitor configuration.
type Config struct {
	Janitor struct {
		// Syslog holds common syslog configuration.
		Syslog cmd.SyslogConfig
		// DebugAddr controls the bind address for prometheus metrics, etc.
		DebugAddr string
		// Features holds potential Feature flags.
		Features map[string]bool
		// Common database connection configuration.
		cmd.DBConfig

		// Certificates describes a cleanup job for the certificates table.
		Certificates struct {
			CleanupConfig
		}

		// CertificateStatus describes a cleanup job for the certificateStatus table.
		CertificateStatus struct {
			CleanupConfig
		}

		// CertificatesPerName describes a cleanup job for the certificatesPerName table.
		CertificatesPerName struct {
			CleanupConfig
		}
	}
}

// Valid checks that each of the cleanup job configurations are valid or returns
// an error.
func (c Config) Valid() error {
	if err := c.Janitor.Certificates.Valid(); err != nil {
		return err
	}
	if err := c.Janitor.CertificateStatus.Valid(); err != nil {
		return err
	}
	if err := c.Janitor.CertificatesPerName.Valid(); err != nil {
		return err
	}
	return nil
}
