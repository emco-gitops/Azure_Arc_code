package main

import (
	"context"
	"emcogit"
	"fmt"
	"os"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

func convertToCommitFile(ref interface{}) []gitprovider.CommitFile {
	var exists bool
	switch ref.(type) {
	case []gitprovider.CommitFile:
		exists = true
	default:
		exists = false
	}
	var rf []gitprovider.CommitFile
	// Create rf is doesn't exist
	if !exists {
		rf = []gitprovider.CommitFile{}
	} else {
		rf = ref.([]gitprovider.CommitFile)
	}
	return rf
}

func main() {

	repoName := "git-go-test-five"
	//path := "arc-k8s-demo"
	userName := "chitti-intel"
	gitType := "github"
	githubToken := os.Getenv("GITTOKEN")
	// commitMessage := "Deleting File"
	// Create a new client
	ctx := context.Background()

	// client := emcogithub.CreateRawClient(ctx, githubToken)

	// Delete file from repoS
	// emcogithub.DeleteFile(ctx, *client, userName, repoName, path, commitMessage)

	c, err := emcogit.CreateClient(githubToken, gitType)
	if err != nil {
		fmt.Println(err)
	}

	// response := emcogithub.CreateRepo(ctx, c, repoName, userName, "Test Repo New")
	// if response != nil {
	// 	fmt.Println(response)
	// }
	// emcogithub.DeleteRepo(ctx, c, userName, repoName)
	files := []gitprovider.CommitFile{}
	// files = emcogithub.Iterate("/home/ubuntu/go_projects/src/new-k8s-test-code/arc-k8s-demo/cluster-apps")
	// response := emcogithub.CommitFiles(ctx, c, userName, repoName, "main", "New app files added", files)

	// content1 := "Hi I am test file 2"
	// treeEntries := make([]*rawgithub.TreeEntry, 0)
	// treeEntries = append(treeEntries, &rawgithub.TreeEntry{
	// 	Path: &path,
	// 	Mode: &githubNewFileMode,
	// })

	// files := []gitprovider.CommitFile{}
	//files = emcogithub.Delete("Test-file2", files)
	files = emcogit.Add("test-repo/Test-file-feb7", "Hi I am a test file", files, gitType).([]gitprovider.CommitFile)
	// files = emcogit.Delete(path, files, gitType)

	response := emcogit.CommitFiles(ctx, c, userName, repoName, "main", "New Commit", files, gitType)
	if response != nil {
		fmt.Println(response)
	}
}
