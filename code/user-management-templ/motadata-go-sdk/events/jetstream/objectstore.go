package jetstream

import (
	"context"
	"io"

	"github.com/nats-io/nats.go/jetstream"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// ObjectStore wraps a JetStream Object Store bucket.
type ObjectStore struct {
	os     jetstream.ObjectStore
	bucket string
}

// CreateObjectStore creates a new Object Store bucket.
func CreateObjectStore(ctx context.Context, js jetstream.JetStream, cfg jetstream.ObjectStoreConfig) (*ObjectStore, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	os, err := js.CreateObjectStore(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &ObjectStore{os: os, bucket: cfg.Bucket}, nil
}

// GetObjectStore retrieves an existing Object Store bucket.
func GetObjectStore(ctx context.Context, js jetstream.JetStream, bucket string) (*ObjectStore, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	os, err := js.ObjectStore(ctx, bucket)
	if err != nil {
		return nil, err
	}

	return &ObjectStore{os: os, bucket: bucket}, nil
}

// DeleteObjectStore deletes an Object Store bucket.
func DeleteObjectStore(ctx context.Context, js jetstream.JetStream, bucket string) error {
	if js == nil {
		return utils.ErrJetStreamNotEnabled
	}

	return js.DeleteObjectStore(ctx, bucket)
}

// ObjectStoreNames returns the names of all Object Store buckets.
func ObjectStoreNames(ctx context.Context, js jetstream.JetStream) ([]string, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	lister := js.ObjectStoreNames(ctx)
	var names []string
	for name := range lister.Name() {
		names = append(names, name)
	}
	if err := lister.Error(); err != nil {
		return names, err
	}

	return names, nil
}

// Get retrieves an object by name.
func (s *ObjectStore) Get(ctx context.Context, name string) (jetstream.ObjectResult, error) {
	return s.os.Get(ctx, name)
}

// GetInfo retrieves metadata for an object.
func (s *ObjectStore) GetInfo(ctx context.Context, name string, opts ...jetstream.GetObjectInfoOpt) (*jetstream.ObjectInfo, error) {
	return s.os.GetInfo(ctx, name, opts...)
}

// Put stores an object with the given metadata and data reader.
func (s *ObjectStore) Put(ctx context.Context, meta jetstream.ObjectMeta, reader io.Reader) (*jetstream.ObjectInfo, error) {
	return s.os.Put(ctx, meta, reader)
}

// PutBytes stores an object as raw bytes.
func (s *ObjectStore) PutBytes(ctx context.Context, name string, data []byte) (*jetstream.ObjectInfo, error) {
	return s.os.PutBytes(ctx, name, data)
}

// Delete removes an object by name.
func (s *ObjectStore) Delete(ctx context.Context, name string) error {
	return s.os.Delete(ctx, name)
}

// List returns information about all objects in the store.
func (s *ObjectStore) List(ctx context.Context, opts ...jetstream.ListObjectsOpt) ([]*jetstream.ObjectInfo, error) {
	return s.os.List(ctx, opts...)
}

// Watch watches for changes to objects in the store.
func (s *ObjectStore) Watch(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.ObjectWatcher, error) {
	return s.os.Watch(ctx, opts...)
}

// Seal seals the Object Store, preventing further modifications.
func (s *ObjectStore) Seal(ctx context.Context) error {
	return s.os.Seal(ctx)
}

// Status returns the status of the Object Store.
func (s *ObjectStore) Status(ctx context.Context) (jetstream.ObjectStoreStatus, error) {
	return s.os.Status(ctx)
}

// BucketName returns the name of the bucket.
func (s *ObjectStore) BucketName() string {
	return s.bucket
}
