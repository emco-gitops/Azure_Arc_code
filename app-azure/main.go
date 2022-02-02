package main

import (
	"azurearc"
	"fmt"
	"log"
)

const (
	clientIdValue           string = "a61fca65-1f43-46cb-bd52-48f891203272"
	resourceName            string = "https://management.core.windows.net/"
	clientSecretValue       string = "3R07Q~RMwTldmA.xi2ARx6jvBUa_BQDF4Jw5R"
	tenantId                string = "1e452d66-a99c-413b-b87f-bc0c04219888"
	subscriptionId          string = "a68d0074-3737-48ac-9af0-b4a2e93f5e11"
	arcClusterResourceGroup string = "flux-demo-rg"
	arcCluster              string = "flux-demo-arc"
	repositoryUrlName       string = "https://github.com/Azure/arc-k8s-demo"
	gitConfigurationName    string = "gitops-demo"
	operatorScope           string = "cluster"
	gitPath                 string = ""
)

func main() {
	//Rest api to get the access token
	accessTokenValue, err := azurearc.GetAccessToken(clientIdValue, clientSecretValue, tenantId)
	if err != nil {
		log.Fatal(err)
	}
	// Install microsoft.flux extension
	extensionResponse, err := azurearc.InstallFluxExtension(accessTokenValue, subscriptionId, arcClusterResourceGroup, arcCluster)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(extensionResponse))

	// PUT request for creating git configuration
	// PUT request body
	// configResponse, err := azurearc.CreateGitConfiguration(accessTokenValue, repositoryUrlName, gitConfigurationName, operatorScope, subscriptionId, arcClusterResourceGroup, arcCluster, "main", gitPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf(string(configResponse))

	// // PUT request for creating fluxV2 configuration
	// // PUT request body
	// configResponse, err := azurearc.CreateFluxConfiguration(accessTokenValue, repositoryUrlName, gitConfigurationName, operatorScope, subscriptionId, arcClusterResourceGroup, arcCluster, "main", "", 60, 60)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf(string(configResponse))

	deleteResponse, err := azurearc.DeleteFluxConfiguration(accessTokenValue, subscriptionId, arcClusterResourceGroup, arcCluster, gitConfigurationName)
	fmt.Printf(string(deleteResponse))

	// deleteResponse, err := azurearc.DeleteGitConfiguration(accessTokenValue, subscriptionId, arcClusterResourceGroup, arcCluster, gitConfigurationName)
	// fmt.Printf(string(deleteResponse))
}
