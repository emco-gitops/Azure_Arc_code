package main

import (
	"fmt"
	"os"

	rawgitlab "github.com/xanzy/go-gitlab"
)

func main() {
	// Create a new client
	gitLabToken := os.Getenv("GITLABTOKEN")
	c, err := rawgitlab.NewClient(gitLabToken, "")
	if err != nil {
		fmt.Println(err)
	}
	resp, err := c.RepositoryFiles.DeleteFile("Test-Repo", "README.md", nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
