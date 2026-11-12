# Roadmap

Goal: build the answer to the Bolt question ([STORY.md](STORY.md)) as a real distributed database, one small step at a time.

## How each step works

1. **Read** the chapters listed for the step.
2. **Notes**: write them in your own words in [notes/](notes/).
3. **Build**: Claude writes the code, one PR per step.
4. **Walkthrough**: the PR explains the code against the theory. Any Go idiom gets explained the first time it appears, with the Python equivalent.
5. **Prove**: run the step's proof yourself. It's a test that would fail if the theory were misunderstood.
6. **Quiz**: 2–3 questions before moving on. Answers go in the notes.

**Books**
- **DDIA**: *Designing Data-Intensive Applications*, 2nd ed. (Kleppmann & Riccomini, 2026). Cited by chapter title, because the numbers moved from the 1st edition.
- **DI**: *Database Internals* (Petrov, 2019). Cited by chapter number.

**Language**: Go.

## Part 1: one machine

| # | Build | Read | Proof | Status |
|---|---|---|---|---|
| 1 | Append-only log file + in-memory map of key → file offset | DDIA *Storage and Retrieval* (log-structured storage, hash index); DI 1 | restart the process, data still there | next |
| 2 | Split the log into segment files, merge old ones (compaction) | DDIA *Storage and Retrieval* | disk usage drops, nothing lost | |
| 3 | Checksums + crash recovery | DI 3 (File Formats), DI 5 (Transaction Processing and Recovery) | kill -9 mid-write, recovers cleanly | |
| 4 | In-memory sorted buffer flushed to sorted files on disk (LSM tree), range scans | DDIA *Storage and Retrieval* (SSTables, LSM-trees); DI 7 (Log-Structured Storage) | range query correct across files | |
| 5 | B-trees: theory and notes, building one is optional | DI 2, 4, 6 | write up when to pick LSM vs B-tree | |
| 6 | Hourly time index + top-N by duration: the Bolt problem on one node | — | matches brute force on random data | |

## Part 2: many machines

This is the part the Bolt interview was lost on.

| # | Build | Read | Proof | Status |
|---|---|---|---|---|
| 7 | Run as a server, with nodes talking over the network | DDIA *The Trouble with Distributed Systems*; DI 8, 9 (Failure Detection) | a client works against 3 processes | |
| 8 | Sharding + a coordinator that queries every shard and merges the results | DDIA *Sharding* | same answers as the single node | |
| 9 | Leader/follower replication, confirm a write once 2 of 3 copies have it | DDIA *Replication*; DI 11 (Replication and Consistency) | kill a follower, no confirmed write lost | |
| 10 | Failover: elect a new leader | DI 10 (Leader Election); DDIA *Consistency and Consensus* | kill the leader, the cluster keeps going | |
| 11 | Raft (consensus algorithm) | DI 14 (Consensus); DDIA *Consistency and Consensus*; the Raft paper | MIT 6.5840-style tests pass | |

After step 8, the Bolt question has a real answer. After step 9, so do its replication follow-ups.

## Notes only (no build)

Read these alongside the steps:
- DDIA: *Data Models and Query Languages*, *Encoding and Evolution*, *Transactions*, *Batch Processing*, *Stream Processing*.
- DI: 12 (Anti-Entropy and Dissemination), 13 (Distributed Transactions).
