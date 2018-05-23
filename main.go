package main

import (
	"bufio"
	logger "log"
	"os"
	"path"
	"strings"

	git "github.com/dlresende/logit/git"
	log "github.com/dlresende/logit/log"
)

func main() {
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	repository, err := git.Init("/tmp/logit")
	if err != nil {
		logger.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(log.ChopEvent)
	for scanner.Scan() {
		logEventStr := scanner.Text()
		logEvent := log.Parse(logEventStr)
		commitTitle := logEvent.Message[:strings.Index(logEvent.Message, "\n")]
		commitDescription := logEvent.Level + "\n" + logEvent.Message
		author := path.Base(file.Name())[:15]
		branch := path.Base(file.Name())
		repository.Commit(commitTitle+"\n\n"+commitDescription, author, branch, logEvent.When)
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
}
