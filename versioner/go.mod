module github.com/felipepimentel/daggerverse/versioner

go 1.21

require (
	dagger.io/dagger v0.9.8
	github.com/99designs/gqlgen v0.17.62
	github.com/Khan/genqlient v0.7.0
	github.com/vektah/gqlparser/v2 v2.5.21
)

require (
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace dagger.io/dagger => github.com/dagger/dagger/sdk/go v0.9.8
