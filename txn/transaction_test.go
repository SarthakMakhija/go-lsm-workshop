package txn

import (
	"go-lsm-workshop/kv"
	"go-lsm-workshop/state"
	"go-lsm-workshop/table"
	"go-lsm-workshop/test_utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadonlyTransactionWithEmptyState(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	transaction := NewReadonlyTransaction(oracle, storageState)
	_, ok := transaction.Get([]byte("paxos"))

	assert.False(t, ok)
}

func TestReadonlyTransactionWithAnExistingKey(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	commitTimestamp := uint64(5)
	oracle.nextTimestamp = commitTimestamp + 1

	batch := kv.NewBatch()
	_ = batch.Put([]byte("consensus"), []byte("raft"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	transaction := NewReadonlyTransaction(oracle, storageState)
	value, ok := transaction.Get([]byte("consensus"))

	assert.True(t, ok)
	assert.Equal(t, "raft", value.String())
}

func TestReadonlyTransactionWithAnExistingKeyButWithATimestampHigherThanCommitTimestamp(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	//simulate a readonly transaction starting first
	oracle.nextTimestamp = uint64(5)
	oracle.commitTimestampMark.Finish(uint64(4))
	transaction := NewReadonlyTransaction(oracle, storageState)

	commitTimestamp := uint64(6)
	batch := kv.NewBatch()
	_ = batch.Put([]byte("raft"), []byte("consensus algorithm"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	_, ok := transaction.Get([]byte("raft"))

	assert.False(t, ok)
}

func TestReadonlyTransactionWithScan(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	commitTimestamp := uint64(5)
	oracle.nextTimestamp = commitTimestamp + 1

	batch := kv.NewBatch()
	_ = batch.Put([]byte("consensus"), []byte("raft"))
	_ = batch.Put([]byte("storage"), []byte("NVMe"))
	_ = batch.Put([]byte("kv"), []byte("distributed"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	transaction := NewReadonlyTransaction(oracle, storageState)
	iterator, _ := transaction.Scan(kv.NewInclusiveKeyRange(kv.RawKey("draft"), kv.RawKey("quadrant")))

	assert.Equal(t, "kv", iterator.Key().RawString())
	assert.Equal(t, "distributed", iterator.Value().String())

	_ = iterator.Next()

	assert.False(t, iterator.IsValid())
}

func TestReadonlyTransactionWithScanHavingSameKeyWithMultipleTimestamps(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	batch := kv.NewBatch()
	_ = batch.Put([]byte("consensus"), []byte("unknown"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, 4)))

	commitTimestamp := uint64(5)
	oracle.nextTimestamp = commitTimestamp + 1

	batch = kv.NewBatch()
	_ = batch.Put([]byte("consensus"), []byte("VSR"))
	_ = batch.Put([]byte("storage"), []byte("NVMe"))
	_ = batch.Put([]byte("kv"), []byte("distributed"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	transaction := NewReadonlyTransaction(oracle, storageState)
	iterator, _ := transaction.Scan(kv.NewInclusiveKeyRange(kv.RawKey("bolt"), kv.RawKey("quadrant")))

	assert.Equal(t, "consensus", iterator.Key().RawString())
	assert.Equal(t, "VSR", iterator.Value().String())

	_ = iterator.Next()

	assert.Equal(t, "kv", iterator.Key().RawString())
	assert.Equal(t, "distributed", iterator.Value().String())

	_ = iterator.Next()

	assert.False(t, iterator.IsValid())
}

func TestAttemptsToCommitAnEmptyReadwriteTransaction(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	oracle.commitTimestampMark.Finish(2)
	transaction := NewReadwriteTransaction(oracle, storageState)

	_, err := transaction.Commit()

	assert.Error(t, err)
	assert.Equal(t, EmptyTransactionErr, err)
}

func TestGetsAnExistingKeyInAReadwriteTransaction(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	transaction := NewReadwriteTransaction(oracle, storageState)
	_ = transaction.Set([]byte("HDD"), []byte("Hard disk"))
	future, _ := transaction.Commit()
	future.Wait()

	anotherTransaction := NewReadwriteTransaction(oracle, storageState)
	_ = anotherTransaction.Set([]byte("SSD"), []byte("Solid state drive"))
	future, _ = anotherTransaction.Commit()
	future.Wait()

	readonlyTransaction := NewReadonlyTransaction(oracle, storageState)

	value, ok := readonlyTransaction.Get([]byte("HDD"))
	assert.Equal(t, true, ok)
	assert.Equal(t, "Hard disk", value.String())

	value, ok = readonlyTransaction.Get([]byte("SSD"))
	assert.Equal(t, true, ok)
	assert.Equal(t, "Solid state drive", value.String())

	_, ok = readonlyTransaction.Get([]byte("non-existing"))
	assert.Equal(t, false, ok)
}

func TestGetsTheValueFromAKeyInAReadwriteTransactionFromBatch(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	transaction := NewReadwriteTransaction(oracle, storageState)
	_ = transaction.Set([]byte("HDD"), []byte("Hard disk"))

	value, ok := transaction.Get([]byte("HDD"))
	assert.Equal(t, true, ok)
	assert.Equal(t, "Hard disk", value.String())

	future, _ := transaction.Commit()
	future.Wait()
}

func TestTracksReadsInAReadwriteTransactionWithGet(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	transaction := NewReadwriteTransaction(oracle, storageState)
	_ = transaction.Set([]byte("HDD"), []byte("Hard disk"))
	transaction.Get([]byte("SSD"))

	future, _ := transaction.Commit()
	future.Wait()

	assert.Equal(t, 1, len(transaction.reads))
	assert.Equal(t, kv.RawKey("SSD"), transaction.reads[0])
}

func TestReadwriteTransactionWithScanHavingMultipleTimestampsOfSameKey(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	batch := kv.NewBatch()
	_ = batch.Put([]byte("consensus"), []byte("unknown"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, 4)))

	commitTimestamp := uint64(5)
	oracle.nextTimestamp = commitTimestamp + 1

	batch = kv.NewBatch()
	_ = batch.Put([]byte("consensus"), []byte("VSR"))
	_ = batch.Put([]byte("storage"), []byte("NVMe"))
	_ = batch.Put([]byte("kv"), []byte("distributed"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	transaction := NewReadwriteTransaction(oracle, storageState)
	iterator, _ := transaction.Scan(kv.NewInclusiveKeyRange(kv.RawKey("bolt"), kv.RawKey("quadrant")))

	assert.Equal(t, "consensus", iterator.Key().RawString())
	assert.Equal(t, "VSR", iterator.Value().String())

	_ = iterator.Next()

	assert.Equal(t, "kv", iterator.Key().RawString())
	assert.Equal(t, "distributed", iterator.Value().String())

	_ = iterator.Next()

	assert.False(t, iterator.IsValid())
}

func TestReadwriteTransactionWithScanHavingDeletedKey(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	commitTimestamp := uint64(5)
	oracle.nextTimestamp = commitTimestamp + 1

	batch := kv.NewBatch()
	batch.Delete([]byte("quadrant"))
	_ = batch.Put([]byte("consensus"), []byte("VSR"))
	_ = batch.Put([]byte("storage"), []byte("NVMe"))
	_ = batch.Put([]byte("kv"), []byte("distributed"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	transaction := NewReadwriteTransaction(oracle, storageState)
	iterator, _ := transaction.Scan(kv.NewInclusiveKeyRange(kv.RawKey("bolt"), kv.RawKey("rocks")))

	assert.Equal(t, "consensus", iterator.Key().RawString())
	assert.Equal(t, "VSR", iterator.Value().String())

	_ = iterator.Next()

	assert.Equal(t, "kv", iterator.Key().RawString())
	assert.Equal(t, "distributed", iterator.Value().String())

	_ = iterator.Next()
	assert.False(t, iterator.IsValid())
}

func TestTracksReadsInAReadwriteTransactionWithScan(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	commitTimestamp := uint64(5)
	oracle.nextTimestamp = commitTimestamp + 1

	batch := kv.NewBatch()
	batch.Delete([]byte("quadrant"))
	_ = batch.Put([]byte("consensus"), []byte("VSR"))
	_ = batch.Put([]byte("storage"), []byte("NVMe"))
	_ = batch.Put([]byte("kv"), []byte("distributed"))
	assert.Nil(t, storageState.Set(kv.NewTimestampedBatchFrom(*batch, commitTimestamp)))
	oracle.commitTimestampMark.Finish(commitTimestamp)

	transaction := NewReadwriteTransaction(oracle, storageState)
	_ = transaction.Set([]byte("hdd"), []byte("Hard disk"))

	iterator, _ := transaction.Scan(kv.NewInclusiveKeyRange(kv.RawKey("bolt"), kv.RawKey("tiger-beetle")))

	assert.Equal(t, "consensus", iterator.Key().RawString())
	assert.Equal(t, "VSR", iterator.Value().String())

	_ = iterator.Next()

	assert.Equal(t, "hdd", iterator.Key().RawString())
	assert.Equal(t, "Hard disk", iterator.Value().String())

	_ = iterator.Next()

	assert.Equal(t, "kv", iterator.Key().RawString())
	assert.Equal(t, "distributed", iterator.Value().String())

	_ = iterator.Next()

	assert.Equal(t, "storage", iterator.Key().RawString())
	assert.Equal(t, "NVMe", iterator.Value().String())

	_ = iterator.Next()
	assert.False(t, iterator.IsValid())

	allTrackedReads := transaction.reads
	assert.Equal(t, 4, len(allTrackedReads))

	assert.Equal(t, "consensus", string(allTrackedReads[0]))
	assert.Equal(t, "hdd", string(allTrackedReads[1]))
	assert.Equal(t, "kv", string(allTrackedReads[2]))
	assert.Equal(t, "storage", string(allTrackedReads[3]))
}

func TestReferencesToSSTableInTransactionGet(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	ssTableBuilder := table.NewSSTableBuilder(4096)
	ssTableBuilder.Add(kv.NewStringKeyWithTimestamp("consensus", 0), kv.NewStringValue("paxos"))
	ssTableBuilder.Add(kv.NewStringKeyWithTimestamp("distributed", 0), kv.NewStringValue("TiKV"))
	ssTableBuilder.Add(kv.NewStringKeyWithTimestamp("etcd", 0), kv.NewStringValue("bbolt"))

	ssTable, err := ssTableBuilder.Build(1, rootPath)
	assert.Nil(t, err)

	storageState.SetSSTableAtLevel(ssTable, 0)

	readonlyTransaction := NewReadonlyTransaction(oracle, storageState)
	value, ok := readonlyTransaction.Get([]byte("consensus"))

	assert.True(t, ok)
	assert.Equal(t, "paxos", value.String())
	assert.Equal(t, int64(0), ssTable.TotalReferences())
}

func TestReferencesToSSTableInTransactionScan(t *testing.T) {
	rootPath := test_utility.SetupADirectoryWithTestName(t)
	storageState, _ := state.NewStorageState(rootPath)
	oracle := NewOracle(NewExecutor(storageState))

	defer func() {
		test_utility.CleanupDirectoryWithTestName(t)
		storageState.Close()
		oracle.Close()
	}()

	ssTableBuilder := table.NewSSTableBuilder(4096)
	ssTableBuilder.Add(kv.NewStringKeyWithTimestamp("consensus", 0), kv.NewStringValue("paxos"))
	ssTableBuilder.Add(kv.NewStringKeyWithTimestamp("distributed", 0), kv.NewStringValue("TiKV"))
	ssTableBuilder.Add(kv.NewStringKeyWithTimestamp("etcd", 0), kv.NewStringValue("bbolt"))

	ssTable, err := ssTableBuilder.Build(1, rootPath)
	assert.Nil(t, err)

	storageState.SetSSTableAtLevel(ssTable, 0)

	readonlyTransaction := NewReadonlyTransaction(oracle, storageState)
	iterator, _ := readonlyTransaction.Scan(kv.NewInclusiveKeyRange(kv.RawKey("draft"), kv.RawKey("quadrant")))
	iterator.Close()

	assert.Equal(t, int64(0), ssTable.TotalReferences())
}
