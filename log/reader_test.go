package log_test

import (
	. "logit/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
		advance, token, err := ChopLogEvent([]byte(logFileContent), false)

		Expect(string(token)).To(Equal(`=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/
`))
		Expect(advance).To(Equal(187))
		Expect(err).To(BeNil())
	})

	It("Should chop multi-line RabbitMQ log events starting with new lines", func() {
		logFileContent := `

=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/
`
		advance, token, err := ChopLogEvent([]byte(logFileContent), false)

		Expect(string(token)).To(Equal(`=INFO REPORT==== 5-May-2018::16:26:52 ===
Starting RabbitMQ 3.6.15 on Erlang 19.3.6.4
Copyright (C) 2007-2018 Pivotal Software, Inc.
Licensed under the MPL.  See http://www.rabbitmq.com/
`))
		Expect(advance).To(Equal(189))
		Expect(err).To(BeNil())
	})

	// https://golang.org/pkg/bufio/#SplitFunc
	It("Should ask for more text when multi-line RabbitMQ log event is not complete", func() {
		logFileContent := `=INFO REPORT==== 5-May-2018::16:26:52 ===
`
		advance, token, err := ChopLogEvent([]byte(logFileContent), false)

		Expect(token).To(BeNil())
		Expect(advance).To(Equal(0))
		Expect(err).To(BeNil())
	})

})
