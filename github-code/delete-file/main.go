package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

func main() {
	repoName := "New-test-repo"
	path := "arc-k8s-demo/namespaces/team-a.yaml"
	userName := "chitti-intel"
	githubToken := os.Getenv("GITTOKEN")
	commitMessage := "Deleting File"
	// Create a new client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Delete file from repoS
	deleteFile(ctx, *client, userName, repoName, path, commitMessage)
}

/*
	Function to delete particular file from github repo
	params : context, github client, user name, repo name, path, commit mesage
	returm : nil
*/
func deleteFile(ctx context.Context, client github.Client, userName string, repoName string, path string, commitMessage string) {

	// Get the file contents and extract sha
	fileContents, _, _, err := client.Repositories.GetContents(ctx, userName, repoName, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*fileContents.SHA)

	sha := *fileContents.SHA
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		SHA:     &sha,
	}

	//Delete file
	deleteResponse, _, err := client.Repositories.DeleteFile(ctx, userName, repoName, path, repositoryContentsOptions)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(deleteResponse)
}
