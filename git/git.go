package git

import (
	"log"
	"time"

	"gopkg.in/src-d/go-billy.v4/osfs"
	goGit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type Repository struct {
	gitRepo *goGit.Repository
}

func Init(gitDir string) (*Repository, error) {
	fs := osfs.New(gitDir)

	var err error
	dot, err := fs.Chroot(".git")
	if err != nil {
		return nil, err
	}

	storer, err := filesystem.NewStorage(dot)
	if err != nil {
		return nil, err
	}

	var gitRepo *goGit.Repository
	var workTree *goGit.Worktree
	var signature *object.Signature
	var options *goGit.CommitOptions

	gitRepo, err = goGit.Open(storer, fs)
	if err != nil {
		log.Printf("Git repo does not exist. Creating new one at %v...", gitDir)
		gitRepo, err = goGit.Init(storer, fs)
		if err != nil {
			return nil, err
		}

		workTree, err = gitRepo.Worktree()
		if err != nil {
			return nil, err
		}

		signature = &object.Signature{Name: "logit", Email: "logit"}
		options = &goGit.CommitOptions{All: false, Author: signature, Committer: signature}
		log.Println("Creating first commit...")
		_, err = workTree.Commit("Initial commit", options)
		if err != nil {
			return nil, err
		}
	}

	return &Repository{gitRepo}, nil
}

func (repo *Repository) Commit(message, author, branch string, when time.Time) {
	log.Printf("Looking for branch %s...\n", branch)
	var err error
	var workTree *goGit.Worktree
	var signature *object.Signature
	var options *goGit.CommitOptions
	_, err = repo.gitRepo.Branch(branch)
	if err != nil {

		workTree, err = repo.gitRepo.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Checking out branch %v...\n", "master")
		err = workTree.Checkout(&goGit.CheckoutOptions{
			Branch: plumbing.Master, Create: false, Force: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		var masterHead *plumbing.Reference
		masterHead, err = repo.gitRepo.Head()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Branch doesn't exist. Creating new one...")
		err = repo.gitRepo.CreateBranch(&config.Branch{Name: branch, Merge: plumbing.ReferenceName("refs/heads/" + branch)})
		if err != nil {
			log.Fatal(err)
		}

		_, err = repo.gitRepo.Branch(branch)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Creating reference for branch %v...\n", branch)
		h := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+branch), masterHead.Hash())
		if err = repo.gitRepo.Storer.SetReference(h); err != nil {
			log.Fatal(err)
		}
	}

	workTree, err = repo.gitRepo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Checking out branch %v...\n", branch)
	err = workTree.Checkout(&goGit.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch), Create: false, Force: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	signature = &object.Signature{Name: author, Email: author, When: when}
	options = &goGit.CommitOptions{All: false, Author: signature, Committer: signature}
	_, err = workTree.Commit(message, options)
	if err != nil {
		log.Fatal(err)
	}
}
