package log

import (
	golog "log"
	"regexp"
)

func ChopLogEvent(data []byte, atEOF bool) (advance int, token []byte, err error) {

	rabbitmqPattern := regexp.MustCompile(`^(?:[[:space:]]|[A-Za-z]|\*)*(=(.*)====\s(.*)\s===\n((?:.+\n)+)\n)`)
	text := string(data)

	golog.Printf("Data chunk:\n%v\n", text)

	if rabbitmqPattern.MatchString(text) {
		groups := rabbitmqPattern.FindStringSubmatch(text)
		logEvent := groups[1]
		golog.Printf("Found match:\n%v\n", logEvent)
		return len([]byte(groups[0])), []byte(logEvent), nil
	}

	golog.Println("No match found")
	return 0, nil, nil
}
