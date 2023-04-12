package gossip

import "github.com/VictoriaMetrics/fastcache"

// Option is a functional option that may be configured when instantiating a new instance of this gossip protocol.
type Option func(protocol *Protocol)

// WithEvents registers a batch of callbacks onto a single gossip protocol instance.
func WithEvents(events Events) Option {
	return func(protocol *Protocol) {
		protocol.events = events
	}
}

// WithCacheSize sets new cache of the specified size for nodes, that have been already seen
func WithCacheSize(maxBytes int) Option {
	return func(protocol *Protocol) {
		protocol.seen = fastcache.New(maxBytes)
	}
}
