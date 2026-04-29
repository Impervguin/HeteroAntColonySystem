package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ListFilesResponse struct {
	Files []string `json:"files"`
}

func GetFiles(apiServer, apiBase string) (*ListFilesResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s%s/tsp/files", apiServer, apiBase))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var files ListFilesResponse
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	return &files, nil
}
