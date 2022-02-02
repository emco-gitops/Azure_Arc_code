package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fluxcd/go-git-providers/gitlab"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

const (
	gitLabDomain = "gitlab.com"
)

func main() {
	// Create a new client
	ctx := context.Background()
	gitLabToken := os.Getenv("GITLABTOKEN")
	c, err := gitlab.NewClient(gitLabToken, "")
	//checkErr(err)
	fmt.Println(err)
	// Get public information about the flux repository.
	userRef := gitprovider.UserRef{
		Domain:    gitLabDomain,
		UserLogin: "chitti-intel",
	}

	desc := "Test repo created using go-git-provider package"
	repoName := "Go-Git-Test-Repo-two"
	userRepoRef := gitprovider.UserRepositoryRef{
		UserRef:        userRef,
		RepositoryName: repoName,
	}
	userRepoInfo := gitprovider.RepositoryInfo{
		Description: &desc,
		Visibility:  gitprovider.RepositoryVisibilityVar(gitprovider.RepositoryVisibilityPublic),
	}

	// Check that the repository doesn't exist
	//_, err = c.UserRepositories().Get(ctx, userRepoRef)
	//Expect(errors.Is(err, gitprovider.ErrNotFound)).To(BeTrue())

	// Create it
	userRepo, err := c.UserRepositories().Create(ctx, userRepoRef, userRepoInfo, &gitprovider.RepositoryCreateOptions{
		AutoInit:        gitprovider.BoolVar(true),
		LicenseTemplate: gitprovider.LicenseTemplateVar(gitprovider.LicenseTemplateApache2),
	})
	//Expect(err).ToNot(HaveOccurred())

	fmt.Println(userRepo)
	fmt.Println(err)

}
