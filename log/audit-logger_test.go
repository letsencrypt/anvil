// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package log

import (
	"errors"
	"fmt"
	"log/syslog"
	"net"
	"testing"
	"time"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/jmhodges/clock"
	"github.com/letsencrypt/boulder/test"
)

const stdoutLevel = 7

func setup(t *testing.T) *impl {
	// Write all logs to UDP on a high port so as to not bother the system
	// which is running the test, particularly for Emerg()
	writer, err := syslog.Dial("udp", "127.0.0.1:65530", syslog.LOG_INFO|syslog.LOG_LOCAL0, "")
	test.AssertNotError(t, err, "Could not construct syslog object")

	logger, err := New(writer, stdoutLevel)
	test.AssertNotError(t, err, "Could not construct syslog object")
	impl, ok := logger.(*impl)
	if !ok {
		t.Fatalf("Wrong type returned from New: %T", logger)
	}
	return impl
}

func TestConstruction(t *testing.T) {
	t.Parallel()
	_ = setup(t)
}

func TestSingleton(t *testing.T) {
	t.Parallel()
	log1 := Get()
	test.AssertNotNil(t, log1, "Logger shouldn't be nil")

	log2 := Get()
	test.AssertEquals(t, log1, log2)

	audit := setup(t)

	// Should not work
	err := Set(audit)
	test.AssertError(t, err, "Can't re-set")

	// Verify no change
	log4 := Get()

	// Verify that log4 != log3
	test.AssertNotEquals(t, log4, audit)

	// Verify that log4 == log2 == log1
	test.AssertEquals(t, log4, log2)
	test.AssertEquals(t, log4, log1)
}

func TestConstructionNil(t *testing.T) {
	t.Parallel()
	_, err := New(nil, stdoutLevel)
	test.AssertError(t, err, "Nil shouldn't be permitted.")
}

func TestEmit(t *testing.T) {
	t.Parallel()
	log := setup(t)

	log.AuditNotice("test message")
}

func TestEmitEmpty(t *testing.T) {
	t.Parallel()
	log := setup(t)

	log.AuditNotice("")
}

func ExampleAuditLogger() {
	impl := setup(nil)

	bw, ok := impl.w.(*bothWriter)
	if !ok {
		fmt.Printf("Wrong type of impl's writer: %T\n", impl.w)
		return
	}
	bw.clk = clock.NewFake()
	impl.AuditErr(errors.New("Error Audit"))
	impl.Warning("Warning Audit")
	// Output: [31m[1mE000000 log.test [AUDIT] Error Audit[0m
	// [33mW000000 log.test Warning Audit[0m
}

func TestSyslogMethods(t *testing.T) {
	t.Parallel()
	audit := setup(t)

	audit.AuditNotice("audit-logger_test.go: audit-notice")
	audit.AuditErr(errors.New("audit-logger_test.go: audit-err"))
	audit.Crit("audit-logger_test.go: critical")
	audit.Debug("audit-logger_test.go: debug")
	audit.Emerg("audit-logger_test.go: emerg")
	audit.Err("audit-logger_test.go: err")
	audit.Info("audit-logger_test.go: info")
	audit.Notice("audit-logger_test.go: notice")
	audit.Warning("audit-logger_test.go: warning")
}

func TestPanic(t *testing.T) {
	t.Parallel()
	audit := setup(t)
	defer audit.AuditPanic()
	panic("Test panic")
	// Can't assert anything here or golint gets angry
}

func TestAuditObject(t *testing.T) {
	t.Parallel()

	log := NewMock()

	// Test a simple object
	log.AuditObject("Prefix", "String")
	if len(log.GetAllMatching("[AUDIT]")) != 1 {
		t.Errorf("Failed to audit log simple object")
	}

	// Test a system object
	log.Clear()
	log.AuditObject("Prefix", t)
	if len(log.GetAllMatching("[AUDIT]")) != 1 {
		t.Errorf("Failed to audit log system object")
	}

	// Test a complex object
	log.Clear()
	type validObj struct {
		A string
		B string
	}
	var valid = validObj{A: "B", B: "C"}
	log.AuditObject("Prefix", valid)
	if len(log.GetAllMatching("[AUDIT]")) != 1 {
		t.Errorf("Failed to audit log complex object")
	}

	// Test logging an unserializable object
	log.Clear()
	type invalidObj struct {
		A chan string
	}

	var invalid = invalidObj{A: make(chan string)}
	log.AuditObject("Prefix", invalid)
	if len(log.GetAllMatching("[AUDIT]")) != 1 {
		t.Errorf("Failed to audit log unserializable object %v", log.GetAllMatching("[AUDIT]"))
	}
}

func TestTransmission(t *testing.T) {
	t.Parallel()

	l, err := newUDPListener("127.0.0.1:0")
	test.AssertNotError(t, err, "Failed to open log server")
	defer l.Close()

	fmt.Printf("Going to %s\n", l.LocalAddr().String())
	writer, err := syslog.Dial("udp", l.LocalAddr().String(), syslog.LOG_INFO|syslog.LOG_LOCAL0, "")
	test.AssertNotError(t, err, "Failed to find connect to log server")

	impl, err := New(writer, stdoutLevel)
	test.AssertNotError(t, err, "Failed to construct audit logger")

	data := make([]byte, 128)

	impl.AuditNotice("audit-logger_test.go: audit-notice")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.AuditErr(errors.New("audit-logger_test.go: audit-err"))
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.Crit("audit-logger_test.go: critical")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.Debug("audit-logger_test.go: debug")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.Emerg("audit-logger_test.go: emerg")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.Err("audit-logger_test.go: err")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.Info("audit-logger_test.go: info")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")

	impl.Warning("audit-logger_test.go: warning")
	_, _, err = l.ReadFrom(data)
	test.AssertNotError(t, err, "Failed to find packet")
}

func newUDPListener(addr string) (*net.UDPConn, error) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	l.SetDeadline(time.Now().Add(100 * time.Millisecond))
	l.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	l.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	return l.(*net.UDPConn), nil
}
