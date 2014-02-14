package mockfs

import (
	"bytes"
	"io/ioutil"
	"path"
	"reflect"
	"testing"
)

var testFS = FileSystem(
	Dir("foo",
		File("bar", "BAR"),
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
	}{
		{path: "foo/bar", contents: "BAR"},
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
		if !reflect.DeepEqual(file.(mockDir), testFS.(mockDir)) {
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
	tmpFile, err := ioutil.TempFile("", "mockfs")
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
