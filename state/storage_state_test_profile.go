//go:build test

package state

import (
	"go-lsm-workshop/manifest"
	"go-lsm-workshop/memory"
	"go-lsm-workshop/table"
	"os"
	"slices"
	"time"
)

// NewStorageState creates new instance of StorageState for testing.
func NewStorageState(rootPath string) (*StorageState, error) {
	return NewStorageStateWithOptions(StorageOptions{
		MemTableSizeInBytes:   1 << 20,
		SSTableSizeInBytes:    1 << 20,
		Path:                  rootPath,
		MaximumMemtables:      5,
		FlushMemtableDuration: 50 * time.Millisecond,
		CompactionOptions: CompactionOptions{
			StrategyOptions: SimpleLeveledCompactionOptions{
				Level0FilesCompactionTrigger:    6,
				MaxLevels:                       totalLevels,
				NumberOfSSTablesRatioPercentage: 200,
			},
		},
	})
}

// DeleteManifest deletes Manifest file, only for testing.
func (storageState *StorageState) DeleteManifest() {
	storageState.manifest.Delete()
}

// DeleteWALDirectory deletes WAL directory path, only for testing.
func (storageState *StorageState) DeleteWALDirectory() {
	if len(storageState.WALDirectoryPath()) > 0 {
		_ = os.RemoveAll(storageState.WALDirectoryPath())
	}
}

// HasImmutableMemtables returns true if there are immutable tables, it is only for testing.
func (storageState *StorageState) HasImmutableMemtables() bool {
	return len(storageState.immutableMemtables) > 0
}

// TotalImmutableMemtables returns the total number of immutable memtables, it is only for testing.
func (storageState *StorageState) TotalImmutableMemtables() int {
	return len(storageState.immutableMemtables)
}

// TotalSSTablesAtLevel returns the total number of SSTables at the given level, it is only for testing.
func (storageState *StorageState) TotalSSTablesAtLevel(level int) int {
	if level == 0 {
		return len(storageState.l0SSTableIds)
	}
	return len(storageState.levels[level-1].SSTableIds)
}

// SSTableReferenceCountAtLevel returns a slice of references of table.SSTable at the given level, it is only for testing.
func (storageState *StorageState) SSTableReferenceCountAtLevel(level int) ([]int64, int) {
	var totalReferenceCount []int64
	if level == 0 {
		for _, ssTableId := range storageState.l0SSTableIds {
			totalReferenceCount = append(totalReferenceCount, storageState.ssTables[ssTableId].TotalReferences())
		}
	} else {
		for _, ssTableId := range storageState.levels[level-1].SSTableIds {
			totalReferenceCount = append(totalReferenceCount, storageState.ssTables[ssTableId].TotalReferences())
		}
	}
	return totalReferenceCount, len(totalReferenceCount)
}

// SetSSTableAtLevel sets SSTable at the given level, only for testing.
func (storageState *StorageState) SetSSTableAtLevel(ssTable *table.SSTable, level int) {
	if level == 0 {
		storageState.l0SSTableIds = append(storageState.l0SSTableIds, ssTable.Id())
	} else {
		existingLevel := storageState.levels[level-1]
		if existingLevel == nil {
			existingLevel = &Level{LevelNumber: level}
		}
		existingLevel.SSTableIds = append(existingLevel.SSTableIds, ssTable.Id())
		storageState.levels[level-1] = existingLevel
	}
	storageState.ssTables[ssTable.Id()] = ssTable
}

// ForceFlushNextImmutableMemtable flushes the next immutable memtable to Level 0, it is only for testing.
func (storageState *StorageState) ForceFlushNextImmutableMemtable() error {
	return storageState.forceFlushNextImmutableMemtable()
}

// forceFreezeCurrentMemtable freezes the current memtable, it is only for testing.
func (storageState *StorageState) forceFreezeCurrentMemtable() {
	storageState.immutableMemtables = append(storageState.immutableMemtables, storageState.currentMemtable)
	storageState.currentMemtable = memory.NewMemtable(
		storageState.idGenerator.NextId(),
		storageState.options.MemTableSizeInBytes,
		storageState.walPath,
	)
	_ = storageState.manifest.Add(manifest.NewMemtableCreated(storageState.currentMemtable.Id()))
}

// hasSSTableWithId returns true if there is an SSTable for the given SSTableId, false otherwise, it is only for testing.
func (storageState *StorageState) hasSSTableWithId(id uint64) bool {
	_, ok := storageState.ssTables[id]
	return ok
}

// sortedMemtableIds returns the sorted memtableIds,  it is only for testing.
func (storageState *StorageState) sortedMemtableIds() []uint64 {
	ids := make([]uint64, 0, 1+len(storageState.immutableMemtables))
	ids = append(ids, storageState.currentMemtable.Id())
	for _, immutableMemtable := range storageState.immutableMemtables {
		ids = append(ids, immutableMemtable.Id())
	}
	slices.Sort(ids)
	return ids
}
