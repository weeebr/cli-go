package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"cli-go/_internal/config"
)

// Client represents a Jira API client
type Client struct {
	BaseURL        string
	Email          string
	APIToken       string
	DefaultProject string
	HTTPClient     *http.Client
}

// NewClient creates a new Jira API client
func NewClient(baseURL, email, apiToken, defaultProject string) *Client {
	// Load config to get timeout setting
	cfg, err := config.LoadConfig()
	timeout := 30 * time.Second // Default fallback
	if err == nil {
		timeout = time.Duration(cfg.Network.TimeoutSeconds) * time.Second
	}

	return &Client{
		BaseURL:        baseURL,
		Email:          email,
		APIToken:       apiToken,
		DefaultProject: defaultProject,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// createAuthHeader creates Basic Auth header
func (c *Client) createAuthHeader() string {
	auth := c.Email + ":" + c.APIToken
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// makeRequest makes an HTTP request to the Jira API
func (c *Client) makeRequest(method, endpoint string) ([]byte, error) {
	url := c.BaseURL + "/rest/api/3/" + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// GetIssue retrieves an issue by key
func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	normalizedKey := NormalizeIssueKey(issueKey, c.DefaultProject)
	endpoint := fmt.Sprintf("issue/%s?fields=summary,status,assignee,reporter,description,customfield_10087,customfield_10093,customfield_10077", url.PathEscape(normalizedKey))

	data, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse issue data: %v", err)
	}

	return &issue, nil
}

// GetComments retrieves comments for an issue
func (c *Client) GetComments(issueKey string) ([]Comment, error) {
	normalizedKey := NormalizeIssueKey(issueKey, c.DefaultProject)
	endpoint := fmt.Sprintf("issue/%s/comment", url.PathEscape(normalizedKey))

	data, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var result struct {
		Comments []Comment `json:"comments"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse comments data: %v", err)
	}

	return result.Comments, nil
}

// GetIssueWithComments retrieves an issue with its comments
func (c *Client) GetIssueWithComments(issueKey string) (*Issue, error) {
	issue, err := c.GetIssue(issueKey)
	if err != nil {
		return nil, err
	}

	comments, err := c.GetComments(issueKey)
	if err != nil {
		// If comments fail, still return the issue
		return issue, nil
	}

	issue.Fields.Comments.Comments = comments
	return issue, nil
}

// GetChangelog retrieves changelog for an issue
func (c *Client) GetChangelog(issueKey string) (*Issue, error) {
	normalizedKey := NormalizeIssueKey(issueKey, c.DefaultProject)
	endpoint := fmt.Sprintf("issue/%s?expand=changelog", url.PathEscape(normalizedKey))

	data, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse changelog data: %v", err)
	}

	return &issue, nil
}

// GetCurrentUser retrieves current user information
func (c *Client) GetCurrentUser() (*User, error) {
	data, err := c.makeRequest("GET", "myself")
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %v", err)
	}

	return &user, nil
}

// SearchUsers searches for users by query
func (c *Client) SearchUsers(query string) ([]User, error) {
	encodedQuery := url.QueryEscape(query)
	endpoint := fmt.Sprintf("user/search?query=%s", encodedQuery)

	data, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse user search data: %v", err)
	}

	return users, nil
}

// SearchJQL searches issues using JQL
func (c *Client) SearchJQL(jql string, maxResults int) (*SearchResults, error) {
	if maxResults <= 0 {
		maxResults = 10
	}

	encodedJQL := url.QueryEscape(jql)
	endpoint := fmt.Sprintf("search/jql?jql=%s&maxResults=%d&fields=key,summary", encodedJQL, maxResults)

	data, err := c.makeRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var results SearchResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %v", err)
	}

	return &results, nil
}

// GetUserActivityJQL retrieves recent activity for a user
func (c *Client) GetUserActivityJQL(userAccountID string, maxResults int) (*SearchResults, error) {
	var jql string
	if userAccountID != "" {
		jql = fmt.Sprintf("updated >= -7d AND (assignee = \"%s\" OR reporter = \"%s\" OR watcher = \"%s\") ORDER BY updated DESC",
			userAccountID, userAccountID, userAccountID)
	} else {
		jql = "updated >= -7d AND (assignee = currentUser() OR reporter = currentUser() OR watcher = currentUser()) ORDER BY updated DESC"
	}

	return c.SearchJQL(jql, maxResults)
}

// GetUserModifiedIssues retrieves issues modified by a user
func (c *Client) GetUserModifiedIssues(userAccountID string, maxResults int) (*SearchResults, error) {
	var jql string
	if userAccountID != "" {
		jql = fmt.Sprintf("updated >= -7d AND (assignee = \"%s\" OR reporter = \"%s\") ORDER BY updated DESC",
			userAccountID, userAccountID)
	} else {
		jql = "updated >= -7d AND (assignee = currentUser() OR reporter = currentUser()) ORDER BY updated DESC"
	}

	return c.SearchJQL(jql, maxResults)
}

// GetIssueWithChangelog retrieves an issue with its changelog
func (c *Client) GetIssueWithChangelog(issueKey string) (*Issue, error) {
	return c.GetChangelog(issueKey)
}

// GetUserActivity retrieves user activity (viewed, created, updated issues)
func (c *Client) GetUserActivity(userName string) (*SearchResults, *SearchResults, *SearchResults, error) {
	// Get current user if no userName provided
	var userAccountID string
	if userName == "" {
		user, err := c.GetCurrentUser()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get current user: %v", err)
		}
		userAccountID = user.AccountID
	} else {
		// Search for user by name
		users, err := c.SearchUsers(userName)
		if err != nil || len(users) == 0 {
			return nil, nil, nil, fmt.Errorf("user not found: %s", userName)
		}
		userAccountID = users[0].AccountID
	}

	// Get viewed issues (recently updated)
	viewed, err := c.GetUserActivityJQL(userAccountID, 10)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get viewed issues: %v", err)
	}

	// Get created issues
	created, err := c.GetUserModifiedIssues(userAccountID, 10)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get created issues: %v", err)
	}

	// Get updated issues (same as modified for now)
	updated, err := c.GetUserModifiedIssues(userAccountID, 10)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get updated issues: %v", err)
	}

	return viewed, created, updated, nil
}

// SearchIssues searches issues using JQL with default parameters
func (c *Client) SearchIssues(jql string) (*SearchResults, error) {
	return c.SearchJQL(jql, 10) // Default to 10 results
}

// GetUserViewedIssues gets issues recently viewed by a user
func (c *Client) GetUserViewedIssues(username string) (*SearchResults, error) {
	// Search for issues recently updated where user is assignee, reporter, or watcher
	jql := fmt.Sprintf("updated >= -7d AND (assignee = \"%s\" OR reporter = \"%s\" OR watcher = \"%s\") ORDER BY updated DESC",
		username, username, username)
	return c.SearchJQL(jql, 10)
}

// GetUserCreatedIssues gets issues created by a user
func (c *Client) GetUserCreatedIssues(username string) (*SearchResults, error) {
	// Search for issues created by user in the last 7 days
	jql := fmt.Sprintf("created >= -7d AND reporter = \"%s\" ORDER BY created DESC", username)
	return c.SearchJQL(jql, 10)
}

// GetUserUpdatedIssues gets issues updated by a user
func (c *Client) GetUserUpdatedIssues(username string) (*SearchResults, error) {
	// Search for issues updated by user in the last 7 days
	jql := fmt.Sprintf("updated >= -7d AND (assignee = \"%s\" OR reporter = \"%s\") ORDER BY updated DESC",
		username, username)
	return c.SearchJQL(jql, 10)
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(summary, description, issueType string) (*Issue, error) {
	// This is a placeholder implementation
	// In a real implementation, you'd make a POST request to create an issue
	return nil, fmt.Errorf("CreateIssue not implemented")
}

// UpdateIssue updates an existing issue
func (c *Client) UpdateIssue(issueKey, summary, description string) (*Issue, error) {
	// This is a placeholder implementation
	// In a real implementation, you'd make a PUT request to update an issue
	return nil, fmt.Errorf("UpdateIssue not implemented")
}

// AddComment adds a comment to an issue
func (c *Client) AddComment(issueKey, comment string) error {
	// This is a placeholder implementation
	// In a real implementation, you'd make a POST request to add a comment
	return fmt.Errorf("AddComment not implemented")
}

// AssignIssue assigns an issue to a user
func (c *Client) AssignIssue(issueKey, assignee string) error {
	// This is a placeholder implementation
	// In a real implementation, you'd make a PUT request to assign an issue
	return fmt.Errorf("AssignIssue not implemented")
}

// TransitionIssue transitions an issue to a new status
func (c *Client) TransitionIssue(issueKey, transition string) error {
	// This is a placeholder implementation
	// In a real implementation, you'd make a POST request to transition an issue
	return fmt.Errorf("TransitionIssue not implemented")
}
