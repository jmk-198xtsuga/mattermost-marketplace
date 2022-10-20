package repository

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-github/v28/github"
)

//go:embed repository_format.json
var json_repository_format string

//go:embed release_format.json
var json_release_format string

func makeHandleRepositoryGet(t *testing.T) func(http.ResponseWriter, *http.Request) {
	if t == nil {
		return nil
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			const pathRepository = "/repos"

			repoRegex := "^" + pathRepository + "(/.+){0,2}$"
			requestPath := r.URL.Path
			matchSimple, _ := regexp.MatchString(repoRegex, pathRepository)
			matchUserParam, _ := regexp.MatchString(repoRegex, pathRepository+"/user")
			matchUserParamRepoParam, _ := regexp.MatchString(repoRegex, pathRepository+"/user/foo")
			if !(matchSimple || matchUserParam || matchUserParamRepoParam) {
				t.Errorf("Regex Match Fail")
			}
			if match, _ := regexp.MatchString(repoRegex, requestPath); !match {
				t.Errorf("Request path was not \"" + pathRepository + "/{owner}/{repo}\"")
			}
			if match, _ := regexp.MatchString("application/([^\\s]+\\+)json", r.Header.Get("Accept")); !match {
				t.Errorf("Request does not accept JSON response")
			}
			w.WriteHeader(http.StatusOK)
			requestPathParts := strings.Split(r.URL.Path, "/")
			requestOwnerName := requestPathParts[len(requestPathParts)-2]
			requestRepoName := requestPathParts[len(requestPathParts)-1]
			response := fmt.Sprintf(json_repository_format, 1, requestRepoName,
				requestOwnerName+"/"+requestRepoName, requestOwnerName,
				1, requestOwnerName, 1)
			//fmt.Print(response)
			w.Write([]byte(response))
		}
	}
}

func makeHandleRepositoryList(t *testing.T) func(http.ResponseWriter, *http.Request) {
	if t == nil {
		return nil
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			var repositoryNames [12]string = [...]string{
				"mattermost-plugin-github",
				"mattermost-plugin-autolink",
				"mattermost-plugin-zoom",
				"mattermost-plugin-jira",
				"mattermost-plugin-welcomebot",
				"mattermost-plugin-jenkins",
				"mattermost-plugin-antivirus",
				"mattermost-plugin-custom-attributes",
				"mattermost-plugin-aws-SNS",
				"mattermost-plugin-gitlab",
				"mattermost-plugin-nps",
				"mattermost-plugin-webex",
			}

			repoRegex := "^/users/.+/repos$"
			requestPath := r.URL.Path
			matchUserParam, _ := regexp.MatchString(repoRegex, "/users/user/repos")
			if !(matchUserParam) {
				t.Errorf("Regex Match Fail")
			}
			if match, _ := regexp.MatchString(repoRegex, requestPath); !match {
				t.Errorf("Request path was not \"/users/{user}/repos\"")
			}
			if match, _ := regexp.MatchString("application/([^\\s]+\\+)json", r.Header.Get("Accept")); !match {
				t.Errorf("Request does not accept JSON response")
			}
			w.WriteHeader(http.StatusOK)
			requestPathParts := strings.Split(r.URL.Path, "/")
			requestOwnerName := requestPathParts[len(requestPathParts)-2]
			response := "[\n    "
			response += strings.ReplaceAll(fmt.Sprintf(json_repository_format, 1, repositoryNames[0],
				requestOwnerName+"/"+repositoryNames[0], requestOwnerName,
				1, requestOwnerName, 1), "\n", "\n    ")
			for idx, name := range repositoryNames[1:] {
				response += ",\n    "
				response += strings.ReplaceAll(fmt.Sprintf(json_repository_format, idx+2, name,
					requestOwnerName+"/"+name, requestOwnerName,
					1, requestOwnerName, 1), "\n", "\n    ")
			}
			response += "\n]"
			//fmt.Print(response)
			w.Write([]byte(response))
		}
	}
}

func makeHandleReleaseList(t *testing.T) func(http.ResponseWriter, *http.Request) {
	if t == nil {
		return nil
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			releaseNames := [...]string{"alpha", "beta", "stable"}
			releaseTags := [...]string{"v0.1.0", "v0.9.0", "v1.0.0"}
			repoRegex := "^/repos(/.+){2}/releases$"
			requestPath := r.URL.Path
			matchUserParam, _ := regexp.MatchString(repoRegex, "/repos/user/foo/releases")
			if !(matchUserParam) {
				t.Errorf("Regex Match Fail")
			}
			if match, _ := regexp.MatchString(repoRegex, requestPath); !match {
				t.Errorf("Request path was not \"/repos/{owner}/{repo}/releases\"")
			}
			if match, _ := regexp.MatchString("application/([^\\s]+\\+)json", r.Header.Get("Accept")); !match {
				t.Errorf("Request does not accept JSON response")
			}
			w.WriteHeader(http.StatusOK)
			requestPathParts := strings.Split(r.URL.Path, "/")
			requestOwnerName := requestPathParts[len(requestPathParts)-3]
			response := "[\n    "
			for idx, name := range releaseNames {
				if idx > 0 {
					response += ","
				}
				response += "\n    "
				response += strings.ReplaceAll(fmt.Sprintf(json_release_format, idx+1, name,
					releaseTags[idx], requestOwnerName, 1), "\n", "\n    ")
			}
			response += "\n]"
			//fmt.Print(response)
			w.Write([]byte(response))
		}
	}
}

func MockGitHubClientRepositoryGet(t *testing.T) *github.Repository {
	var repository *github.Repository = nil
	server := httptest.NewServer(http.HandlerFunc(makeHandleRepositoryGet(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.Get(ctx, "mattermost", "mattermost-plugin-github")

	if err != nil {
		if t != nil {
			t.Errorf("Repositiories.Get returned error: %v", err)
		}
	} else if resp.StatusCode != http.StatusOK {
		if t != nil {
			t.Errorf("Repositories.Get did not get an OK HTTP response")
		}
	} else if got == nil {
		if t != nil {
			t.Errorf("Repositories.Get returned nil")
		}
	} else {
		if (*got.Owner.Login != "mattermost") || (*got.Name != "mattermost-plugin-github") {
			if t != nil {
				t.Errorf("Repository mock response failed")
			}
		} else {
			repository = got
		}
	}
	return repository
}

func MockGitHubRepository(owner string, repo string) *github.Repository {
	var repository *github.Repository = new(github.Repository)
	repository.ID = new(int64)
	*repository.ID = 1
	repository.Name = new(string)
	*repository.Name = repo
	repository.FullName = new(string)
	*repository.FullName = strings.Join([]string{owner, repo}, "/")
	repository.Owner = new(github.User)
	repository.Owner.ID = new(int64)
	*repository.Owner.ID = 1
	repository.Owner.Login = new(string)
	*repository.Owner.Login = owner
	return repository
}

func MockGitHubClientRepositoryList(t *testing.T) []*github.Repository {
	var repositories []*github.Repository = nil
	server := httptest.NewServer(http.HandlerFunc(makeHandleRepositoryList(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.List(ctx, "mattermost", nil)

	if err != nil {
		if t != nil {
			t.Errorf("Repositiories.List returned error: %v", err)
		}
	} else if resp.StatusCode != http.StatusOK {
		if t != nil {
			t.Errorf("Repositories.List did not get an OK HTTP response")
		}
	} else if got == nil {
		if t != nil {
			t.Errorf("Repositories.List returned nil")
		}
	} else {
		if len(got) == 0 {
			if t != nil {
				t.Errorf("Repository list mock response failed")
			}
		} else if *got[0].Owner.Login != "mattermost" {
			if t != nil {
				t.Errorf("Repository list mock response failed")
			}
		} else {
			repositories = got
		}
	}
	return repositories
}

func MockGitHubRepositoryList(owner string, repos []string) []*github.Repository {
	var repositories []*github.Repository
	for i, r := range repos {
		var repo *github.Repository = MockGitHubRepository(owner, r)
		*repo.ID = int64(i + 1)
		repositories = append(repositories, repo)
	}
	return repositories
}

func MockGitHubClientReleaseList(t *testing.T) []*github.RepositoryRelease {
	var releases []*github.RepositoryRelease = nil
	server := httptest.NewServer(http.HandlerFunc(makeHandleReleaseList(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.ListReleases(ctx, "mattermost", "mattermost-plugin-github", nil)

	if err != nil {
		if t != nil {
			t.Errorf("Repositiories.ListReleases returned error: %v", err)
		}
	} else if resp.StatusCode != http.StatusOK {
		if t != nil {
			t.Errorf("Repositories.ListReleases did not get an OK HTTP response")
		}
	} else if got == nil {
		if t != nil {
			t.Errorf("Repositories.ListReleases returned nil")
		}
	} else {
		if len(got) == 0 {
			if t != nil {
				t.Errorf("Release list mock response failed")
			}
		} else if *got[len(got)-1].Name != "stable" {
			if t != nil {
				t.Errorf("Release list mock response failed")
			}
		} else if *got[len(got)-1].TagName != "v1.0.0" {
			if t != nil {
				t.Errorf("Release list mock response failed")
			}
		} else {
			releases = got
		}
	}
	return releases
}

func MockGitHubReleaseList(owner string, repo string, tags []string) []*github.RepositoryRelease {
	var releases []*github.RepositoryRelease
	for i, t := range tags {
		var release *github.RepositoryRelease = new(github.RepositoryRelease)
		release.ID = new(int64)
		*release.ID = int64(i + 1)
		release.Name = new(string)
		*release.Name = t
		release.TagName = new(string)
		*release.TagName = t
		release.Author = new(github.User)
		release.Author.ID = new(int64)
		*release.Author.ID = 1
		release.Author.Login = new(string)
		*release.Author.Login = owner
		releases = append(releases, release)
	}
	return releases
}
