package main

import (
	logger "log"
	"os"

	"github.com/dlresende/logit/git"
	"github.com/dlresende/logit/log"
	"github.com/dlresende/logit/log/rabbitmq"
)

func main() {
	logFile, err := os.Open(os.Args[1])
	defer logFile.Close()
	if err != nil {
		logger.Fatal(err)
	}

	repository, err := git.Init("/tmp/logit")
	if err != nil {
		logger.Fatal(err)
	}

	log.Read(rabbitmqparser.NewRabbitMQLogParser(), logFile, repository)
}
