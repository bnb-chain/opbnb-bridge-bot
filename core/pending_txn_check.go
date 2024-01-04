package core

type PendingTxnCheck struct {
	inner map[uint]uint64 // #{withdrawalId=>nonce}
}

// NewPendingTxsManager creates a new PendingTxnCheck
func NewPendingTxsManager() *PendingTxnCheck {
	return &PendingTxnCheck{inner: make(map[uint]uint64)}
}

// IsPendingTxn checks whether there is pending transaction for the specific event id.
func (c *PendingTxnCheck) IsPendingTxn(id uint) bool {
	_, ok := c.inner[id]
	return ok
}

// AddPendingTxn adds a pending item.
func (c *PendingTxnCheck) AddPendingTxn(id uint, nonce uint64) {
	c.inner[id] = nonce
}

// Prune removes the transactions with staled nonce.
func (c *PendingTxnCheck) Prune(chainNonce uint64) {
	for id, nonce := range c.inner {
		if nonce <= chainNonce {
			delete(c.inner, id)
		}
	}
}
