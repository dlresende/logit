package rabbitmqparser_test

import (
	"time"

	"github.com/dlresende/logit/log"
	rabbitmq "github.com/dlresende/logit/log/rabbitmq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	It("should parse RabbitMQ log events", func() {
		rabbitmqEvent := `=WARNING REPORT==== 5-May-2018::16:27:16 ===
global: rabbit@6b44fa337952635b266d4cda4f1fb908 failed to connect to rabbit@8063b04b5406b66334951c641fe9f634

`
		rabbitmqParser := rabbitmq.NewRabbitMQLogParser()

		Expect(rabbitmqParser.Parse(rabbitmqEvent)).To(Equal(log.Event{
			Level:   "WARNING REPORT",
			Message: "global: rabbit@6b44fa337952635b266d4cda4f1fb908 failed to connect to rabbit@8063b04b5406b66334951c641fe9f634\n",
			When:    time.Date(2018, time.May, 5, 16, 27, 16, 0, time.UTC)}))
	})

	It("should parse RabbitMQ log events with multiple lines", func() {
		rabbitmqEvent := `=INFO REPORT==== 5-May-2018::16:26:53 ===
FHC read buffering:  OFF
FHC write buffering: ON

`
		rabbitmqParser := rabbitmq.NewRabbitMQLogParser()

		Expect(rabbitmqParser.Parse(rabbitmqEvent)).To(Equal(log.Event{
			Level:   "INFO REPORT",
			Message: "FHC read buffering:  OFF\nFHC write buffering: ON\n",
			When:    time.Date(2018, time.May, 5, 16, 26, 53, 0, time.UTC)}))
	})
})

var _ = Describe("Log/Reader", func() {
	It("Should chop multi-line RabbitMQ log events", func() {
		logFileContent := `=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/

=INFO REPORT==== 5-May-2018::16:26:52 ===
node           : rabbit@6b44fa337952635b266d4cda4f1fb908
home dir       : /var/vcap/store/rabbitmq
config file(s) : /var/vcap/jobs/rabbitmq-server/bin/../etc/rabbitmq.config
cookie hash    : 9XNMH4g9Js9lG9NYVrNgfw==
log            : /var/vcap/sys/log/rabbitmq-server/rabbit@6b44fa337952635b266d4cda4f1fb908.log
sasl log       : /var/vcap/sys/log/rabbitmq-server/rabbit@6b44fa337952635b266d4cda4f1fb908-sasl.log
database dir   : /var/vcap/store/rabbitmq/mnesia/db

`

		rabbitmqParser := rabbitmq.NewRabbitMQLogParser()

		advance, token, err := rabbitmqParser.ChopEvent([]byte(logFileContent), false)

		Expect(string(token)).To(Equal(`=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/

`))
		Expect(advance).To(Equal(188))
		Expect(err).To(BeNil())
	})

	It("Should chop multi-line RabbitMQ log events starting with new lines", func() {
		logFileContent := `

=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/

`
		rabbitmqParser := rabbitmq.NewRabbitMQLogParser()

		advance, token, err := rabbitmqParser.ChopEvent([]byte(logFileContent), false)

		Expect(string(token)).To(Equal(`=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/

`))
		Expect(advance).To(Equal(190))
		Expect(err).To(BeNil())
	})

	// https://golang.org/pkg/bufio/#SplitFunc
	It("Should ask for more text when multi-line RabbitMQ log event is not complete", func() {
		logFileContent := `=INFO REPORT==== 5-May-2018::16:26:52 ===
`

		rabbitmqParser := rabbitmq.NewRabbitMQLogParser()

		advance, token, err := rabbitmqParser.ChopEvent([]byte(logFileContent), false)

		Expect(token).To(BeNil())
		Expect(advance).To(Equal(0))
		Expect(err).To(BeNil())
	})

	It("Should ignore text that does not match a log event", func() {
		logFileContent := `
**********************************************************
*** Publishers will be blocked until this alarm clears ***
**********************************************************

=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/

`

		rabbitmqParser := rabbitmq.NewRabbitMQLogParser()

		advance, token, err := rabbitmqParser.ChopEvent([]byte(logFileContent), false)

		Expect(string(token)).To(Equal(`=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/

`))
		Expect(advance).To(Equal(367))
		Expect(err).To(BeNil())
	})
})
