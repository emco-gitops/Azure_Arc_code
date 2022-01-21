package emcogitlab

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fluxcd/go-git-providers/gitlab"
	"github.com/fluxcd/go-git-providers/gitprovider"
	rawgitlab "github.com/xanzy/go-gitlab"
)

const (
	gitlabDomain = "gitlab.com"
)

// create gitlab client
/*
	Function to create gitlab client
	params : gitlab token
	return : gitporvider client , error
*/
func CreateClient(gitLabToken string) (gitprovider.Client, error) {
	// Create a new client
	c, err := gitlab.NewClient(gitLabToken, "", gitprovider.WithDestructiveAPICalls(true))

	if err != nil {
		return nil, err
	}
	return c, nil
}

//create raw gitlab client
/*
	Function to create raw gitlab client
	params : gitlab token
	return : raw gitlab client (github.com/xanzy/go-gitlab), error
*/
func CreateRawClient(gitLabToken string) (*rawgitlab.Client, error) {
	c, err := rawgitlab.NewClient(gitLabToken, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
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

			// fmt.Printf("File Name: %s\n", info.Name())
			//Get the contents of the file
			content := GetContent(path)
			//get the relative path of the file
			relativePath := path[length+1:]
			// append to array of files
			files = append(files, gitprovider.CommitFile{Path: &relativePath, Content: &content})
		}
		return nil
	})

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
	Internal function to create a repo refercnce
	params : user name, repo name
	return : repo reference
*/

func getRepoRef(userName string, repoName string) gitprovider.UserRepositoryRef {
	// Create the user reference
	userRef := gitprovider.UserRef{
		Domain:    gitlabDomain,
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
	Function to delete repo
	params : context, gitprovider client , user name, repo name
	return : nil
*/
func DeleteRepo(ctx context.Context, c gitprovider.Client, userName string, repoName string) error {

	// create repo reference
	userRepoRef := getRepoRef(userName, repoName)
	// get the reference of the repo to be deleted
	userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)

	if err != nil {
		return err
	}
	//delete repo
	fmt.Println(userRepo.Delete(ctx))

	return nil
}

/*
	Function to delete file in repository
	params : rawgitlab client, username, repo name, branch, path, message
	return : error
*/

func DeleteFile(c rawgitlab.Client, userName string, repoName string, branch string, path string, message string) error {

	// get the project name
	project := userName + "/" + repoName

	// create the deleteoptions
	opt := rawgitlab.DeleteFileOptions{&branch, nil, nil, nil, &message, nil}
	_, err := c.RepositoryFiles.DeleteFile(project, path, &opt)
	if err != nil {
		return err
	}

	return nil
}

/*
	Function to update file in repository
	params : rawgitlab client, username, repo name, branch, path, message, content
	return : error
*/
func UpdateFile(c rawgitlab.Client, userName string, repoName string, branch string, path string, message string, content string) error {
	// get the project name
	project := userName + "/" + repoName

	//options to update the file
	opt := rawgitlab.UpdateFileOptions{&branch, nil, nil, nil, nil, &content, &message, nil}
	_, _, err := c.RepositoryFiles.UpdateFile(project, path, &opt)
	if err != nil {
		return err
	}

	return nil
}

/*
	Function to create a commit with multiple files and actions
	params : rawgitlab Client , username, repo name, branch, message, actions with files
	return : error
*/
func CreateCommit(c rawgitlab.Client, userName string, repoName string, branch string, message string, files []*rawgitlab.CommitActionOptions) error {
	// get the project name
	project := userName + "/" + repoName

	// options to create the commit
	opt := rawgitlab.CreateCommitOptions{&branch, &message, nil, nil, nil, files, nil, nil, nil, nil}

	_, _, err := c.Commits.CreateCommit(project, &opt, nil)
	if err != nil {
		return err
	}

	return nil
}

/*
	Add file to the commit (generic)
	params : file content, file action (update, add), files (array of type CommitActionOptions)
	return : error
*/

func AddFileToCommit(content string, action rawgitlab.FileActionValue, path string, files []*rawgitlab.CommitActionOptions) []*rawgitlab.CommitActionOptions {

	// append the new file to the files array
	files = append(files, &rawgitlab.CommitActionOptions{&action, &path, nil, &content, nil, nil, nil})

	return files
}

/*
	Add file (create) to the commit
	params : content, path, files array
	return : files array
*/
func Add(content string, path string, files []*rawgitlab.CommitActionOptions) []*rawgitlab.CommitActionOptions {
	var action rawgitlab.FileActionValue
	action = "create"
	files = append(files, &rawgitlab.CommitActionOptions{&action, &path, nil, &content, nil, nil, nil})

	return files
}

/*
	Add delete file to the commit
	params : path, files array
	return : files arrauy
*/
func Delete(path string, files []*rawgitlab.CommitActionOptions) []*rawgitlab.CommitActionOptions {
	var action rawgitlab.FileActionValue
	action = "delete"
	content := "Deleting file"
	files = append(files, &rawgitlab.CommitActionOptions{&action, &path, nil, &content, nil, nil, nil})

	return files
}

/*
	Add update file to the commit
	params : content, path, files array
	return : files array
*/
func Update(content string, path string, files []*rawgitlab.CommitActionOptions) []*rawgitlab.CommitActionOptions {
	var action rawgitlab.FileActionValue
	action = "update"
	files = append(files, &rawgitlab.CommitActionOptions{&action, &path, nil, &content, nil, nil, nil})

	return files
}
