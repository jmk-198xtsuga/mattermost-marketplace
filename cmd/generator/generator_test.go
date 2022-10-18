package main

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

func _TestGitHubRepositoryGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(makeHandleRepositoryGet(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.Get(ctx, "mattermost", "mattermost-plugin-github")

	if err != nil {
		t.Errorf("Repositiories.Get returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Repositories.Get did not get an OK HTTP response")
	}

	if got == nil {
		t.Errorf("Repositories.Get returned nil")
	} else {
		if (*got.Owner.Login != "mattermost") || (*got.Name != "mattermost-plugin-github") {
			t.Errorf("Repository mock response failed")
		}
	}
}

func _TestGitHubRepositoryList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(makeHandleRepositoryList(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.List(ctx, "mattermost", nil)

	if err != nil {
		t.Errorf("Repositiories.List returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Repositories.List did not get an OK HTTP response")
	}

	if got == nil {
		t.Errorf("Repositories.List returned nil")
	} else {
		if len(got) == 0 {
			t.Errorf("Repository list mock response failed")
		} else if *got[0].Owner.Login != "mattermost" {
			t.Errorf("Repository list mock response failed")
		}
	}
}

func _TestGitHubReleaseList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(makeHandleReleaseList(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.ListReleases(ctx, "mattermost", "mattermost-plugin-github", nil)

	if err != nil {
		t.Errorf("Repositiories.ListReleases returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Repositories.ListReleases did not get an OK HTTP response")
	}

	if got == nil {
		t.Errorf("Repositories.ListReleases returned nil")
	} else {
		if len(got) == 0 {
			t.Errorf("Release list mock response failed")
		} else if *got[len(got)-1].Name != "stable" {
			t.Errorf("Release list mock response failed")
		} else if *got[len(got)-1].TagName != "v1.0.0" {
			t.Errorf("Release list mock response failed")
		}
	}
}

func TestGitHubMockGeneration(t *testing.T) {
	_TestGitHubRepositoryGet(t)
	_TestGitHubRepositoryList(t)
	_TestGitHubReleaseList(t)
}
