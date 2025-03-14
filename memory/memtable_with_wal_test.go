package memory

import (
	"go-lsm-workshop/kv"
	"go-lsm-workshop/log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemtableWithWALWithASingleKey(t *testing.T) {
	directoryPath := "."
	walDirectoryPath := filepath.Join(directoryPath, "wal")
	assert.Nil(t, os.MkdirAll(walDirectoryPath, os.ModePerm))

	defer func() {
		_ = os.RemoveAll(walDirectoryPath)
	}()

	memTable := NewMemtable(1, testMemtableSize, log.NewWALPath(directoryPath))
	_ = memTable.Set(kv.NewStringKeyWithTimestamp("consensus", 5), kv.NewStringValue("raft"))

	value, ok := memTable.Get(kv.NewStringKeyWithTimestamp("consensus", 5))
	assert.True(t, ok)
	assert.Equal(t, kv.NewStringValue("raft"), value)
}

func TestMemtableWithWALWithMultipleKeys(t *testing.T) {
	directoryPath := "."
	walDirectoryPath := filepath.Join(directoryPath, "wal")
	assert.Nil(t, os.MkdirAll(walDirectoryPath, os.ModePerm))

	defer func() {
		_ = os.RemoveAll(walDirectoryPath)
	}()

	memTable := NewMemtable(2, testMemtableSize, log.NewWALPath(directoryPath))
	_ = memTable.Set(kv.NewStringKeyWithTimestamp("consensus", 5), kv.NewStringValue("raft"))
	_ = memTable.Set(kv.NewStringKeyWithTimestamp("storage", 6), kv.NewStringValue("NVMe"))

	value, ok := memTable.Get(kv.NewStringKeyWithTimestamp("consensus", 6))
	assert.True(t, ok)
	assert.Equal(t, kv.NewStringValue("raft"), value)

	value, ok = memTable.Get(kv.NewStringKeyWithTimestamp("storage", 6))
	assert.True(t, ok)
	assert.Equal(t, kv.NewStringValue("NVMe"), value)
}

func TestMemtableRecoveryFromWAL(t *testing.T) {
	directoryPath := "."
	walDirectoryPath := filepath.Join(directoryPath, "wal")
	assert.Nil(t, os.MkdirAll(walDirectoryPath, os.ModePerm))

	defer func() {
		_ = os.RemoveAll(walDirectoryPath)
	}()

	memTable := NewMemtable(3, testMemtableSize, log.NewWALPath(directoryPath))
	_ = memTable.Set(kv.NewStringKeyWithTimestamp("consensus", 5), kv.NewStringValue("raft"))
	_ = memTable.Set(kv.NewStringKeyWithTimestamp("storage", 6), kv.NewStringValue("NVMe"))

	memTable.wal.Close()

	recoveredMemTable, maxTimestamp, err := RecoverFromWAL(3, testMemtableSize, walDirectoryPath)
	assert.Nil(t, err)

	value, ok := recoveredMemTable.Get(kv.NewStringKeyWithTimestamp("consensus", 5))
	assert.True(t, ok)
	assert.Equal(t, kv.NewStringValue("raft"), value)

	value, ok = recoveredMemTable.Get(kv.NewStringKeyWithTimestamp("storage", 6))
	assert.True(t, ok)
	assert.Equal(t, kv.NewStringValue("NVMe"), value)

	assert.Equal(t, uint64(6), maxTimestamp)
}
