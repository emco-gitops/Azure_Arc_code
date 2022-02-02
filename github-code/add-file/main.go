package main

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
)

const (
	githubDomain = "github.com"
)

func main() {

	ctx := context.Background()
	githubToken := os.Getenv("GITTOKEN")
	userName := "chitti-intel"
	repoName := "Azure-test-repo"
	branch := "main"
	commitMessage := "App files added"
	// Create a new client
	c, err := github.NewClient(gitprovider.WithOAuth2Token(githubToken))
	if err != nil {
		fmt.Println(err)
	}
	// get the files array
	files := iterate("/home/ubuntu/go_projects/src/new-k8s-test-code")
	// commit the files to the github repo
	commitFiles(ctx, c, userName, repoName, branch, commitMessage, files)
}

/*
iterate through all the files and create a file array with path and contents
params : path string
return : []gitprovider.CommitFile
*/
func iterate(path string) []gitprovider.CommitFile {
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
			content := getContent(path)
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
func getContent(filePath string) string {

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

func commitFiles(ctx context.Context, c gitprovider.Client, userName string, repoName string, branch string, commitMessage string, files []gitprovider.CommitFile) {

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

	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)
	//Commit file to this repo
	_, err = userRepo.Commits().Create(ctx, branch, commitMessage, files)

	if err != nil {
		fmt.Println("Error in CommitFiles func")
		fmt.Println(err)
	}

}
