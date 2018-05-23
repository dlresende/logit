package rabbitmqparser

import (
	logger "log"
	"regexp"
	"time"

	log "github.com/dlresende/logit/log"
)

const RabbitMQLogEvent = `^(?:[[:space:]]|[A-Za-z]|\*)*(=(.*)====\s(.*)\s===\n((?:.+\n)+)\n)`

type RabbitMQLogParser struct{}

func (p *RabbitMQLogParser) Parse(logEvent string) log.Event {
	rabbitmqPattern := regexp.MustCompile(RabbitMQLogEvent)

	if rabbitmqPattern.MatchString(logEvent) {
		logData := rabbitmqPattern.FindStringSubmatch(logEvent)
		logEventTime, err := time.Parse("2-Jan-2006::15:4:5", logData[3])
		if err != nil {
			logger.Fatalln(err)
		}

		return log.Event{Level: logData[2], Message: logData[4], When: logEventTime}
	}

	return log.Event{}
}

func (p *RabbitMQLogParser) ChopEvent(data []byte, atEOF bool) (advance int, token []byte, err error) {

	rabbitmqPattern := regexp.MustCompile(RabbitMQLogEvent)
	text := string(data)

	if rabbitmqPattern.MatchString(text) {
		groups := rabbitmqPattern.FindStringSubmatch(text)
		logEvent := groups[1]
		return len([]byte(groups[0])), []byte(logEvent), nil
	}

	return 0, nil, nil
}

func NewRabbitMQLogParser() *RabbitMQLogParser {
	return &RabbitMQLogParser{}
}
