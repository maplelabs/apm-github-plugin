/* Package gitprovider gives utilites to deal with cloud vcs provider APIs */
package gitprovider

import (
	"time"

	"github.com/maplelabs/github-audit/logger"
)

var (
	log logger.Logger
)

func init() {
	log = logger.GetLogger()
}

// GitProvider represents git cloud provider which will be used to pull metrics
type GitProvider interface {
	// CheckCredentials checks user credentials
	CheckCredentials() error

	// GetCommits fetches commits using APIs
	GetCommits(from time.Time, to time.Time, branch string) ([]byte, error)

	// GetPullRequests fetches pull requests using APIs
	GetPullRequests(int) ([]byte, error)

	// GetIssues fetches issues using APIs
	GetIssues(to time.Time) ([]byte, error)
}

// NewGitProvider returns a new git provider based on git cloud type
func NewGitProvider(host string, repoOwner string, repoName string, userName string, accessToken string) GitProvider {
	if host == "github" {
		return NewGithubClient(repoOwner, repoName, userName, accessToken)
	}
	return nil
}
