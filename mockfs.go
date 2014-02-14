// package mockfs provides a mock filesystem object that implements the
// http.FileSystem interface.
//
// To use it, call the MockFileSystem function with one or more http.File
// objects. The MockFile and Mockdir helper functions create files and
// directories respectively. You do not have to use MockFile, if you'd like
// to write your own mock or stub or even use a real file you can. For
// example:
//
//     testFS := mockfs.FileSystem(
//             mockfs.File("/robots.txt", "User-agent: *\nDisallow: /"),
//             mockfs.Dir("/misc",
//                     mockfs.File("hello.txt", "Hello")
//                     os.Open("/path/to/some/real/file.txt")
//             )
//     );
//
//     file, err := testFS.Open("/robots.txt")
//     file, err := testFS.Open("/misc/file.txt")
package mockfs

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

// Creates a mock filesystem containing the files.
func FileSystem(files ...http.File) http.FileSystem {
	return mockDir{"", files}
}

//  a mock file with string contents.
func File(name, contents string) http.File {
	b := []byte(contents)
	return mockFile{name, int64(len(b)), bytes.NewReader(b)}
}

type mockFile struct {
	name string
	size int64
	*bytes.Reader
}

func (f mockFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f mockFile) Readdir(int) ([]os.FileInfo, error) {
	return nil, errors.New("Not dir")
}

func (f mockFile) Read(b []byte) (int, error) {
	return f.Reader.Read(b)
}

func (f mockFile) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}

func (f mockFile) Close() error {
	_, err := f.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	return nil
}

func (f mockFile) Name() string {
	return f.name
}

func (f mockFile) Size() int64 {
	return int64(f.Reader.Len())
}

func (f mockFile) Mode() os.FileMode {
	panic("unimplemented")
}

func (f mockFile) ModTime() time.Time {
	panic("unimplemented")
}

func (f mockFile) IsDir() bool {
	return false
}

func (f mockFile) Sys() interface{} {
	return nil
}

// Creates a mock directory containing the files.
func Dir(name string, files ...http.File) http.File {
	return mockDir{name, files}
}

type mockDir struct {
	name  string
	files []http.File
}

func (d mockDir) Open(name string) (http.File, error) {
	parts := strings.SplitN(name, "/", 2)
	file, err := d.find(parts[0])
	if len(parts) == 1 {
		return file, err
	}
	if subDir, ok := file.(mockDir); ok {
		return subDir.Open(parts[1])
	}
	return nil, os.ErrNotExist
}

func (d mockDir) find(name string) (http.File, error) {
	if name == "" || name == "." {
		return d, nil
	}
	for _, file := range d.files {
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if stat.Name() == name {
			return file, nil
		}
	}
	return nil, os.ErrNotExist
}

func (d mockDir) Stat() (os.FileInfo, error) {
	return d, nil
}

func (d mockDir) Readdir(count int) ([]os.FileInfo, error) {
	var err error
	r := make([]os.FileInfo, len(d.files))
	for i, f := range d.files {
		r[i], err = f.Stat()
		if err != nil {
			return r[0 : i-1], err
		}
	}
	return r, nil
}

func (d mockDir) Read(b []byte) (int, error) {
	return 0, errors.New("Not regular file")
}

func (d mockDir) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("Not regular file")
}

func (d mockDir) Close() error {
	return nil
}

func (d mockDir) Name() string {
	return d.name
}

func (d mockDir) Size() int64 {
	return 0
}

func (d mockDir) Mode() os.FileMode {
	panic("unimplemented")
}

func (d mockDir) ModTime() time.Time {
	panic("unimplemented")
}

func (d mockDir) IsDir() bool {
	return true
}

func (d mockDir) Sys() interface{} {
	return nil
}
