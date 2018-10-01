<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://godoc.org/github.com/gobuffalo/buffalo-goth"><img src="https://godoc.org/github.com/gobuffalo/buffalo-goth?status.svg" alt="GoDoc" /></a>
<a href="https://travis-ci.org/gobuffalo/buffalo-goth"><img src="https://travis-ci.org/gobuffalo/buffalo-goth.svg?branch=master" alt="Build Status" /></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/buffalo-goth"><img src="https://goreportcard.com/badge/github.com/gobuffalo/buffalo-goth" alt="Go Report Card" /></a>
</p>

# Goth Generator for Buffalo

In Buffalo `v0.9.4` the built in generator for [github.com/markbates/goth](https://github.com/markbates/goth) was removed in favor of this plugin.

## Installation

```bash
$ go get -u github.com/gobuffalo/buffalo-goth
```

## Usage

Generate Users, Routes

```bash
$ buffalo generate goth-auth facebook twitter linkedin etc...
```

Generate Routes only

```bash
$ buffalo generate goth facebook twitter linkedin etc...
```

For more detailed usage visit [https://gobuffalo.io/docs/generators#goth](https://gobuffalo.io/docs/generators#goth).
