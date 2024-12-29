// A generated module for Versioner functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"dagger/versioner/internal/dagger"
)

// VersionConfig holds configuration for version management
type VersionConfig struct {
	// Strategy for versioning (default: "semver")
	Strategy string
	// Prefix for version tags (default: "v")
	TagPrefix string
	// Files to check for version patterns
	VersionFiles []string
	// Custom version pattern regex
	VersionPattern string
	// Whether to validate commit messages (default: true)
	ValidateCommits bool
	// Whether to auto-generate changelog (default: true)
	GenerateChangelog bool
	// Whether to create git tags (default: true)
	CreateTags bool
	// Branch configuration
	BranchConfig *BranchConfig
	// Commit message configuration
	CommitConfig *CommitConfig
}

// BranchConfig holds branch-specific configuration
type BranchConfig struct {
	// Main branch name (default: "main")
	MainBranch string
	// Development branch name (default: "develop")
	DevelopBranch string
	// Release branch prefix (default: "release/")
	ReleaseBranchPrefix string
	// Hotfix branch prefix (default: "hotfix/")
	HotfixBranchPrefix string
	// Feature branch prefix (default: "feature/")
	FeatureBranchPrefix string
}

// CommitType represents a type of commit and its version increment
type CommitType struct {
	// Type of commit (e.g., "feat", "fix")
	Type string
	// Version increment for this type
	Increment VersionIncrement
}

// CommitConfig holds commit message configuration
type CommitConfig struct {
	// Types of changes that trigger version increments
	Types []CommitType
	// Scopes allowed in commit messages
	Scopes []string
	// Whether to require scope (default: false)
	RequireScope bool
	// Whether to allow custom scopes (default: true)
	AllowCustomScopes bool
	// Breaking change indicators
	BreakingChangeIndicators []string
}

// VersionIncrement represents how a version number should be incremented
type VersionIncrement string

const (
	Major VersionIncrement = "major"
	Minor VersionIncrement = "minor"
	Patch VersionIncrement = "patch"
	None  VersionIncrement = "none"
)

// Versioner manages semantic versioning for projects
type Versioner struct {
	// Configuration for versioning
	Config *VersionConfig
}

// getDefaultConfig returns default configuration
func (v *Versioner) getDefaultConfig() *VersionConfig {
	return &VersionConfig{
		Strategy:          "semver",
		TagPrefix:         "v",
		VersionFiles:      []string{"pyproject.toml", "package.json", "VERSION"},
		VersionPattern:    `version\s*=\s*["']?(\d+\.\d+\.\d+(?:-[a-zA-Z0-9.-]+)?(?:\+[a-zA-Z0-9.-]+)?)["']?`,
		ValidateCommits:   true,
		GenerateChangelog: true,
		CreateTags:        true,
		BranchConfig: &BranchConfig{
			MainBranch:          "main",
				DevelopBranch:       "develop",
				ReleaseBranchPrefix: "release/",
				HotfixBranchPrefix:  "hotfix/",
				FeatureBranchPrefix: "feature/",
		},
		CommitConfig: &CommitConfig{
			Types: []CommitType{
				{Type: "feat", Increment: Minor},
				{Type: "fix", Increment: Patch},
				{Type: "perf", Increment: Patch},
				{Type: "refactor", Increment: None},
				{Type: "style", Increment: None},
				{Type: "docs", Increment: None},
				{Type: "test", Increment: None},
				{Type: "build", Increment: None},
				{Type: "ci", Increment: None},
				{Type: "chore", Increment: None},
			},
			Scopes: []string{
				"core",
				"deps",
				"docs",
				"tests",
				"build",
				"ci",
			},
			RequireScope:            false,
			AllowCustomScopes:       true,
			BreakingChangeIndicators: []string{"BREAKING CHANGE:", "BREAKING-CHANGE:", "!"},
		},
	}
}

// WithConfig sets the configuration
func (v *Versioner) WithConfig(config *VersionConfig) *Versioner {
	v.Config = config
	return v
}

// GetCurrentVersion retrieves the current version from source code
func (v *Versioner) GetCurrentVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	config := v.Config
	if config == nil {
		config = v.getDefaultConfig()
	}

	// First try to find version in project files
	for _, file := range config.VersionFiles {
		fmt.Printf("Checking file: %s\n", file)
		if contents, err := source.File(file).Contents(ctx); err == nil {
			fmt.Printf("Found file: %s\n", file)
			re := regexp.MustCompile(config.VersionPattern)
			matches := re.FindStringSubmatch(contents)
			if len(matches) > 1 && v.isValidSemVer(matches[1]) {
				fmt.Printf("Found version %s in %s\n", matches[1], file)
				return matches[1], nil
			}
			fmt.Printf("No valid version found in %s\n", file)
		} else {
			fmt.Printf("Error reading %s: %v\n", file, err)
		}
	}

	fmt.Println("No version found in files, returning default 0.1.0")
	return "0.1.0", nil
}

// BumpVersion increments the version based on commit messages
func (v *Versioner) BumpVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	// Initialize config if nil
	if v.Config == nil {
		v.Config = v.getDefaultConfig()
	}

	// Get current version from files
	currentVersion, err := v.GetCurrentVersion(ctx, source)
	if err != nil {
		return "", fmt.Errorf("failed to get current version: %w", err)
	}

	// Initialize container
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	// Get all commits
	output, err := container.WithExec([]string{
		"git", "log", "--format=%B%n-hash-%n%H",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get git history: %w", err)
	}

	// Analyze commits to determine version increment
	increment := v.analyzeCommits(strings.Split(output, "\n"), v.Config.CommitConfig)
	if increment == None {
		increment = Patch
	}

	// Bump version according to semver rules
	newVersion := v.incrementVersion(currentVersion, increment)

	// Update version in files
	if err := v.updateVersionInFiles(ctx, source, newVersion); err != nil {
		return "", fmt.Errorf("failed to update version in files: %w", err)
	}

	// Create git tag
	if v.Config.CreateTags {
		tag := fmt.Sprintf("%s%s", v.Config.TagPrefix, newVersion)
		container = container.
			WithExec([]string{"git", "config", "--global", "user.email", "versioner@dagger.io"}).
			WithExec([]string{"git", "config", "--global", "user.name", "Dagger Versioner"})

		// Add and commit version changes
		container = container.
			WithExec([]string{"git", "add", "."}).
			WithExec([]string{"git", "commit", "-m", fmt.Sprintf("chore(release): bump version to %s", newVersion)})

		// Create and push tag
		_, err = container.
			WithExec([]string{"git", "tag", "-a", tag, "-m", fmt.Sprintf("Release %s", tag)}).
			Stdout(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to create git tag: %w", err)
		}
	}

	return newVersion, nil
}

// hasTag checks if a specific tag exists
func (v *Versioner) hasTag(ctx context.Context, container *dagger.Container, tag string) bool {
	// Try to get the tag
	_, err := container.WithExec([]string{
		"git", "rev-parse", "--verify", tag,
	}).Stdout(ctx)
	
	// Return true only if the tag exists
	return err == nil
}

// isValidSemVer checks if a version string follows semantic versioning
func (v *Versioner) isValidSemVer(version string) bool {
	semverPattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	match, _ := regexp.MatchString(semverPattern, version)
	return match
}

// analyzeCommits determines version increment based on commit messages
func (v *Versioner) analyzeCommits(commits []string, config *CommitConfig) VersionIncrement {
	increment := None

	for _, commit := range commits {
		if commit == "" {
			continue
		}

		re := regexp.MustCompile(`^(\w+)(?:\(([^)]+)\))?!?: (.+)`)
		matches := re.FindStringSubmatch(commit)
		if len(matches) < 4 {
			continue
		}

		commitType := matches[1]

		// Get increment type for commit type
		for _, ct := range config.Types {
			if ct.Type == commitType {
				if ct.Increment > increment {
					increment = ct.Increment
				}
				break
			}
		}

		// Check for breaking changes
		for _, indicator := range config.BreakingChangeIndicators {
			if strings.Contains(commit, indicator) {
				return Major
			}
		}
	}

	return increment
}

// incrementVersion bumps the version number according to semver rules
func (v *Versioner) incrementVersion(version string, increment VersionIncrement) string {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return version
	}

	major := parts[0]
	minor := parts[1]
	patch := parts[2]

	switch increment {
	case Major:
		majorNum := v.parseInt(major)
		return fmt.Sprintf("%d.0.0", majorNum+1)
	case Minor:
		minorNum := v.parseInt(minor)
		return fmt.Sprintf("%s.%d.0", major, minorNum+1)
	case Patch:
		patchNum := v.parseInt(patch)
		return fmt.Sprintf("%s.%s.%d", major, minor, patchNum+1)
	default:
		return version
	}
}

// parseInt safely converts string to int
func (v *Versioner) parseInt(s string) int {
	num := 0
	fmt.Sscanf(s, "%d", &num)
	return num
}

// updateVersionInFiles updates version in all configured files
func (v *Versioner) updateVersionInFiles(ctx context.Context, source *dagger.Directory, version string) error {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Check if pyproject.toml exists and update it
	if _, err := source.File("pyproject.toml").Contents(ctx); err == nil {
		fmt.Printf("Updating version to %s in pyproject.toml\n", version)
		container = container.
			WithExec([]string{"apk", "add", "--no-cache", "python3", "py3-pip", "poetry"})
		
		_, err := container.
			WithExec([]string{"poetry", "version", version}).
			Stdout(ctx)
		if err != nil {
			fmt.Printf("Error updating pyproject.toml: %v\n", err)
			return fmt.Errorf("failed to update version in pyproject.toml: %w", err)
		}
		fmt.Printf("Successfully updated pyproject.toml to version %s\n", version)
		return nil
	}

	fmt.Println("No recognized version files found")
	return nil
}

// createTag creates a git tag for the new version
func (v *Versioner) createTag(ctx context.Context, source *dagger.Directory, tag string) error {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Install git
	container = container.WithExec([]string{
		"apk", "add", "--no-cache", "git",
	})

	// Configure git
	container = container.WithExec([]string{
		"git", "config", "--global", "user.email", "versioner@dagger.io",
	}).WithExec([]string{
		"git", "config", "--global", "user.name", "Dagger Versioner",
	})

	// Create tag
	_, err := container.WithExec([]string{
		"git", "tag", "-a", tag, "-m", fmt.Sprintf("Release %s", tag),
	}).Stdout(ctx)

	return err
}

// generateChangelog generates or updates CHANGELOG.md
func (v *Versioner) generateChangelog(ctx context.Context, source *dagger.Directory, version string) error {
	config := v.Config
	if config == nil {
		config = v.getDefaultConfig()
	}

	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Install git
	container = container.WithExec([]string{
		"apk", "add", "--no-cache", "git",
	})

	// Get all commits if this is the first version
	var output string
	var err error
	
	if v.hasTag(ctx, container, fmt.Sprintf("%s%s", config.TagPrefix, version)) {
		// Get commits since last version
		output, err = container.WithExec([]string{
			"git", "log", "--format=%B%n-hash-%n%H",
			fmt.Sprintf("%s%s..HEAD", config.TagPrefix, version),
		}).Stdout(ctx)
	} else {
		// Get all commits for first version
		output, err = container.WithExec([]string{
			"git", "log", "--format=%B%n-hash-%n%H",
		}).Stdout(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to get git history: %w", err)
	}

	// Parse commits and generate changelog
	changelog := v.formatChangelog(strings.Split(output, "\n"), version)

	// Write changelog to file
	_, err = container.WithNewFile(
		"CHANGELOG.md",
		changelog,
		dagger.ContainerWithNewFileOpts{
			Permissions: 0644,
		},
	).WithExec([]string{
		"git", "add", "CHANGELOG.md",
	}).WithExec([]string{
		"git", "commit", "-m", fmt.Sprintf("docs: update changelog for version %s", version),
	}).Stdout(ctx)

	return err
}

// formatChangelog formats commits into changelog entries
func (v *Versioner) formatChangelog(commits []string, version string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s%s\n\n", v.Config.TagPrefix, version))

	// Group commits by type
	typeGroups := make(map[string][]string)
	for _, commit := range commits {
		if commit == "" {
			continue
		}

		re := regexp.MustCompile(`^(\w+)(?:\(([^)]+)\))?!?: (.+)`)
		matches := re.FindStringSubmatch(commit)
		if len(matches) < 4 {
			continue
		}

		commitType := matches[1]
		scope := matches[2]
		message := matches[3]

		entry := message
		if scope != "" {
			entry = fmt.Sprintf("**%s**: %s", scope, message)
		}

		typeGroups[commitType] = append(typeGroups[commitType], entry)
	}

	// Write grouped commits
	for _, commitType := range []string{"feat", "fix", "perf", "refactor", "docs", "style", "test", "build", "ci", "chore"} {
		entries := typeGroups[commitType]
		if len(entries) > 0 {
			sb.WriteString(fmt.Sprintf("### %s\n\n", strings.Title(commitType)))
			for _, entry := range entries {
				sb.WriteString(fmt.Sprintf("* %s\n", entry))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ValidateCommitMessage checks if a commit message follows the conventional commit format
func (v *Versioner) ValidateCommitMessage(message string) error {
	config := v.Config
	if config == nil {
		config = v.getDefaultConfig()
	}

	// Check if empty
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	// Parse conventional commit format
	re := regexp.MustCompile(`^(\w+)(?:\(([^)]+)\))?!?: (.+)`)
	matches := re.FindStringSubmatch(message)
	if len(matches) < 4 {
		return fmt.Errorf("commit message must follow format: type(scope): description")
	}

	commitType := matches[1]
	scope := matches[2]
	description := matches[3]

	// Validate type
	validType := false
	for _, ct := range config.CommitConfig.Types {
		if ct.Type == commitType {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid commit type: %s", commitType)
	}

	// Validate scope if required
	if config.CommitConfig.RequireScope && scope == "" {
		return fmt.Errorf("commit scope is required")
	}

	// Validate scope if custom scopes not allowed
	if !config.CommitConfig.AllowCustomScopes && scope != "" {
		validScope := false
		for _, s := range config.CommitConfig.Scopes {
			if s == scope {
				validScope = true
				break
			}
		}
		if !validScope {
			return fmt.Errorf("invalid commit scope: %s", scope)
		}
	}

	// Validate description
	if len(description) < 10 {
		return fmt.Errorf("commit description must be at least 10 characters")
	}

	return nil
}

// PushChanges commits and pushes version changes
func (v *Versioner) PushChanges(ctx context.Context, source *dagger.Directory, version string) error {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Install git
	container = container.WithExec([]string{
		"apk", "add", "--no-cache", "git",
	})

	// Configure git
	container = container.WithExec([]string{
		"git", "config", "--global", "user.email", "versioner@dagger.io",
	}).WithExec([]string{
		"git", "config", "--global", "user.name", "Dagger Versioner",
	})

	// Add all changes
	_, err := container.WithExec([]string{
		"git", "add", ".",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Commit changes
	_, err = container.WithExec([]string{
		"git", "commit", "-m", fmt.Sprintf("chore(release): bump version to %s", version),
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Push changes and tags
	_, err = container.WithExec([]string{
		"git", "push", "origin", "HEAD", "--tags",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	return nil
}

// GetVersionHistory returns a list of all versions with their commit messages
func (v *Versioner) GetVersionHistory(ctx context.Context, source *dagger.Directory) (string, error) {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Install git
	container = container.WithExec([]string{
		"apk", "add", "--no-cache", "git",
	})

	// Get all tags with their messages
	output, err := container.WithExec([]string{
		"git", "tag", "-l", "--format=%(tag) %(subject)",
	}).Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to get version history: %w", err)
	}

	return output, nil
}
