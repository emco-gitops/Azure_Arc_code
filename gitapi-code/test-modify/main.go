package main

import (
"fmt"
"bytes"
"net/http"
"io/ioutil"
"encoding/json"
"log"
"os"
//"encoding/json"
"bufio"
//"path/filepath"
"encoding/base64"

)


const (
   gitAccessToken string="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"
   gitRepoName string="git-go-test-four"
)

type Body struct {
   Message string `json:"message"`
   Content string `json:"content"`
   Sha string `json:"sha"`
}


func main() {

// obtain the sha of the file
path := "/cluster-apps/arc-k8s-demo.yaml" 
shaValue, err := getSha(path)
if err!= nil{
   log.Fatal(err)
}
fmt.Println(shaValue)

// modify the file contents

//encode the modified file contents
encoded := encode("/home/ubuntu/go_projects/src/k8s-test-code/modified-file")

//send the file
response,err := sendFile(path, shaValue, encoded)
if err!= nil{
   log.Fatal(err)
}

fmt.Println(string(response))
}

//function to obtain the sha of file
func getSha(path string) (string,error) {
url := "https://api.github.com/repos/chitti-intel/" + gitRepoName + "/contents" + path
req, _ := http.NewRequest("GET", url, nil)
//authorizationString := "token " + gitAccessToken
//req.Header.Add("Authorization",authorizationString)
res, _ := http.DefaultClient.Do(req)

defer res.Body.Close()

body, _ := ioutil.ReadAll(res.Body)
m := make(map[string]interface{})
err := json.Unmarshal(body, &m)
if err != nil {
    return "",err
}

sha := fmt.Sprintf("%v", m["sha"])
return sha,nil

}



//encode contents of the file to base64
func encode(filePath string) string {

    // Open file on disk.
    f, _ := os.Open(filePath)

    // Read entire file into byte slice.
    reader := bufio.NewReader(f)
    content, _ := ioutil.ReadAll(reader)

    // Encode as base64.
    encoded := base64.StdEncoding.EncodeToString(content)

    return encoded
}


//write file to the gitrepo
func sendFile(path string, sha string, contents string)(string,error){
     //create a new client
     client := http.Client{}

     urlPut := "https://api.github.com/repos/chitti-intel/" + gitRepoName + "/contents" + path

     dataBody:= Body{"Test file creation from go", contents, sha}
     data, err:= json.Marshal(dataBody)
     if err != nil {
         fmt.Println("Error in parsing post body")
         return "",err
     }
     req,err:=http.NewRequest("PUT", urlPut, bytes.NewBuffer(data))
     if err!=nil{
       //Handle Error
       return "",err
     }

     authorizationString := "token " + gitAccessToken
     req.Header.Set("Accept", "application/vnd.github.v3+json")
     req.Header.Add("Authorization", authorizationString)

     res, err:=client.Do(req)
     if err != nil {
         return "",err
     }
     responseData, err := ioutil.ReadAll(res.Body)
     if err != nil {
         return "",err
     }
     return string(responseData),nil

}
