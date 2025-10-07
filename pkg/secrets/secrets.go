package secrets

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"sync"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type S struct {
	client     *secretmanager.Client
	crc32Table *crc32.Table
	ttl        time.Duration

	mu    sync.Mutex
	cache map[string]entry
}

type entry struct {
	value  []byte
	expiry time.Time
}

func New(ctx context.Context, ttl time.Duration) (*S, error) {
	if ttl < 0 {
		return nil, errors.New("ttl must be non-negative")
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	return &S{
		client:     client,
		crc32Table: crc32.MakeTable(crc32.Castagnoli),
		ttl:        ttl,
		cache:      make(map[string]entry),
	}, nil
}

func (s *S) Close() error {
	return s.client.Close()
}

func (s *S) Read(ctx context.Context, verisonName string) ([]byte, error) {
	if value, ok := s.readCache(verisonName); ok {
		return value, nil
	}

	result, err := s.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: verisonName})
	if err != nil {
		return nil, fmt.Errorf("reading secret: %w", err)
	}

	if result.Payload.DataCrc32C == nil {
		return nil, errors.New("secret is missing checksum")
	}
	if checksum := int64(crc32.Checksum(result.Payload.Data, s.crc32Table)); checksum != *result.Payload.DataCrc32C {
		return nil, errors.New("checksum mismatch")
	}

	s.writeCache(verisonName, result.Payload.Data)
	return result.Payload.Data, nil
}

func (s *S) readCache(name string) ([]byte, bool) {
	if s.ttl == 0 {
		return nil, false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.cache[name]
	if !ok {
		return nil, false
	}
	if e.expiry.Before(time.Now()) {
		delete(s.cache, name)
		return nil, false
	}
	return e.value, true
}

func (s *S) writeCache(name string, value []byte) {
	if s.ttl == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[name] = entry{
		value:  value,
		expiry: time.Now().Add(s.ttl),
	}
}
