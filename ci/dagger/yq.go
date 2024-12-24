package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"main/internal/dagger"
)

func (c *Ci) Yq(ctx context.Context, testDataSrc *dagger.Directory) error {
	var eg errgroup.Group

	//Read a key
	eg.Go(func() error {
		val, err := dag.
			Yq(testDataSrc).
			Get(ctx, ".foo.bar", "test.yaml")
		if err != nil {
			return err
		}

		if val != "qux" {
			return fmt.Errorf("expression '.foo.bar' should return 'qux'")
		}

		return nil
	})

	//eg.Go(func() error {
	//	val, err := dag.
	//		Pipeline("Edit the file dans read it").
	//		Yq(testDataSrc).
	//		Set(".foo.bar=\"super_qux\"", "test.yaml").
	//		Get(ctx, ".foo.bar", "test.yaml")
	//	if err != nil {
	//		return err
	//	}
	//
	//	if val != "super_qux" {
	//		return fmt.Errorf("expression '.foo.bar' should return 'super_qux' but")
	//	}
	//
	//	return err
	//})

	return eg.Wait()
}
