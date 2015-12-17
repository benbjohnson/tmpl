# tmpl

This program is a command line interface to Go's `text/template` library. It
can be used by passing in a set of JSON-encoded data and a list of template
paths ending in a `.tmpl` extension. The templates are processed and their
results are saved to the filename with the `.tmpl` extension removed.


## Getting Started

To install `tmpl`, simply run:

```sh
$ go get github.com/benbjohnson/tmpl
```

Then run `tmpl` against your desired templates:

```sh
$ tmpl -data '["foo","bar","baz"]' a.go.tmpl b.go.tmpl
```

You will now have templates generated at `a.go` and `b.go`.
