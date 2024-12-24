# Auto-detection

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.2-0f0f19.svg?style=flat-square)

A module which analyze a project in order to extract information for lazy mode module

## Features

 * Node:
   * Extract information from package.json
     * Application name and version
     * Engine version
     * Define if this is a package or not
   * Define if there is some tests in the project
   * Detect the package manager
 * OCI:
   * Detect if a dockerfile or containerfile is present in the repository

## Prerequisite
### Node

The module expects to find some information in the `package.json`:
* Fields:
   * `.name`
   * `.engines.node`
* Scripts:
   * `test` (required if test are find): expect a command to run tests
   * `build` (required): expect to command to build / transpile the code
   * `clean` (required): cleanup the project
   * `lint` (optional): expect a command to check the lint not fixing it

## Examples

```go
nodeAnalyzer := dag.
   Autodetection().
   Node(
      src,
      dagger.AutodetectionNodeOpts{
        PatternExclusions: []string{"node_modules"},
      },
   )

engineVersion, err := nodeAnalyzer.GetEngineVersion(ctx)
isYarn, err := nodeAnalyzer.IsYarn(ctx)
isNpm, err := nodeAnalyzer.IsNpm(ctx)
appVersion, err := nodeAnalyzer.GetVersion(ctx)
nodeAutoSetup.Version = appVersion
appName, err := nodeAnalyzer.GetName(ctx)
isTest, err := nodeAnalyzer.IsTest(ctx)
isPackage, err := nodeAnalyzer.IsPackage(ctx)
istLint, err := nodeAnalyzer.Is(ctx, "lint")
```

### OCI
```go
isOci, err := dag.
	Autodetection().
    Oci(
        src,
        dagger.AutodetectionOciOpts{
            PatternExclusions: []string{"node_modules"},
        },
   ).
   IsOci(ctx)
```

more example in the `/ci/node.go`

## To Do

- [ ] Add golang
- [ ] Add python
