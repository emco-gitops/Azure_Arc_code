package main

import (
	"context"
	"emcogithub"
	"fmt"
	"os"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

func main() {

	repoName := "Azure-test-repo-three"
	// path := "arc-k8s-demo/namespaces/team-b.yaml"
	userName := "chitti-intel"
	githubToken := os.Getenv("GITTOKEN")
	// commitMessage := "Deleting File"
	// Create a new client
	ctx := context.Background()

	// client := emcogithub.CreateRawClient(ctx, githubToken)

	// Delete file from repoS
	// emcogithub.DeleteFile(ctx, *client, userName, repoName, path, commitMessage)

	c, err := emcogithub.CreateClient(githubToken)
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
	files = emcogithub.Add("Test-file", "Hi I am a test file", files)

	response := emcogithub.CommitFiles(ctx, c, userName, repoName, "main", "New Commit", files)
	if response != nil {
		fmt.Println(response)
	}
}
