#### Memtable Assignments

1. Assignment 1

```go
    memtable.entries.Put(key, value)
```

2. Assignment 2

```go
    memtable.Set(key, kv.EmptyValue)
```

3. Assignment 3

```go
    memtable.entries.Put(key, value)
```

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

#### Bloom filter Assignments

1. Assignment 1

```go
for index := 0; index < len(positions); index++ {
    position := positions[index]
    filter.bitVector.Set(uint(position))
}
```

2. Assignment 2

```go
for index := 0; index < len(positions); index++ {
    position := positions[index]
    if !filter.bitVector.Test(uint(position)) {
        return false
    }
}
```