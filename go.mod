module github.com/felipepimentel/daggerverse/python

go 1.21.5

require (
	dagger.io/dagger v0.9.8
	github.com/99designs/gqlgen v0.17.44
	github.com/Khan/genqlient v0.7.0
	github.com/felipepimentel/daggerverse/python-builder v0.0.0-00010101000000-000000000000
	github.com/felipepimentel/daggerverse/python-publisher v0.0.0-00010101000000-000000000000
	github.com/felipepimentel/daggerverse/python-versioner v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.5.11
	golang.org/x/sync v0.6.0
)

require (
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alexflint/go-arg v1.4.2 // indirect
	github.com/alexflint/go-scalar v1.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/felipepimentel/daggerverse/python-builder => ../python-builder
	github.com/felipepimentel/daggerverse/python-publisher => ../python-publisher
	github.com/felipepimentel/daggerverse/python-versioner => ../python-versioner
)
