package log

import (
	"regexp"
)

func ChopLogEvent(data []byte, atEOF bool) (advance int, token []byte, err error) {

	rabbitmqPattern := regexp.MustCompile(`^(?:[[:space:]]|[A-Za-z]|\*)*(=(.*)====\s(.*)\s===\n((?:.+\n)+)\n)`)
	text := string(data)

	if rabbitmqPattern.MatchString(text) {
		groups := rabbitmqPattern.FindStringSubmatch(text)
		logEvent := groups[1]
		return len([]byte(groups[0])), []byte(logEvent), nil
	}

	return 0, nil, nil
}
