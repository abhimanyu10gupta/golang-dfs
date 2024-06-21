package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "my best picture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "a009f671f9564204575c776b260c87415ad821da"
	expectedPathName := "a009f/671f9/56420/4575c/776b2/60c87/415ad/821da"

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s found %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.Filename != expectedPathName {
		t.Errorf("have %s found %s", pathKey.Filename, expectedOriginalKey)
	}

}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "myspecialpicture"

	data := []byte("some jpeg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	key := "myspecialpicture"

	data := []byte("some jpeg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("want %s hav %s", data, b)
	}
	fmt.Println(string(b))
	s.Delete(key)
}
