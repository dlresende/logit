package log

import (
	"fmt"
	"regexp"
	"time"
)

type LogEvent struct {
	Level, Message string
	When           time.Time
}

func Parse(logEvent string) LogEvent {
	fmt.Printf(logEvent)
	rabbitmqPattern := regexp.MustCompile(`(?s)^=(.*)====\s(.*)\s===\n(.*)`)
	if rabbitmqPattern.MatchString(logEvent) {
		logData := rabbitmqPattern.FindStringSubmatch(logEvent)
		logEventTime, _ := time.Parse("2-Jan-2006::15:4:5", logData[2])
		return LogEvent{logData[1], logData[3], logEventTime}
	}
	return LogEvent{}
}
