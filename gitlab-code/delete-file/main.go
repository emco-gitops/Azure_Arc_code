package main

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

func main() {
	// Create a new client
	gitLabToken := os.Getenv("GITLABTOKEN")
	c, err := gitlab.NewClient(gitLabToken, nil)
	if err != nil {
		fmt.Println(err)
	}

	branch := "main"
	message := "Delete file"
	userName := "chitti-intel"
	repoName := "Test-Repo"
	path := "README.md"
	project := username + "/" + repoName
	opt := gitlab.DeleteFileOptions{&branch, nil, nil, nil, &message, nil}
	resp, err := c.RepositoryFiles.DeleteFile(project, path, &opt)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
