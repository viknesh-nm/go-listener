# go-listener

Build go projects automatically when files get modified

## Install

```go
    go get -u github.com/viknesh-nm/go-listener
```

## Usage

- Update `PATH` with `GOPATH/bin`

    ```
    Usage:
        go-listener [flags]

    Flags:
        -b, --build            build only mode that generates the binary file
        -c, --command string   extra additional commands that needs to be run
        -h, --help             help for go-listener
        -n, --name string      name of the project
        -p, --path string      directory to be watch
    ```
