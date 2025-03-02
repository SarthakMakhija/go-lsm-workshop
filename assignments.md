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

#### SSTableBuilder Assignments

1. Assignment 1

```go
    if builder.blockBuilder.Add(key, value) {
		return
	}
```

2. Assignment 2

```go
    encodedBlock := builder.blockBuilder.Build().Encode()
	builder.blockMetaList.Add(block.Meta{
		BlockStartingOffset: uint32(len(builder.allBlocksData)),
		StartingKey:         builder.startingKey,
		EndingKey:           builder.endingKey,
	})
	builder.allBlocksData = append(builder.allBlocksData, encodedBlock...)
```

3. Assignment 3

```go
    buffer.Write(builder.allBlocksData)          //data blocks
	buffer.Write(builder.blockMetaList.Encode()) //metadata section block.MetaList.Encode()
	buffer.Write(blockMetaStartingOffset())      //4 bytes to indicate where the meta section starts from
	
    buffer.Write(encodedFilter)             //bloom filter section bloom.Filter.Encode()
	buffer.Write(bloomFilterStartingOffset) //4 bytes to indicate where the bloom filter section starts from
```

#### BlockBuilder Assignments

1. Assignment 1

```go
    binary.LittleEndian.PutUint16(keyValueBuffer[:], uint16(key.EncodedSizeInBytes()))
	copy(keyValueBuffer[ReservedKeySize:], key.EncodedBytes())

	binary.LittleEndian.PutUint16(keyValueBuffer[ReservedKeySize+key.EncodedSizeInBytes():], uint16(value.SizeInBytes()))
	copy(keyValueBuffer[ReservedKeySize+key.EncodedSizeInBytes()+ReservedValueSize:], value.Bytes())
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

#### SSTable Assignments

1. Assignment 1

```go
    table.file.Read(int64(startingOffset), buffer)
```

2. Assignment 2

```go
    table.blockMetaList.MaybeBlockMetaContaining(key)
```

3. Assignment 3

```go
    mid := low + (high-low)/2
    meta := metaList.list[mid]
    switch key.CompareKeysWithDescendingTimestamp(meta.StartingKey) {
    case -1:
        high = mid - 1
    case 0:
        return meta, mid
    case 1:
        possibleIndex = mid
        low = mid + 1
    }
```

#### Transactions Assignments

1. Assignment 1

```go
    beginTimestamp := oracle.nextTimestamp - 1
```

2. Assignment 2

```go
    commitTimestamp := oracle.nextTimestamp
```

3. Assignment 3

```go
    for _, committedTransaction := range oracle.readyToCommitTransactions {
		if committedTransaction.commitTimestamp <= transaction.beginTimestamp {
			continue
		}
		for _, key := range transaction.reads {
			if committedTransaction.transaction.batch.Contains(key) {
				return true
			}
		}
	}
```

4. Assignment 4:

```go
    transaction.trackReads(key)
```

5. Assignment 5:

```go
    transaction.oracle.executor.submit(kv.NewTimestampedBatchFrom(*transaction.batch, commitTimestamp)
```

6. Assignment 6:

```go
    executor.state.Set(executionRequest.batch)
    executionRequest.future.MarkDoneAsError(err)
    executionRequest.future.MarkDoneAsOk()
```

#### Compaction

