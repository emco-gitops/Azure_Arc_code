package main

import (
	"fmt"

	rawgitlab "github.com/xanzy/go-gitlab"
)

func main() {
	// Create a new client
	gitLabToken := "glpat-VsaWds-rWtbx6eM5ejBm"
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
