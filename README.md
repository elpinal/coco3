# Coco3

Coco3 is a shell written in Go, providing Vim-like key mappings and painless
syntax.

## Install

To install, use `go get`:

```bash
$ go get -u github.com/elpinal/coco3
```

-u flag stands for "update".

## Documentation

To get help, type ":help" or execute "help" built-in command.

## Motivation

Most shells have too complicated syntax. Some shells have an insufficient
editing facility. These facts lead to developing Coco3 shell.

## Note

If you are curious, try the <b>extra mode</b>!

```
$ coco3 -extra
```

The behavior of the <b>extra mode</b> is not documented, so read the source code in `extra` directory.

## Contribution

1. Fork ([https://github.com/elpinal/coco3/fork](https://github.com/elpinal/coco3/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[elpinal](https://github.com/elpinal)
