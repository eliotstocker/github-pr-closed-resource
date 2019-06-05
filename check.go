package resource

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Check (business logic)
func Check(request CheckRequest, manager Github) (CheckResponse, error) {
	var response CheckResponse

	pulls, err := manager.ListClosedPullRequests()
	if err != nil {
		return nil, fmt.Errorf("failed to get last commits: %s", err)
	}

	for _, p := range pulls {
        // Filter out commits that are too old.
        if !p.ClosedAt.Time.IsZero() {
            if !p.ClosedAt.Time.After(request.Version.ClosedDate) {
                continue
            }
        } else if !p.MergedAt.Time.IsZero() {
            if !p.MergedAt.Time.After(request.Version.ClosedDate) {
                continue
            }
        } else {
            continue
        }

		// Filter to only requested PRs
		if len(request.Source.Filter) > 0 && !contains(request.Source.Filter, p.Number) {
		    continue
		}

		if request.Source.DisableForks && p.IsCrossRepository {
            continue
        }

		response = append(response, NewVersion(p))
	}

	// Sort the commits by date
	sort.Sort(response)

	// If there are no new but an old version = return the old
	if len(response) == 0 && request.Version.PR != "" {
		response = append(response, request.Version)
	}
	// If there are new versions and no previous = return just the latest
	if len(response) != 0 && request.Version.PR == "" {
		response = CheckResponse{response[len(response)-1]}
	}
	return response, nil
}

// IsInsidePath checks whether the child path is inside the parent path.
//
// /foo/bar is inside /foo, but /foobar is not inside /foo.
// /foo is inside /foo, but /foo is not inside /foo/
func IsInsidePath(parent, child string) bool {
	if parent == child {
		return true
	}

	// we add a trailing slash so that we only get prefix matches on a
	// directory separator
	parentWithTrailingSlash := parent
	if !strings.HasSuffix(parentWithTrailingSlash, string(filepath.Separator)) {
		parentWithTrailingSlash += string(filepath.Separator)
	}

	return strings.HasPrefix(child, parentWithTrailingSlash)
}

// CheckRequest ...
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// CheckResponse ...
type CheckResponse []Version

func (r CheckResponse) Len() int {
	return len(r)
}

func (r CheckResponse) Less(i, j int) bool {
	return r[j].ClosedDate.After(r[i].ClosedDate)
}

func (r CheckResponse) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func contains(arr []int, find int) bool {
   for i := range arr {
      if arr[i] == find {
         return true
      }
   }
   return false
}
