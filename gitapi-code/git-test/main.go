package main

import (
    "fmt"
    "bytes"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "encoding/json"
)

const (
   gitAccessToken string="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"
   newGitRepoName string="git-go-test-five"
)

type Repo struct {
   Name string `json:"name"`
}

func main() {

     //create a new repo
     client := http.Client{}

     urlPost := "https://api.github.com/user/repos"

     dataBody:= Repo{newGitRepoName}
     data, err:= json.Marshal(dataBody)
     if err != nil {
         fmt.Println("Error in parsing post body")
         log.Fatal("bad", err)
     }
     req,err:=http.NewRequest("POST", urlPost, bytes.NewBuffer(data))
     if err!=nil{
       //Handle Error
     }

     authorizationString := "token " + gitAccessToken 
     req.Header.Set("Accept", "application/vnd.github.v3+json")
     req.Header.Add("Authorization", authorizationString)
     res, err:=client.Do(req)
     if err != nil {
         fmt.Print(err.Error())
         os.Exit(1)
     }
     responseData, err := ioutil.ReadAll(res.Body)
     if err != nil {
         log.Fatal(err)
     }
     fmt.Printf(string(responseData))

}
