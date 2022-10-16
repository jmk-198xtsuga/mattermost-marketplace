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

const pathRepository = "/repos"

func makeHandleRepositoryGet(t *testing.T) func(http.ResponseWriter, *http.Request) {
	if t == nil {
		return nil
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte(response))
		}
	}
}

func TestGitHubRepositories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(makeHandleRepositoryGet(t)))
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	ctx := context.Background()
	got, resp, err := client.Repositories.Get(ctx, "octocat", "Hello-World")

	if err != nil {
		t.Errorf("Repositiories.Get returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Repositories.Get did not get an OK HTTP response")
	}

	if got == nil {
		t.Errorf("Repositories.Get returned nil")
	} else {
		if *got.Name != "Hello-World" || *got.Owner.Login != "octocat" {
			t.Errorf("Mock response failed")
		}
	}
}
