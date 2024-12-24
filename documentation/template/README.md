# Documentation

## Writing documentation

This documentation is based on the [Diátaxis] documentation system. Follow the principles it prescribes.

Start by choosing the kind of documentation to write based on the [Diátaxis compass](https://diataxis.fr/compass/).

Then write one or more pages using [AsciiDoc] syntax with the added features of [Asciidoctor Kroki] and [Hugo] (always prefer AsciiDoc syntax over Hugo shortcodes) in the `content` folder.

## Building the documentation

### Option 1: npm

NOTE: If the documentation contains diagrams, modify the `kroki-server-url` attribute in the the `markup.yaml` file of the `config` folder to point to a running [Kroki] server.

Install Node.js and npm.

Execute the following commands to build the documentation:

```shell
npm clean-install
npm run build
```

To preview the documentation, continually rebuilt each time it is modified, execute the following command, then visit [localhost:1313](http://localhost:1313):

```shell
npm run serve
```

### Option 2: Dagger

Install [Dagger].

Execute the following command to build the documentation:

```shell
dagger call -m documentation builder --directory '.' build directory --output 'public'
```

To preview the documentation, execute the following command, then visit [localhost:8080](http://localhost:8080):

```shell
dagger call -m documentation builder --directory '.' build server up
```

[Diátaxis]: https://diataxis.fr/
[AsciiDoc]: https://docs.asciidoctor.org/asciidoc/latest/
[Asciidoctor Kroki]: https://github.com/asciidoctor/asciidoctor-kroki
[Hugo]: https://gohugo.io/
[Kroki]: https://kroki.io/
[Dagger]: https://dagger.io/
