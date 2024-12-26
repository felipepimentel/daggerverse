// A generated module for Shellchecktest functions
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
	"regexp"

	"github.com/sourcegraph/conc/pool"
)

// Shellchecktest is a module for testing shell scripts.
type Shellchecktest struct{}

// All runs all tests.
func (m *Shellchecktest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckDirectory)

	return p.Wait()
}

// CheckDirectory runs a test on a directory.
func (m *Shellchecktest) CheckDirectory(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Shellcheck().Check(dir).Stderr(ctx)

	if err != nil {
		re := regexp.MustCompile("exit code: 123")
		if re.MatchString(err.Error()) {
			return nil
		}
	}

	return err
}
