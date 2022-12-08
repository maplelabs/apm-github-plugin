package dataprocessor

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/maplelabs/github-audit/metricformator"
)

const (
	COMMIT      = "commit"
	PULLREQUEST = "pull_request"
	ISSUE       = "issue"
	GITHUB      = "github"
)

// GithubProcessor process data from github APIs
type GithubProcessor struct {
	// Repository Name
	RepoName string

	// Repository URL
	RepoURL string

	// Current time in milliseconds
	CurrentTimeInMS int64

	// Metricformator instance to customise processed data
	MetricFormator *metricformator.MetricFormator
}

// NewGithubProcessor provides new instance of github api processor
func NewGithubProcessor(repoName string, repoURL string) GithubProcessor {
	var gp GithubProcessor
	gp.RepoName = repoName
	gp.RepoURL = repoURL
	gp.CurrentTimeInMS = time.Now().UnixNano() / 1000000
	gp.MetricFormator = metricformator.NewMetricFormator()
	return gp
}

// Commit represents commit document
type Commit struct {
	// DocumentType is "commit"
	DocumentType string `json:"document_type"`

	// RepoType is "github" , represents the git provider
	RepoType string `json:"repo_type"`

	// RepoName is repository name
	RepoName string `json:"repo_name"`

	// RepoURL is repository url
	RepoURL string `json:"repo_url"`

	// CommitURL is the github api url to commit
	CommitURL string `json:"commit_url"`

	// CreatedAt represents at what time this commit was created
	CreatedAt time.Time `json:"created_at"`

	// Message represents commit message
	Message string `json:"message"`

	// Committer provides info related to user who commited changes
	Committer User `json:"committer"`

	// Sha represents commit sha
	Sha string `json:"sha"`

	// time in milliseconds
	Time int64 `json:"time"`
}

// User represents a git user
type User struct {
	// ID of the user
	ID string `json:"id"`

	// User contains user name
	User string `json:"user"`
}

// PullRequest represents pull reequest document
type PullRequest struct {
	// DocumentType is "pull_request"
	DocumentType string `json:"document_type"`

	// RepoType is "github" , represents the git provider
	RepoType string `json:"repo_type"`

	// RepoName is repository name
	RepoName string `json:"repo_name"`

	// RepoURL is repository url
	RepoURL string `json:"repo_url"`

	// CreatedAt represents at what time this pull request is created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt represents at what time this pull request is updated
	UpdatedAt time.Time `json:"updated_at"`

	// ClosedAt represents at what time this pull request is closed
	ClosedAt time.Time `json:"closed_at"`

	// State represents the state of pull request
	State string `json:"state"`

	// PullRequestNo represents pull request number
	PullRequestNo string `json:"pull_request_no"`

	// Title represents pull request title
	Title string `json:"title"`

	// URL is the api url for pull request
	URL string `json:"url"`

	// MergedAt represents at what time this pull request is merged
	MergedAt time.Time `json:"merged_at"`

	// MergeCommitSha represents pull request sha
	MergeCommitSha string `json:"merge_commit_sha"`

	// Reviewers holds list of reviewrs
	Reviewers []User `json:"reviewers"`

	// RequestFromRepo shows from where this pull request is raised
	RequestFromRepo RequestFromRepository `json:"request_from_repo"`

	// MergeToRepo shows the base repository where pull request will be merged
	MergeToRepo MergeToRepository `json:"merge_to_repo"`

	// time in milliseconds
	Time int64 `json:"time"`
}

// RequestFromRepository represents from where pull request is raised
type RequestFromRepository struct {
	// Name of repository
	Name string `json:"name"`

	// URL is api url to repository
	URL string `json:"url"`

	// Private is whether this repository is private or public
	Private bool `json:"private"`

	// Sha is the sha to last commit to repository
	Sha string `json:"sha"`

	// Branch represents the branch to repo
	Branch string `json:"branch"`

	// ByUser repesents the user initiating request
	ByUser User `json:"by_user"`
}

// RequestFromRepository represents to where pull request will be merged
type MergeToRepository struct {
	// Name of repository
	Name string `json:"name"`

	// URL is api url to repository
	URL string `json:"url"`

	// Private is whether this repository is private or public
	Private bool `json:"private"`

	// Sha is the sha to last commit to repository
	Sha string `json:"sha"`

	// Branch represents the branch to repo
	Branch string `json:"branch"`
}

// Issue represents issue document
type Issue struct {
	// DocumentType is "pull_request"
	DocumentType string `json:"document_type"`

	// RepoType is "github" , represents the git provider
	RepoType string `json:"repo_type"`

	// RepoName is repository name
	RepoName string `json:"repo_name"`

	// RepoURL is repository url
	RepoURL string `json:"repo_url"`

	// IssueNo is issue number
	IssueNo string `json:"issue_no"`

	// State represents the state of issue
	State string `json:"state"`

	// Title represents issue title
	Title string `json:"title"`

	// CreatedAt represents at what time this issue is created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt represents at what time this issue is updated
	UpdatedAt time.Time `json:"updated_at"`

	// ClosedAt represents at what time this issue is closed
	ClosedAt string `json:"closed_at"`

	// URL is api url to issue
	URL string `json:"url"`

	// CreatedBy shows the user who created it
	CreatedBy User `json:"created_by"`

	// Assignees represents issue assignees
	Assignees []User `json:"assignees"`

	// time in milliseconds
	Time int64 `json:"time"`
}

// ProcessCommits prepares commit output documents
func (g GithubProcessor) ProcessCommits(data []byte, tags map[string]string) ([]interface{}, error) {
	var commits []github.RepositoryCommit
	commitDocuments := make([]interface{}, 0)
	err := json.Unmarshal(data, &commits)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling commits for repository %v", err, g.RepoName)
		return commitDocuments, err
	}
	for _, c := range commits {
		var commit Commit
		commit.RepoName = g.RepoName
		commit.RepoURL = g.RepoURL
		commit.DocumentType = COMMIT
		commit.Message = c.Commit.GetMessage()
		commit.RepoType = GITHUB
		commit.CommitURL = c.GetURL()
		commit.Sha = c.GetSHA()
		commit.CreatedAt = c.Commit.Committer.GetDate().Local()
		commit.Committer.ID = strconv.FormatInt(c.Committer.GetID(), 10)
		commit.Committer.User = c.Commit.Author.GetName()
		commit.Time = g.CurrentTimeInMS
		commitDocuments = append(commitDocuments, commit)
	}
	b, _ := json.Marshal(commitDocuments)
	b = g.MetricFormator.CustomizeMetrics(b)
	finalDocs := AddTags(b, tags)
	return finalDocs, err
}

// ProcessPullRequests prepares pull request output documents
func (g GithubProcessor) ProcessPullRequests(data []byte, tags map[string]string) ([]interface{}, error) {
	var pullRequests []github.PullRequest
	prDocuments := make([]interface{}, 0)
	err := json.Unmarshal(data, &pullRequests)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling pull requests for repository %v", err, g.RepoName)
		return prDocuments, err
	}
	for _, p := range pullRequests {
		var pr PullRequest
		pr.PullRequestNo = strconv.Itoa(p.GetNumber())
		pr.DocumentType = PULLREQUEST
		pr.RepoType = GITHUB
		pr.RepoName = g.RepoName
		pr.RepoURL = g.RepoURL
		pr.CreatedAt = p.GetCreatedAt().Local()
		pr.UpdatedAt = p.GetUpdatedAt().Local()
		pr.ClosedAt = p.GetClosedAt().Local()
		pr.State = p.GetState()
		pr.URL = p.GetURL()
		pr.Title = p.GetTitle()
		pr.MergedAt = p.GetMergedAt().Local()
		pr.MergeCommitSha = p.GetMergeCommitSHA()
		pr.Time = g.CurrentTimeInMS
		var reqFromRepo RequestFromRepository
		reqFromRepo.Branch = p.Head.GetRef()
		reqFromRepo.ByUser.ID = strconv.FormatInt(p.Head.User.GetID(), 10)
		reqFromRepo.ByUser.User = p.Head.User.GetLogin()
		reqFromRepo.Name = p.Head.Repo.GetFullName()
		reqFromRepo.Private = p.Head.Repo.GetPrivate()
		reqFromRepo.URL = p.Head.Repo.GetURL()
		reqFromRepo.Sha = p.Head.GetSHA()
		pr.RequestFromRepo = reqFromRepo
		var mergeToRepo MergeToRepository
		mergeToRepo.Name = p.Base.Repo.GetFullName()
		mergeToRepo.Branch = p.Base.GetRef()
		mergeToRepo.Private = p.Base.Repo.GetPrivate()
		mergeToRepo.Sha = p.Base.GetSHA()
		mergeToRepo.URL = p.Base.Repo.GetURL()
		pr.MergeToRepo = mergeToRepo
		var reviewers []User
		for _, rr := range p.RequestedReviewers {
			var u User
			u.ID = strconv.FormatInt(rr.GetID(), 10)
			u.User = rr.GetLogin()
			reviewers = append(reviewers, u)
		}
		pr.Reviewers = reviewers
		prDocuments = append(prDocuments, pr)
	}

	b, _ := json.Marshal(prDocuments)
	b = g.MetricFormator.CustomizeMetrics(b)
	finalDocs := AddTags(b, tags)
	return finalDocs, nil
}

// ProcessIssues prepares issue output documents
func (g GithubProcessor) ProcessIssues(data []byte, tags map[string]string) ([]interface{}, error) {
	var issues []github.Issue
	issueDocuments := make([]interface{}, 0)
	err := json.Unmarshal(data, &issues)
	if err != nil {
		log.Errorf("error[%v] in unmarshalling pull requests for repository %v", err, g.RepoName)
		return issueDocuments, err
	}
	for _, i := range issues {
		// only capturing issues
		if i.PullRequestLinks == nil {
			var issue Issue
			issue.DocumentType = ISSUE
			issue.RepoType = GITHUB
			issue.RepoName = g.RepoName
			issue.RepoURL = g.RepoURL
			issue.IssueNo = strconv.Itoa(i.GetNumber())
			issue.Title = i.GetTitle()
			issue.URL = i.GetURL()
			issue.State = i.GetState()
			issue.Time = g.CurrentTimeInMS
			issue.CreatedAt = i.GetCreatedAt().Local()
			issue.UpdatedAt = i.GetUpdatedAt().Local()
			issue.CreatedBy.ID = strconv.FormatInt(i.User.GetID(), 10)
			issue.CreatedBy.User = i.User.GetLogin()
			var assignees []User
			for _, rr := range i.Assignees {
				var u User
				u.ID = strconv.FormatInt(rr.GetID(), 10)
				u.User = rr.GetLogin()
				assignees = append(assignees, u)
			}
			issue.Assignees = assignees
			issueDocuments = append(issueDocuments, issue)
		}
	}
	b, _ := json.Marshal(issueDocuments)
	b = g.MetricFormator.CustomizeMetrics(b)
	finalDocs := AddTags(b, tags)
	return finalDocs, err
}
