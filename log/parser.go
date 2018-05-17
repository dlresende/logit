package log

import (
	"log"
	"regexp"
	"strings"
	"time"
)

type Event struct {
	Level, Message string
	When           time.Time
}

func Parse(logEvent string) Event {
	rabbitmqPattern := regexp.MustCompile(`^(?:[[:space:]]|[A-Za-z]|\*)*(=(.*)====\s(.*)\s===\n((?:.+\n)+)\n)`)

	if rabbitmqPattern.MatchString(logEvent) {
		logData := rabbitmqPattern.FindStringSubmatch(logEvent)
		logEventTime, err := time.Parse("2-Jan-2006::15:4:5", logData[3])
		if err != nil {
			log.Fatalln(err)
		}

		commitTitle := logData[4][:strings.Index(logData[4], "\n")]
		return Event{commitTitle, logData[1], logEventTime}
	}

	return Event{}
}
