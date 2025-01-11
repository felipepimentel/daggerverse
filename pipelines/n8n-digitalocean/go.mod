module dagger/n8n-digitalocean

go 1.21

require (
	dagger.io/dagger v0.9.3
	dagger/digitalocean v0.0.0
	dagger/docker v0.0.0
	dagger/n8n v0.0.0
	github.com/99designs/gqlgen v0.17.31
	github.com/Khan/genqlient v0.6.0
	github.com/vektah/gqlparser/v2 v2.5.6
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
)

replace (
	dagger/digitalocean => ../../libraries/digitalocean
	dagger/docker => ../../libraries/docker
	dagger/n8n => ../n8n
)
