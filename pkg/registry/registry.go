package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	NodeArchitecture = "amd64" // replace with your node's architecture
)

type Manifest struct {
	Architecture string `json:"architecture"`
}

func CheckImageArchitecture(image string) bool {
	parts := strings.Split(image, ":")
	repo, tag := parts[0], parts[1]

	var url string
	if strings.Contains(repo, ".azurecr.io") {
		// It's an Azure Container Registry image
		url = fmt.Sprintf("https://%s/v2/%s/manifests/%s", repo, getRepoName(repo), tag)

		// Get an ACR access token
		accessToken, err := GetAcrAccessToken(repo)
		if err != nil {
			return false
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		// It's a Docker Hub image
		url = fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", repo, tag)
	}

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return false
	}

	return manifest.Architecture == NodeArchitecture
}

func getRepoName(repo string) string {
	parts := strings.Split(repo, "/")
	return parts[len(parts)-1]
}
