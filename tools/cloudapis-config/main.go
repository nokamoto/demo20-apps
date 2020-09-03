package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nokamoto/demo20-apis/cloud/api"
)

const (
	gitURL      = "GIT_URL"
	gitRevision = "GIT_REVISION"
)

func env(env, value string) string {
	v := os.Getenv(env)
	if len(v) == 0 {
		return value
	}
	return v
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func ishex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s filename\n", os.Args[0])
		os.Exit(1)
	}
	url := env(gitURL, "https://github.com/nokamoto/demo20-apis.git")
	revision := env(gitRevision, "master")
	output := os.Args[1]

	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	assert(err)

	repo.Head()

	workTree, err := repo.Worktree()
	assert(err)

	var opt git.CheckoutOptions
	if ishex(revision) {
		opt.Hash = plumbing.NewHash(revision)
	} else {
		opt.Branch = plumbing.NewBranchReferenceName(revision)
	}

	assert(workTree.Checkout(&opt))

	head, err := repo.Head()
	assert(err)

	commit, err := repo.CommitObject(head.Hash())
	assert(err)

	tree, err := commit.Tree()
	assert(err)

	var files []*object.File
	err = tree.Files().ForEach(func(file *object.File) error {
		if strings.HasSuffix(file.Name, ".json") {
			files = append(files, file)
		}
		return nil
	})
	assert(err)

	cfg := api.AuthzConfig{
		Config: make(map[string]*api.Authz),
	}
	for _, file := range files {
		content, err := file.Contents()
		assert(err)

		var xs api.AuthzConfig
		assert(jsonpb.UnmarshalString(content, &xs))

		for k, v := range xs.GetConfig() {
			cfg.Config[k] = v
		}
	}

	m := jsonpb.Marshaler{
		Indent: "  ",
	}
	json, err := m.MarshalToString(&cfg)
	assert(err)
	assert(ioutil.WriteFile(output, []byte(json), 0644))
}
