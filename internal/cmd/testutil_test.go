package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/salmonumbrella/threads-go/internal/config"
	"github.com/salmonumbrella/threads-go/internal/iocontext"
	"github.com/salmonumbrella/threads-go/internal/secrets"
)

type stubStore struct{}

func (s *stubStore) Set(string, secrets.Credentials) error { return errors.New("not implemented") }
func (s *stubStore) Get(string) (*secrets.Credentials, error) {
	return nil, errors.New("not implemented")
}
func (s *stubStore) Delete(string) error     { return errors.New("not implemented") }
func (s *stubStore) List() ([]string, error) { return nil, errors.New("not implemented") }
func (s *stubStore) Keys() ([]string, error) { return nil, errors.New("not implemented") }

func newTestFactory(t *testing.T) *Factory {
	t.Helper()
	io := &iocontext.IO{
		Out:    &bytes.Buffer{},
		ErrOut: &bytes.Buffer{},
		In:     &bytes.Buffer{},
	}
	cfg := config.Default()
	f, err := NewFactory(context.Background(), FactoryOptions{
		IO:     io,
		Config: cfg,
		Store: func() (secrets.Store, error) {
			return &stubStore{}, nil
		},
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	return f
}
