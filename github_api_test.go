package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUsername is the GitHub username used for testing.
const TestUsername = "octocat"

func TestFetchGitHubProfile(t *testing.T) {
	profile, err := fetchGitHubProfile(TestUsername)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, TestUsername, profile.Login)
	assert.NotEmpty(t, profile.Name)
}

func TestFetchRepositories(t *testing.T) {
	repos, err := fetchRepositories(TestUsername)
	assert.NoError(t, err)
	assert.NotEmpty(t, repos)
	for _, repo := range repos {
		assert.NotEmpty(t, repo.Name)
	}
}

func TestFetchRepositoryContents(t *testing.T) {
	contents, err := fetchRepositoryContents(TestUsername, "Hello-World", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, contents)
	for _, item := range contents {
		assert.NotEmpty(t, item.Name)
		assert.NotEmpty(t, item.Type)
	}
}

func TestFetchFileContent(t *testing.T) {
	content, err := fetchFileContent(TestUsername, "Hello-World", "README")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "Hello World!")
}

func TestSearchRepositories(t *testing.T) {
	repos, err := searchRepositories(TestUsername, "Hello-World")
	assert.NoError(t, err)
	assert.NotEmpty(t, repos)
	assert.Contains(t, repos[0].Name, "Hello-World")
}
