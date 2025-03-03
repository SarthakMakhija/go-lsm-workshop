package compact

import (
	"go-lsm-workshop/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCompactionTaskForSimpleLayeredCompactionWithNoCompaction(t *testing.T) {
	compactionOptions := state.SimpleLeveledCompactionOptions{
		NumberOfSSTablesRatioPercentage: 200,
		MaxLevels:                       2,
		Level0FilesCompactionTrigger:    2,
	}
	snapshot := state.StorageStateSnapshot{
		L0SSTableIds: []uint64{1},
		Levels: []*state.Level{
			{LevelNumber: 1, SSTableIds: []uint64{2, 3}},
			{LevelNumber: 2, SSTableIds: []uint64{4, 5, 6, 7}},
		},
	}

	compaction := NewSimpleLeveledCompaction(compactionOptions)
	_, ok := compaction.CompactionDescription(snapshot)

	assert.False(t, ok)
}

func TestGenerateCompactionTaskForSimpleLayeredCompactionWithCompactionForLevel0And1(t *testing.T) {
	compactionOptions := state.SimpleLeveledCompactionOptions{
		NumberOfSSTablesRatioPercentage: 200,
		MaxLevels:                       2,
		Level0FilesCompactionTrigger:    2,
	}
	snapshot := state.StorageStateSnapshot{
		L0SSTableIds: []uint64{1, 2},
		Levels: []*state.Level{
			{LevelNumber: 1, SSTableIds: nil},
		},
	}

	compaction := NewSimpleLeveledCompaction(compactionOptions)
	compactionDescription, ok := compaction.CompactionDescription(snapshot)

	assert.True(t, ok)
	assert.Equal(t, -1, compactionDescription.UpperLevel)
	assert.Equal(t, 1, compactionDescription.LowerLevel)
	assert.Equal(t, []uint64{1, 2}, compactionDescription.UpperLevelSSTableIds)
	assert.Equal(t, []uint64(nil), compactionDescription.LowerLevelSSTableIds)
}

func TestGenerateCompactionTaskForSimpleLayeredCompactionWithCompactionForLevel1And2(t *testing.T) {
	compactionOptions := state.SimpleLeveledCompactionOptions{
		NumberOfSSTablesRatioPercentage: 200,
		MaxLevels:                       2,
		Level0FilesCompactionTrigger:    2,
	}
	snapshot := state.StorageStateSnapshot{
		L0SSTableIds: []uint64{1},
		Levels: []*state.Level{
			{LevelNumber: 1, SSTableIds: []uint64{2, 3}},
			{LevelNumber: 2, SSTableIds: []uint64{4}},
		},
	}

	compaction := NewSimpleLeveledCompaction(compactionOptions)
	compactionDescription, ok := compaction.CompactionDescription(snapshot)

	assert.True(t, ok)
	assert.Equal(t, 1, compactionDescription.UpperLevel)
	assert.Equal(t, 2, compactionDescription.LowerLevel)
	assert.Equal(t, []uint64{2, 3}, compactionDescription.UpperLevelSSTableIds)
	assert.Equal(t, []uint64{4}, compactionDescription.LowerLevelSSTableIds)
}
