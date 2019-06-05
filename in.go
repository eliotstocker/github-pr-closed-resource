package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Get (business logic)
func Get(request GetRequest, github Github, outputDir string) (*GetResponse, error) {
	if request.Params.SkipDownload {
		return &GetResponse{Version: request.Version}, nil
	}

	pull, err := github.GetPullRequest(request.Version.PR, request.Version.Commit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pull request: %s", err)
	}

	// Create the metadata
	var metadata Metadata
	metadata.Add("pr", strconv.Itoa(pull.Number))
	metadata.Add("url", pull.URL)
	metadata.Add("head_name", pull.HeadRefName)
	metadata.Add("head_sha", pull.Tip.OID)
	metadata.Add("base_name", pull.BaseRefName)
	metadata.Add("message", pull.Tip.Message)
	metadata.Add("author", pull.Tip.Author.User.Login)

	// Write version and metadata for reuse in PUT

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %s", err)
	}
	b, err := json.Marshal(request.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal version: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, "version.json"), b, 0644); err != nil {
		return nil, fmt.Errorf("failed to write version: %s", err)
	}
	b, err = json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, "metadata.json"), b, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, "commit"), []byte(request.Version.Commit), 0644); err != nil {
        return nil, fmt.Errorf("failed to write commit: %s", err)
    }
    if err := ioutil.WriteFile(filepath.Join(outputDir, "pr"), []byte(request.Version.PR), 0644); err != nil {
        return nil, fmt.Errorf("failed to write pr: %s", err)
    }

	return &GetResponse{
		Version:  request.Version,
		Metadata: metadata,
	}, nil
}

// GetParameters ...
type GetParameters struct {
	SkipDownload bool `json:"skip_download"`
}

// GetRequest ...
type GetRequest struct {
	Source  Source        `json:"source"`
	Version Version       `json:"version"`
	Params  GetParameters `json:"params"`
}

// GetResponse ...
type GetResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata,omitempty"`
}
