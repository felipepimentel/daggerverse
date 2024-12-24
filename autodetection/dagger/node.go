package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"main/internal/dagger"
	"os"
	"regexp"
	"slices"
)

var defaultNodeExclude = []string{
	"node_modules",
	".tsconfig",
}

var defaultNodePatterns = map[string]PatternMatch{
	"test": {
		Patterns: []string{
			".+\\.(test|spec)\\.js",
			".+\\.(test|spec)\\.jsx",
			".+\\.(test|spec)\\.ts",
			".+\\.(test|spec)\\.tsx",
			"(.+/)*(__)*tests*(__)*/.+",
		},
	},
	"yarn": {
		Patterns: []string{
			".*yarn.lock",
		},
	},
	"npm": {
		Patterns: []string{
			".*package-lock.json",
		},
	},
}

type packageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Workspaces      []string          `json:"workspaces"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Repository      *repository       `json:"repository"`
	Engines         *engines          `json:"engines"`
	PublishConfig   *publishConfig    `json:"publishConfig,omitempty"`
}

type repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type engines struct {
	Node string `json:"node"`
}

type publishConfig struct {
	Registry string `json:"registry"`
}

type NodeAnalyzer struct {
	Matches    []string
	PkgJsonRep string
}

func newNodeAnalyzer(ctx context.Context, dir *dagger.Directory, patternExclusions []string) (*NodeAnalyzer, error) {
	anlzr, err := newAnalyzer(
		dir,
		append(patternExclusions, defaultNodeExclude...),
		defaultNodePatterns,
	)
	if err != nil {
		return nil, err
	}

	err = anlzr.run(ctx)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(analyzeFolder + "/package.json")
	if err != nil {
		return nil, err
	}

	return &NodeAnalyzer{
		Matches:    anlzr.getMatch(),
		PkgJsonRep: string(content),
	}, nil
}

func (n NodeAnalyzer) toPkgJson() (*packageJson, error) {
	pkgJson := packageJson{}
	err := json.Unmarshal([]byte(n.PkgJsonRep), &pkgJson)
	if err != nil {
		return nil, err
	}

	return &pkgJson, nil
}

func (n *NodeAnalyzer) IsTest() bool {
	return slices.Contains(n.Matches, "test")
}

func (n *NodeAnalyzer) IsYarn() bool {
	return slices.Contains(n.Matches, "yarn")
}

func (n *NodeAnalyzer) IsNpm() bool {
	return slices.Contains(n.Matches, "npm")
}

func (n *NodeAnalyzer) IsPackage() (bool, error) {
	info, err := n.toPkgJson()
	if err != nil {
		return false, err
	}

	return info.PublishConfig != nil, nil
}

func (n *NodeAnalyzer) Is(
	// Define if a script is present or not in the package.json
	scriptName string,
) (bool, error) {
	scriptNames, err := n.GetScriptNames()
	if err != nil {
		return false, err
	}

	return slices.Contains(scriptNames, scriptName), nil
}

func (n *NodeAnalyzer) GetEngineVersion() (string, error) {
	info, err := n.toPkgJson()
	if err != nil {
		return "", err
	}

	nodeEngineVersion := regexp.
		MustCompile(`(\d+\.\d+\.\d+)`).
		FindStringSubmatch(info.Engines.Node)
	if len(nodeEngineVersion) != 2 {
		return "", fmt.Errorf("not able to parse the node engine version: '%s'", info.Engines.Node)
	}

	return nodeEngineVersion[1], nil
}

func (n *NodeAnalyzer) GetVersion() (string, error) {
	info, err := n.toPkgJson()
	if err != nil {
		return "", err
	}

	return info.Version, nil
}

func (n *NodeAnalyzer) GetName() (string, error) {
	info, err := n.toPkgJson()
	if err != nil {
		return "", err
	}

	return info.Name, nil
}

func (n *NodeAnalyzer) GetScriptNames() ([]string, error) {
	info, err := n.toPkgJson()
	if err != nil {
		return nil, err
	}

	return maps.Keys(info.Scripts), nil
}

func (n *NodeAnalyzer) GetWorkspaces() ([]string, error) {
	info, err := n.toPkgJson()
	if err != nil {
		return nil, err
	}

	return info.Workspaces, nil
}
