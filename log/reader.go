package log

import (
	"bufio"
	logger "log"
	"os"
	"path"
	"strings"
	"time"

	git "github.com/dlresende/logit/git"
)

type Event struct {
	Level, Message string
	When           time.Time
}

func Read(parser Parser, logFile *os.File, repository *git.Repository) {
	scanner := bufio.NewScanner(logFile)
	scanner.Split(parser.ChopEvent)

	for scanner.Scan() {
		logEntry := scanner.Text()
		logEvent := parser.Parse(logEntry)
		commitTitle := logEvent.Message[:strings.Index(logEvent.Message, "\n")]
		commitDescription := logEvent.Level + "\n" + logEvent.Message
		author := path.Base(logFile.Name())[:15]
		branch := path.Base(logFile.Name())
		commitMessage := commitTitle + "\n\n" + commitDescription
		repository.Commit(commitMessage, author, branch, logEvent.When)
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
}

type Parser interface {
	ChopEvent(data []byte, atEOF bool) (advance int, token []byte, err error)
	Parse(logEvent string) Event
}
