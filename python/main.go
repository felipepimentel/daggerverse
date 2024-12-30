package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"dagger/python/internal/dagger"

	"github.com/spf13/viper"
)

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value string
}

// PyPIConfig holds PyPI deployment configuration
type PyPIConfig struct {
	// Registry URL (default: https://upload.pypi.org/legacy/)
	Registry string
	// Token for authentication
	Token *dagger.Secret
	// Skip existing versions (default: false)
	SkipExisting bool
	// Allow dirty versions (default: false)
	AllowDirty bool
	// Additional publish arguments
	ExtraArgs []string
	// Environment variables for publishing
	Env []KeyValue
	// Repository name in Poetry config (default: "pypi")
	RepositoryName string
	// Skip build before publishing (default: false)
	SkipBuild bool
	// Skip verification before publishing (default: false)
	SkipVerify bool
}

// TestConfig holds pytest configuration options
type TestConfig struct {
	// Verbose output (default: true)
	Verbose bool
	// Number of parallel workers (default: auto)
	Workers int
	// Coverage configuration
	Coverage *CoverageConfig
	// Additional pytest arguments
	ExtraArgs []string
	// Environment variables for tests
	Env []KeyValue
	// Test markers to select
	Markers []string
	// Test paths to run (default: ".")
	TestPaths []string
	// Skip installing test dependencies (default: false)
	SkipInstall bool
	// JUnit XML report path
	JUnitXML string
	// Maximum test duration in seconds (0 for no limit)
	MaxTestTime int
	// Stop on first failure (default: false)
	FailFast bool
}

// CoverageConfig holds coverage reporting configuration
type CoverageConfig struct {
	// Enable coverage reporting (default: true)
	Enabled bool
	// Coverage report formats (default: ["term", "xml"])
	Formats []string
	// Minimum coverage percentage (default: 0)
	MinCoverage int
	// Coverage output directory (default: "coverage")
	OutputDir string
	// Paths to include in coverage
	Include []string
	// Paths to exclude from coverage
	Exclude []string
	// Show missing lines in report (default: true)
	ShowMissing bool
	// Branch coverage (default: false)
	Branch bool
	// Context lines in report (default: 0)
	Context int
}

// BuildConfig holds Poetry build configuration
type BuildConfig struct {
	// Additional build arguments
	BuildArgs []string
	// Additional dependencies to install
	ExtraDependencies []string
	// Poetry configuration options
	PoetryConfig []KeyValue
	// Environment variables for build
	Env []KeyValue
	// Cache configuration
	Cache *CacheConfig
	// Poetry dependency groups to install (default: ["dev"])
	DependencyGroups []string
	// Poetry optional dependencies to install
	OptionalDependencies []string
	// Skip installing dependencies (default: false)
	SkipDependencies bool
	// Install only selected dependency groups (default: false)
	OnlyGroups bool
	// Skip installing the root package (default: true)
	SkipRoot bool
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	// Enable pip cache (default: true)
	PipCache bool
	// Enable poetry cache (default: true)
	PoetryCache bool
	// Custom cache volume names
	PipCacheVolume string
	PoetryCacheVolume string
}

// LintConfig holds configuration for code linting
type LintConfig struct {
	// Enable ruff linting (default: true)
	Enabled bool
	// Linting rules to enable
	Select []string
	// Linting rules to ignore
	Ignore []string
	// Files or directories to exclude
	Exclude []string
	// Maximum line length (default: 88)
	LineLength int
	// Fix issues automatically (default: false)
	Fix bool
	// Show source code for each error (default: true)
	ShowSource bool
	// Format output (options: text, json, junit)
	Format string
	// Output file path
	Output string
	// Additional ruff arguments
	ExtraArgs []string
}

// FormatConfig holds configuration for code formatting
type FormatConfig struct {
	// Enable black formatting (default: true)
	Enabled bool
	// Line length (default: 88)
	LineLength int
	// Skip string normalization (default: false)
	SkipStringNormalization bool
	// Target specific Python versions
	TargetVersion []string
	// Include files/patterns
	Include []string
	// Exclude files/patterns
	Exclude []string
	// Check only (no changes) (default: false)
	Check bool
	// Show diff of changes (default: true)
	ShowDiff bool
	// Additional black arguments
	ExtraArgs []string
}

// DocsConfig holds configuration for documentation generation
type DocsConfig struct {
	// Enable documentation generation (default: true)
	Enabled bool
	// Documentation tool to use (options: sphinx, mkdocs)
	Tool string
	// Documentation source directory (default: "docs")
	SourceDir string
	// Documentation output directory (default: "site")
	OutputDir string
	// Documentation format (options: html, pdf, epub)
	Format string
	// Documentation theme
	Theme string
	// Project name
	ProjectName string
	// Project version
	Version string
	// Project author
	Author string
	// Additional extensions to enable
	Extensions []string
	// Additional dependencies to install
	ExtraDependencies []string
	// Environment variables for documentation build
	Env []KeyValue
	// Additional build arguments
	ExtraArgs []string
}

// GitConfig holds configuration for Git operations
type GitConfig struct {
	// Repository URL to clone
	Repository string
	// Branch or tag to checkout (default: main)
	Ref string
	// Depth of git history to clone (default: 1)
	Depth int
	// Whether to fetch all history (default: false)
	FetchAll bool
	// Whether to fetch all tags (default: false)
	FetchTags bool
	// Whether to fetch submodules (default: false)
	Submodules bool
	// Authentication token for private repositories
	Token *dagger.Secret
	// SSH key for authentication
	SSHKey *dagger.Secret
	// Known hosts file content for SSH authentication
	KnownHosts string
	// Additional git configuration
	Config []KeyValue
	// Environment variables for git operations
	Env []KeyValue
}

// Python represents a Python module with Poetry support
type Python struct {
	// PythonVersion specifies the Python version to use (default: "3.12")
	PythonVersion string
	// PackagePath specifies the path to the package within the source (default: ".")
	PackagePath string
	// PyPIConfig holds the PyPI deployment configuration
	PyPIConfig *PyPIConfig
	// TestConfig holds the test configuration
	TestConfig *TestConfig
	// BuildConfig holds the build configuration
	BuildConfig *BuildConfig
	// LintConfig holds the linting configuration
	LintConfig *LintConfig
	// FormatConfig holds the formatting configuration
	FormatConfig *FormatConfig
	// DocsConfig holds the documentation configuration
	DocsConfig *DocsConfig
	// GitConfig holds the Git configuration
	GitConfig *GitConfig
}

// validatePythonVersion validates Python version
func (m *Python) validatePythonVersion(version string) error {
	if version == "" {
		return nil
	}
	// Basic version format validation
	if !strings.HasPrefix(version, "3.") {
		return fmt.Errorf("unsupported Python version: %s (only Python 3.x is supported)", version)
	}
	return nil
}

// WithPythonVersion sets the Python version to use
func (m *Python) WithPythonVersion(version string) *Python {
	if err := m.validatePythonVersion(version); err != nil {
		panic(err) // Panic is appropriate here as this is a builder pattern
	}
	m.PythonVersion = version
	return m
}

// WithPackagePath sets the package path within the source
func (m *Python) WithPackagePath(path string) *Python {
	m.PackagePath = path
	return m
}

// WithPyPIConfig sets the PyPI deployment configuration
func (m *Python) WithPyPIConfig(config *PyPIConfig) *Python {
	m.PyPIConfig = config
	return m
}

// WithTestConfig sets the test configuration
func (m *Python) WithTestConfig(config *TestConfig) *Python {
	m.TestConfig = config
	return m
}

// WithBuildConfig sets the build configuration
func (m *Python) WithBuildConfig(config *BuildConfig) *Python {
	m.BuildConfig = config
	return m
}

// WithLintConfig sets the linting configuration
func (m *Python) WithLintConfig(config *LintConfig) *Python {
	m.LintConfig = config
	return m
}

// WithFormatConfig sets the formatting configuration
func (m *Python) WithFormatConfig(config *FormatConfig) *Python {
	m.FormatConfig = config
	return m
}

// WithDocsConfig sets the documentation configuration
func (m *Python) WithDocsConfig(config *DocsConfig) *Python {
	m.DocsConfig = config
	return m
}

// WithGitConfig sets the Git configuration
func (m *Python) WithGitConfig(config *GitConfig) *Python {
	m.GitConfig = config
	return m
}

// getBaseImage returns the Python base image with the configured version
func (m *Python) getBaseImage() string {
	version := m.PythonVersion
	if version == "" {
		version = "3.12"
	}
	if err := m.validatePythonVersion(version); err != nil {
		panic(err)
	}
	return fmt.Sprintf("python:%s-slim", version)
}

// getWorkdir returns the working directory path
func (m *Python) getWorkdir(basePath string) string {
	if m.PackagePath == "" {
		return basePath
	}
	return filepath.Join(basePath, m.PackagePath)
}

// getDefaultPyPIConfig returns default PyPI configuration
func (m *Python) getDefaultPyPIConfig() *PyPIConfig {
	return &PyPIConfig{
		Registry:       "https://upload.pypi.org/legacy/",
		SkipExisting:  false,
		AllowDirty:    false,
		ExtraArgs:     []string{},
		Env:           []KeyValue{},
		RepositoryName: "pypi",
		SkipBuild:     false,
		SkipVerify:    false,
	}
}

// validatePyPIConfig validates PyPI configuration
func (m *Python) validatePyPIConfig(config *PyPIConfig) error {
	if config == nil {
		return nil
	}
	if config.Registry != "" {
		if _, err := url.Parse(config.Registry); err != nil {
			return fmt.Errorf("invalid registry URL: %w", err)
		}
	}
	return nil
}

// getEnvOrSecret tries to get a value from multiple sources in order:
// 1. Command line argument
// 2. Environment variable
// 3. .env file
// 4. Default value
func (m *Python) getEnvOrSecret(key string, defaultValue string) (*dagger.Secret, error) {
	// Initialize viper
	v := viper.New()
	
	// Set up viper to read from environment
	v.SetEnvPrefix("")
	v.AutomaticEnv()
	
	// Try to read from .env file if it exists
	v.SetConfigFile(".env")
	_ = v.ReadInConfig() // Ignore error if file doesn't exist
	
	// Get value from any source
	value := v.GetString(key)
	if value == "" {
		value = os.Getenv(key)
	}
	if value == "" {
		value = defaultValue
	}
	
	if value == "" {
		return nil, fmt.Errorf("%s is required", key)
	}
	
	return dag.SetSecret(key, value), nil
}

// Publish builds, tests and publishes the Python package to a registry
func (m *Python) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	config := m.PyPIConfig
	if config == nil {
		config = m.getDefaultPyPIConfig()
	}

	if err := m.validatePyPIConfig(config); err != nil {
		return "", err
	}

	// Find pyproject.toml location
	projectPath, err := m.findPyProjectToml(source)
	if err != nil {
		// If not found, use default package path
		projectPath = m.PackagePath
	}

	// Run tests before publishing unless verification is skipped
	if !config.SkipVerify {
		if _, err := m.Test(ctx, source); err != nil {
			return "", fmt.Errorf("tests failed: %w", err)
		}
	}

	// Build the package unless skipped
	var container *dagger.Container
	if !config.SkipBuild {
		container = m.Build(source)
	} else {
		container = m.BuildEnv(source)
	}

	// Add environment variables
	for _, kv := range config.Env {
		container = container.WithEnvVariable(kv.Key, kv.Value)
	}

	// Configure Poetry authentication
	if token != nil {
		container = container.
			WithSecretVariable("POETRY_HTTP_BASIC_PYPI_USERNAME", dag.SetSecret("PYPI_USERNAME", "__token__")).
			WithSecretVariable("POETRY_HTTP_BASIC_PYPI_PASSWORD", token).
			WithSecretVariable("POETRY_PYPI_TOKEN", token)
	}

	// Build publish command
	args := []string{"poetry", "publish"}

	// Add skip existing flag
	if config.SkipExisting {
		args = append(args, "--skip-existing")
	}

	// Add allow dirty flag
	if config.AllowDirty {
		args = append(args, "--allow-dirty")
	}

	// Add extra arguments
	args = append(args, config.ExtraArgs...)

	// Run publish command
	return container.WithWorkdir(filepath.Join("/app", projectPath)).WithExec(args).Stdout(ctx)
}

// Build creates a Python package using Poetry
func (m *Python) Build(source *dagger.Directory) *dagger.Container {
	container := m.BuildEnv(source).
		WithExec([]string{
			"poetry", "build",
			"--no-interaction",
		})
	
	return container.WithDirectory("/dist", container.Directory("/app/dist"))
}

// getDefaultTestConfig returns default test configuration
func (m *Python) getDefaultTestConfig() *TestConfig {
	return &TestConfig{
		Verbose: true,
		Workers: 0, // auto
		Coverage: &CoverageConfig{
			Enabled: true,
			Formats: []string{"term", "xml"},
			MinCoverage: 0,
			OutputDir: "coverage",
			Include: []string{},
			Exclude: []string{},
			ShowMissing: true,
			Branch: false,
			Context: 0,
		},
		Env: []KeyValue{},
		Markers: []string{},
		TestPaths: []string{"."},
		SkipInstall: false,
		JUnitXML: "",
		MaxTestTime: 0,
		FailFast: false,
	}
}

// Test runs pytest with the specified configuration
func (m *Python) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	config := m.TestConfig
	if config == nil {
		config = m.getDefaultTestConfig()
	}

	// Find pyproject.toml location
	projectPath, err := m.findPyProjectToml(source)
	if err != nil {
		// If not found, use default package path
		projectPath = m.PackagePath
	}

	// Build environment unless skipped
	var container *dagger.Container
	if !config.SkipInstall {
		container = m.BuildEnv(source)
	} else {
		container = dag.Container().
			From(m.getBaseImage()).
			WithDirectory("/app", source).
			WithWorkdir(filepath.Join("/app", projectPath))
	}

	// Add environment variables
	for _, kv := range config.Env {
		container = container.WithEnvVariable(kv.Key, kv.Value)
	}

	args := []string{"poetry", "run", "pytest"}

	if config.Verbose {
		args = append(args, "--verbose", "--color=yes")
	}

	if config.Workers > 0 {
		args = append(args, fmt.Sprintf("-n=%d", config.Workers))
	}

	if config.Coverage != nil && config.Coverage.Enabled {
		args = append(args, fmt.Sprintf("--cov=%s", projectPath))
		
		for _, format := range config.Coverage.Formats {
			switch format {
			case "xml":
				args = append(args, "--cov-report=xml")
			case "html":
				args = append(args, fmt.Sprintf("--cov-report=html:%s/html", config.Coverage.OutputDir))
			case "term":
				args = append(args, "--cov-report=term")
			}
		}

		if config.Coverage.MinCoverage > 0 {
			args = append(args, fmt.Sprintf("--cov-fail-under=%d", config.Coverage.MinCoverage))
		}

		if len(config.Coverage.Include) > 0 {
			args = append(args, "--cov-include=" + strings.Join(config.Coverage.Include, ","))
		}

		if len(config.Coverage.Exclude) > 0 {
			args = append(args, "--cov-exclude=" + strings.Join(config.Coverage.Exclude, ","))
		}

		if config.Coverage.ShowMissing {
			args = append(args, "--cov-report=term-missing")
		}

		if config.Coverage.Branch {
			args = append(args, "--cov-branch")
		}

		if config.Coverage.Context > 0 {
			args = append(args, fmt.Sprintf("--cov-context=%d", config.Coverage.Context))
		}

		args = append(args, "--no-cov-on-fail")
	}

	// Add test markers
	for _, marker := range config.Markers {
		args = append(args, "-m", marker)
	}

	// Add JUnit XML report
	if config.JUnitXML != "" {
		args = append(args, fmt.Sprintf("--junitxml=%s", config.JUnitXML))
	}

	// Add max test time
	if config.MaxTestTime > 0 {
		args = append(args, fmt.Sprintf("--maxtime=%d", config.MaxTestTime))
	}

	// Add fail fast
	if config.FailFast {
		args = append(args, "--exitfirst")
	}

	// Add test paths
	if len(config.TestPaths) > 0 {
		args = append(args, config.TestPaths...)
	} else {
		args = append(args, projectPath)
	}

	// Add any extra arguments
	args = append(args, config.ExtraArgs...)

	return container.WithExec(args).Stdout(ctx)
}

// getDefaultBuildConfig returns default build configuration
func (m *Python) getDefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		BuildArgs: []string{},
		ExtraDependencies: []string{},
		PoetryConfig: []KeyValue{},
		Env: []KeyValue{},
		Cache: &CacheConfig{
			PipCache: true,
			PoetryCache: true,
			PipCacheVolume: "pip-cache",
			PoetryCacheVolume: "poetry-cache",
		},
		DependencyGroups: []string{"dev"},
		OptionalDependencies: []string{},
		SkipDependencies: false,
		OnlyGroups: false,
		SkipRoot: true,
	}
}

// BuildEnv creates a container with Python and Poetry installed
func (m *Python) BuildEnv(source *dagger.Directory) *dagger.Container {
	config := m.BuildConfig
	if config == nil {
		config = m.getDefaultBuildConfig()
	}

	// Find pyproject.toml location
	projectPath, err := m.findPyProjectToml(source)
	if err != nil {
		// If not found, use default package path
		projectPath = m.PackagePath
	}

	// Initialize container
	container := dag.Container().From(m.getBaseImage())

	// Add environment variables
	for _, kv := range config.Env {
		container = container.WithEnvVariable(kv.Key, kv.Value)
	}

	// Setup caches
	if config.Cache != nil {
		if config.Cache.PipCache {
			pipCache := dag.CacheVolume(config.Cache.PipCacheVolume)
			container = container.WithMountedCache("/root/.cache/pip", pipCache)
		}
		if config.Cache.PoetryCache {
			poetryCache := dag.CacheVolume(config.Cache.PoetryCacheVolume)
			container = container.WithMountedCache("/root/.cache/pypoetry", poetryCache)
		}
	}

	// Add source code
	container = container.
		WithDirectory("/app", source).
		WithWorkdir(filepath.Join("/app", projectPath))

	// Install base dependencies
	container = container.WithExec([]string{
		"pip", "install",
		"--no-cache-dir",
		"--upgrade",
		"pip",
		"poetry",
	})

	// Install extra dependencies if any
	if len(config.ExtraDependencies) > 0 {
		container = container.WithExec(append(
			[]string{"pip", "install", "--no-cache-dir"},
			config.ExtraDependencies...,
		))
	}

	// Configure Poetry
	for _, kv := range config.PoetryConfig {
		container = container.WithExec([]string{
			"poetry", "config",
			kv.Key, kv.Value,
		})
	}

	// Skip dependency installation if requested
	if config.SkipDependencies {
		return container
	}

	// Install project dependencies
	installArgs := []string{
		"poetry", "install",
		"--no-interaction",
	}

	// Handle root package installation
	if config.SkipRoot {
		installArgs = append(installArgs, "--no-root")
	}

	// Handle dependency groups
	if len(config.DependencyGroups) > 0 {
		if config.OnlyGroups {
			installArgs = append(installArgs, "--only")
		} else {
			installArgs = append(installArgs, "--with")
		}
		installArgs = append(installArgs, strings.Join(config.DependencyGroups, ","))
	}

	// Handle optional dependencies
	if len(config.OptionalDependencies) > 0 {
		installArgs = append(installArgs, "--extras")
		installArgs = append(installArgs, strings.Join(config.OptionalDependencies, ","))
	}

	// Add build arguments
	installArgs = append(installArgs, config.BuildArgs...)
	
	return container.WithExec(installArgs)
}

// getDefaultLintConfig returns default linting configuration
func (m *Python) getDefaultLintConfig() *LintConfig {
	return &LintConfig{
		Enabled: true,
		Select: []string{
			"E", "F", "W",  // pycodestyle, pyflakes, warnings
			"I",            // isort
			"N",            // pep8-naming
			"UP",           // pyupgrade
			"RUF",          // ruff-specific
		},
		Ignore: []string{},
		Exclude: []string{},
		LineLength: 88,
		Fix: false,
		ShowSource: true,
		Format: "text",
		Output: "",
		ExtraArgs: []string{},
	}
}

// getDefaultFormatConfig returns default formatting configuration
func (m *Python) getDefaultFormatConfig() *FormatConfig {
	return &FormatConfig{
		Enabled: true,
		LineLength: 88,
		SkipStringNormalization: false,
		TargetVersion: []string{},
		Include: []string{},
		Exclude: []string{},
		Check: false,
		ShowDiff: true,
		ExtraArgs: []string{},
	}
}

// Lint runs code linting using ruff
func (m *Python) Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	config := m.LintConfig
	if config == nil {
		config = m.getDefaultLintConfig()
	}

	if !config.Enabled {
		return "Linting skipped", nil
	}

	// Find pyproject.toml location
	projectPath, err := m.findPyProjectToml(source)
	if err != nil {
		// If not found, use default package path
		projectPath = m.PackagePath
	}

	container := m.BuildEnv(source)

	// Install ruff if not in extra dependencies
	container = container.WithExec([]string{
		"pip", "install", "--no-cache-dir", "ruff",
	})

	args := []string{"ruff", "check"}

	// Add selected rules
	if len(config.Select) > 0 {
		args = append(args, "--select", strings.Join(config.Select, ","))
	}

	// Add ignored rules
	if len(config.Ignore) > 0 {
		args = append(args, "--ignore", strings.Join(config.Ignore, ","))
	}

	// Add excluded patterns
	if len(config.Exclude) > 0 {
		args = append(args, "--exclude", strings.Join(config.Exclude, ","))
	}

	// Add line length
	if config.LineLength > 0 {
		args = append(args, "--line-length", fmt.Sprintf("%d", config.LineLength))
	}

	// Add fix flag
	if config.Fix {
		args = append(args, "--fix")
	}

	// Add show source flag
	if config.ShowSource {
		args = append(args, "--show-source")
	}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add extra arguments
	args = append(args, config.ExtraArgs...)

	// Add target path
	args = append(args, projectPath)

	return container.WithExec(args).Stdout(ctx)
}

// Format runs code formatting using black
func (m *Python) Format(ctx context.Context, source *dagger.Directory) (string, error) {
	config := m.FormatConfig
	if config == nil {
		config = m.getDefaultFormatConfig()
	}

	if !config.Enabled {
		return "Formatting skipped", nil
	}

	// Find pyproject.toml location
	projectPath, err := m.findPyProjectToml(source)
	if err != nil {
		// If not found, use default package path
		projectPath = m.PackagePath
	}

	container := m.BuildEnv(source)

	// Install black if not in extra dependencies
	container = container.WithExec([]string{
		"pip", "install", "--no-cache-dir", "black",
	})

	args := []string{"black"}

	// Add line length
	if config.LineLength > 0 {
		args = append(args, "--line-length", fmt.Sprintf("%d", config.LineLength))
	}

	// Add string normalization flag
	if config.SkipStringNormalization {
		args = append(args, "--skip-string-normalization")
	}

	// Add target versions
	if len(config.TargetVersion) > 0 {
		args = append(args, "--target-version", strings.Join(config.TargetVersion, ","))
	}

	// Add included patterns
	if len(config.Include) > 0 {
		args = append(args, "--include", strings.Join(config.Include, "|"))
	}

	// Add excluded patterns
	if len(config.Exclude) > 0 {
		args = append(args, "--exclude", strings.Join(config.Exclude, "|"))
	}

	// Add check flag
	if config.Check {
		args = append(args, "--check")
	}

	// Add diff flag
	if config.ShowDiff {
		args = append(args, "--diff")
	}

	// Add extra arguments
	args = append(args, config.ExtraArgs...)

	// Add target path
	args = append(args, projectPath)

	return container.WithExec(args).Stdout(ctx)
}

// getDefaultDocsConfig returns default documentation configuration
func (m *Python) getDefaultDocsConfig() *DocsConfig {
	return &DocsConfig{
		Enabled: true,
		Tool: "sphinx",
		SourceDir: "docs",
		OutputDir: "site",
		Format: "html",
		Theme: "sphinx_rtd_theme",
		ProjectName: "",
		Version: "",
		Author: "",
		Extensions: []string{
			"sphinx.ext.autodoc",
			"sphinx.ext.napoleon",
			"sphinx.ext.viewcode",
			"sphinx.ext.intersphinx",
		},
		ExtraDependencies: []string{
			"sphinx",
			"sphinx-rtd-theme",
			"myst-parser",
		},
		Env: []KeyValue{},
		ExtraArgs: []string{},
	}
}

// BuildDocs generates documentation using Sphinx or MkDocs
func (m *Python) BuildDocs(ctx context.Context, source *dagger.Directory) (string, error) {
	config := m.DocsConfig
	if config == nil {
		config = m.getDefaultDocsConfig()
	}

	if !config.Enabled {
		return "Documentation generation skipped", nil
	}

	// Find pyproject.toml location
	projectPath, err := m.findPyProjectToml(source)
	if err != nil {
		// If not found, use default package path
		projectPath = m.PackagePath
	}

	container := m.BuildEnv(source)

	// Install documentation dependencies
	container = container.WithExec(append(
		[]string{"pip", "install", "--no-cache-dir"},
		config.ExtraDependencies...,
	))

	// Add environment variables
	for _, kv := range config.Env {
		container = container.WithEnvVariable(kv.Key, kv.Value)
	}

	var args []string
	switch config.Tool {
	case "mkdocs":
		args = []string{"mkdocs", "build"}
		if config.OutputDir != "" {
			args = append(args, "--site-dir", config.OutputDir)
		}
	default: // sphinx
		args = []string{
			"sphinx-build",
			"-b", config.Format,
			filepath.Join(projectPath, config.SourceDir),
			filepath.Join(projectPath, config.OutputDir, config.Format),
		}
		if config.ProjectName != "" {
			args = append(args, "-D", fmt.Sprintf("project=%s", config.ProjectName))
		}
		if config.Version != "" {
			args = append(args, "-D", fmt.Sprintf("version=%s", config.Version))
		}
		if config.Author != "" {
			args = append(args, "-D", fmt.Sprintf("author=%s", config.Author))
		}
		if config.Theme != "" {
			args = append(args, "-D", fmt.Sprintf("html_theme=%s", config.Theme))
		}
	}

	// Add extra arguments
	args = append(args, config.ExtraArgs...)

	return container.WithExec(args).Stdout(ctx)
}

// getDefaultGitConfig returns default Git configuration
func (m *Python) getDefaultGitConfig() *GitConfig {
	return &GitConfig{
		Ref: "main",
		Depth: 1,
		FetchAll: false,
		FetchTags: false,
		Submodules: false,
		Config: []KeyValue{},
		Env: []KeyValue{},
	}
}

// validateGitConfig validates Git configuration
func (m *Python) validateGitConfig(config *GitConfig) error {
	if config == nil {
		return nil
	}
	if config.Repository == "" {
		return fmt.Errorf("repository URL is required")
	}
	if _, err := url.Parse(config.Repository); err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}
	if config.Depth < 0 {
		return fmt.Errorf("depth must be non-negative")
	}
	return nil
}

// Checkout clones a Git repository and returns its directory
func (m *Python) Checkout(ctx context.Context) (*dagger.Directory, error) {
	config := m.GitConfig
	if config == nil {
		config = m.getDefaultGitConfig()
	}

	if err := m.validateGitConfig(config); err != nil {
		return nil, err
	}

	// Start with git container
	container := dag.Container().
		From("alpine/git:latest")

	// Add environment variables
	for _, kv := range config.Env {
		container = container.WithEnvVariable(kv.Key, kv.Value)
	}

	// Add git configuration
	for _, kv := range config.Config {
		container = container.WithExec([]string{
			"git", "config", "--global",
			kv.Key, kv.Value,
		})
	}

	// Setup authentication if needed
	if config.Token != nil {
		repoURL, err := url.Parse(config.Repository)
		if err != nil {
			return nil, fmt.Errorf("invalid repository URL: %w", err)
		}

		// Add token to URL
		token, err := config.Token.Plaintext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
		repoURL.User = url.UserPassword("git", token)
		config.Repository = repoURL.String()
	} else if config.SSHKey != nil {
		// Create SSH directory
		container = container.WithExec([]string{
			"mkdir", "-p", "/root/.ssh",
		})

		// Write SSH key
		sshKey, err := config.SSHKey.Plaintext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get SSH key: %w", err)
		}
		container = container.WithNewFile("/root/.ssh/id_rsa", sshKey, dagger.ContainerWithNewFileOpts{
			Permissions: 0600,
		})

		if config.KnownHosts != "" {
			container = container.WithNewFile("/root/.ssh/known_hosts", config.KnownHosts, dagger.ContainerWithNewFileOpts{
				Permissions: 0600,
			})
		}
	}

	// Prepare clone command
	cloneArgs := []string{
		"clone",
		"--branch", config.Ref,
	}

	if !config.FetchAll {
		if config.Depth > 0 {
			cloneArgs = append(cloneArgs, "--depth", fmt.Sprintf("%d", config.Depth))
		}
		if !config.FetchTags {
			cloneArgs = append(cloneArgs, "--no-tags")
		}
	}

	if !config.Submodules {
		cloneArgs = append(cloneArgs, "--no-recurse-submodules")
	}

	cloneArgs = append(cloneArgs, config.Repository, ".")

	// Clone repository
	container = container.
		WithWorkdir("/src").
		WithExec(cloneArgs)

	// Return the cloned directory
	return container.Directory("."), nil
}

// CI runs the Continuous Integration pipeline (test and build)
func (m *Python) CI(ctx context.Context, source *dagger.Directory) (string, error) {
	// Run tests first
	if _, err := m.Test(ctx, source); err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Then build
	container := m.Build(source)
	if container == nil {
		return "", fmt.Errorf("build failed: container is nil")
	}

	return "CI pipeline completed successfully", nil
}

// CD runs the Continuous Delivery pipeline (publish to PyPI)
func (m *Python) CD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	// Get token from provided token or environment
	publishToken := token
	if publishToken == nil {
		var err error
		publishToken, err = m.getEnvOrSecret("PYPI_TOKEN", "")
		if err != nil {
			return "", err
		}
	}

	// Publish to PyPI
	return m.Publish(ctx, source, publishToken)
}

// CICD runs the complete CI/CD pipeline including version management
func (m *Python) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	// First, handle versioning
	version, err := m.bumpVersion(ctx, source)
	if err != nil {
		return "", fmt.Errorf("error bumping version: %v", err)
	}

	// Run tests
	testOutput, err := m.Test(ctx, source)
	if err != nil {
		return "", fmt.Errorf("error running tests: %v", err)
	}
	fmt.Println("Tests output:", testOutput)

	// Run linting
	lintOutput, err := m.Lint(ctx, source)
	if err != nil {
		return "", fmt.Errorf("error running linter: %v", err)
	}
	fmt.Println("Lint output:", lintOutput)

	// Run formatting
	formatOutput, err := m.Format(ctx, source)
	if err != nil {
		return "", fmt.Errorf("error running formatter: %v", err)
	}
	fmt.Println("Format output:", formatOutput)

	// Build package
	buildOutput := m.Build(source)
	fmt.Println("Build output:", buildOutput)

	// Update version in pyproject.toml
	container := m.BuildEnv(source)
	container = container.WithExec([]string{
		"poetry", "version", version,
	})

	// Publish to PyPI
	publishOutput, err := m.Publish(ctx, source, token)
	if err != nil {
		return "", fmt.Errorf("error publishing to PyPI: %v", err)
	}
	fmt.Println("Publish output:", publishOutput)

	return version, nil
}

func (m *Python) bumpVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	// Setup container with Node.js and required tools
	container := dag.Container().
		From("node:lts-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("GIT_AUTHOR_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "github-actions[bot]@users.noreply.github.com")

	// Install required packages
	container = container.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "openssh-client"})

	// Install semantic-release and all required plugins
	container = container.WithExec([]string{
		"npm", "install", "-g",
		"semantic-release",
		"@semantic-release/commit-analyzer",
		"@semantic-release/release-notes-generator",
		"@semantic-release/changelog",
		"@semantic-release/git",
		"@semantic-release/github",
	})

	// Configure Git user
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Run semantic-release with explicit configuration
	output, err := container.
		WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithEnvVariable("GH_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithExec([]string{
			"npx", "semantic-release",
			"--branches", "main",
			"--ci", "false", // Disable CI detection since we're running in Dagger
			"--debug", // Enable debug logging
		}).Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("error running semantic-release: %v", err)
	}

	// Extract version from output
	version := strings.TrimSpace(output)
	if version == "" {
		return "", fmt.Errorf("no version found in semantic-release output")
	}

	return version, nil
}

// findPyProjectToml recursively searches for pyproject.toml file
func (m *Python) findPyProjectToml(source *dagger.Directory) (string, error) {
	// First try the package path
	if m.PackagePath != "" {
		if _, err := source.File(filepath.Join(m.PackagePath, "pyproject.toml")).Contents(context.Background()); err == nil {
			return m.PackagePath, nil
		}
	}

	// Then try root directory
	if _, err := source.File("pyproject.toml").Contents(context.Background()); err == nil {
		return ".", nil
	}

	// Finally, search recursively
	entries, err := source.Entries(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to list directory entries: %w", err)
	}

	for _, entry := range entries {
		if entry == "pyproject.toml" {
			return ".", nil
		}
		
		// Check if it's a directory
		if _, err := source.Directory(entry).ID(context.Background()); err == nil {
			// Recursively search in subdirectory
			subdir := source.Directory(entry)
			if path, err := m.findPyProjectToml(subdir); err == nil {
				return filepath.Join(entry, path), nil
			}
		}
	}

	return "", fmt.Errorf("pyproject.toml not found in the source directory")
}
