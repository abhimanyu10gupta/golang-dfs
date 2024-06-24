package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "my best picture"
	pathKey := CASPathTransformFunc(key)
	expectedFilename := "a009f671f9564204575c776b260c87415ad821da"
	expectedPathName := "a009f/671f9/56420/4575c/776b2/60c87/415ad/821da"

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s found %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.Filename != expectedFilename {
		t.Errorf("have %s found %s", pathKey.Filename, expectedFilename)
	}

}

func TestStore(t *testing.T) {

	s := newStore()
	id := generateID()
	defer tearDown(t, s)

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("foo_%d", i)

		data := []byte("some jpeg bytes")

		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("want %s hav %s", data, b)
		}
		fmt.Println(string(b))
		if err := s.Delete(id, key); err != nil {
			t.Error(err)
		}
		if ok := s.Has(id, key); ok {
			t.Errorf("expected to not have %s key", key)
		}
	}

}
