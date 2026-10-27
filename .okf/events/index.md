# events

* [Computation cancelled](/events/computation-cancelled.md) - Domain event: a running computation was interrupted by client, navigation or shutdown, carrying its partial state and freeing its slot.
* [Computation completed](/events/computation-completed.md) - Domain event: a background computation finished successfully, its result became available, its retention window started and its computation slot was released.
* [Job expired](/events/job-expired.md) - Domain event: a finished job passed its retention window (or the count ceiling evicted it) and polling it now yields not-found.
