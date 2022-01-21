package main

import (
	"emcogitlab"
	"fmt"

	rawgitlab "github.com/xanzy/go-gitlab"
)

func main() {
	// //Create a new client
	// ctx := context.Background()
	// gitLabToken := "glpat-VsaWds-rWtbx6eM5ejBm"
	// repoName := "Test-Repo"
	// userName := "chitti-intel"
	// c, err := emcogitlab.CreateClient(gitLabToken)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // create a new repo
	// fmt.Println(emcogitlab.CreateRepo(ctx, c, repoName, userName, "Repo to test out go gitlab code"))

	// // adding files to repo
	// files := emcogitlab.Iterate("/home/ubuntu/go_projects/src/k8s-test-code/arc-k8s-demo")
	// fmt.Println(emcogitlab.CommitFiles(ctx, c, userName, repoName, "main", "Initial commit", files))

	// Create a new client
	repoName := "Test-Repo"
	userName := "chitti-intel"
	gitLabToken := "glpat-VsaWds-rWtbx6eM5ejBm"
	c, err := emcogitlab.CreateRawClient(gitLabToken)
	if err != nil {
		fmt.Println(err)
	}
	// Delete an existing file
	// resp := emcogitlab.DeleteFile(*c, userName, repoName, "main", "README.md", "Deleteing file")
	// if resp != nil {
	// 	fmt.Println(resp)
	// }

	// //Delete Repo
	// fmt.Println(emcogitlab.DeleteRepo(ctx, c, userName, repoName))
	// content := emcogitlab.GetContent("/home/ubuntu/go_projects/src/k8s-test-code/modified-test-file")
	// resp := emcogitlab.UpdateFile(*c, userName, repoName, "main", "README.md", "Updated Readme", content)
	// if resp != nil {
	// 	fmt.Println(resp)
	// }

	// Create a new commit with new files
	files := []*rawgitlab.CommitActionOptions{}
	files = emcogitlab.Delete("Test-file-3", files)
	files = emcogitlab.Update("Hi there I am file 4 updated", "Test-file-4", files)
	resp := emcogitlab.CreateCommit(*c, userName, repoName, "main", "New commit", files)
	if resp != nil {
		fmt.Println(resp)
	}
}
