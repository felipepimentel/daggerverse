module github.com/felipepimentel/daggerverse/python

go 1.21.5
require (
	dagger.io/dagger v0.9.8
	github.com/99designs/gqlgen v0.17.62
	github.com/Khan/genqlient v0.6.0
	github.com/felipepimentel/daggerverse/python-builder v0.0.0-00010101000000-000000000000
	github.com/felipepimentel/daggerverse/python-publisher v0.0.0-00010101000000-000000000000
	github.com/felipepimentel/daggerverse/python-versioner v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.5.21
	golang.org/x/sync v0.10.0
)

require (
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/agnivade/levenshtein v1.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/urfave/cli/v2 v2.27.5 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/felipepimentel/daggerverse/python-builder => ../python-builder
	github.com/felipepimentel/daggerverse/python-publisher => ../python-publisher
	github.com/felipepimentel/daggerverse/python-versioner => ../python-versioner
)
