package resource

import (
	"errors"
	"strconv"
	"time"

	"github.com/shurcooL/githubv4"
)

// Source represents the configuration for the resource.
type Source struct {
	Repository          string   `json:"repository"`
	AccessToken         string   `json:"access_token"`
	V3Endpoint          string   `json:"v3_endpoint"`
	V4Endpoint          string   `json:"v4_endpoint"`
	Paths               []string `json:"paths"`
	IgnorePaths         []string `json:"ignore_paths"`
	DisableCISkip       bool     `json:"disable_ci_skip"`
	SkipSSLVerification bool     `json:"skip_ssl_verification"`
	DisableForks        bool     `json:"disable_forks"`
	GitCryptKey         string   `json:"git_crypt_key"`
	Filter              []int    `json:"filter"`
}

// Validate the source configuration.
func (s *Source) Validate() error {
	if s.AccessToken == "" {
		return errors.New("access_token must be set")
	}
	if s.Repository == "" {
		return errors.New("repository must be set")
	}
	if s.V3Endpoint != "" && s.V4Endpoint == "" {
		return errors.New("v4_endpoint must be set together with v3_endpoint")
	}
	if s.V4Endpoint != "" && s.V3Endpoint == "" {
		return errors.New("v3_endpoint must be set together with v4_endpoint")
	}
	return nil
}

// Metadata output from get/put steps.
type Metadata []*MetadataField

// Add a MetadataField to the Metadata.
func (m *Metadata) Add(name, value string) {
	*m = append(*m, &MetadataField{Name: name, Value: value})
}

// MetadataField ...
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Version communicated with Concourse.
type Version struct {
	PR            string    `json:"pr"`
	Commit        string    `json:"commit"`
	ClosedDate    time.Time `json:"committed,omitempty"`
}

// NewVersion constructs a new Version.
func NewVersion(p *PullRequest) Version {
    var closed time.Time
    if !p.ClosedAt.Time.IsZero() {
        closed = p.ClosedAt.Time
    } else if !p.MergedAt.Time.IsZero() {
        closed = p.MergedAt.Time
    }
	return Version{
		PR:            strconv.Itoa(p.Number),
		Commit:        p.Tip.OID,
		ClosedDate:    closed,
	}
}

// PullRequest represents a pull request and includes the tip (commit).
type PullRequest struct {
	PullRequestObject
	Tip CommitObject
}

// PullRequestObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/pullrequest/
type PullRequestObject struct {
	ID          string
	Number      int
	Title       string
	URL         string
	BaseRefName string
	HeadRefName string
	Repository  struct {
		URL string
	}
	ClosedAt    githubv4.DateTime
	MergedAt    githubv4.DateTime
	IsCrossRepository bool
}

// CommitObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/commit/
type CommitObject struct {
	ID            string
	OID           string
	CommittedDate githubv4.DateTime
	Message       string
	Author        struct {
		User struct {
			Login string
		}
	}
}
