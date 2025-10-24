// cmd/gcp/gcp.go
package gcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Configuration structures
type Project struct {
	Name         string   `json:"name"`
	ID           string   `json:"id"`
	Environments []string `json:"environments"`
}

type Service struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	Projects []Project `json:"projects"`
	Services []Service `json:"services"`
}

type CacheData struct {
	Project string `json:"project"`
	Env     string `json:"env"`
	Service string `json:"service"`
}

var (
	// CLI flags
	projectFlag string
	envFlag     string
	serviceFlag string
	listFlag    bool
	repeatFlag  bool

	// Configuration
	config    Config
	configDir string
	cacheFile string
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorDim    = "\033[2m"
)

// bellSkipper implements an io.WriteCloser that skips the terminal bell character
type bellSkipper struct{}

// Write implements io.Writer by filtering out the bell character
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // Bell character
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements io.Closer
func (bs *bellSkipper) Close() error {
	return nil
}

// GcpCmd represents the gcp command
var GcpCmd = &cobra.Command{
	Use:   "gcp [project] [env] [service]",
	Short: "Interactive Google Cloud Console launcher",
	Long: `Opens GCP Console pages for different projects, environments, and services.
	
Examples:
  gcp                          # Interactive mode
  gcp air prod k8s             # Direct mode with partial matches
  gcp --repeat                 # Use last selection
  gcp --list                   # List available options`,
	RunE: runGcpCommand,
}

func init() {
	// Setup flags
	GcpCmd.Flags().StringVarP(&projectFlag, "project", "p", "", "Project name (partial match supported)")
	GcpCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment name")
	GcpCmd.Flags().StringVarP(&serviceFlag, "service", "s", "", "Service name (partial match supported)")
	GcpCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available projects and services")
	GcpCmd.Flags().BoolVarP(&repeatFlag, "repeat", "r", false, "Use last selection")

	// Initialize configuration paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError: Unable to determine home directory%s\n", colorRed, colorReset)
		os.Exit(1)
	}

	configDir = filepath.Join(homeDir, ".config", "sun-cli")
	cacheFile = filepath.Join(configDir, "gcp-cache.json")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: Unable to create config directory: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Load configuration
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "%sWarning: Unable to load config: %v%s\n", colorYellow, err, colorReset)
		fmt.Fprintf(os.Stderr, "%sUsing default configuration%s\n", colorYellow, colorReset)
		config = getDefaultConfig()
	}
}

// loadConfig loads configuration from JSON file
func loadConfig() error {
	configFile := filepath.Join(configDir, "gcp-config.json")

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config file
		defaultConfig := getDefaultConfig()
		return saveConfig(configFile, defaultConfig)
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// saveConfig saves configuration to JSON file
func saveConfig(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getDefaultConfig returns default configuration
func getDefaultConfig() Config {
	return Config{
		Projects: []Project{
			{
				Name:         "AirAsia MOVE",
				ID:           "airasia-move-project-id",
				Environments: []string{"prod", "staging", "dev"},
			},
			{
				Name:         "ARRK Engineering",
				ID:           "arrk-engineering-project-id",
				Environments: []string{"prod", "dev"},
			},
			{
				Name:         "Personal Sandbox",
				ID:           "my-sandbox-project-id",
				Environments: []string{"test"},
			},
		},
		Services: []Service{
			{Name: "Kubernetes Workloads", Path: "kubernetes/workload"},
			{Name: "Cloud SQL (MySQL)", Path: "sql/instances"},
			{Name: "Logs Explorer", Path: "logs/query"},
			{Name: "Monitoring Dashboards", Path: "monitoring/dashboards"},
			{Name: "Cloud Storage", Path: "storage/browser"},
			{Name: "Cloud Run", Path: "run"},
			{Name: "Cloud Functions", Path: "functions/list"},
			{Name: "IAM & Admin", Path: "iam-admin/iam"},
			{Name: "Compute Engine", Path: "compute/instances"},
			{Name: "BigQuery", Path: "bigquery"},
		},
	}
}

// runGcpCommand executes the main command logic
func runGcpCommand(cmd *cobra.Command, args []string) error {
	// Handle list flag
	if listFlag {
		return listOptions()
	}

	// Handle repeat flag
	if repeatFlag {
		return repeatLastSelection()
	}

	// Print welcome banner
	printBanner()

	// Parse arguments
	if len(args) > 0 {
		projectFlag = args[0]
	}
	if len(args) > 1 {
		envFlag = args[1]
	}
	if len(args) > 2 {
		serviceFlag = args[2]
	}

	// Selection loop with back navigation
	var project *Project
	var env string
	var service *Service
	var err error

	// Step 1: Select project
	for {
		project, err = selectProject(projectFlag)
		if err != nil {
			return err
		}
		break
	}

	// Step 2: Select environment (with back to project)
	for {
		env, err = selectEnvironment(project, envFlag)
		if err != nil {
			if err.Error() == "go_back" {
				// User wants to go back to project selection
				projectFlag = "" // Reset to show interactive prompt
				project, err = selectProject(projectFlag)
				if err != nil {
					return err
				}
				continue
			}
			return err
		}
		break
	}

	// Step 3: Select service (with back to environment)
	for {
		service, err = selectService(serviceFlag)
		if err != nil {
			if err.Error() == "go_back" {
				// User wants to go back to environment selection
				envFlag = "" // Reset to show interactive prompt
				continue
			}
			return err
		}
		break
	}

	// Build and open URL
	url := buildURL(project, env, service)
	printSummary(project, env, service, url)

	if err := openBrowser(url); err != nil {
		return err
	}

	// Cache selection
	saveCache(project.Name, env, service.Name)

	return nil
}

// selectProject handles project selection with improved partial matching
func selectProject(filter string) (*Project, error) {
	if filter != "" {
		// Find matching project (case-insensitive, partial match)
		matched := findMatchingProject(filter)
		if matched == nil {
			// Show available projects to help user
			fmt.Printf("%sNo project matching '%s'. Available projects:%s\n", colorYellow, filter, colorReset)
			for _, p := range config.Projects {
				fmt.Printf("  • %s\n", p.Name)
			}
			return nil, fmt.Errorf("no project matching '%s'", filter)
		}
		fmt.Printf("%s✓ Matched project:%s %s\n", colorGreen, colorReset, matched.Name)
		return matched, nil
	}

	// Interactive selection with fuzzy search
	fmt.Printf("\n%s%s📁 Select a GCP Project:%s\n", colorBold, colorBlue, colorReset)

	projectNames := make([]string, len(config.Projects))
	for i, p := range config.Projects {
		projectNames[i] = p.Name
	}
	sort.Strings(projectNames)

	searcher := func(input string, index int) bool {
		project := projectNames[index]
		name := strings.Replace(strings.ToLower(project), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		// Fuzzy search: check if all characters in input appear in order
		if input == "" {
			return true
		}

		inputIdx := 0
		for _, char := range name {
			if inputIdx < len(input) && char == rune(input[inputIdx]) {
				inputIdx++
			}
		}
		return inputIdx == len(input)
	}

	prompt := promptui.Select{
		Label:             "Project",
		Items:             projectNames,
		Size:              len(projectNames),
		HideHelp:          false,
		Stdout:            &bellSkipper{},
		CursorPos:         0,
		StartInSearchMode: true,
		Searcher:          searcher,
		Templates: &promptui.SelectTemplates{
			Active:   "▸ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "{{ \"✓\" | green }} {{ . | green }}",
			Help:     "{{ \"Type to search\" | faint }} {{ \"[↑↓ to move, enter to select, / to search, esc to cancel]\" | faint }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("project selection cancelled: %w", err)
	}

	return findProjectByName(result), nil
}

// selectEnvironment handles environment selection with validation
func selectEnvironment(project *Project, filter string) (string, error) {
	if filter != "" {
		// Validate environment (case-insensitive match)
		for _, env := range project.Environments {
			if strings.EqualFold(env, filter) {
				fmt.Printf("%s✓ Matched environment:%s %s\n", colorGreen, colorReset, env)
				return env, nil
			}
		}
		// Environment not found
		fmt.Printf("%sInvalid environment '%s' for project '%s'. Available:%s\n",
			colorYellow, filter, project.Name, colorReset)
		for _, env := range project.Environments {
			fmt.Printf("  • %s\n", env)
		}
		return "", fmt.Errorf("invalid environment '%s' for project '%s'. Valid: %v",
			filter, project.Name, project.Environments)
	}

	// Interactive selection with fuzzy search and back option
	fmt.Printf("\n%s%s🌎 Select an Environment:%s\n", colorBold, colorBlue, colorReset)

	// Add "← Go Back" option
	envOptions := make([]string, len(project.Environments)+1)
	envOptions[0] = "← Go Back"
	copy(envOptions[1:], project.Environments)

	searcher := func(input string, index int) bool {
		// Don't filter the back option
		if index == 0 {
			return strings.Contains(strings.ToLower("back"), strings.ToLower(input))
		}

		env := envOptions[index]
		name := strings.ToLower(env)
		input = strings.ToLower(input)

		// Simple fuzzy search for environments
		if input == "" {
			return true
		}

		inputIdx := 0
		for _, char := range name {
			if inputIdx < len(input) && char == rune(input[inputIdx]) {
				inputIdx++
			}
		}
		return inputIdx == len(input)
	}

	templates := &promptui.SelectTemplates{
		Active:   "▸ {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "{{ \"✓\" | green }} {{ . | green }}",
		Help:     "{{ \"Type to search\" | faint }} {{ \"[↑↓ to move, enter to select, / to search, esc to cancel]\" | faint }}",
	}

	// Customize template for back option
	templates.FuncMap = promptui.FuncMap
	templates.Active = `{{if eq . "← Go Back"}}▸ {{ . | yellow }}{{else}}▸ {{ . | cyan }}{{end}}`
	templates.Inactive = `{{if eq . "← Go Back"}}  {{ . | faint }}{{else}}  {{ . }}{{end}}`

	prompt := promptui.Select{
		Label:             "Environment",
		Items:             envOptions,
		Size:              len(envOptions),
		HideHelp:          false,
		Stdout:            &bellSkipper{},
		CursorPos:         1, // Start on first real environment, not back option
		StartInSearchMode: len(project.Environments) > 4,
		Searcher:          searcher,
		Templates:         templates,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("environment selection cancelled: %w", err)
	}

	// Check if user selected go back
	if result == "← Go Back" {
		return "", fmt.Errorf("go_back")
	}

	return result, nil
}

// selectService handles service selection with improved partial matching
func selectService(filter string) (*Service, error) {
	if filter != "" {
		// Find matching service (case-insensitive, partial match)
		matched := findMatchingService(filter)
		if matched == nil {
			// Show available services to help user
			fmt.Printf("%sNo service matching '%s'. Available services:%s\n", colorYellow, filter, colorReset)
			for _, s := range config.Services {
				fmt.Printf("  • %s\n", s.Name)
			}
			return nil, fmt.Errorf("no service matching '%s'", filter)
		}
		fmt.Printf("%s✓ Matched service:%s %s\n", colorGreen, colorReset, matched.Name)
		return matched, nil
	}

	// Interactive selection with fuzzy search and back option
	fmt.Printf("\n%s%s🧩 Select a Service:%s\n", colorBold, colorBlue, colorReset)

	serviceNames := make([]string, len(config.Services))
	for i, s := range config.Services {
		serviceNames[i] = s.Name
	}
	sort.Strings(serviceNames)

	// Add "← Go Back" option
	serviceOptions := make([]string, len(serviceNames)+1)
	serviceOptions[0] = "← Go Back"
	copy(serviceOptions[1:], serviceNames)

	searcher := func(input string, index int) bool {
		// Don't filter the back option
		if index == 0 {
			return strings.Contains(strings.ToLower("back"), strings.ToLower(input))
		}

		service := serviceOptions[index]
		name := strings.Replace(strings.ToLower(service), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		// Fuzzy search: check if all characters in input appear in order
		if input == "" {
			return true
		}

		// Also check for common abbreviations
		abbreviations := map[string][]string{
			"k8s":     {"kubernetes"},
			"gke":     {"kubernetes"},
			"sql":     {"sql", "database"},
			"db":      {"sql", "database"},
			"storage": {"storage", "gcs"},
			"gcs":     {"storage"},
			"logs":    {"logs", "explorer"},
			"iam":     {"iam", "admin"},
			"compute": {"compute", "engine", "vm"},
			"vm":      {"compute", "engine"},
			"bq":      {"bigquery"},
			"pubsub":  {"pub", "sub"},
			"run":     {"run"},
			"fn":      {"functions"},
			"lambda":  {"functions"},
		}

		// Check abbreviation match
		if matches, ok := abbreviations[input]; ok {
			for _, match := range matches {
				if strings.Contains(name, match) {
					return true
				}
			}
		}

		// Fuzzy character matching
		inputIdx := 0
		for _, char := range name {
			if inputIdx < len(input) && char == rune(input[inputIdx]) {
				inputIdx++
			}
		}
		return inputIdx == len(input)
	}

	templates := &promptui.SelectTemplates{
		Active:   "▸ {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "{{ \"✓\" | green }} {{ . | green }}",
		Help:     "{{ \"Type to search (e.g., k8s, sql, logs)\" | faint }} {{ \"[↑↓ to move, enter to select, / to toggle search]\" | faint }}",
	}

	// Customize template for back option
	templates.FuncMap = promptui.FuncMap
	templates.Active = `{{if eq . "← Go Back"}}▸ {{ . | yellow }}{{else}}▸ {{ . | cyan }}{{end}}`
	templates.Inactive = `{{if eq . "← Go Back"}}  {{ . | faint }}{{else}}  {{ . }}{{end}}`

	prompt := promptui.Select{
		Label:             "Service",
		Items:             serviceOptions,
		Size:              10,
		HideHelp:          false,
		Stdout:            &bellSkipper{},
		CursorPos:         1, // Start on first real service, not back option
		StartInSearchMode: true,
		Searcher:          searcher,
		Templates:         templates,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("service selection cancelled: %w", err)
	}

	// Check if user selected go back
	if result == "← Go Back" {
		return nil, fmt.Errorf("go_back")
	}

	return findServiceByName(result), nil
}

// buildURL constructs the GCP Console URL
func buildURL(project *Project, env string, service *Service) string {
	baseURL := "https://console.cloud.google.com"
	url := fmt.Sprintf("%s/%s?project=%s-%s", baseURL, service.Path, project.ID, env)

	// Add environment parameter for services that support it
	if needsEnvParam(service.Path) {
		url = fmt.Sprintf("%s&environment=%s", url, env)
	}

	return url
}

// needsEnvParam checks if service needs environment parameter
func needsEnvParam(servicePath string) bool {
	envServices := []string{"kubernetes/", "run", "functions/"}
	for _, s := range envServices {
		if strings.Contains(servicePath, s) {
			return true
		}
	}
	return false
}

// openBrowser opens URL in default browser (cross-platform)
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("%s⚠️  Cannot auto-open browser. Please visit:%s\n", colorYellow, colorReset)
		fmt.Printf("%s%s%s\n", colorBlue, url, colorReset)
		return nil
	}

	return nil
}

// findMatchingProject finds a project by partial, case-insensitive name match
func findMatchingProject(filter string) *Project {
	filter = strings.ToLower(strings.TrimSpace(filter))

	// First pass: exact match (case-insensitive)
	for i := range config.Projects {
		if strings.EqualFold(config.Projects[i].Name, filter) {
			return &config.Projects[i]
		}
	}

	// Second pass: starts with filter
	for i := range config.Projects {
		if strings.HasPrefix(strings.ToLower(config.Projects[i].Name), filter) {
			return &config.Projects[i]
		}
	}

	// Third pass: contains filter anywhere
	for i := range config.Projects {
		if strings.Contains(strings.ToLower(config.Projects[i].Name), filter) {
			return &config.Projects[i]
		}
	}

	// Fourth pass: check individual words (for acronyms like "air" matching "AirAsia")
	for i := range config.Projects {
		words := strings.Fields(strings.ToLower(config.Projects[i].Name))
		for _, word := range words {
			if strings.HasPrefix(word, filter) {
				return &config.Projects[i]
			}
		}
	}

	return nil
}

// findMatchingService finds a service by partial, case-insensitive name match
func findMatchingService(filter string) *Service {
	filter = strings.ToLower(strings.TrimSpace(filter))

	// First pass: exact match (case-insensitive)
	for i := range config.Services {
		if strings.EqualFold(config.Services[i].Name, filter) {
			return &config.Services[i]
		}
	}

	// Second pass: starts with filter
	for i := range config.Services {
		if strings.HasPrefix(strings.ToLower(config.Services[i].Name), filter) {
			return &config.Services[i]
		}
	}

	// Third pass: contains filter anywhere
	for i := range config.Services {
		if strings.Contains(strings.ToLower(config.Services[i].Name), filter) {
			return &config.Services[i]
		}
	}

	// Fourth pass: check individual words and common abbreviations
	abbreviations := map[string]string{
		"k8s":     "kubernetes",
		"gke":     "kubernetes",
		"sql":     "sql",
		"db":      "sql",
		"storage": "storage",
		"gcs":     "storage",
		"logs":    "logs",
		"iam":     "iam",
		"compute": "compute",
		"vm":      "compute",
		"bq":      "bigquery",
		"pubsub":  "pub/sub",
		"run":     "run",
		"fn":      "functions",
		"lambda":  "functions",
	}

	// Check if filter is a known abbreviation
	if expanded, ok := abbreviations[filter]; ok {
		for i := range config.Services {
			if strings.Contains(strings.ToLower(config.Services[i].Name), expanded) {
				return &config.Services[i]
			}
		}
	}

	// Check individual words
	for i := range config.Services {
		words := strings.Fields(strings.ToLower(config.Services[i].Name))
		for _, word := range words {
			if strings.HasPrefix(word, filter) {
				return &config.Services[i]
			}
		}
	}

	return nil
}

// findProjectByName finds a project by exact name
func findProjectByName(name string) *Project {
	for i := range config.Projects {
		if config.Projects[i].Name == name {
			return &config.Projects[i]
		}
	}
	return nil
}

// findServiceByName finds a service by exact name
func findServiceByName(name string) *Service {
	for i := range config.Services {
		if config.Services[i].Name == name {
			return &config.Services[i]
		}
	}
	return nil
}

// saveCache saves the selection to cache file
func saveCache(project, env, service string) {
	cache := CacheData{
		Project: project,
		Env:     env,
		Service: service,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(cacheFile, data, 0644)
}

// loadCache loads cached selection
func loadCache() (*CacheData, error) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// repeatLastSelection repeats the last selection
func repeatLastSelection() error {
	cache, err := loadCache()
	if err != nil {
		return fmt.Errorf("no cached selection found")
	}

	fmt.Printf("%s🔄 Using last selection...%s\n", colorYellow, colorReset)

	// Set flags from cache
	projectFlag = cache.Project
	envFlag = cache.Env
	serviceFlag = cache.Service

	// Execute the main logic directly instead of calling runGcpCommand
	printBanner()

	// Select project
	project, err := selectProject(projectFlag)
	if err != nil {
		return err
	}

	// Select environment
	env, err := selectEnvironment(project, envFlag)
	if err != nil {
		return err
	}

	// Select service
	service, err := selectService(serviceFlag)
	if err != nil {
		return err
	}

	// Build and open URL
	url := buildURL(project, env, service)
	printSummary(project, env, service, url)

	if err := openBrowser(url); err != nil {
		return err
	}

	// Cache selection (already cached but update timestamp)
	saveCache(project.Name, env, service.Name)

	return nil
}

// listOptions lists all available projects and services
func listOptions() error {
	fmt.Printf("\n%s%s📋 Available Configurations%s\n\n", colorBold, colorBlue, colorReset)

	// List projects
	fmt.Printf("%sProjects:%s\n", colorBold, colorReset)
	for _, p := range config.Projects {
		fmt.Printf("  • %s → %s\n", p.Name, p.ID)
		fmt.Printf("    %sEnvironments: %v%s\n", colorDim, p.Environments, colorReset)
	}

	// List services
	fmt.Printf("\n%sServices:%s\n", colorBold, colorReset)
	for _, s := range config.Services {
		fmt.Printf("  • %s → %s\n", s.Name, s.Path)
	}

	fmt.Printf("\n%sConfig location: %s%s\n", colorDim, filepath.Join(configDir, "gcp-config.json"), colorReset)
	fmt.Printf("\n%sTip: Use partial names like 'air' for AirAsia or 'k8s' for Kubernetes%s\n",
		colorDim, colorReset)

	return nil
}

// printBanner prints welcome banner
func printBanner() {
	fmt.Printf("\n%s%s╔════════════════════════════════════════╗%s\n", colorBold, colorBlue, colorReset)
	fmt.Printf("%s%s║  Google Cloud Console Launcher         ║%s\n", colorBold, colorBlue, colorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════╝%s\n", colorBold, colorBlue, colorReset)
}

// printSummary prints selection summary
func printSummary(project *Project, env string, service *Service, url string) {
	fmt.Printf("\n%s%s✓ Configuration%s\n", colorGreen, colorBold, colorReset)
	fmt.Printf("%s  Project:     %s%s %s(%s)%s\n", colorDim, colorReset, project.Name, colorDim, project.ID, colorReset)
	fmt.Printf("%s  Environment: %s%s\n", colorDim, colorReset, env)
	fmt.Printf("%s  Service:     %s%s\n", colorDim, colorReset, service.Name)
	fmt.Printf("\n%s🚀 Opening: %s%s\n\n", colorBlue, url, colorReset)
}
