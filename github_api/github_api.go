package github_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"ghexplorer/config"

	"io"
	"net/http"
	"net/url"
)

// GitHubProfile is GitHub profile struct
type GitHubProfile struct {
	Name        string `json:"name"`
	Login       string `json:"login"`
	Description string `json:"bio"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
}

// Repository is GitHub profile repository struct
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// FileInfo is GitHub profile repository file info struct
type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// FetchGitHubProfile fetch GitHub profile
func FetchGitHubProfile(username string) (*GitHubProfile, error) {
	customUrl := fmt.Sprintf("%s/users/%s", config.GithubAPIBaseURL, username)
	resp, err := http.Get(customUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch profile: %s", resp.Status)
	}

	var profile GitHubProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// fetchProfileReadme fetches the user's profile README.md content
// func fetchProfileReadme(username string) (string, error) {
// 	content, err := fetchFileContent(username, username, "README.md")
// 	if err != nil {
// 		return "", err
// 	}
// 	return content, nil
// }

// FetchRepositories fetch GitHub profile repositories with pagination
func FetchRepositories(username string) ([]*Repository, error) {
	var allRepos []*Repository
	page := 1
	perPage := 100 // Maximum allowed by GitHub API

	for {
		customUrl := fmt.Sprintf("%s/users/%s/repos?page=%d&per_page=%d", config.GithubAPIBaseURL, username, page, perPage)
		resp, err := http.Get(customUrl)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			err := resp.Body.Close()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("failed to fetch repositories: %s", resp.Status)
		}

		var repos []*Repository
		err = json.NewDecoder(resp.Body).Decode(&repos)
		errResp := resp.Body.Close()
		if errResp != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)

		// Check if we've received less than perPage items, meaning this is the last page
		if len(repos) < perPage {
			break
		}

		page++
	}

	return allRepos, nil
}

// FetchRepositoryContents fetch GitHub profile repository contents
func FetchRepositoryContents(username, repo, path string) ([]*FileInfo, error) {
	customUrl := fmt.Sprintf("%s/repos/%s/%s/contents/%s", config.GithubAPIBaseURL, username, repo, path)
	resp, err := http.Get(customUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repository contents: %s", resp.Status)
	}

	var contents []*FileInfo
	err = json.NewDecoder(resp.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// FetchFileContent fetch GitHub profile repository file contents
func FetchFileContent(username, repo, path string) (string, error) {
	customUrl := fmt.Sprintf("%s/repos/%s/%s/contents/%s", config.GithubAPIBaseURL, username, repo, path)
	resp, err := http.Get(customUrl)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file content: %s", resp.Status)
	}

	var fileContent struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	err = json.NewDecoder(resp.Body).Decode(&fileContent)
	if err != nil {
		return "", err
	}

	if fileContent.Encoding == "base64" {
		decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
		if err != nil {
			return "", err
		}
		return string(decodedContent), nil
	}
	return fileContent.Content, nil
}

// SearchRepositories perform searching through GitHub profile repositories
func SearchRepositories(username, query string) ([]*Repository, error) {
	customUrl := fmt.Sprintf("%s/search/repositories?q=%s+user:%s", config.GithubAPIBaseURL, url.QueryEscape(query), username)
	resp, err := http.Get(customUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search repositories: %s", resp.Status)
	}

	var searchResult struct {
		Items []*Repository `json:"items"`
	}
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}
	return searchResult.Items, nil
}
