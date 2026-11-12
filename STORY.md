# Why Tay-db exists

## The interview

In 2025 I interviewed for a Software Engineer role at Bolt. One of the rounds was system design, and this was the question:

> We want to build a small analytics system for fraud detection on orders. It needs to provide the following API:
>
> - `/insertOrder(order: Order): void`: add an order to the storage. The order can be considered a blob of 1-10KBs in size, with an `orderId`, `beginTime` and `finishTime` as distinguished fields.
> - `/getLongestNOrdersByDuration(n: int, startTime: datetime, endTime: datetime): Order[]`: retrieve the longest n orders that started between `startTime` and `endTime`, as measured by duration (`finishTime - beginTime`).
> - `/getShortestNOrdersByDuration(n: int, startTime: datetime, endTime: datetime): Order[]`: retrieve the shortest n orders that started between `startTime` and `endTime`, as measured by duration (`finishTime - beginTime`).
>
> You are not allowed to use any external systems at all! So no MySQL, no Redis, no S3, etc. Just a bunch of machines in "the cloud", with disks and decent amounts of memory and compute.

It looked like a database question, and it was. "No external systems" meant I had to design the database myself.

I started well. I asked about data volumes, read and write patterns, and retention. I suggested using a heap to find the longest and shortest orders. I knew the data had to survive a machine dying, so I said "replication".

Then I designed the whole thing on one machine.

## The response

Then this arrived:

> Hey, Abdulmateen!
>
> Thank you again for taking the time to interview with us for the Software Engineer position. We enjoyed our conversation during the System Design interview.
>
> After careful consideration, we've decided not to move forward with your application at this time.
>
> We want to provide some specific, constructive feedback, as we were impressed by several aspects of your interview.
>
> **What Went Well**
>
> We really liked your friendly communication style and your methodical approach to gathering non-functional requirements. You did a good job asking about data volumes, read/write patterns, and retention, which is a critical starting point. We also thought your suggestion to use a heap to find the longest and shortest orders was a great algorithmic insight.
>
> **Areas for Growth**
>
> The main gap we identified was in scaling the solution beyond a single machine. The design you proposed was a good solution for a single node, but it didn't address the challenges of a large-scale, distributed system, which is a key requirement for this role.
>
> In future system design interviews, we'd recommend focusing on these areas:
>
> **Data Partitioning (Sharding):** The design needed a strategy for partitioning the data (sharding) to handle a volume that wouldn't fit on a single database. Proactively identifying this bottleneck and suggesting a solution is key.
>
> **Distributed Design:** We were looking for a design that used multiple nodes for both storage and processing, rather than centralizing the logic on one machine.
>
> **Replication Details:** You correctly identified the need for reliability and suggested replication. The next step, which we were looking for, was to describe how that replication would be implemented (e.g., number of replicas, leader/follower model, consistency trade-offs).
>
> You have a strong foundation, especially in algorithms and initial requirements gathering. With more focus on distributed architecture and scalability patterns, we're confident you'll be a strong candidate for future roles.
>
> I'd also be happy to discuss this over a call if that is what you'd like. Let me know if this is something you'd be open to, but no pressure.
>
> We wish you the very best in your job search.
>
> — Bolt Recruiting

## What I took from it

I had the answer in my hands and didn't use it.

I asked how much data there would be. The numbers I got back already said it wouldn't fit on one machine. That was the moment to say "so we split it across many machines". I didn't say it. I kept designing for one box.

Their three points are really one question: **what happens when the data doesn't fit on one machine?**

- How do you split it? (sharding)
- How do many machines answer one query together? (distributed design)
- How do you keep copies, and what happens when one copy dies? (replication)

The heap was right. It just belonged on every machine and at the merge step, not on a single server.

## What this repo is

Tay-db is the database I should have been able to design in that room.

I'm building it little by little, and learning the theory as I go:

- **Designing Data-Intensive Applications** by Martin Kleppmann
- **Database Internals** by Alex Petrov

Each step goes the same way: read a chapter, write notes in my own words, build the smallest thing that uses the idea, and prove it with a test.

It starts on one machine: a log file, then crash recovery, then sorted files on disk, then the Bolt problem on a single node.

Then it grows past one machine, which is the part I missed: sharding, queries spread across nodes, replication, and leader election.

The finish line is the Bolt question, answered properly, running on real code I wrote and understand.

And then, maybe, I'll interview again.
