# Git hooks

Installation:

```
$ git config core.hooksPath contrib/githooks
```

## pre-commit

The hook automatically runs `gofmt`, `goimports`, and `misspell`
to correctly format the `.go` files included in the commit, provided
that all the aforementioned commands are installed and available
in the user's search `$PATH` environment variable:

```
$ go install golang.org/x/tools/cmd/goimports@latest
$ go install github.com/golangci/misspell/cmd/misspell@master
```

It also runs `go mod tidy` and `golangci-lint` if available.
