package main

// Import key modules.
import (
  //"context"
  "log"
  "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
  //"github.com/Azure/azure-sdk-for-go/sdk/azcore"
  "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
  "github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
)

// Define key global variables.
var (
  subscriptionId = "0035266d-c66d-40f6-bfdb-02cc4146eb70"
)

// Define the function to create a resource group.

func main() {
  cred, err := azidentity.NewDefaultAzureCredential(nil)
  if err != nil {
    log.Fatalf("Authentication failure: %+v", err)
  }

  // Azure SDK Azure Resource Management clients accept the credential as a parameter
  client := armresources.NewResourcesClient(arm.NewDefaultConnection(cred, nil), subscriptionId)

  log.Printf("Authenticated to subscription", client)
  //log.Printf("Authenticated to subscription")
}
