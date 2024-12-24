package main

type KindConfig struct {
	Kind string `yaml:"kind"`

	ApiVersion string `yaml:"apiVersion"`

	Networking Networking `yaml:"networking"`
}

type Networking struct {
	ApiServerPort int `yaml:"apiServerPort"`
}
