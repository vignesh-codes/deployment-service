package svc

import (
	"context"
	adapter "deployment-service/apps/repository/adapter"
	"deployment-service/logger"
	model_build "deployment-service/models/model.build"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"go.uber.org/zap"
)

type BuildService struct {
	repository *adapter.Repository
}

func (svc BuildService) CreateNewRepoScout(payload model_build.RepoScout) (map[string]interface{}, error) {
	fmt.Println("payload is ", payload)
	result, err := svc.repository.MongoDB.InsertOne("REPO_SCOUTS", payload)
	if err != nil {
		logger.Logger.Error("Error while inserting new repo scout", zap.Any(logger.KEY_ERROR, err.Error()))
	}
	return map[string]interface{}{
		"result": payload,
		"data":   result,
	}, nil
}

func (svc BuildService) GetAllRepoScouts(namespace string) ([]model_build.RepoScout, error) {
	// Prepare the query filter to match the specified namespace
	filter := bson.D{
		{"namespace", namespace},
	}

	// Fetch all documents that match the filter from the REPO_SCOUTS collection
	cursor, err := svc.repository.MongoDB.GetAll("REPO_SCOUTS", filter)
	defer cursor.Close(context.TODO()) // Close the cursor after we're done

	var result = []model_build.RepoScout{}
	for cursor.Next(context.TODO()) {
		var repoScout model_build.RepoScout
		if err := cursor.Decode(&repoScout); err != nil {
			fmt.Println("Error decoding document:", err)
			return nil, err
		}
		result = append(result, repoScout)
	}

	// Check for cursor iteration errors
	if err := cursor.Err(); err != nil {
		fmt.Println("Error iterating over cursor:", err)
		return nil, err
	}

	if err != nil {
		// Log the error if fetching fails
		logger.Logger.Error("Error while fetching all repo scouts", zap.Any(logger.KEY_ERROR, err.Error()))
		return nil, err
	}

	// Return the fetched result
	return result, nil
}

func (svc BuildService) GetReleaseInfo(repoURLs []string, topK int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// Iterate over each repository URL
	for _, repoURL := range repoURLs {
		// Build the GitHub API URL for releases
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases", repoURL)
		// Send HTTP GET request to GitHub API
		resp, err := http.Get(apiURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch releases for %s: %w", repoURL, err)
		}
		defer resp.Body.Close()

		// Check the HTTP status code
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, repoURL)
		}

		// Log the body for debugging
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for %s: %w", repoURL, err)
		}

		// Print the body to see what was returned
		// fmt.Println("Response Body:", string(body))

		// Decode the response
		var releases []model_build.ReleaseInfo
		if err := json.Unmarshal(body, &releases); err != nil {
			return nil, fmt.Errorf("failed to decode response for %s: %w", repoURL, err)
		}
		// Select the top 3 releases
		topReleases := []map[string]interface{}{}
		for i, release := range releases {
			if i >= topK {
				break
			}

			// Create a map for each release with the necessary fields
			topReleases = append(topReleases, map[string]interface{}{
				"html_url":     release.HtmlURL,
				"tag_name":     release.TagName,
				"created_at":   release.CreatedAt,
				"published_at": release.PublishedAt,
			})
		}

		// Append to the result with the formatted repo info
		result = append(result, map[string]interface{}{
			"repo_url": "https://github.com/" + repoURL, // Keep the repository URL for reference
			"releases": topReleases,                     // Append the top releases for the repo
		})
	}

	return result, nil
}
