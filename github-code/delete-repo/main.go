package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

const (
	githubDomain = "github.com"
)

func main() {
	// Create a new client
	ctx := context.Background()
	githubToken := os.Getenv("GITTOKEN")
	userName := "chitti-intel"
	repoName := "New-test-repoV3"
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
	err = userRepo.Delete(ctx)

        if err != nil {
		fmt.Println(err)
        }
	fmt.Println("Repo Deleted")
}
