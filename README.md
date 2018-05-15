# LogIt

LogIt is a tool that parses log events and commit each event to Git. 

Debugging distributed systems is hard, mainly because it is difficult to sort
log events into a timeline that would allow them to be easily correlated.

LogIt parses all the log events from a file and commit them to a specific
branch in Git. LogIt leverages Git tools to show log events in a timeline
(branches side by side), making it easier to draw cause effect relations
between log events.

## Usage

```shell
$ go get gopkg.in/src-d/go-git.v4

$ git run main.go rabbitmq@host1.log
$ git run main.go rabbitmq@host2.log
$ git run main.go rabbitmq@host3.log

$ GIT_DIR=/tmp/logit tig --all
```
