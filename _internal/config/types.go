package config

// RepoConfig represents a repository configuration
type RepoConfig struct {
	Owner      string `json:"owner" yaml:"owner"`
	Repo       string `json:"repo" yaml:"repo"`
	Path       string `json:"path" yaml:"path"`
	MainBranch string `json:"mainBranch" yaml:"main_branch"`
	Main       bool   `json:"main" yaml:"main"`
}

// Config represents the unified configuration
type Config struct {
	Ringier struct {
		DefaultProjectKey string `json:"defaultProjectKey" yaml:"default_project_key"`
		DefaultUser       string `json:"defaultUser" yaml:"default_user"`
	} `json:"ringier" yaml:"ringier"`

	Jira struct {
		BaseURL        string `json:"baseURL" yaml:"base_url"`
		Email          string `json:"email" yaml:"email"`
		DefaultProject string `json:"defaultProject" yaml:"default_project"`
	} `json:"jira" yaml:"jira"`

	Repositories []RepoConfig `json:"repositories" yaml:"repositories"`

	History struct {
		DefaultDays  int    `json:"defaultDays" yaml:"default_days"`
		AuthorFilter string `json:"authorFilter" yaml:"author_filter"`
	} `json:"history" yaml:"history"`

	Tests []string `json:"tests" yaml:"tests"`

	// AI configuration
	AI struct {
		Models struct {
			OpenAI    string `json:"openai" yaml:"openai"`
			Anthropic string `json:"anthropic" yaml:"anthropic"`
			Google    string `json:"google" yaml:"google"`
			Groq      string `json:"groq" yaml:"groq"`
		} `json:"models" yaml:"models"`
		Timeouts struct {
			Default int `json:"default" yaml:"default"`
		} `json:"timeouts" yaml:"timeouts"`
	} `json:"ai" yaml:"ai"`

	// Network configuration
	Network struct {
		TimeoutSeconds int `json:"timeoutSeconds" yaml:"timeout_seconds"`
		RetryAttempts  int `json:"retryAttempts" yaml:"retry_attempts"`
	} `json:"network" yaml:"network"`

	// Prompts configuration
	Prompts struct {
		BaseDir string `json:"baseDir" yaml:"base_dir"`
	} `json:"prompts" yaml:"prompts"`

	// Cache configuration
	Cache struct {
		BaseDir   string `json:"baseDir" yaml:"base_dir"`
		MaxSizeMB int    `json:"maxSizeMB" yaml:"max_size_mb"`
		TTLDays   int    `json:"ttlDays" yaml:"ttl_days"`
	} `json:"cache" yaml:"cache"`

	// Server configuration
	Server struct {
		Port         int    `json:"port" yaml:"port"`
		ViteURL      string `json:"viteURL" yaml:"vite_url"`
		APIPort      int    `json:"apiPort" yaml:"api_port"`
		FrontendType string `json:"frontendType" yaml:"frontend_type"`
		Environment  string `json:"environment" yaml:"environment"`
	} `json:"server" yaml:"server"`

	// Credentials configuration
	Credentials struct {
		File string `json:"file" yaml:"file"`
	} `json:"credentials" yaml:"credentials"`
}
