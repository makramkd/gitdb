package gitdb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	git "gopkg.in/src-d/go-git.v4"
)

// DB is an interface to the git db
type DB interface {
	// Save the given data to the database, creating a new commit with a random msg
	Save(data []byte, filename string) error

	// Save the given data to the database, creating a new commit wth the given commit msg
	SaveWithMessage(data []byte, filename string, commitMessage string) error

	// Open the db at the given path
	Open(path string) error

	// Read a file from the db and return it's contents
	Read(filename string) ([]byte, error)

	// Version returns the version of the database (hash)
	Version() string
}

var instance DB
var once sync.Once

// Errors that could be encountered
var (
	ErrorNoDB = errors.New("DB not initialized")
)

// TODO: use interfaces instead so we can mock
type db struct {
	m    *sync.Mutex
	repo *git.Repository
	path string
}

// GetInstance gets an instance of the git db
func GetInstance() DB {
	once.Do(func() {
		instance = &db{
			m: &sync.Mutex{},
		}
	})
	return instance
}

func (d *db) Save(data []byte, filename string) error {
	return d.save(data, filename, "")
}

func (d *db) SaveWithMessage(data []byte, filename string, commitMessage string) error {
	return d.save(data, filename, commitMessage)
}

func (d *db) save(data []byte, filename string, commitMessage string) error {
	// synchronize
	d.m.Lock()
	defer d.m.Unlock()
	// check if repo is available
	if d.repo == nil {
		return ErrorNoDB
	}
	// (over)write file to path
	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", d.path, filename), data, os.ModePerm); err != nil {
		return err
	}
	// commit with random msg
	worktree, err := d.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = worktree.Add(filename)
	if err != nil {
		return err
	}
	m := commitMessage
	if commitMessage == "" {
		m = "some message"
	}
	_, err = worktree.Commit(m, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "gitdb",
			Email: "gitdb@github.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *db) Open(path string) error {
	// try to init, if error returned, then it already exists
	if repo, err := git.PlainInit(path, false); err != nil {
		// already exists, simply open
		repo, err := git.PlainOpen(path)
		if err != nil {
			return err
		}
		d.repo = repo
		d.path = path
	} else {
		d.repo = repo
		d.path = path
	}
	return nil
}

func (d *db) Read(filename string) ([]byte, error) {
	// naive file read
	d.m.Lock()
	defer d.m.Unlock()

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (d *db) Version() string {
	ref, err := d.repo.Head()
	if err != nil {
		return "error"
	}
	return ref.String()
}
