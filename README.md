# Internals of key-value storage engines: LSM-trees and beyond 

Workshop on **Internals of key-value storage engines: LSM-trees and beyond** based on the [go-lsm](https://github.com/SarthakMakhija/go-lsm) repository.

_This code neither compiles, nor runs. The code will run after all the assignments are done._

### About the workshop

This workshop offers a deep dive into the practical aspects of building a an embedded key-value storage engine using the Log-Structured Merge-tree (LSM-tree) architecture. 
Participants will gain a comprehensive understanding of the fundamental concepts, starting with the intricacies of block devices, file I/O, and disk I/O patterns, along with a review of B+Trees for comparison.  
The workshop delves into the theoretical underpinnings of RUM conjecture and the LSM-tree itself, before transitioning into hands-on implementation. 

Attendees will learn to build core components such as Memtables, Write-Ahead Logs (WAL), and SSTables enhanced with Bloom filters.  
Finally, the workshop covers advanced topics like transaction management and compaction strategies, providing a complete roadmap for creating a robust and efficient storage engine. 

This hands-on experience will equip participants with the practical skills and theoretical knowledge necessary to tackle real-world storage challenges.

### Workshop Content

<img width="636" alt="Workshop content" src="https://github.com/user-attachments/assets/4d0c55de-28c0-42e1-b419-ef20b56cfb6d" />

### Prerequisites for attending

To get the most out of this hands-on workshop, participants should come prepared with:

- **Practical Golang experience**: You should be comfortable writing and running Go code.
- **Go 1.24 installation**: Please ensure you have Go version 1.24 installed on your machine before the workshop begins.
