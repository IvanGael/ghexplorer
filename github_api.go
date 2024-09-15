package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const githubAPIBaseURL = "https://api.github.com"

func fetchGitHubProfile(username string) (*GitHubProfile, error) {
	url := fmt.Sprintf("%s/users/%s", githubAPIBaseURL, username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile GitHubProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func fetchRepositories(username string) ([]*Repository, error) {
	url := fmt.Sprintf("%s/users/%s/repos", githubAPIBaseURL, username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []*Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func fetchRepositoryContents(username, repo, path string) ([]*FileInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBaseURL, username, repo, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contents []*FileInfo
	err = json.NewDecoder(resp.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func fetchFileContent(username, repo, path string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBaseURL, username, repo, path)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

func searchRepositories(username, query string) ([]*Repository, error) {
	url := fmt.Sprintf("%s/search/repositories?q=%s+user:%s", githubAPIBaseURL, url.QueryEscape(query), username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResult struct {
		Items []*Repository `json:"items"`
	}
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}

	return searchResult.Items, nil
}
