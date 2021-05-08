# Install

### Ready binaries

Go to the [releases](https://github.com/vkcom/nocolor/releases) page and download the required version.

### With `go get`

To build with `go get` you need [Golang 1.16+](https://golang.org/).

After installation, just run the following command:

```
$ go get -v github.com/vkcom/nocolor
```

And **NoColor** will be install to `$GOPATH/bin/nocolor`, which usually expands to `$HOME/go/bin/nocolor`.

For convenience, you can add this path to the **PATH**.

### From source

To build from source you need [Golang 1.16+](https://golang.org/).

After installation, clone this repository:

```
$ git clone https://github.com/vkcom/nocolor
```

And in the **NoColor** source folder, run the following command:

```
$ make build
```

Optionally, you can pass the name of the binary:

```
$ make build BIN_NAME=nocolor.bin
```

As a result, you will receive a binary file ready to run.

## Next steps

- [Getting Started with NoColor](https://github.com/vkcom/nocolor/blob/master/docs/usage.md)
- [Description of the color concept](https://github.com/vkcom/nocolor/blob/master/docs/concept_of_colors.md)

