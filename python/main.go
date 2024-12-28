package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger/python/internal/dagger"
)

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
	Env map[string]string
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
	Env map[string]string
}

// CoverageConfig holds coverage reporting configuration
type CoverageConfig struct {
	// Enable coverage reporting (default: true)
	Enabled bool
	// Coverage report formats (default: ["term", "xml"])
	Formats []string
	// Minimum coverage percentage (default: 0)
	MinCoverage float64
	// Coverage output directory (default: "coverage")
	OutputDir string
}

// BuildConfig holds Poetry build configuration
type BuildConfig struct {
	// Additional build arguments
	BuildArgs []string
	// Additional dependencies to install
	ExtraDependencies []string
	// Poetry configuration options
	PoetryConfig map[string]string
	// Environment variables for build
	Env map[string]string
	// Cache configuration
	Cache *CacheConfig
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
}

// WithPythonVersion sets the Python version to use
func (m *Python) WithPythonVersion(version string) *Python {
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

// getBaseImage returns the Python base image with the configured version
func (m *Python) getBaseImage() string {
	version := m.PythonVersion
	if version == "" {
		version = "3.12"
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

// getPyPIRegistry returns the configured PyPI registry URL with a default
func (m *Python) getPyPIRegistry() string {
	if m.PyPIConfig == nil || m.PyPIConfig.Registry == "" {
		return "https://upload.pypi.org/legacy/"
	}
	return m.PyPIConfig.Registry
}

// getDefaultPyPIConfig returns default PyPI configuration
func (m *Python) getDefaultPyPIConfig() *PyPIConfig {
	return &PyPIConfig{
		Registry: "https://upload.pypi.org/legacy/",
		SkipExisting: false,
		AllowDirty: false,
		ExtraArgs: []string{},
		Env: make(map[string]string),
		RepositoryName: "pypi",
		SkipBuild: false,
		SkipVerify: false,
	}
}

// Publish builds, tests and publishes the Python package to a registry
func (m *Python) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	config := m.PyPIConfig
	if config == nil {
		config = m.getDefaultPyPIConfig()
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
	for key, value := range config.Env {
		container = container.WithEnvVariable(key, value)
	}

	// Configure Poetry for publishing
	container = container.WithExec([]string{
		"poetry", "config",
		fmt.Sprintf("repositories.%s.url", config.RepositoryName),
		config.Registry,
	})

	// Use provided token or fallback to PyPIConfig token
	publishToken := token
	if publishToken == nil {
		publishToken = config.Token
	}

	// Add authentication if token is provided
	if publishToken != nil {
		container = container.WithSecretVariable(
			fmt.Sprintf("POETRY_PYPI_TOKEN_%s", strings.ToUpper(config.RepositoryName)),
			publishToken,
		)
	} else {
		return "", fmt.Errorf("PyPI token is required for publishing. Use --token flag or configure PyPIConfig")
	}

	// Prepare publish command
	publishCmd := []string{
		"poetry", "publish",
		"--repository", config.RepositoryName,
		"--no-interaction",
	}

	// Add build flag if not skipping build
	if !config.SkipBuild {
		publishCmd = append(publishCmd, "--build")
	}

	// Add optional flags
	if config.SkipExisting {
		publishCmd = append(publishCmd, "--skip-existing")
	}
	if config.AllowDirty {
		publishCmd = append(publishCmd, "--allow-dirty")
	}

	// Add any extra arguments
	publishCmd = append(publishCmd, config.ExtraArgs...)

	// Execute publish command
	return container.WithExec(publishCmd).Stdout(ctx)
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
		},
		Env: make(map[string]string),
	}
}

// Test runs the test suite using pytest with coverage reporting
func (m *Python) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	config := m.TestConfig
	if config == nil {
		config = m.getDefaultTestConfig()
	}

	container := m.BuildEnv(source)

	// Add environment variables
	for key, value := range config.Env {
		container = container.WithEnvVariable(key, value)
	}

	args := []string{"poetry", "run", "pytest"}

	if config.Verbose {
		args = append(args, "--verbose", "--color=yes")
	}

	if config.Workers > 0 {
		args = append(args, fmt.Sprintf("-n=%d", config.Workers))
	}

	if config.Coverage != nil && config.Coverage.Enabled {
		args = append(args, fmt.Sprintf("--cov=%s", m.PackagePath))
		
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
			args = append(args, fmt.Sprintf("--cov-fail-under=%f", config.Coverage.MinCoverage))
		}

		args = append(args, "--no-cov-on-fail")
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
		PoetryConfig: map[string]string{
			"virtualenvs.in-project": "true",
		},
		Env: make(map[string]string),
		Cache: &CacheConfig{
			PipCache: true,
			PoetryCache: true,
			PipCacheVolume: "pip-cache",
			PoetryCacheVolume: "poetry-cache",
		},
	}
}

// BuildEnv prepares a Python development environment with Poetry
func (m *Python) BuildEnv(source *dagger.Directory) *dagger.Container {
	config := m.BuildConfig
	if config == nil {
		config = m.getDefaultBuildConfig()
	}

	container := dag.Container().From(m.getBaseImage())

	// Add environment variables
	for key, value := range config.Env {
		container = container.WithEnvVariable(key, value)
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
		WithWorkdir(m.getWorkdir("/app"))

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
	for key, value := range config.PoetryConfig {
		container = container.WithExec([]string{
			"poetry", "config",
			key, value,
		})
	}

	// Install project dependencies
	installArgs := []string{
		"poetry", "install",
		"--no-interaction",
		"--no-root",
		"--with", "dev",
	}
	installArgs = append(installArgs, config.BuildArgs...)
	
	return container.WithExec(installArgs)
}
