package main

import (
	"context"
	"fmt"
	"io/fs"
	"main/internal/dagger"
	"path/filepath"
	"regexp"
)

const (
	analyzeFolder = "/tmp/analyze"
)

type analyzer struct {
	PatternExclusions []string
	PatternMatches    map[string]PatternMatch
	dir               *dagger.Directory
}

type PatternMatch struct {
	Match    bool
	Patterns []string
}

func newAnalyzer(dir *dagger.Directory, patternExclusions []string, patternMatches map[string]PatternMatch) (*analyzer, error) {
	if patternMatches == nil {
		return nil, fmt.Errorf("pattern has to be set")
	}

	return &analyzer{
		PatternExclusions: patternExclusions,
		PatternMatches:    patternMatches,
		dir:               dir,
	}, nil
}

func (a *analyzer) run(ctx context.Context) error {
	_, err := dag.
		Container().
		From("alpine:latest").
		WithMountedDirectory(analyzeFolder, a.dir).
		Directory(analyzeFolder).
		Export(ctx, analyzeFolder)
	if err != nil {
		return err
	}

	return filepath.WalkDir(analyzeFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, exclusion := range a.PatternExclusions {
			re := regexp.MustCompile(exclusion)
			if re.MatchString(path) {
				return nil
			}
		}

		for k, patternMatch := range a.PatternMatches {
			for _, pattern := range patternMatch.Patterns {
				re := regexp.MustCompile(pattern)
				if re.MatchString(d.Name()) {
					patternMatch.Match = true
					a.PatternMatches[k] = patternMatch
					break
				}
			}
		}

		return nil
	})
}

func (a *analyzer) getMatch() []string {
	var matched []string

	for k, v := range a.PatternMatches {
		if v.Match {
			matched = append(matched, k)
		}
	}
	return matched
}
