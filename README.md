# MemryDB

**MemryDB** is a high-performance, **in-memory** key-value storage engine built from the ground up in Go. It focuses on extreme speed and thread-safety by leveraging an optimized sharding architecture to minimize lock contention.



## 🚀 Features

* **In-Memory Performance:** Optimized for ultra-low latency data access entirely in RAM.
* **Thread-Safe Sharding:** Utilizes a `SharedMap` with multiple independent shards, each protected by its own `sync.RWMutex`.
* **Zero-Allocation Hasher:** Includes a custom FNV-1 implementation that is **4x faster** than the Go standard library.
* **Minimal GC Pressure:** Architected to keep internal logic on the stack, reducing Garbage Collector overhead and preventing "Stop the World" latency spikes.
