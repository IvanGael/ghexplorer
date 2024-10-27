package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// fetchGitHubProfile fetch github profile
func fetchGitHubProfile(username string) (*GitHubProfile, error) {
	url := fmt.Sprintf("%s/users/%s", githubAPIBaseURL, username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
func fetchProfileReadme(username string) (string, error) {
	content, err := fetchFileContent(username, username, "README.md")
	if err != nil {
		return "", err
	}
	return content, nil
}

// fetchRepositories fetch github profile repositories with pagination
func fetchRepositories(username string) ([]*Repository, error) {
	var allRepos []*Repository
	page := 1
	perPage := 100 // Maximum allowed by GitHub API

	for {
		url := fmt.Sprintf("%s/users/%s/repos?page=%d&per_page=%d", githubAPIBaseURL, username, page, perPage)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch repositories: %s", resp.Status)
		}

		var repos []*Repository
		err = json.NewDecoder(resp.Body).Decode(&repos)
		resp.Body.Close()

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

// fetchRepositoryContents fetch github profile repository contents
func fetchRepositoryContents(username, repo, path string) ([]*FileInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBaseURL, username, repo, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

// fetchFileContent fetch github profile repository file contents
func fetchFileContent(username, repo, path string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBaseURL, username, repo, path)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

// searchRepositories perform searching through github profile repositories
func searchRepositories(username, query string) ([]*Repository, error) {
	url := fmt.Sprintf("%s/search/repositories?q=%s+user:%s", githubAPIBaseURL, url.QueryEscape(query), username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
