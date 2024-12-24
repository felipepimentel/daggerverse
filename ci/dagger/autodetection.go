package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"main/internal/dagger"
)

func (c *Ci) Autodetection(ctx context.Context, testDataSrc *dagger.Directory) error {
	var eg errgroup.Group

	eg.Go(func() error {
		nodeAnalyzer := dag.
			Autodetection().
			Node(testDataSrc.Directory("myapi"))

		detectTest, err := nodeAnalyzer.IsTest(ctx)
		if err != nil {
			return err
		}
		if !detectTest {
			return fmt.Errorf("should detect test")
		}

		detectPackage, err := nodeAnalyzer.IsPackage(ctx)
		if err != nil {
			return err
		}
		if detectPackage {
			return fmt.Errorf("should not detect package")
		}

		detectNpm, err := nodeAnalyzer.IsNpm(ctx)
		if err != nil {
			return err
		}
		if !detectNpm {
			return fmt.Errorf("should detect npm")
		}

		detectYarn, err := nodeAnalyzer.IsYarn(ctx)
		if err != nil {
			return err
		}
		if detectYarn {
			return fmt.Errorf("should detect yarn")
		}

		detectLint, err := nodeAnalyzer.Is(ctx, "lint")
		if err != nil {
			return err
		}
		if detectLint {
			return fmt.Errorf("should not detect lint")
		}

		return nil
	})

	return eg.Wait()
}
