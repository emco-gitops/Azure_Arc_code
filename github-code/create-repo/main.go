package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
         log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

const (
	githubDomain = "github.com"
)

func main() {
	// Create a new client
	ctx := context.Background()
	githubToken := os.Getenv("GITTOKEN")
	userName := "chitti-intel"
	repoName := "Azure-test-repo-four"
	desc := "This repo contains azure arc and git golang code"
	c, err := github.NewClient(gitprovider.WithOAuth2Token(githubToken))

	if err != nil {
		fmt.Println(err)
	}

	createRepo(ctx, c, repoName, userName, desc)

}

/*
	Function to create a new Repo in github
	params : context, github client, Repository Name, User Name, description
	return : nil
*/

func createRepo(ctx context.Context, c gitprovider.Client, repoName string, userName string, desc string) {

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

	// Create repoinfo reference
	userRepoInfo := gitprovider.RepositoryInfo{
		Description: &desc,
		Visibility:  gitprovider.RepositoryVisibilityVar(gitprovider.RepositoryVisibilityPublic),
	}

	// Check that the repository doesn't exist
	//_, err = c.UserRepositories().Get(ctx, userRepoRef)
	//Expect(errors.Is(err, gitprovider.ErrNotFound)).To(BeTrue())

	// Create the repository
	userRepo, err := c.UserRepositories().Create(ctx, userRepoRef, userRepoInfo, &gitprovider.RepositoryCreateOptions{
		AutoInit:        gitprovider.BoolVar(true),
		LicenseTemplate: gitprovider.LicenseTemplateVar(gitprovider.LicenseTemplateApache2),
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userRepo)
        log.Info("Repo Created",log.Fields{})
}
