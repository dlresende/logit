package main

import (
	"bufio"
	"log"
	"os"
	"path"

	git "github.com/dlresende/logit/git"
	l "github.com/dlresende/logit/log"
)

func main() {
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	repository, err := git.Init("/tmp/logit")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(l.ChopLogEvent)
	for scanner.Scan() {
		logEventStr := scanner.Text()
		logEvent := l.Parse(logEventStr)
		repository.Commit(logEvent.Level+"\n\n"+logEvent.Message, path.Base(file.Name()), path.Base(file.Name())[:15], logEvent.When)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
