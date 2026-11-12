# Roadmap

Goal: build the answer to the Bolt question ([STORY.md](STORY.md)) as a real distributed database, one small step at a time.

## How each step works

It runs like a class.

1. **Lesson**: Claude teaches the theory in chat, with one concrete example and the papers and book chapters behind it.
2. **Questions**: Claude asks, you answer in your own words. No code until the idea has clicked.
3. **Notes**: yours, in [notes/](notes/).
4. **Code session**: Claude builds the step, you read it with a walkthrough. Any Go idiom gets explained the first time it appears, with the Python equivalent.
5. **Prove**: you run the step's proof yourself. It's a test that would fail if the theory were wrong.
6. **Quiz**: 2–3 questions before the next step.

**Books** (for going deeper after the lesson)
- **DDIA**: *Designing Data-Intensive Applications*, 2nd ed. (Kleppmann & Riccomini, 2026). Cited by chapter title, because the numbers moved from the 1st edition.
- **DI**: *Database Internals* (Petrov, 2019). Cited by chapter number.

**Language**: Go.

## Part 1: one machine

| # | Build | Theory | Proof | Status |
|---|---|---|---|---|
| 1 | Append-only log file + in-memory map of key → file offset | DDIA *Storage and Retrieval*; DI 1. Paper: *Bitcask* (Sheehy & Smith, 2010) | restart the process, data still there | in class |
| 2 | Split the log into segment files, merge old ones (compaction) | DDIA *Storage and Retrieval*; *Bitcask* (merging) | disk usage drops, nothing lost | |
| 3 | Checksums + crash recovery | DI 3, DI 5. Paper (skim): *ARIES* (Mohan et al., 1992) | kill -9 mid-write, recovers cleanly | |
| 4 | In-memory sorted buffer flushed to sorted files on disk (LSM tree), range scans | DDIA *Storage and Retrieval*; DI 7. Papers: *The Log-Structured Merge-Tree* (O'Neil et al., 1996), *Bigtable* (Chang et al., 2006) | range query correct across files | |
| 5 | B-trees: theory and notes, building one is optional | DI 2, 4, 6. Paper: *The Ubiquitous B-Tree* (Comer, 1979) | write up when to pick LSM vs B-tree | |
| 6 | Hourly time index + top-N by duration: the Bolt problem on one node | — | matches brute force on random data | |

## Part 2: many machines

This is the part the Bolt interview was lost on.

| # | Build | Theory | Proof | Status |
|---|---|---|---|---|
| 7 | Run as a server, with nodes talking over the network | DDIA *The Trouble with Distributed Systems*; DI 8, 9. Paper: *Time, Clocks, and the Ordering of Events* (Lamport, 1978) | a client works against 3 processes | |
| 8 | Sharding + a coordinator that queries every shard and merges the results | DDIA *Sharding*. Paper: *Dynamo* (DeCandia et al., 2007) | same answers as the single node | |
| 9 | Leader/follower replication, confirm a write once 2 of 3 copies have it | DDIA *Replication*; DI 11 | kill a follower, no confirmed write lost | |
| 10 | Failover: elect a new leader | DI 10; DDIA *Consistency and Consensus* | kill the leader, the cluster keeps going | |
| 11 | Raft (consensus algorithm) | DI 14; DDIA *Consistency and Consensus*. Paper: *In Search of an Understandable Consensus Algorithm* (Ongaro & Ousterhout, 2014) | MIT 6.5840-style tests pass | |

After step 8, the Bolt question has a real answer. After step 9, so do its replication follow-ups.

## Notes only (no build)

Read these alongside the steps:
- DDIA: *Data Models and Query Languages*, *Encoding and Evolution*, *Transactions*, *Batch Processing*, *Stream Processing*.
- DI: 12 (Anti-Entropy and Dissemination), 13 (Distributed Transactions).
