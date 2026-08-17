package jetstream

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// KVStore wraps a JetStream KeyValue bucket with a simplified API for
// key-value operations.
//
// JetStream KeyValue stores are backed by a dedicated stream with
// subject-per-key mapping. Each Put or Update creates a new message
// (revision) in the underlying stream, enabling full history tracking
// and change notification via Watch.
//
// Optimistic concurrency is supported through revisions: every write
// returns a monotonically increasing uint64 revision number. The Update
// method accepts a lastRevision parameter and only succeeds if the key's
// current revision matches, preventing lost-update problems in concurrent
// environments. Use Create for insert-if-absent semantics (fails if the
// key already exists) and Put for unconditional upsert.
type KVStore struct {
	// kv is the underlying JetStream KeyValue handle.
	kv jetstream.KeyValue
	// bucket is the name of this KV bucket (maps to a stream name internally).
	bucket string
}

// CreateKVStore creates a new KeyValue store bucket with the given
// configuration. The KeyValueConfig controls the bucket name, TTL for
// automatic key expiration, history depth (number of revisions kept per
// key), max value size, storage backend, and replication factor.
// Returns ErrJetStreamNotEnabled when js is nil.
func CreateKVStore(ctx context.Context, js jetstream.JetStream, cfg jetstream.KeyValueConfig) (*KVStore, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	kv, err := js.CreateKeyValue(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &KVStore{kv: kv, bucket: cfg.Bucket}, nil
}

// GetKVStore retrieves an existing KeyValue store bucket by name.
// Returns an error if the bucket does not exist on the server.
func GetKVStore(ctx context.Context, js jetstream.JetStream, bucket string) (*KVStore, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return nil, err
	}

	return &KVStore{kv: kv, bucket: bucket}, nil
}

// DeleteKVStore deletes a KeyValue store bucket and all of its keys,
// revisions, and underlying stream data. This operation is irreversible.
func DeleteKVStore(ctx context.Context, js jetstream.JetStream, bucket string) error {
	if js == nil {
		return utils.ErrJetStreamNotEnabled
	}

	return js.DeleteKeyValue(ctx, bucket)
}

// KVStoreNames returns the names of all KeyValue store buckets visible to
// the current account. It drains the server-side lister into a slice for
// convenience; for very large numbers of buckets, consider using the
// underlying JetStream lister directly to process names incrementally.
func KVStoreNames(ctx context.Context, js jetstream.JetStream) ([]string, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	// drainLister - collect all bucket names into a single slice.
	lister := js.KeyValueStoreNames(ctx)
	var names []string
	for name := range lister.Name() {
		names = append(names, name)
	}
	if err := lister.Error(); err != nil {
		return names, err
	}

	return names, nil
}

// Get retrieves the latest revision of the value for a key. The returned
// KeyValueEntry includes the value bytes, the revision number (useful for
// subsequent Update calls with optimistic concurrency), the operation type,
// and a timestamp.
func (s *KVStore) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	return s.kv.Get(ctx, key)
}

// Put stores a value for a key unconditionally, creating the key if it
// does not exist or overwriting the current value if it does. Returns the
// new revision number assigned to this write.
func (s *KVStore) Put(ctx context.Context, key string, value []byte) (uint64, error) {
	return s.kv.Put(ctx, key, value)
}

// Create stores a value only if the key does not already exist (insert-if-
// absent semantics). Returns the initial revision number on success, or an
// error if the key is already present. This is useful for initializing
// shared state without risking an overwrite from a concurrent writer.
func (s *KVStore) Create(ctx context.Context, key string, value []byte) (uint64, error) {
	return s.kv.Create(ctx, key, value)
}

// Update stores a value only if the key's current revision matches
// lastRevision (compare-and-swap / optimistic concurrency control). This
// prevents lost updates when multiple writers operate on the same key
// concurrently. The typical pattern is: Get the current entry, read its
// Revision(), compute the new value, then call Update with that revision.
// If another writer modified the key in between, Update returns an error
// and the caller can retry.
func (s *KVStore) Update(ctx context.Context, key string, value []byte, lastRevision uint64) (uint64, error) {
	return s.kv.Update(ctx, key, value, lastRevision)
}

// Delete soft-deletes a key by placing a delete marker in the revision
// history. The key will no longer appear in Get or Keys results, but its
// history is retained (up to the configured history depth) and can still
// be observed via Watch. Use Purge to hard-delete all revisions.
func (s *KVStore) Delete(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	return s.kv.Delete(ctx, key, opts...)
}

// Keys returns all non-deleted keys currently in the bucket. Internally
// this replays the stream to build the key set, so it may be slow for
// very large buckets.
func (s *KVStore) Keys(ctx context.Context, opts ...jetstream.WatchOpt) ([]string, error) {
	return s.kv.Keys(ctx, opts...)
}

// History returns the full revision history of values for a key, ordered
// from oldest to newest. The number of retained revisions is controlled
// by the bucket's History setting in KeyValueConfig. Each entry includes
// the value, revision number, operation type, and timestamp.
func (s *KVStore) History(ctx context.Context, key string, opts ...jetstream.WatchOpt) ([]jetstream.KeyValueEntry, error) {
	return s.kv.History(ctx, key, opts...)
}

// Watch watches for changes to keys matching the given pattern (e.g.,
// "config.*" or "users.>") and delivers change notifications through the
// returned KeyWatcher's Updates() channel. The watcher first replays all
// existing matching entries (initial values), then switches to live
// updates. Use jetstream.IgnoreDeletes() or jetstream.UpdatesOnly() watch
// options to filter the notification stream.
func (s *KVStore) Watch(ctx context.Context, keys string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return s.kv.Watch(ctx, keys, opts...)
}

// WatchAll watches for changes to all keys in the bucket. It behaves
// identically to Watch with a wildcard pattern: it replays all current
// values first, then delivers live updates for every key mutation.
func (s *KVStore) WatchAll(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return s.kv.WatchAll(ctx, opts...)
}

// Purge hard-deletes all revisions of a key from the underlying stream,
// freeing storage. Unlike Delete (which places a soft-delete marker),
// Purge permanently removes the data.
func (s *KVStore) Purge(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	return s.kv.Purge(ctx, key, opts...)
}

// PurgeDeletes removes all soft-delete markers from the bucket, reclaiming
// storage occupied by tombstones. This is a maintenance operation typically
// run periodically to keep the underlying stream compact.
func (s *KVStore) PurgeDeletes(ctx context.Context, opts ...jetstream.KVPurgeOpt) error {
	return s.kv.PurgeDeletes(ctx, opts...)
}

// Status returns the status of the KeyValue bucket, including the number
// of keys, total bytes used, the underlying stream info, and the bucket's
// configuration.
func (s *KVStore) Status(ctx context.Context) (jetstream.KeyValueStatus, error) {
	return s.kv.Status(ctx)
}

// BucketName returns the name of the bucket as specified during creation.
func (s *KVStore) BucketName() string {
	return s.bucket
}
