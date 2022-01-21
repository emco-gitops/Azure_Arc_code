package main

import (
    "fmt"
    "bytes"
    "io/ioutil"
    "log"
    "os"
    "bufio"
    "path/filepath"
    "encoding/base64"

)

const (
   gitAccessToken string="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"
   gitRepoName string="git-go-test-five"
)

type Body struct {
   Message string `json:"message"`
   Content string `json:"content"`
}

func main() {

     iterate("/home/ubuntu/go_projects/src/k8s-test-code")

}

//iterate through all the files 
func iterate(path string) {
    length:=len([]rune(path))
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatalf(err.Error())
        }
        if info.IsDir()==false{
         fmt.Printf("File Name: %s\n", info.Name())
         //Encode the file to base64
         encodedValue := encode(path)
         // Print encoded data to console
         fmt.Println("ENCODED: " + encodedValue)
         //send the file to the github repo
         response,err := sendFile(path[length:],encodedValue)
         if err != nil {
            log.Fatal(err)
         }
         fmt.Println(response)   
        }
         return nil
    })
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
func sendFile(path string, contents string)(string,error){
     //create a new client
     client := http.Client{}

     urlPut := "https://api.github.com/repos/chitti-intel/" + gitRepoName + "/contents" + path

     dataBody:= Body{"Test file creation from go",contents}
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
