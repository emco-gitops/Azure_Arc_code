package main

import (
        "context"
        "fmt"
	"github.com/fluxcd/go-git-providers/gitlab"
        "github.com/fluxcd/go-git-providers/gitprovider"
    	"io/ioutil"
    	"log"
    	"os"
    	"bufio"
    	"path/filepath"

)
const (
	gitLabDomain = "gitlab.com"
)


func main() {

	iterate("/home/ubuntu/go_projects/src/k8s-test-code/arc-k8s-demo")
}

//iterate through all the files
func iterate(path string) {
    files := []gitprovider.CommitFile{}
    length:=len([]rune(path))
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error{
        if err != nil {
            log.Fatalf(err.Error())
        }
        if info.IsDir()==false{
         fmt.Printf("File Name: %s\n", info.Name())
         //Encode the file to base64
         content := getContent(path)
         //send the file to the github repo
	 relativePath := path[length+1:]
	 fmt.Println(relativePath)
         files = append(files,gitprovider.CommitFile{Path : &relativePath, Content: &content})
        }
         return nil
    })
    fmt.Println(len(files))
    sendFile(files)
}

//get contents of the file
func getContent(filePath string) string {

    // Open file on disk.
    f, _ := os.Open(filePath)

    // Read entire file into byte slice.
    reader := bufio.NewReader(f)
    content, _ := ioutil.ReadAll(reader)

    return string(content)
}


func sendFile(files []gitprovider.CommitFile) {
        // Create a new client
        ctx := context.Background()
	gitLabToken := "glpat-VsaWds-rWtbx6eM5ejBm"
	c, err := gitlab.NewClient(gitLabToken,"")
        userRef := gitprovider.UserRef{
                Domain:    gitLabDomain,
                UserLogin: "chitti-intel",
        }

        repoName := "Go-Git-Test-Repo-two"
        userRepoRef := gitprovider.UserRepositoryRef{
                UserRef:        userRef,
                RepositoryName: repoName,
        }

        userRepo, err := c.UserRepositories().Get(ctx, userRepoRef)
        //Commit file to this repo
        _, err = userRepo.Commits().Create(ctx, "main", "app files added", files)

        fmt.Println(err)

}
