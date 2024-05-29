package table

import (
	"bytes"
	"encoding/binary"
	"go-lsm/table/block"
	"go-lsm/txn"
)

type SSTableBuilder struct {
	blockBuilder  *block.Builder
	blockMetaList *BlockMetaList
	startingKey   txn.Key
	data          []byte
	blockSize     uint
}

func NewSSTableBuilder(blockSize uint) *SSTableBuilder {
	return &SSTableBuilder{
		blockBuilder:  block.NewBlockBuilder(blockSize),
		blockMetaList: NewBlockMetaList(),
		blockSize:     blockSize,
	}
}

func (builder *SSTableBuilder) Add(key txn.Key, value txn.Value) {
	if builder.startingKey.IsEmpty() {
		builder.startingKey = key
	}
	if builder.blockBuilder.Add(key, value) {
		return
	}
	builder.finishBlock()
	builder.startNewBlock(key)
	builder.blockBuilder.Add(key, value)
}

// Build
// TODO: Bloom
func (builder *SSTableBuilder) Build(id uint64, filePath string) (SSTable, error) {
	builder.finishBlock()

	blockMetaOffset := make([]byte, block.Uint32Size)
	binary.LittleEndian.PutUint32(blockMetaOffset, uint32(len(builder.data)))

	buffer := new(bytes.Buffer)
	buffer.Write(builder.data)
	buffer.Write(builder.blockMetaList.encode())
	buffer.Write(blockMetaOffset)

	file, err := Create(filePath, buffer.Bytes())
	if err != nil {
		return SSTable{}, err
	}
	//TODO: Block cache + bloom fields
	return SSTable{
		id:              id,
		file:            file,
		blockMetaList:   builder.blockMetaList,
		blockMetaOffset: uint32(len(builder.data)),
		blockSize:       builder.blockSize,
	}, nil
}

func (builder *SSTableBuilder) finishBlock() {
	encodedBlock := builder.blockBuilder.Build().Encode()
	builder.blockMetaList.add(BlockMeta{
		offset:      uint32(len(builder.data)),
		startingKey: builder.startingKey,
	})
	builder.data = append(builder.data, encodedBlock...)
}

func (builder *SSTableBuilder) startNewBlock(key txn.Key) {
	builder.blockBuilder = block.NewBlockBuilder(builder.blockSize)
	builder.startingKey = key
}
