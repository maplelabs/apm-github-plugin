package gitprovider

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

// GithubClient represents new Github client to access github APIs
type GithubClient struct {
	// Client is github access client
	Client *github.Client

	// RepositoryName for accessing data
	RepositoryName string

	// RepositoryOwner for accessing data
	RepositoryOwner string

	// Username for authentication
	Username string

	// Accesstoken for authetication
	Accesstoken string

	// ctx for request
	ctx context.Context
}

// NewGithubClient returns a new github api client
func NewGithubClient(repoOwner string, repoName string, userName string, accessToken string) *GithubClient {
	gc := new(GithubClient)
	gc.RepositoryName = repoName
	gc.RepositoryOwner = repoOwner
	gc.Username = userName
	gc.Accesstoken = accessToken
	ctx := context.Background()
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		gc.Client = github.NewClient(tc)
	} else {
		gc.Client = github.NewClient(nil)
	}
	gc.ctx = ctx
	return gc
}

// CheckCredentials checks credentails for the user
func (gc *GithubClient) CheckCredentials() error {
	// ignoring response as we are only concerned with authentication
	_, _, err := gc.Client.Users.Get(gc.ctx, "")
	if err != nil {
		log.Errorf("error[%v] in authenticating credentials for user %v", err, gc.Username)
		return err
	}
	return nil
}

// GetCommits fetches commits for the user
func (gc *GithubClient) GetCommits(from time.Time, to time.Time, branch string) ([]byte, error) {
	log.Debug("commit to be fetched from branch %v for repository %v after %v to %v", branch, gc.RepositoryName, from, to)
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		SHA:         branch,
		Since:       from,
		Until:       to,
	}
	var allCommits []*github.RepositoryCommit
	for {
		commits, resp, err := gc.Client.Repositories.ListCommits(gc.ctx, gc.RepositoryOwner, gc.RepositoryName, opt)
		if err != nil {
			log.Errorf("error[%v] in fetching commits for repository %v", err, gc.RepositoryName)
			return nil, err
		}
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	allCommitsByte, err := json.Marshal(allCommits)
	return allCommitsByte, err
}

// GetPullRequests fetches pull request for the user
func (gc *GithubClient) GetPullRequests(fromNo int) ([]byte, error) {
	log.Debug("pull requests to be fetched from pull_request no. %v repository %v", fromNo, gc.RepositoryName)
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}
	var allPullRequests []*github.PullRequest
	for {
		pullRequests, resp, err := gc.Client.PullRequests.List(gc.ctx, gc.RepositoryOwner, gc.RepositoryName, opt)
		if err != nil {
			log.Errorf("error[%v] in fetching pull requests for repository %v", err, gc.RepositoryName)
			return nil, err
		}
		for _, pr := range pullRequests {
			if pr.GetNumber() > fromNo {
				allPullRequests = append(allPullRequests, pr)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	allPullRequestsByte, err := json.Marshal(allPullRequests)
	return allPullRequestsByte, err
}

// GetIssues fetches issues for the user
func (gc *GithubClient) GetIssues(from time.Time) ([]byte, error) {
	log.Debug("issues to be fetched after %v for repository %v", from, gc.RepositoryName)
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
		Since:       from,
	}
	var allIssues []*github.Issue
	for {
		issues, resp, err := gc.Client.Issues.ListByRepo(gc.ctx, gc.RepositoryOwner, gc.RepositoryName, opt)
		if err != nil {
			log.Errorf("error[%v] in fetching issues for repository %v", err, gc.RepositoryName)
			return nil, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	allIssuesByte, err := json.Marshal(allIssues)
	return allIssuesByte, err
}
