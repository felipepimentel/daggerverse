# Presentation

## Writing the slides

Write the content of the slides using [AsciiDoc] syntax with the added features of [Asciidoctor reveal.js] and [Asciidoctor Kroki] in the `index.adoc` file.

## Building the presentation

### Prerequisites

Authentication to GitHub Packages is needed to install the dependencies. Install [GitHub CLI] and execute `gh auth login --scopes 'read:packages'`.

### Option 1: npm

NOTE: If the presentation contains diagrams, modify the `kroki-server-url` attribute in the document header (at the top of the `index.adoc` file) to point to a running [Kroki] server.

Install Node.js and npm.

Execute the following commands to build the presentation:

```shell
npm set "//npm.pkg.github.com/:_authToken=$(gh auth token)"
npm clean-install
npm run build
```

To continually rebuild the presentation each time the slides are modified, use the following command:

```shell
npm run autobuild
```

To preview the presentation, execute the following command (it can be executed in parallel with the automatic build), then visit [localhost:1234](http://localhost:1234):

```shell
npm run serve
```

### Option 2: Dagger

Install [Dagger].

Execute the following command to build the presentation:

```shell
dagger call -m presentation builder --directory '.' --npmrc 'cmd:echo //npm.pkg.github.com/:_authToken=$(gh auth token)' build directory --output 'dist'
```

To preview the presentation, execute the following command, then visit [localhost:8080](http://localhost:8080):

```shell
dagger call -m presentation builder --directory '.' --npmrc 'cmd:echo //npm.pkg.github.com/:_authToken=$(gh auth token)' build server up
```

[AsciiDoc]: https://docs.asciidoctor.org/asciidoc/latest/
[Asciidoctor reveal.js]: https://docs.asciidoctor.org/reveal.js-converter/latest/
[Asciidoctor Kroki]: https://github.com/asciidoctor/asciidoctor-kroki
[Kroki]: https://kroki.io/
[GitHub CLI]: https://cli.github.com/
[Dagger]: https://dagger.io/
