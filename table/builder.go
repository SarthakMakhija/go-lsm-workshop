package table

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-lsm-workshop/kv"
	"go-lsm-workshop/table/block"
	"go-lsm-workshop/table/bloom"
	"path/filepath"
)

// SSTableBuilder allows building SSTable in a step-by-step manner.
type SSTableBuilder struct {
	blockBuilder       *block.Builder
	blockMetaList      *block.MetaList
	bloomFilterBuilder *bloom.FilterBuilder
	startingKey        kv.Key
	endingKey          kv.Key
	allBlocksData      []byte
	blockSize          uint
}

// NewSSTableBuilderWithDefaultBlockSize creates a new instance of SSTableBuilder with block.DefaultBlockSize = 4Kb.
func NewSSTableBuilderWithDefaultBlockSize() *SSTableBuilder {
	return NewSSTableBuilder(block.DefaultBlockSize)
}

// NewSSTableBuilder creates a new instance of SSTableBuilder with the given block size.
// The specified block size will be used to limit the size of each block that will be a part of the final SSTable.
func NewSSTableBuilder(blockSize uint) *SSTableBuilder {
	return &SSTableBuilder{
		blockBuilder:       block.NewBlockBuilder(blockSize),
		blockMetaList:      block.NewBlockMetaList(),
		bloomFilterBuilder: bloom.NewBloomFilterBuilder(),
		blockSize:          blockSize,
	}
}

// Add adds the key/value pair in the current block builder.
// Add involves:
// 1) Keeping a track of the starting key and ending key of the current block.
// 2) Adding the key to the bloom.FilterBuilder
// 3) Adding the key/value pair to the current block.Builder.
// 4) Finishing the current block, if it is full and starting a new block (or block.Builder).
func (builder *SSTableBuilder) Add(key kv.Key, value kv.Value) {
	if builder.startingKey.IsRawKeyEmpty() {
		builder.startingKey = key
	}
	builder.endingKey = key
	builder.bloomFilterBuilder.Add(key)

	//Assignment 1
	//Step1: Add the key-value pair to the block builder.

	builder.finishBlock()
	builder.startNewBlockBuilder(key)
	builder.blockBuilder.Add(key, value)
}

// Build builds the SSTable using the given id and rootPath.
// It involves encoding the SSTable, writing the entire table to persistent storage and creating an in-memory representation
// in the form of SSTable with a reference to its File.
// The encoding looks like:
/**
  ----------------------------------------------------------------------------------------------------------------------------------------------------------
| data block | data block |...| data block | metadata section | 4 bytes for meta starting offset | bloom filter section | 4 bytes for bloom starting offset |
|										   |				  |									 |		                |                                   |			                                        |
 ----------------------------------------------------------------------------------------------------------------------------------------------------------
*/
func (builder *SSTableBuilder) Build(id uint64, rootPath string) (*SSTable, error) {
	blockMetaStartingOffset := func() []byte {
		blockMetaStartingOffset := make([]byte, block.Uint32Size)
		binary.LittleEndian.PutUint32(blockMetaStartingOffset, uint32(len(builder.allBlocksData)))
		return blockMetaStartingOffset
	}
	bloomStartingOffset := func(buffer *bytes.Buffer) []byte {
		bloomStartingOffset := make([]byte, block.Uint32Size)
		binary.LittleEndian.PutUint32(bloomStartingOffset, uint32(buffer.Len()))
		return bloomStartingOffset
	}

	builder.finishBlock()
	buffer := new(bytes.Buffer)

	//Assignment 3
	//Step1: Write allBlocksData in the buffer
	//Step2: Write block meta-list
	//Step3: Write the starting offset of the block meta-section

	filter := builder.bloomFilterBuilder.Build(bloom.FalsePositiveRate)
	encodedFilter, err := filter.Encode()
	if err != nil {
		return nil, err
	}

	bloomFilterStartingOffset := bloomStartingOffset(buffer)

	//Step4: Write the bloom filter
	//Step5: Write the starting offset of the bloom filter

	file, err := CreateAndWrite(SSTableFilePath(id, rootPath), buffer.Bytes())
	if err != nil {
		return nil, err
	}

	startingKey, _ := builder.blockMetaList.StartingKeyOfFirstBlock()
	endingKey, _ := builder.blockMetaList.EndingKeyOfLastBlock()
	return &SSTable{
		id:                      id,
		file:                    file,
		blockMetaList:           builder.blockMetaList,
		bloomFilter:             filter,
		blockMetaStartingOffset: uint32(len(builder.allBlocksData)),
		blockSize:               builder.blockSize,
		startingKey:             startingKey,
		endingKey:               endingKey,
	}, nil
}

// EstimatedSize returns an estimate of the size of the encoded data of all the blocks.
func (builder SSTableBuilder) EstimatedSize() int {
	return len(builder.allBlocksData)
}

// finishBlock finishes the current block. It involves:
// 1) Encoding the current block.
// 2) Storing the block.Meta in the block meta-list.
// 3) Collecting the encoded data of the current block in allBlocksData.
func (builder *SSTableBuilder) finishBlock() {
	//Assignment 2
	//Step1: Encode the block using blockBuilder.
	//Step2: Store the block.Meta in the block meta-list.
	//Step3: Collect the encoded data of the current block in allBlocksData.
}

// startNewBlockBuilder creates a new instance of SSTableBuilder.
func (builder *SSTableBuilder) startNewBlockBuilder(key kv.Key) {
	builder.blockBuilder = block.NewBlockBuilder(builder.blockSize)
	builder.startingKey = key
	builder.endingKey = key
}

// SSTableFilePath returns the SSTable filepath which consists of rootPath/id.sst.
func SSTableFilePath(id uint64, rootPath string) string {
	return filepath.Join(rootPath, fmt.Sprintf("%v.sst", id))
}
