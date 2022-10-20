package repository

import "testing"

func TestGitHubMockGeneration(t *testing.T) {
	MockGitHubRepositoryGet(t)
	MockGitHubRepositoryList(t)
	MockGitHubReleaseList(t)
}
