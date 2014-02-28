// Copyright 2014 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package fakehttpfs

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

var now = time.Now()
var testFS = FileSystem(
	Dir("foo",
		File("bar", "BAR", now),
		Dir("baz",
			Dir("baz",
				Dir("baz",
					File("baz", "BAZ"),
				),
			),
		),
	),
	File("hello", "hello"),
)

func TestFilesystem(t *testing.T) {

	fileTests := []struct {
		path     string
		contents string
		modTime  time.Time
	}{
		{path: "foo/bar", contents: "BAR", modTime: now},
		{path: "foo/baz/baz/baz/baz", contents: "BAZ"},
		{path: "hello", contents: "hello"},
		{path: "/hello", contents: "hello"},
		{path: "./hello", contents: "hello"},
		{path: ".//hello", contents: "hello"},
	}

	for _, test := range fileTests {
		file, err := testFS.Open(test.path)
		if err != nil {
			t.Errorf("expected %s to not error, got %v", test.path, err)
		}
		if file == nil {
			t.Errorf("expected %s to be a file, got nil", test.path)
		} else {
			b := new(bytes.Buffer)
			b.ReadFrom(file)
			file.Close()
			if s := b.String(); s != test.contents {
				t.Errorf("expected %s to contain %q, got %q", test.path, test.contents, s)
			}
			stat, err := file.Stat()
			if err != nil {
				t.Fatalf("expected stat to not error, got %v", err)
			}
			if mt := stat.ModTime(); mt != test.modTime {
				t.Errorf("expected %s modtime to be %v, got %v", test.path, test.modTime, mt)

			}
		}
	}
}

func TestDirs(t *testing.T) {
	dirTests := []string{
		"foo",
		"foo/baz",
		"foo/baz/baz",
	}

	for _, path := range dirTests {
		file, err := testFS.Open(path)
		if err != nil {
			t.Errorf("expected %s to not error, got %v", path, err)
		}
		stat, err := file.Stat()
		if err != nil {
			t.Errorf("expected %s stat to not error, got %v", path, err)
		}
		if !stat.IsDir() {
			t.Errorf("expected %s to be a directory, got %v", path, file)
		}
	}
}

func TestSelf(t *testing.T) {
	selfTests := []string{
		"/",
		"",
		".",
	}

	for _, path := range selfTests {
		file, err := testFS.Open(path)
		if err != nil {
			t.Errorf("expected %s to not error, got %v", path, err)
		}
		if !reflect.DeepEqual(file.(*dir), testFS.(*dir)) {
			t.Errorf("expected %s to be fs, got %v", path, file)
		}
	}
}

func TestErrors(t *testing.T) {
	errTests := []string{
		"oops",
		"foo/oops",
		"foo/baz/baz/oops",
		"foo/baz/baz/baz/baz/oops",
		"hello/oops",
		// we don't do .. cleaning
		"../hello",
		// we don't support trailing slashes
		"hello/",
		"/hello/",
	}

	for _, path := range errTests {
		file, err := testFS.Open(path)
		if file != nil {
			t.Errorf("expecred open(%q) to be nil, got %v", path, file)
		}
		if err == nil {
			t.Errorf("expecred open(%q) to return error, got nil", path)
		}
	}
}

func TestOtherFiles(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "fakehttpfs")
	name := path.Base(tmpFile.Name())
	if err != nil {
		panic(err)
	}
	fs := FileSystem(
		tmpFile,
	)
	file, err := fs.Open(name)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if file != tmpFile {
		t.Errorf("expecred open(%q) to be tmpfile, got %v", name, file)
	}
}

var testDir = Dir("foo",
	File("one", ""),
	File("two", ""),
	File("three", ""),
	File("four", ""),
	File("five", ""),
)

func TestReaddirAll(t *testing.T) {
	fileInfos, err := testDir.Readdir(0)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expected := []string{"one", "two", "three", "four", "five"}
	if names := names(fileInfos); !reflect.DeepEqual(names, expected) {
		t.Errorf("dir.Readdir(0) did not return all elements, expected\n%#v\ngot\n%#v", expected, names)
	}
}

func TestReaddir(t *testing.T) {
	tests := []struct {
		names []string
		err   error
	}{
		{[]string{"one", "two"}, nil},
		{[]string{"three", "four"}, nil},
		{[]string{"five"}, io.EOF},
		{[]string{}, io.EOF},
		{[]string{}, io.EOF},
	}
	for i, test := range tests {
		fileInfos, err := testDir.Readdir(2)
		if err != test.err {
			t.Errorf("got error %v, expected %v", err, test.err)
		}
		if names := names(fileInfos); !reflect.DeepEqual(names, test.names) {
			t.Errorf("iteration %d of dir.Readdir(2) did not return correct elements, expected\n%#v\ngot\n%#v", i+1, test.names, names)
		}
	}
}

func names(infos []os.FileInfo) []string {
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names
}
