package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	//"os"
	"encoding/json"
)

const (
	clientIdValue           string = "a61fca65-1f43-46cb-bd52-48f891203272"
	resourceName            string = "https://management.core.windows.net/"
	clientSecretValue       string = "3R07Q~RMwTldmA.xi2ARx6jvBUa_BQDF4Jw5R"
	tenantId                string = "1e452d66-a99c-413b-b87f-bc0c04219888"
	subscriptionId          string = "a68d0074-3737-48ac-9af0-b4a2e93f5e11"
	arcClusterResourceGroup string = "flux-demo-rg"
	arcCluster              string = "flux-demo-arc"
	repositoryUrlName       string = "https://github.com/chitti-intel/Azure-test-repo-two"
	gitConfigurationName    string = "gitops-demo"
	operatorScope           string = "cluster"
)

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
}

type Properties struct {
	RepositoryUrl         string `json:"repositoryUrl"`
	OperatorNamespace     string `json:"operatorNamespace"`
	OperatorInstanceName  string `json:"operatorInstanceName"`
	OperatorType          string `json:"operatorType"`
	OperatorParams        string `json:"operatorParams"`
	OperatorScope         string `json:"operatorScope"`
	SshKnownHostsContents string `json:"sshKnownHostsContents"`
}

type PropertiesFlux struct {
	Scope          string             `json:"scope"`
	Namespace      string             `json:"namespace"`
	SourceKind     string             `json:"sourceKind"`
	Suspend        bool               `json:"suspend"`
	GitRepository  RepoProperties     `json:"gitRepository"`
	Kustomizations KustomizationsUnit `json:"kustomizations"`
}

type RepoProperties struct {
	Url           string  `json:"url"`
	RepositoryRef RepoRef `json:"repositoryRef"`
}

type RepoRef struct {
	Branch string `json:"branch"`
}

type KustomizationsUnit struct {
	FirstKustomization KustomizationProperties `json:"kustomization-1"`
}

type KustomizationProperties struct {
	Path                  string `json:"path"`
	TimeoutInSeconds      int    `json:"timeoutInSeconds"`
	SyncIntervalInSeconds int    `json:"syncIntervalInSeconds"`
	Prune                 bool   `json:"prune"`
	Force                 bool   `json:"force"`
}

type Requestbody struct {
	Properties Properties `json:"properties"`
}

type RequestbodyFlux struct {
	Properties PropertiesFlux `json:"properties"`
}

type FluxExtension struct {
	Identity   IndentityProp `json:"identity"`
	Properties ExtensionProp `json:"properties"`
}

type IndentityProp struct {
	Type string `json:"type"`
}

type ExtensionProp struct {
	ExtensionType           string `json:"extensionType"`
	AutoUpgradeMinorVersion bool   `json:"autoUpgradeMinorVersion"`
}

func main() {
	//Rest api to get the access token
	accessTokenValue, err := getAccessToken(clientIdValue, clientSecretValue, tenantId)
	if err != nil {
		log.Fatal(err)
	}
	// Install microsoft.flux extension
	extensionResponse, err := installFluxExtension(accessTokenValue, subscriptionId, arcClusterResourceGroup, arcCluster)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(extensionResponse))

	// // PUT request for creating git configuration
	// // PUT request body
	// configResponse, err := createGitConfiguration(accessTokenValue, repositoryUrlName, gitConfigurationName, operatorScope, subscriptionId, arcClusterResourceGroup, arcCluster)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf(string(configResponse))

	// PUT request for creating fluxV2 configuration
	// PUT request body
	configResponse, err := createFluxConfiguration(accessTokenValue, repositoryUrlName, gitConfigurationName, operatorScope, subscriptionId, arcClusterResourceGroup, arcCluster)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(configResponse))

}

//func to get the access token for azure arc
func getAccessToken(clientId string, clientSecret string, tenantIdValue string) (string, error) {
	//Rest api to get the access token
	client := http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Add("client_id", clientId)
	data.Add("resource", "https://management.core.windows.net/")
	data.Add("client_secret", clientSecret)

	urlPost := "https://login.microsoftonline.com/" + tenantIdValue + "/oauth2/token"

	req, err := http.NewRequest("POST", urlPost, bytes.NewBufferString(data.Encode()))
	if err != nil {
		//Handle Error
		fmt.Print(err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	res, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}

	// Unmarshall the response body into json and get token value
	newToken := Token{}
	json.Unmarshal(responseData, &newToken)

	return newToken.AccessToken, nil
}

//create a git configuration for the mentioned user repo
func createGitConfiguration(accessToken string, repositoryUrl string, gitConfiguration string, operatorScopeType string, subscriptionIdValue string, arcClusterResourceGroupName string, arcClusterName string) (string, error) {
	// PUT request for creating git configuration
	// PUT request body
	client := http.Client{}
	properties := Requestbody{Properties{repositoryUrl, gitConfiguration, gitConfiguration, "flux", "--git-branch=main --sync-garbage-collection --git-path=arc-k8s-demo", operatorScopeType, ""}}
	dataProperties, err := json.Marshal(properties)
	if err != nil {
		fmt.Println("Error in parsing properties")
		return "", err
	}

	urlPut := "https://management.azure.com/subscriptions/" + subscriptionIdValue + "/resourceGroups/" + arcClusterResourceGroupName + "/providers/Microsoft.Kubernetes/connectedClusters/" + arcClusterName + "/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/" + gitConfiguration + "?api-version=2021-03-01"
	reqPut, err := http.NewRequest(http.MethodPut, urlPut, bytes.NewBuffer(dataProperties))

	if err != nil {
		//Handle Error
		fmt.Println("Error in the put request")
		return "", err
	}
	// Add request header
	authorizationString := "Bearer " + accessToken
	reqPut.Header.Set("Content-Type", "application/json; charset=UTF-8")
	reqPut.Header.Add("Authorization", authorizationString)
	fmt.Println(reqPut)
	resPut, err := client.Do(reqPut)
	if err != nil {
		//Handle Error
		fmt.Println("Error in the put request from server")
		return "", err
	}
	responseDataPut, err := ioutil.ReadAll(resPut.Body)
	if err != nil {
		fmt.Println("Error in the response of put request")
		return "", err
	}
	fmt.Printf(string(responseDataPut))
	return string(responseDataPut), nil

}

//create a FluxV2 configuration for the mentioned user repo
func createFluxConfiguration(accessToken string, repositoryUrl string, gitConfiguration string, operatorScopeType string, subscriptionIdValue string, arcClusterResourceGroupName string, arcClusterName string) (string, error) {
	// PUT request for creating git configuration
	// PUT request body
	client := http.Client{}
	// properties := Requestbody{Properties{repositoryUrl, gitConfiguration, gitConfiguration, "flux", "--git-branch=main --sync-garbage-collection --git-path=arc-k8s-demo", operatorScopeType, ""}}
	properties := RequestbodyFlux{PropertiesFlux{operatorScopeType, gitConfiguration, "GitRepository", false, RepoProperties{repositoryUrl, RepoRef{"main"}}, KustomizationsUnit{KustomizationProperties{"./Arc-test-folder", 600, 600, true, false}}}}
	dataProperties, err := json.Marshal(properties)
	if err != nil {
		fmt.Println("Error in parsing properties")
		return "", err
	}

	urlPut := "https://management.azure.com/subscriptions/" + subscriptionIdValue + "/resourceGroups/" + arcClusterResourceGroupName + "/providers/Microsoft.Kubernetes/connectedClusters/" + arcClusterName + "/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/" + gitConfiguration + "?api-version=2022-01-01-preview"
	reqPut, err := http.NewRequest(http.MethodPut, urlPut, bytes.NewBuffer(dataProperties))

	if err != nil {
		//Handle Error
		fmt.Println("Error in the put request")
		return "", err
	}
	// Add request header
	authorizationString := "Bearer " + accessToken
	reqPut.Header.Set("Content-Type", "application/json; charset=UTF-8")
	reqPut.Header.Add("Authorization", authorizationString)
	fmt.Println(reqPut)
	resPut, err := client.Do(reqPut)
	if err != nil {
		//Handle Error
		fmt.Println("Error in the put request from server")
		return "", err
	}
	responseDataPut, err := ioutil.ReadAll(resPut.Body)
	if err != nil {
		fmt.Println("Error in the response of put request")
		return "", err
	}
	fmt.Printf(string(responseDataPut))
	return string(responseDataPut), nil

}

//install microsoft.flux extension on the cluster
func installFluxExtension(accessToken string, subscriptionIdValue string, arcClusterResourceGroupName string, arcClusterName string) (string, error) {
	// PUT request for installing microsoft.flux extension
	// PUT request body
	client := http.Client{}
	properties := FluxExtension{IndentityProp{"SystemAssigned"}, ExtensionProp{"microsoft.flux", true}}
	dataProperties, err := json.Marshal(properties)
	if err != nil {
		fmt.Println("Error in parsing properties for extension install")
		return "", err
	}

	// urlPut := "https://management.azure.com/subscriptions/" + subscriptionIdValue + "/resourceGroups/" + arcClusterResourceGroupName + "/providers/Microsoft.Kubernetes/connectedClusters/" + arcClusterName + "/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/" + gitConfiguration + "?api-version=2021-03-01"
	urlPut := "https://management.azure.com/subscriptions/" + subscriptionIdValue + "/resourceGroups/" + arcClusterResourceGroupName + "/providers/Microsoft.Kubernetes/connectedClusters/" + arcClusterName + "/providers/Microsoft.KubernetesConfiguration/extensions/flux?api-version=2021-09-01"
	reqPut, err := http.NewRequest(http.MethodPut, urlPut, bytes.NewBuffer(dataProperties))

	if err != nil {
		//Handle Error
		fmt.Println("Error in the put request for extension install")
		return "", err
	}
	// Add request header
	authorizationString := "Bearer " + accessToken
	reqPut.Header.Set("Content-Type", "application/json; charset=UTF-8")
	reqPut.Header.Add("Authorization", authorizationString)
	fmt.Println(reqPut)
	resPut, err := client.Do(reqPut)
	if err != nil {
		//Handle Error
		fmt.Println("Error in the put request for extension install from server")
		return "", err
	}
	responseDataPut, err := ioutil.ReadAll(resPut.Body)
	if err != nil {
		fmt.Println("Error in the response of put request for extension install")
		return "", err
	}
	fmt.Printf(string(responseDataPut))
	return string(responseDataPut), nil
}
