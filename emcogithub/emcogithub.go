package emcogithub

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
	rawgithub "github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

const (
	githubDomain = "github.com"
)

type commitType struct {
	k rawgithub.Commit
	c *gitprovider.CommitClient
}

/*
	Function to create gitprovider githubClient
	params : github token
	return : gitprovider github client, error
*/
func CreateClient(githubToken string) (gitprovider.Client, error) {
	c, err := github.NewClient(gitprovider.WithOAuth2Token(githubToken), gitprovider.WithDestructiveAPICalls(true))
	if err != nil {
		return nil, err
	}
	return c, nil
}

/*
	Function to create raw githubClient
	params : github token
	return : raw gitprovider github client
*/

func CreateRawClient(ctx context.Context, githubToken string) *rawgithub.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := rawgithub.NewClient(tc)

	return client
}

/*
	Function to create a new Repo in github
	params : context, github client, Repository Name, User Name, description
	return : nil
*/

func CreateRepo(ctx context.Context, c gitprovider.Client, repoName string, userName string, desc string) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)

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
		return err
	}
	fmt.Println(userRepo)

	return nil
}

/*
iterate through all the files and create a file array with path and contents
params : path string
return : []gitprovider.CommitFile
*/
func Iterate(path string) []gitprovider.CommitFile {
	files := []gitprovider.CommitFile{}
	length := len([]rune(path))
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		// check if the path is a directory
		if info.IsDir() == false {

			fmt.Printf("File Name: %s\n", info.Name())
			//Get the contents of the file
			content := GetContent(path)
			//get the relative path of the file
			relativePath := path[length+1:]
			fmt.Println(relativePath)
			// append to array of files
			files = append(files, gitprovider.CommitFile{Path: &relativePath, Content: &content})
		}
		return nil
	})
	fmt.Println(len(files))
	return files
}

/*
	Function to get the contents of the file
	params : filepath (string)
	return : contents (string)
*/
func GetContent(filePath string) string {

	// Open file on disk.
	f, _ := os.Open(filePath)

	// Read entire file into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	return string(content)
}

/*
	Function to commit multiple files to the github repo
	params : context, github client, User Name, Repo Name, Branch Name, Commit Message, files ([]gitprovider.CommitFile)
	return : nil
*/

func CommitFiles(ctx context.Context, c gitprovider.Client, userName string, repoName string, branch string, commitMessage string, files []gitprovider.CommitFile) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)

	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)
	//Commit file to this repo
	_, err = userRepo.Commits().Create(ctx, branch, commitMessage, files)

	if err != nil {
		return err
	}
	return nil
}

/*
	Function to delete particular file from github repo
	params : context, github client, user name, repo name, path, commit mesage
	returm : nil
*/
func DeleteFile(ctx context.Context, client rawgithub.Client, userName string, repoName string, path string, commitMessage string) error {

	// Get the file contents and extract sha
	fileContents, _, _, err := client.Repositories.GetContents(ctx, userName, repoName, path, &rawgithub.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}
	fmt.Println(*fileContents.SHA)

	sha := *fileContents.SHA
	repositoryContentsOptions := &rawgithub.RepositoryContentFileOptions{
		Message: &commitMessage,
		SHA:     &sha,
	}

	//Delete file
	deleteResponse, _, err := client.Repositories.DeleteFile(ctx, userName, repoName, path, repositoryContentsOptions)
	if err != nil {
		return err
	}
	fmt.Println(deleteResponse)

	return nil
}

/*
	Function to delete repo
	params : context, gitprovider client , user name, repo name
	return : nil
*/
func DeleteRepo(ctx context.Context, c gitprovider.Client, userName string, repoName string) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)
	// get thels reference of the repo to be deleted
	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)

	if err != nil {
		return err
	}
	//delete repo
	fmt.Println(userRepo.Delete(ctx))

	return nil
}

/*
	Internal function to create a repo refercnce
	params : user name, repo name
	return : repo reference
*/

func getRepoRef(userName string, repoName string) gitprovider.UserRepositoryRef {
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

	return userRepoRef
}

/*
	function to create a commit
	params :
	return :
*/

func CreateCommit(ctx context.Context, client rawgithub.Client, c gitprovider.Client, userName string, repoName string, branch string, message string, treeEntries []*rawgithub.TreeEntry) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)

	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)
	//Commit file to this repo

	commits, err := userRepo.Commits().ListPage(ctx, branch, 1, 0)
	if err != nil {
		return err
	}

	latestCommitTreeSHA := commits[0].Get().TreeSha

	tree, _, err := client.Git.CreateTree(ctx, userRepoRef.GetIdentity(), userRepoRef.GetRepository(), latestCommitTreeSHA, treeEntries)
	if err != nil {
		return err
	}

	latestCommitSHA := commits[0].Get().Sha

	nCommit, _, err := client.Git.CreateCommit(ctx, userRepoRef.GetIdentity(), userRepoRef.GetRepository(), &rawgithub.Commit{
		Message: &message,
		Tree:    tree,
		Parents: []*rawgithub.Commit{
			&rawgithub.Commit{
				SHA: &latestCommitSHA,
			},
		},
	})
	if err != nil {
		return err
	}

	ref := "refs/heads/" + branch
	ghRef := &rawgithub.Reference{
		Ref: &ref,
		Object: &rawgithub.GitObject{
			SHA: nCommit.SHA,
		},
	}

	if _, _, err := client.Git.UpdateRef(ctx, userRepoRef.GetIdentity(), userRepoRef.GetRepository(), ghRef, true); err != nil {
		return err
	}

	return nil
}

/*
	Function to Add file to the commit
	params : path , content, files (gitprovider commitfile array)
	return : files (gitprovider commitfile array)
*/
func Add(path string, content string, files []gitprovider.CommitFile) []gitprovider.CommitFile {
	files = append(files, gitprovider.CommitFile{
		Path:    &path,
		Content: &content,
	})

	return files
}

/*
	Function to Delete file from the commit
	params : path, files (gitprovider commitfile array)
	return : files (gitprovider commitfile array)
*/
func Delete(path string, files []gitprovider.CommitFile) []gitprovider.CommitFile {
	files = append(files, gitprovider.CommitFile{
		Path:    &path,
		Content: nil,
	})

	return files
}
