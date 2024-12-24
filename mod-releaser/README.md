# Mod Releaser

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/mod-releaser)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.2-0f0f19.svg?style=flat-square)

A module which publish dagger module

## Usage

*This example publish the module `mod-releaser` and the call is run inside this module.*  
*Also the daggerverse repository is inside the home so adapt the command according where your repository is located.*

```shell
dagger call \
  --git-repo ..
  --component mod-releaser \
  with-git-config --cfg=$HOME/.gitconfig \
  minor \
  publish \
  repository \
  export --path $HOME/daggerverse/
```

# TODO

 * [ ] Implement ssh agent socket when socket is available as input
 * [ ] Automatic decision to bump major / minor / patch