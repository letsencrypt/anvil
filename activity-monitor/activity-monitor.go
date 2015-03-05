// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/codegangsta/cli"
	"github.com/letsencrypt/boulder"
	"github.com/letsencrypt/boulder/analysisengine"
	"github.com/streadway/amqp"
	"log"
	"net/url"
	"os"
)

const (
	QueueName        = "Monitor"
	AmqpExchange     = "boulder"
	AmqpExchangeType = "topic"
	AmqpInternal     = false
	AmqpDurable      = false
	AmqpDeleteUnused = false
	AmqpExclusive    = false
	AmqpNoWait       = false
	AmqpNoLocal      = false
	AmqpAutoAck      = false
	AmqpMandatory    = false
	AmqpImmediate    = false
)

func startMonitor(AmqpUrl string, logger *boulder.JsonLogger) {

	ae := analysisengine.NewLoggingAnalysisEngine(logger)

	// For convenience at the broker, identifiy ourselves by hostname
	consumerTag, err := os.Hostname()
	if err != nil {
		log.Fatalf("Could not determine hostname")
	}

	conn, err := amqp.Dial(AmqpUrl)
	if err != nil {
		log.Fatalf("Could not connect to AMQP server: %s", err)
		return
	}

	rpcCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not start channel: %s", err)
		return
	}

	err = rpcCh.ExchangeDeclare(
		AmqpExchange,
		AmqpExchangeType,
		AmqpDurable,
		AmqpDeleteUnused,
		AmqpInternal,
		AmqpNoWait,
		nil)
	if err != nil {
		log.Fatalf("Could not declare exchange: %s", err)
		return
	}

	_, err = rpcCh.QueueDeclare(
		QueueName,
		AmqpDurable,
		AmqpDeleteUnused,
		AmqpExclusive,
		AmqpNoWait,
		nil)
	if err != nil {
		log.Fatalf("Could not declare queue: %s", err)
		return
	}

	err = rpcCh.QueueBind(
		QueueName,
		"#", //wildcard
		AmqpExchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("Could not bind queue: %s", err)
		return
	}

	deliveries, err := rpcCh.Consume(
		QueueName,
		consumerTag,
		AmqpAutoAck,
		AmqpExclusive,
		AmqpNoLocal,
		AmqpNoWait,
		nil)
	if err != nil {
		log.Fatalf("Could not subscribe to queue: %s", err)
		return
	}

	// Run forever.
	for d := range deliveries {
		// Pass each message to the Analysis Engine
		ae.ProcessMessage(d)
		// Only ack the delivery we actually handled (ackMultiple=false)
		const ackMultiple = false
		d.Ack(ackMultiple)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "activity-monitor"
	app.Usage = "Monitor Boulder's communications."
	app.Version = "0.0.0"

	// Specify AMQP Server
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "amqp",
			Value: "amqp://guest:guest@localhost:5672",
			Usage: "AMQP Broker String",
		},
		cli.StringFlag{
			Name:  "jsonlog",
			Usage: "JSON logging server and port (e.g., tcp://localhost:515)",
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "Enable debug logging to stdout",
		},
		cli.IntFlag{
			Name:  "level",
			Value: 4,
			Usage: "Minimum Level to log (0-7), 7=Debug",
		},
	}

	app.Action = func(c *cli.Context) {
		logger := boulder.NewJsonLogger("am")

		// Parse SysLog URL if one was provided
		if c.GlobalString("jsonlog") == "" {
			log.Println("No external logging server; defaulting to stdout.")
			logger.EnableStdOut(true)
		} else {
			syslogU, err := url.Parse(c.GlobalString("jsonlog"))
			if err != nil {
				log.Fatalf("Could not parse Syslog URL: %s", err)
				return
			}

			logger.SetEndpoint(syslogU.Scheme, syslogU.Host)
			err = logger.Connect()
			if err != nil {
				log.Fatalf("Could not open remote syslog: %s", err)
				return
			}

			logger.EnableStdOut(c.GlobalBool("stdout"))

		}

		logger.SetLevel(c.GlobalInt("level"))

		startMonitor(c.GlobalString("amqp"), logger)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Could not start: %s", err)
		return
	}
}
