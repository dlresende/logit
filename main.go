package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"time"

	l "bitbucket.org/dlresende/logit/log"

	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
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
		logEventStr := scanner.Text()
		logEvent := l.Parse(logEventStr)
		commit(logEvent.Level+"\n\n"+logEvent.Message, path.Base(file.Name()), logEvent.When)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func commit(message, branch string, when time.Time) {
	gitDir := "/tmp/logit"
	fs := osfs.New(gitDir)
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

	log.Printf("Opening git repo %v...\n", gitDir)
	repo, err = git.Open(storer, fs)
	if err != nil {
		log.Println("Git repo does not exist. Creating new one...")
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
		log.Println("Creating first commit...")
		_, err = workTree.Commit("Initial commit", options)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Looking for branch %s...\n", branch)
	_, err = repo.Branch(branch)
	if err != nil {

		workTree, err = repo.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Checking out branch %v...\n", "master")
		err = workTree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.Master, Create: false, Force: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		var masterHead *plumbing.Reference
		masterHead, err = repo.Head()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Branch doesn't exist. Creating new one...")
		err = repo.CreateBranch(&config.Branch{Name: branch, Merge: plumbing.ReferenceName("refs/heads/" + branch)})
		if err != nil {
			log.Fatal(err)
		}

		_, err = repo.Branch(branch)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Creating reference for branch %v...\n", branch)
		h := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+branch), masterHead.Hash())
		if err = storer.SetReference(h); err != nil {
			log.Fatal(err)
		}
	}

	workTree, err = repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Checking out branch %v...\n", branch)
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch), Create: false, Force: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	author := branch[:15]
	signature = &object.Signature{Name: author, Email: author, When: when}
	options = &git.CommitOptions{All: false, Author: signature, Committer: signature}
	_, err = workTree.Commit(message, options)
	if err != nil {
		log.Fatal(err)
	}
}
