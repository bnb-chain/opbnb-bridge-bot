package core

import "time"

// txNonceAndTimestamp represents the nonce and sending timestamp of transaction.
type txNonceAndTimestamp struct {
	nonce     uint64
	timestamp time.Time
}

// PendingTxnCheck records the nonce and sending timestamp of transactions, in order to prevent re-processing
// the same withdrawals.
type PendingTxnCheck struct {
	inner map[uint]txNonceAndTimestamp // #{withdrawalId=>txNonceAndTimestamp}
}

// NewPendingTxsManager creates a new PendingTxnCheck
func NewPendingTxsManager() *PendingTxnCheck {
	return &PendingTxnCheck{inner: make(map[uint]txNonceAndTimestamp)}
}

// IsPendingTxn checks whether there is pending transaction for the specific event id.
func (c *PendingTxnCheck) IsPendingTxn(id uint) bool {
	_, ok := c.inner[id]
	return ok
}

// AddPendingTxn adds a pending item.
func (c *PendingTxnCheck) AddPendingTxn(id uint, nonce uint64) {
	c.inner[id] = txNonceAndTimestamp{nonce: nonce, timestamp: time.Now()}
}

// Prune removes the transactions with chainNonce and the current timestamp.
func (c *PendingTxnCheck) Prune(chainNonce uint64) {
	const DurationSafelyPrunePendingTxs = 5 * time.Minute
	timeSafelyPrunePendingTxs := time.Now().Add(-1 * DurationSafelyPrunePendingTxs)

	for id, nt := range c.inner {
		if nt.nonce <= chainNonce && nt.timestamp.Before(timeSafelyPrunePendingTxs) {
			delete(c.inner, id)
		}
	}
}
