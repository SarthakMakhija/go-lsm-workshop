package kv

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

const TimestampSize = int(unsafe.Sizeof(uint64(0)))

// Key represents a versioned key. It contains the original key (of the raw key)
// along with the timestamp.
// A small comment on timestamp:
// When an instance of Key is stored in the system, the timestamp is the commit-timestamp,
// generated by txn.Oracle.
// When the user performs a get (/read) operation, the system creates an instance of Key
// with the user provided key and the begin-timestamp of the transaction.
type Key struct {
	key       []byte
	timestamp uint64
}

var EmptyKey = Key{key: nil}

func DecodeFrom(buffer []byte) Key {
	if len(buffer) < TimestampSize {
		panic("buffer too small to decode the key from")
	}
	length := len(buffer)
	return Key{
		key:       buffer[:length-TimestampSize],
		timestamp: binary.LittleEndian.Uint64(buffer[length-TimestampSize:]),
	}
}

// NewKey creates a new instance of the Key.
func NewKey(key []byte, timestamp uint64) Key {
	return Key{
		key:       key,
		timestamp: timestamp,
	}
}

// IsLessThanOrEqualTo returns true if the Key is less than or equal to the other Key.
func (key Key) IsLessThanOrEqualTo(other LessOrEqual) bool {
	otherKey := other.(Key)
	comparison := bytes.Compare(key.key, otherKey.key)
	if comparison > 0 {
		return false
	}
	if comparison < 0 {
		return true
	}
	//comparison == 0
	return key.timestamp <= otherKey.timestamp
}

// IsEqualTo returns true if the Key is equal to the other Key.
func (key Key) IsEqualTo(other Key) bool {
	return bytes.Equal(key.key, other.key) && key.timestamp == other.timestamp
}

// CompareKeysWithDescendingTimestamp compares the two keys.
// It compares the raw keys ([]byte).
// If the comparison result is not zero, it is returned.
// Else, the timestamps of the two keys are compared:
// 1) It returns 0, if the timestamps of the two keys are same.
// 2) It returns -1, if the timestamp of the key is greater than the timestamp of the other key.
// 3) It returns 1, if the timestamp of the key is less than the timestamp of the other key.
// Timestamp plays an important role in ordering of keys. Consider a key "consensus" with timestamps 15 and 13 in memtable,
// and user wants to perform a scan between "consensus" to "decimal" with timestamp as 16. This means we would like to return
// all the keys falling between "consensus" and "decimal" such that: timestamp of the keys in system <= 16.
// However, "consensus" is present with 15 and 13 timestamps. We would only return "consensus" with timestamp 15. If the key
// "consensus" with timestamp 15 is placed before the key "consensus" with timestamp 13, range iteration becomes easier because
// the first key will always have the latest timestamp.
// Instictively, the key "consensus" with timestamp 15 is greater than the same key with timestamp 13, but we would like to place
// the key "consensus" with timestamp 15 before the same key with timestamp 13 in Skiplist.
// Hence, "consensus_15".CompareKeysWithDescendingTimestamp("consensus_15") returns -1.
// Look at the test: TestKeyComparisonLessThanBasedOnTimestamp to understand key comparison.
func (key Key) CompareKeysWithDescendingTimestamp(other Key) int {
	comparison := bytes.Compare(key.key, other.key)
	if comparison != 0 {
		return comparison
	}
	if key.timestamp == other.timestamp {
		return 0
	}
	if key.timestamp > other.timestamp {
		return -1
	}
	return 1
}

// CompareKeys compares the user provided key with timestamp and the instance of the Key existing in the system.
// It is mainly called from external.SkipList.
func CompareKeys(userKey, systemKey Key) int {
	return userKey.CompareKeysWithDescendingTimestamp(systemKey)
}

// IsRawKeyEqualTo returns true if the raw key two keys is the same.
func (key Key) IsRawKeyEqualTo(other Key) bool {
	return bytes.Equal(key.key, other.key)
}

// IsRawKeyGreaterThan returns true if the raw key of key is greater than the raw key of the other.
func (key Key) IsRawKeyGreaterThan(other Key) bool {
	return bytes.Compare(key.key, other.key) > 0
}

// IsRawKeyLesserThan returns true if the raw key of key is lesser than the raw key of the other.
func (key Key) IsRawKeyLesserThan(other Key) bool {
	return bytes.Compare(key.key, other.key) < 0
}

// IsRawKeyEmpty returns true if the raw key is empty.
func (key Key) IsRawKeyEmpty() bool {
	return key.RawSizeInBytes() == 0
}

// EncodedBytes returns the encoded format of the Key.
// The encoded format of Key includes:
//
// | Raw Key| timestamp |
func (key Key) EncodedBytes() []byte {
	if key.IsRawKeyEmpty() {
		return nil
	}
	buffer := make([]byte, key.EncodedSizeInBytes())

	numberOfBytesWritten := copy(buffer, key.key)
	binary.LittleEndian.PutUint64(buffer[numberOfBytesWritten:], key.timestamp)

	return buffer
}

// RawBytes returns the raw key.
func (key Key) RawBytes() []byte {
	return key.key
}

// RawString returns the string representation of raw key.
func (key Key) RawString() string {
	return string(key.RawBytes())
}

// EncodedSizeInBytes returns the length of the encoded key.
func (key Key) EncodedSizeInBytes() int {
	if key.IsRawKeyEmpty() {
		return 0
	}
	return len(key.key) + TimestampSize
}

// RawSizeInBytes returns the size of the raw key.
func (key Key) RawSizeInBytes() int {
	return len(key.RawBytes())
}

// Timestamp returns the timestamp.
func (key Key) Timestamp() uint64 {
	return key.timestamp
}

// RawKey represents the raw key (provided by the user).
type RawKey []byte

// IsLessThanOrEqualTo returns true if RawKey is less than or equal to other key.
func (key RawKey) IsLessThanOrEqualTo(other LessOrEqual) bool {
	otherKey := other.(RawKey)
	return bytes.Compare(key, otherKey) <= 0
}
