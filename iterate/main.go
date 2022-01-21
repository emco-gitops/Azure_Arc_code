package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "io/ioutil"
    "path/filepath"
    "encoding/base64"
)

func main() {
    //currentDirectory, err := os.Getwd()
    //if err != nil {
        //log.Fatal(err)
    //}
    //iterate(currentDirectory)
    iterate("/home/ubuntu/go_projects/src/test-folder")
}

func iterate(path string) {
    length:=len([]rune(path))
    fmt.Printf("Length of Strinf: %d\n",len([]rune(path)))
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatalf(err.Error())
        }
        if info.IsDir()==false{
         fmt.Printf("File Name: %s\n", info.Name())
         fmt.Printf("Path of file: %s\n", path[length:])
         encodedValue := encode(path)
         // Print encoded data to console.
         // ... The base64 image can be used as a data URI in a browser.
         fmt.Println("ENCODED: " + encodedValue)
        }
         return nil
    })
}


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
