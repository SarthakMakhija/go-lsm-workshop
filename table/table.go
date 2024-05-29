package table

import (
	"bytes"
	"encoding/binary"
	"go-lsm/txn"
)

type SSTable struct {
	id              uint64
	blockMetaList   *BlockMetaList
	blockMetaOffset uint32
	file            *File
	blockSize       uint
}

func (table SSTable) readBlock(blockIndex int) (Block, error) {
	startingOffset, endOffset := table.offsetRangeOfBlockAt(blockIndex)
	buffer := make([]byte, endOffset-startingOffset)
	n, err := table.file.Read(int64(startingOffset), buffer)
	if err != nil {
		return Block{}, err
	}
	return decodeToBlock(buffer[:n]), nil
}

func (table SSTable) offsetRangeOfBlockAt(blockIndex int) (uint32, uint32) {
	blockMeta, ok := table.blockMetaList.getAt(blockIndex)
	if !ok {
		panic("block meta not found")
	}
	nextBlockMeta, ok := table.blockMetaList.getAt(blockIndex + 1)
	var endOffset uint32
	if ok {
		endOffset = nextBlockMeta.offset
	} else {
		endOffset = table.blockMetaOffset
	}
	return blockMeta.offset, endOffset
}

type BlockMeta struct {
	offset      uint32
	startingKey txn.Key
}

type BlockMetaList struct {
	list []BlockMeta
}

func NewBlockMetaList() *BlockMetaList {
	return &BlockMetaList{}
}

func (metaList *BlockMetaList) add(block BlockMeta) {
	metaList.list = append(metaList.list, block)
}

func (metaList *BlockMetaList) encode() []byte {
	numberOfBlocks := make([]byte, uint32Size)
	binary.LittleEndian.PutUint32(numberOfBlocks, uint32(len(metaList.list)))

	resultingBuffer := new(bytes.Buffer)
	resultingBuffer.Write(numberOfBlocks)

	for _, blockMeta := range metaList.list {
		buffer := make([]byte, uint32Size+reservedKeySize+blockMeta.startingKey.Size())

		binary.LittleEndian.PutUint32(buffer[:], blockMeta.offset)
		binary.LittleEndian.PutUint16(buffer[uint32Size:], uint16(blockMeta.startingKey.Size()))
		copy(buffer[uint32Size+reservedKeySize:], blockMeta.startingKey.Bytes())

		resultingBuffer.Write(buffer)
	}

	return resultingBuffer.Bytes()
}

func (metaList *BlockMetaList) getAt(index int) (BlockMeta, bool) {
	if index < len(metaList.list) {
		return metaList.list[index], true
	}
	return BlockMeta{}, false
}

func decodeToBlockMetaList(buffer []byte) BlockMetaList {
	numberOfBlocks := binary.LittleEndian.Uint32(buffer[:])
	blockList := make([]BlockMeta, 0, numberOfBlocks)

	buffer = buffer[uint32Size:]
	for index := 0; index < len(buffer); {
		offset := binary.LittleEndian.Uint32(buffer[index:])
		keySize := binary.LittleEndian.Uint16(buffer[index+uint32Size:])
		key := buffer[index+uint32Size+reservedKeySize : index+uint32Size+reservedKeySize+int(keySize)]

		blockList = append(blockList, BlockMeta{
			offset:      offset,
			startingKey: txn.NewKey(key),
		})
		index = index + uint32Size + reservedKeySize + int(keySize)
	}
	return BlockMetaList{
		list: blockList,
	}
}
