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

### Reflections

**06th-07th March 2025**

The workshop was well-received, with participants providing positive feedback regarding the content delivery and the overall flow of information.
I was satisfied with the overall timing, the content delivery and the way the workshop content has shaped up and story building leading to a concept. 

However, there are a few things that I would like to work upon:

1. The "Transactions" section felt somewhat rushed. To address this, the Day 2 schedule will need restructuring to dedicate a full 90 minutes to the "Transactions" topic, ensuring adequate time for exploration and practice.
2. The current assignment structure includes some relatively simple tasks that could be eliminated. This adjustment will allow for the inclusion of more assignments focused on "SSTable" and "Transactions".
3. A dedicated 15-minute segment will be added to the workshop to provide a clear and concise overview of how all the components of the Log-Structured Merge-tree (LSM) work together in [go-lsm](https://github.com/SarthakMakhija/go-lsm).

### Workshop content

| Date                 | PDF Link                                                                                                                                                                                |
|----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 06th-07th March 2025 | [Internals of key-value storage engines: LSM-trees and beyond.pdf](https://github.com/user-attachments/files/19141831/Internals.of.key-value.storage.engines_.LSM-trees.and.beyond.pdf) |


