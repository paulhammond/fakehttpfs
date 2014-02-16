// package fakehttpfs provides a fake filesystem object that implements the
// http.FileSystem interface.
//
// To use it, call the FileSystem function with one or more http.File objects.
// The File and Dir helper functions create files and directories
// respectively. You do not have to use the File helper, if you'd like
// to write your own mock or stub or even use a real file you can. For
// example:
//
//     testFS := fakehttpfs.FileSystem(
//             fakehttpfs.File("/robots.txt", "User-agent: *\nDisallow: /"),
//             fakehttpfs.Dir("/misc",
//                     fakehttpfs.File("hello.txt", "Hello")
//                     os.Open("/path/to/some/real/file.txt")
//             )
//     );
//
//     file, err := testFS.Open("/robots.txt")
//     file, err := testFS.Open("/misc/file.txt")
package fakehttpfs

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Creates a test fake filesystem containing the files.
func FileSystem(files ...http.File) http.FileSystem {
	return &dir{"", files, 0}
}

//  a test fake file with string contents.
func File(name, contents string) http.File {
	b := []byte(contents)
	return file{name, int64(len(b)), bytes.NewReader(b)}
}

type file struct {
	name string
	size int64
	*bytes.Reader
}

func (f file) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f file) Readdir(int) ([]os.FileInfo, error) {
	return nil, errors.New("Not dir")
}

func (f file) Read(b []byte) (int, error) {
	return f.Reader.Read(b)
}

func (f file) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}

func (f file) Close() error {
	_, err := f.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	return nil
}

func (f file) Name() string {
	return f.name
}

func (f file) Size() int64 {
	return int64(f.Reader.Len())
}

func (f file) Mode() os.FileMode {
	panic("unimplemented")
}

func (f file) ModTime() time.Time {
	panic("unimplemented")
}

func (f file) IsDir() bool {
	return false
}

func (f file) Sys() interface{} {
	return nil
}

// Creates a test fake directory containing the files.
func Dir(name string, files ...http.File) http.File {
	return &dir{name, files, 0}
}

type dir struct {
	name     string
	files    []http.File
	position int
}

func (d *dir) Open(name string) (http.File, error) {
	parts := strings.SplitN(name, "/", 2)
	file, err := d.find(parts[0])
	if len(parts) == 1 {
		return file, err
	}
	if subDir, ok := file.(*dir); ok {
		return subDir.Open(parts[1])
	}
	return nil, os.ErrNotExist
}

func (d *dir) find(name string) (http.File, error) {
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

func (d *dir) Stat() (os.FileInfo, error) {
	return d, nil
}

func (d *dir) Readdir(count int) (r []os.FileInfo, err error) {
	if count == 0 {
		r = make([]os.FileInfo, len(d.files))
		for i, f := range d.files {
			r[i], err = f.Stat()
			if err != nil {
				return r[0 : i-1], err
			}
		}
	} else {
		r = make([]os.FileInfo, count)
		for i := 0; i < count; i++ {
			if d.position > len(d.files)-1 {
				return r[0:i], io.EOF
			}
			r[i], err = d.files[d.position].Stat()
			if err != nil {
				return r[0:i], err
			}
			d.position++
		}
	}
	return r, nil
}

func (d *dir) Read(b []byte) (int, error) {
	return 0, errors.New("Not regular file")
}

func (d *dir) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("Not regular file")
}

func (d *dir) Close() error {
	d.position = 0
	return nil
}

func (d *dir) Name() string {
	return d.name
}

func (d *dir) Size() int64 {
	return 0
}

func (d *dir) Mode() os.FileMode {
	panic("unimplemented")
}

func (d *dir) ModTime() time.Time {
	panic("unimplemented")
}

func (d *dir) IsDir() bool {
	return true
}

func (d *dir) Sys() interface{} {
	return nil
}