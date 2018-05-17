package log_test

import (
	"time"

	. "github.com/dlresende/logit/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	It("should parse RabbitMQ log events", func() {
		rabbitmqEvent := `=WARNING REPORT==== 5-May-2018::16:27:16 ===
global: rabbit@6b44fa337952635b266d4cda4f1fb908 failed to connect to rabbit@8063b04b5406b66334951c641fe9f634

`
		Expect(Parse(rabbitmqEvent)).To(Equal(Event{
			Level:   "global: rabbit@6b44fa337952635b266d4cda4f1fb908 failed to connect to rabbit@8063b04b5406b66334951c641fe9f634",
			Message: "=WARNING REPORT==== 5-May-2018::16:27:16 ===\nglobal: rabbit@6b44fa337952635b266d4cda4f1fb908 failed to connect to rabbit@8063b04b5406b66334951c641fe9f634\n\n",
			When:    time.Date(2018, time.May, 5, 16, 27, 16, 0, time.UTC)}))
	})

	It("should parse RabbitMQ log events with multiple lines", func() {
		rabbitmqEvent := `=INFO REPORT==== 5-May-2018::16:26:53 ===
FHC read buffering:  OFF
FHC write buffering: ON

`
		Expect(Parse(rabbitmqEvent)).To(Equal(Event{
			Level:   "FHC read buffering:  OFF",
			Message: "=INFO REPORT==== 5-May-2018::16:26:53 ===\nFHC read buffering:  OFF\nFHC write buffering: ON\n\n",
			When:    time.Date(2018, time.May, 5, 16, 26, 53, 0, time.UTC)}))
	})
})
