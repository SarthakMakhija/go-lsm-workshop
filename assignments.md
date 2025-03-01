#### WAL Assignments

1. Assignment 1

```go
	binary.LittleEndian.PutUint16(buffer, uint16(key.EncodedSizeInBytes()))
	copy(buffer[block.ReservedKeySize:], key.EncodedBytes())

	binary.LittleEndian.PutUint16(buffer[block.ReservedKeySize+key.EncodedSizeInBytes():], uint16(value.SizeInBytes()))
	copy(buffer[block.ReservedKeySize+key.EncodedSizeInBytes()+block.ReservedValueSize:], value.Bytes())
    
    _, err := wal.file.Write(buffer)
```

2. Assignment 2

```go
		keySize := binary.LittleEndian.Uint16(bytes)
		key := bytes[block.ReservedKeySize : uint16(block.ReservedKeySize)+keySize]

		valueSize := binary.LittleEndian.Uint16(bytes[uint16(block.ReservedKeySize)+keySize:])
		value := bytes[uint16(block.ReservedKeySize)+keySize+uint16(block.ReservedValueSize) : uint16(block.ReservedKeySize)+keySize+uint16(block.ReservedValueSize)+valueSize]
```