package main

import (
	"context"
	"dagger/kind/internal/dagger"
	"fmt"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Kind struct {
	DockerSocket *dagger.Socket
	KindSvc      *dagger.Service
	KindPort     int
	Container    *dagger.Container
	ClusterName  string
}

func New(
	ctx context.Context,
	// Docker socket path. E.g. /var/run/docker.sock
	// How to use it:
	// dagger call --docker-sock=/var/run/docker.sock kind --kind-svc=tcp://127.0.0.1:3000
	// +required
	dockerSocket *dagger.Socket,

	// It should be the tcp://127.0.0.1 followed by any port. E.g. tcp://127.0.0.1:3000
	// Before launch this function, make sure that you have configured in your /etc/hosts file
	// an entry for localhost 127.0.0.1 . Otherwise, the alpine container will not be able to connect to the kind cluster.
	// +required
	kindSvc *dagger.Service,

	// +optional
	// +default="dagger-kubernetes-cluster"
	// +description=The name of the kind cluster
	clusterName string,

) *Kind {

	ep, err := kindSvc.Endpoint(ctx)

	if err != nil {

		panic(err)

	}

	port, err := strconv.Atoi(strings.Split(ep, ":")[1])

	if err != nil {

		panic(err)

	}

	if port < 1024 || port > 65535 {

		panic(fmt.Sprintf("Invalid port number: %d, it should be between 1024 and 65535", port))

	}

	kindConfig := &KindConfig{
		Kind:       "Cluster",
		ApiVersion: "kind.x-k8s.io/v1alpha4",
		Networking: Networking{
			ApiServerPort: port,
		},
	}

	yamlFileContent, err := yaml.Marshal(kindConfig)

	container, err := dag.Container().
		From("alpine").
		WithUnixSocket("/var/run/docker.sock", dockerSocket).
		WithNewFile("kind.yaml", string(yamlFileContent)).
		WithExec([]string{"apk", "add", "docker", "kubectl", "k9s", "curl"}).
		WithExec([]string{"curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.25.0/kind-linux-amd64"}).
		WithExec([]string{"chmod", "+x", "./kind"}).
		WithExec([]string{"mv", "./kind", "/usr/local/bin/kind"}).
		WithEnvVariable("BUST", time.Now().String()).
		WithExec([]string{
			"kind", "delete", "cluster",
			"--name", clusterName,
		}).
		WithExec([]string{
			"kind", "create", "cluster",
			"--name", clusterName,
			"--config", "kind.yaml",
			"--wait", "1m",
		}).
		WithServiceBinding("localhost", kindSvc).
		WithExec([]string{
			"kubectl", "config",
			"set-cluster", fmt.Sprintf("kind-%s", clusterName), fmt.Sprintf("--server=https://localhost:%d", port)},
		).Sync(ctx)

	if err != nil {

		panic(err)

	}

	return &Kind{
		DockerSocket: dockerSocket,
		KindSvc:      kindSvc,
		KindPort:     port,
		Container:    container,
		ClusterName:  clusterName,
	}
}

// Loads a container to kind cluster, previously it was saved as a tarball with
// the annotations required by kind. You can use this function into your module
func (m *Kind) LoadContainerOnKind(

	ctx context.Context,

	container *dagger.Container,

	tag string,

) *dagger.Container {

	containerFileTaName := fmt.Sprintf("%s.tar", tag)

	tarball := container.
		// This is the image name that will be loaded in the kind cluster
		WithAnnotation(
			"org.opencontainers.image.ref.name",
			fmt.Sprintf("%s:latest", tag),
		).

		// Kind requires the docker.io/library prefix, otherwise it will load the image
		// This a fake image name in docker.io, it is not a real image.
		// You should user imagePullPolicy: Never in your Kind manifests.
		WithAnnotation(
			"io.containerd.image.name",
			fmt.Sprintf("docker.io/library/%s:latest", tag),
		).
		AsTarball()

	return m.Container.
		WithFile(containerFileTaName, tarball).
		WithEnvVariable("BUST", time.Now().String()).
		WithExec([]string{"kind", "load", "image-archive", containerFileTaName, "--name", m.ClusterName}).
		WithoutFile(containerFileTaName)

}

// Launch k9s terminal
// Example usage:
// dagger call --docker-socket=/var/run/docker.sock --kind-svc=tcp://127.0.0.1:3000 knines
func (m *Kind) Knines(

	ctx context.Context,

) *dagger.Container {

	return m.Container.Terminal(dagger.ContainerTerminalOpts{
		Cmd: []string{"k9s"},
	})

}

// Inspect returns the container that will be launched
// Example usage:
// dagger call --docker-socket=/var/run/docker.sock --kind-svc=tcp://127.0.0.1:3000 inspect
func (m *Kind) Inspect(

	ctx context.Context,

) *dagger.Container {

	return m.Container.Terminal()

}

// Inspect returns the container that will be launched
// Example usage:
// dagger call --docker-socket=/var/run/docker.sock --kind-svc=tcp://127.0.0.1:3000 inspect
func (m *Kind) Test(

	ctx context.Context,

) *dagger.Container {

	return dag.Container().From("node:20").From("node:20-alpine")

}
