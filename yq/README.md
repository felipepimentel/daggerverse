# Yq

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.2-0f0f19.svg?style=flat-square)

A yq module

## Features

* Read a given key
* Edit the yaml


## Examples

### Read a field

```go
dag.
    Yq(testDataSrc).
    Get(ctx, ".foo.bar", "test.yaml")                                                                                          
```

```shell
dagger call -m "github.com/Dudesons/daggerverse/yq" --source ../testdata/yq/ \
  get --expr=".foo" --yaml-file-path=test.yaml
```

### Edit

```go
dag.
    Yq(testDataSrc).
    Set(".foo.bar=\"super_qux\"", "test.yaml").
    Get(ctx, ".foo.bar", "test.yaml")
```

```shell
dagger call -m "github.com/Dudesons/daggerverse/yq" --source ../testdata/yq/ \
  set --expr=".foo.bar=\"toto\"" --yaml-file-path=test.yaml \
  get --expr=".foo" --yaml-file-path=test.yaml
```

### Export to host

```shell
dagger call -m "github.com/Dudesons/daggerverse/yq" --source ../testdata/yq/ \
  set --expr=".foo.bar=\"toto\"" --yaml-file-path=test.yaml \
  state \
  export --path=/tmp/testdata3/yq/
```

### Open a shell

```shell
dagger call -m "github.com/Dudesons/daggerverse/yq" --source ../testdata/yq/ --yaml-file-path=test.yaml \
  shell                                                                                        
```


**more example in the `/ci/yq.go`**
