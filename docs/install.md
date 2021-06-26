# Installation


## Ready binaries — the easiest way

Go to the [Releases](https://github.com/vkcom/nocolor/releases) page and download the latest version for your OS.

Check that it launches correctly:
```bash
nocolor version
```

*(here and then, we suppose that the `nocolor` binary is available by name)*

You're done! Proceed to the [Getting started](/docs/getting_started.md) page.


## With `go get`

Make sure you have [Go](https://golang.org/dl/) version 1.16 or higher installed.  
Run the following command:
```bash
go get -v github.com/vkcom/nocolor
```

NoColor will be installed to `$GOPATH/bin/nocolor`, which usually expands to `$HOME/go/bin/nocolor`.

For convenience, you can add this folder to the **PATH**.


## Build from source

Make sure you have [Go](https://golang.org/dl/) version 1.16 or higher installed.  
Clone this repository and run `make build`:
```bash
git clone https://github.com/vkcom/nocolor
cd nocolor
make build
```

Optionally, you can pass a name of the binary:
```bash
make build BIN_NAME=nocolor.bin
```

A resulting binary will be placed in the `./build` folder.


## What's next?

Proceed to the [Getting started](/docs/getting_started.md) page.

