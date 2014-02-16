# fakehttpfs

fakehttpfs is a go test fake filesystem implementing the
[`net/http` Filesystem interface](http://golang.org/pkg/net/http/#FileSystem).
It is designed to make HTTP/filesystem related tests easier to read and write.

## Usage

Run `go get github.com/paulhammond/fakehttpfs` to install.

Create a fake filesystem using the fakehttpfs.FileSystem function:

```go
testFS := fakehttpfs.FileSystem(
        fakehttpfs.File("/robots.txt", "User-agent: *\nDisallow: /"),
        fakehttpfs.Dir("/misc",
                fakehttpfs.File("hello.txt", "Hello")
        )
);

file, err = testFS.Open("/robots.txt")
file, err = testFS.Open("/misc/file.txt")
```

Full documentation is available through
[godoc](http://godoc.org/github.com/paulhammond/fakehttpfs).

## License

MIT license, see [LICENSE.txt](LICENSE.txt) for details.
