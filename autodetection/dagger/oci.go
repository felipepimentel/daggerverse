package main

import (
	"context"
	"main/internal/dagger"
	"slices"
)

var defaultOciPatterns = map[string]PatternMatch{
	"oci": {
		Patterns: []string{
			".*Dockerfile",
			".*Containerfile",
		},
	},
}

type OciAnalyzer struct {
	Matches []string
}

func newOciAnalyzer(ctx context.Context, dir *dagger.Directory, patternExclusions []string) (*OciAnalyzer, error) {
	anlzr, err := newAnalyzer(
		dir,
		patternExclusions,
		defaultOciPatterns,
	)
	if err != nil {
		return nil, err
	}

	err = anlzr.run(ctx)
	if err != nil {
		return nil, err
	}

	return &OciAnalyzer{
		Matches: anlzr.getMatch(),
	}, nil
}

func (n *OciAnalyzer) IsOci() bool {
	return slices.Contains(n.Matches, "oci")
}
