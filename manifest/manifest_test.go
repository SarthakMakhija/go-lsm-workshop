package manifest

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateANewManifestWithNewMemtableCreatedEvent(t *testing.T) {
	manifestDirectoryPath := filepath.Join(".", "TestCreateANewManifestWithNewMemtableCreatedEvent")
	assert.Nil(t, os.MkdirAll(manifestDirectoryPath, os.ModePerm))

	manifest, _, err := CreateNewOrRecoverFrom(manifestDirectoryPath)
	defer func() {
		_ = os.RemoveAll(manifestDirectoryPath)
	}()

	assert.Nil(t, err)
	assert.Nil(t, manifest.Add(NewMemtableCreated(10)))
}

func TestCreateANewManifestWithNewSSTableFlushedEvent(t *testing.T) {
	manifestDirectoryPath := filepath.Join(".", "TestCreateANewManifestWithNewSSTableFlushedEvent")
	assert.Nil(t, os.MkdirAll(manifestDirectoryPath, os.ModePerm))

	manifest, _, err := CreateNewOrRecoverFrom(manifestDirectoryPath)
	defer func() {
		_ = os.RemoveAll(manifestDirectoryPath)
	}()

	assert.Nil(t, err)
	assert.Nil(t, manifest.Add(NewSSTableFlushed(10)))
}

func TestRecoversAnExistingManifest(t *testing.T) {
	manifestDirectoryPath := filepath.Join(".", "TestRecoversAnExistingManifest")
	assert.Nil(t, os.MkdirAll(manifestDirectoryPath, os.ModePerm))

	manifest, _, err := CreateNewOrRecoverFrom(manifestDirectoryPath)
	defer func() {
		_ = os.RemoveAll(manifestDirectoryPath)
	}()

	assert.Nil(t, err)
	assert.Nil(t, manifest.Add(NewMemtableCreated(10)))
	assert.Nil(t, manifest.Add(NewMemtableCreated(20)))
	assert.Nil(t, manifest.Add(NewSSTableFlushed(10)))

	upperLevel := -1
	lowerLevel := 1
	upperLevelSSTableIds := []uint64{20, 30}
	lowerLevelSSTableIds := []uint64{50, 60}

	compactionDone := NewCompactionDone([]uint64{10, 11}, upperLevel, lowerLevel, upperLevelSSTableIds, lowerLevelSSTableIds)
	assert.Nil(t, manifest.Add(compactionDone))

	manifest, events, err := CreateNewOrRecoverFrom(manifestDirectoryPath)
	assert.Nil(t, err)

	assert.Equal(t, 4, len(events))
	assert.Equal(t, uint64(10), events[0].(*MemtableCreated).MemtableId)
	assert.Equal(t, uint64(20), events[1].(*MemtableCreated).MemtableId)
	assert.Equal(t, uint64(10), events[2].(*SSTableFlushed).SsTableId)

	assert.Equal(t, []uint64{10, 11}, events[3].(*CompactionDone).NewSSTableIds)
	assert.Equal(t, -1, events[3].(*CompactionDone).UpperLevel)
	assert.Equal(t, 1, events[3].(*CompactionDone).LowerLevel)
	assert.Equal(t, []uint64{20, 30}, events[3].(*CompactionDone).UpperLevelSSTableIds)
	assert.Equal(t, []uint64{50, 60}, events[3].(*CompactionDone).LowerLevelSSTableIds)
}
