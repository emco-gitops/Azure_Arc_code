package main

import (
	"context"
	"fmt"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

const (
	githubDomain = "github.com"
)

func main() {
	// Create a new client
	ctx := context.Background()
	githubToken := "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"
	userName := "chitti-intel"
	repoName := "New-test-repo"
	c, err := github.NewClient(gitprovider.WithOAuth2Token(githubToken), gitprovider.WithDestructiveAPICalls(true))
	if err != nil {
		fmt.Println(err)
	}

	deleteRepo(ctx, c, userName, repoName)
}

func deleteRepo(ctx context.Context, c gitprovider.Client, userName string, repoName string) {

	// Create the user reference
	userRef := gitprovider.UserRef{
		Domain:    githubDomain,
		UserLogin: userName,
	}

	// Create the repo reference
	userRepoRef := gitprovider.UserRepositoryRef{
		UserRef:        userRef,
		RepositoryName: repoName,
	}
	// get the reference of the repo to be deleted
	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)

	if err != nil {
		fmt.Println(err)
	}
	//delete repo
	fmt.Println(userRepo.Delete(ctx))
}
