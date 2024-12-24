package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"main/internal/dagger"
	"strings"
)

func (c *Ci) Node(ctx context.Context, testDataSrc *dagger.Directory) error {
	var eg errgroup.Group

	// Lazy mode pipeline with oci build
	eg.Go(func() error {
		refs, err := dag.
			Node().
			WithAutoSetup(
				"testdata-myapi",
				testDataSrc.Directory("myapi"),
			).
			Pipeline(
				ctx,
				dagger.NodePipelineOpts{
					DryRun: true,
					TTL:    "5m",
					IsOci:  true,
				},
			)

		fmt.Println("image: " + refs)

		return err
	})

	// Explicit mode pipeline with oci build
	eg.Go(func() error {
		refs, err := dag.
			Node().
			WithPipelineID("testdata-myapi").
			WithVersion("20.9.0").
			WithSource(testDataSrc.Directory("myapi")).
			WithNpm().
			Install().
			Test().
			Build().
			OciBuild(ctx, nil, dagger.NodeOciBuildOpts{IsTTL: true, TTL: "5m"})

		fmt.Println("image: " + strings.Join(refs, "\n"))

		return err
	})

	//Lazy mode pipeline with package build
	eg.Go(func() error {
		_, err := dag.
			Node().
			WithAutoSetup(
				"testdata-lib",
				testDataSrc.Directory("mylib"),
			).
			Pipeline(
				ctx,
				dagger.NodePipelineOpts{
					DryRun:        true,
					PackageDevTag: "beta",
				},
			)

		return err
	})

	// Explicit mode pipeline with package build
	eg.Go(func() error {
		_, err := dag.
			Node().
			WithPipelineID("testdata-mylib").
			WithVersion("20.9.0").
			WithSource(testDataSrc.Directory("mylib")).
			WithNpm().
			Install().
			Test().
			Build().
			Publish(dagger.NodePublishOpts{DryRun: true, DevTag: "beta"}).
			Do(ctx)

		return err
	})

	return eg.Wait()
}
