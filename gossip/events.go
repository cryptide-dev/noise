package gossip

import "github.com/cryptide-dev/noise"

// Events comprise of callbacks that may be hooked against by a user to handle inbound gossip messages/events that
// occur throughout the lifecycle of this gossip protocol.
type Events struct {
	// OnGossipReceived is called whenever new gossip is received from the network. An error may be return to
	// disconnect the sender sending you data; indicating that the gossip received is invalid.
	OnGossipReceived func(sender noise.ID, data []byte) error
}
