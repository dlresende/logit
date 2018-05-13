package main

import (
	"bufio"
	"log"
	"os"
	"time"

	l "logit/log"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func main() {
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(l.ChopLogEvent)
	for scanner.Scan() {
		// fmt.Println("NEW EVENT")
		logEventStr := scanner.Text()
		// fmt.Println(logEventStr)
		logEvent := l.Parse(logEventStr)
		commit(logEvent.Level+"\n\n"+logEvent.Message, "rabbit", logEvent.When)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func commit(message, branch string, when time.Time) {
	fs := osfs.New("/tmp/logit")
	var err error
	dot, err := fs.Chroot(".git")
	if err != nil {
		log.Fatal(err)
	}
	storer, err := filesystem.NewStorage(dot)
	if err != nil {
		log.Fatal(err)
	}
	var repo *git.Repository
	var workTree *git.Worktree
	var signature *object.Signature
	var options *git.CommitOptions
	var firstCommit plumbing.Hash

	log.Println("Trying to open git repo")
	repo, err = git.Open(storer, fs)
	if err != nil {
		log.Println("Git repo does not exist. Creating new one.")
		repo, err = git.Init(storer, fs)
		if err != nil {
			log.Fatal(err)
		}

		workTree, err = repo.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		signature = &object.Signature{Name: "logit", Email: "logit"}
		options = &git.CommitOptions{All: false, Author: signature, Committer: signature}
		log.Println("Creating first commit")
		firstCommit, err = workTree.Commit("Initial commit", options)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Looking for existing branch")
	_, err = repo.Branch(branch)
	if err != nil {
		log.Println("Branch doesn't exist. Creating new one.")
		err = repo.CreateBranch(&config.Branch{Name: "rabbit", Merge: "refs/heads/rabbit"})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Creating reference for branch")
		h := plumbing.NewHashReference("refs/heads/rabbit", firstCommit)
		if err = storer.SetReference(h); err != nil {
			log.Fatal(err)
		}
	}

	workTree, err = repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Checking branch out")
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/rabbit", Create: false, Force: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	signature = &object.Signature{Name: "logit", Email: "logit", When: when}
	options = &git.CommitOptions{All: false, Author: signature, Committer: signature}
	_, err = workTree.Commit(message, options)
	if err != nil {
		log.Fatal(err)
	}
}
