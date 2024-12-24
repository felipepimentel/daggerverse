package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"main/internal/dagger"
)

func (c *Ci) Infrabox(ctx context.Context, testDataSrc *dagger.Directory) error {
	var eg errgroup.Group

	eg.Go(func() error {
		_, err := dag.
			Infrabox().
			Terragrunt().
			WithSource("/terraform", testDataSrc.Directory("terraform")).
			DisableColor().
			Plan("/terraform/stacks/dev/europe-west1/staging/qux").
			Apply("/terraform/stacks/dev/europe-west1/staging/qux").
			Plan("/terraform/stacks/dev/europe-west1/staging/qux", dagger.InfraboxTfPlanOpts{DetailedExitCode: true}).
			Do(ctx)

		return err
	})

	eg.Go(func() error {
		_, err := dag.
			Infrabox().
			Terragrunt().
			WithSource("/terraform", testDataSrc.Directory("terraform")).
			DisableColor().
			Plan("/terraform/stacks/dev/europe-west1/staging/foo", dagger.InfraboxTfPlanOpts{DetailedExitCode: true}).
			Do(ctx)

		if err == nil {
			return errors.New("it should failed because the stack was not applied")
		}

		return nil
	})

	return eg.Wait()
}
