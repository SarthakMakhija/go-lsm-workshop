//go:build test

package go_lsm_workshop

import "go-lsm-workshop/state"

// StorageState returns the StorageState, it is only for testing.
func (db *Db) StorageState() *state.StorageState {
	return db.storageState
}
