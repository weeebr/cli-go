package jira

// Issue represents a Jira issue
type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary     string      `json:"summary"`
		Description interface{} `json:"description"` // Can be string or ADF object
		Status      struct {
			Name string `json:"name"`
		} `json:"status"`
		Assignee struct {
			DisplayName string `json:"displayName"`
			Email       string `json:"emailAddress"`
		} `json:"assignee"`
		Reporter struct {
			DisplayName string `json:"displayName"`
			Email       string `json:"emailAddress"`
		} `json:"reporter"`
		Comments struct {
			Comments []Comment `json:"comments"`
		} `json:"comments"`
		// Custom fields for testing instructions
		CustomField10087 interface{} `json:"customfield_10087"` // Testing instructions
		CustomField10093 interface{} `json:"customfield_10093"` // Additional custom field
		CustomField10077 interface{} `json:"customfield_10077"` // Additional custom field
	} `json:"fields"`
	Changelog *Changelog `json:"changelog,omitempty"`
}

// Comment represents a Jira comment
type Comment struct {
	ID      string      `json:"id"`
	Author  User        `json:"author"`
	Body    interface{} `json:"body"` // Can be string or ADF object
	Created string      `json:"created"`
	Updated string      `json:"updated"`
}

// User represents a Jira user
type User struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
}

// Changelog represents issue changelog
type Changelog struct {
	Histories []History `json:"histories"`
}

// History represents a changelog history entry
type History struct {
	ID      string `json:"id"`
	Author  User   `json:"author"`
	Created string `json:"created"`
	Items   []Item `json:"items"`
}

// Item represents a changelog item
type Item struct {
	Field      string `json:"field"`
	FieldType  string `json:"fieldtype"`
	FromString string `json:"fromString"`
	ToString   string `json:"toString"`
}

// SearchResults represents JQL search results
type SearchResults struct {
	Issues []Issue `json:"issues"`
	Total  int     `json:"total"`
}

// UserSearchResults represents user search results
type UserSearchResults []User

// NormalizeIssueKey converts numeric shorthand to full issue key
func NormalizeIssueKey(issueKey, defaultProject string) string {
	// If it's already a full key (PROJECT-123), return as-is
	// Check if it contains a dash (PROJECT-123 format)
	for i, char := range issueKey {
		if char == '-' && i > 0 {
			return issueKey
		}
	}

	// If it's just a number, add default project
	if len(issueKey) > 0 && issueKey[0] >= '0' && issueKey[0] <= '9' {
		return defaultProject + "-" + issueKey
	}

	// Return as-is for other cases
	return issueKey
}
