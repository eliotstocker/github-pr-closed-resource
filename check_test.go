package resource_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/telia-oss/github-pr-resource"
	"github.com/telia-oss/github-pr-resource/mocks"
)

var (
	testPullRequests = []*resource.PullRequest{
		createTestPR(1, true, false),
		createTestPR(2, false, false),
		createTestPR(3, true, false),
		createTestPR(4, false, false),
		createTestPR(5, true, true),
		createTestPR(6, false, false),
	}
)

func TestCheck(t *testing.T) {
	tests := []struct {
		description  string
		source       resource.Source
		version      resource.Version
		files        [][]string
		pullRequests []*resource.PullRequest
		expected     resource.CheckResponse
	}{
		{
			description: "check returns the latest version if there is no previous",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.Version{},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[0]),
			},
		},

		{
			description: "check returns the previous version when its still latest",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.NewVersion(testPullRequests[0]),
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[0]),
			},
		},

		{
			description: "check returns all new versions since the last",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			version:      resource.NewVersion(testPullRequests[3]),
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[2]),
				resource.NewVersion(testPullRequests[1]),
				resource.NewVersion(testPullRequests[0]),
			},
		},

		{
			description: "check correctly ignores cross repo pull requests",
			source: resource.Source{
				Repository:   "itsdalmo/test-repository",
				AccessToken:  "oauthtoken",
				DisableForks: true,
			},
			version:      resource.NewVersion(testPullRequests[5]),
			pullRequests: testPullRequests,
			expected: resource.CheckResponse{
				resource.NewVersion(testPullRequests[3]),
				resource.NewVersion(testPullRequests[2]),
				resource.NewVersion(testPullRequests[1]),
				resource.NewVersion(testPullRequests[0]),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			github := mocks.NewMockGithub(ctrl)
			github.EXPECT().ListClosedPullRequests().Times(1).Return(tc.pullRequests, nil)

			input := resource.CheckRequest{Source: tc.source, Version: tc.version}
			output, err := resource.Check(input, github)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got, want := output, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
		})
	}
}
