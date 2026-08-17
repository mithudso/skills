# MVCC, VACUUM & Bloat

> Reference for **postgresql-expert**. Engine internals for multi-version concurrency control, dead-tuple reclamation, and the bloat/wraparound failure modes that follow from them. Verified-as-of 2026-06-30 (PostgreSQL 18 line).

## The MVCC model

PostgreSQL implements **MVCC by keeping multiple physical versions of each row** in the heap. There is no separate undo/rollback segment (unlike Oracle or MySQL InnoDB) — old versions live *in the table itself* until vacuumed.

Every heap tuple carries a 23-byte header. The fields that drive visibility:

| Field | Meaning |
|---|---|
| `xmin` | Transaction id (XID) that **inserted** this tuple version |
| `xmax` | XID that **deleted/updated** it (0 if still live) |
| `cmin`/`cmax` | Command id within the transaction (intra-txn visibility) |
| `ctid` | Physical location `(block, offset)` of this version; an updated row's old version points forward via the update chain |
| `t_infomask` | Hint bits: `HEAP_XMIN_COMMITTED`, `HEAP_XMAX_COMMITTED`, frozen, etc. |

**A write never overwrites in place.** `UPDATE` = insert a new tuple version + set `xmax` on the old. `DELETE` = set `xmax`. The old version remains until VACUUM removes it.

### Snapshots & visibility

A transaction takes a **snapshot**: `(xmin_horizon, xmax, [in-progress XID list])`. A tuple is visible iff its `xmin` committed before the snapshot **and** its `xmax` is null/aborted/after the snapshot. Consequences:

- **Readers never block writers and writers never block readers** — each sees its own consistent version. Only writer-writer conflicts on the *same row* block.
- **Long-running transactions are dangerous.** An open transaction (or an idle-in-transaction session, or a stale replication slot / `hot_standby_feedback`) holds back the **xmin horizon**, so VACUUM cannot remove dead tuples newer than the oldest snapshot anywhere in the system. This is the #1 root cause of runaway bloat.
- Hint bits are set lazily on first read after commit; this is why the *first* `SELECT` after a big write can be unexpectedly slow (it dirties pages writing hint bits).

### HOT (Heap-Only Tuple) updates

If an `UPDATE` does **not** change any indexed column **and** the new version fits on the **same page**, PostgreSQL performs a **HOT update**: the new version is chained via `ctid` on the same page and **no index entry is created**. Indexes point at the original tuple; the chain is followed at read time. HOT updates dramatically reduce index bloat and write amplification.

- Maximize HOT: leave `fillfactor` headroom on update-heavy tables (e.g. `fillfactor=85–90`); avoid indexing hot-updated columns.
- `pg_stat_user_tables.n_tup_hot_upd` vs `n_tup_upd` shows your HOT ratio.

## VACUUM

VACUUM reclaims space from **dead tuples** (versions no longer visible to any snapshot) and performs critical maintenance:

1. **Removes dead tuples** and their index entries; marks freed line pointers reusable (space stays in the file — see bloat).
2. **Freezes** old tuples (anti-wraparound, below).
3. **Updates the visibility map (VM)** — pages where all tuples are visible to everyone. The VM enables **index-only scans** and lets later vacuums skip all-visible pages.
4. **Updates the free space map (FSM)** so inserts reuse freed space.
5. `VACUUM (ANALYZE)` also refreshes planner statistics.

`VACUUM` (lazy) runs online and does **not** return space to the OS — it makes space reusable *within* the table. `VACUUM FULL` rewrites the table compactly and **returns space to the OS** but takes an `ACCESS EXCLUSIVE` lock (table offline) and needs free disk equal to the table size. For online compaction use **`pg_repack`** (extension) or **`pg_squeeze`**.

### Autovacuum

The `autovacuum` launcher spawns workers when dead-tuple/insert thresholds are crossed:

```
vacuum threshold = autovacuum_vacuum_threshold (50)
                 + autovacuum_vacuum_scale_factor (0.2) * reltuples
```

The default `0.2` scale factor means a 100M-row table waits for **20M dead tuples** before vacuuming — far too lax for large tables. **Lower the scale factor per-table** for big/hot tables:

```sql
ALTER TABLE events SET (
  autovacuum_vacuum_scale_factor = 0.02,
  autovacuum_vacuum_insert_scale_factor = 0.02,  -- PG13+: vacuum insert-heavy tables for VM/freeze
  autovacuum_vacuum_cost_limit = 2000            -- let it work faster
);
```

Key knobs: `autovacuum_max_workers`, `autovacuum_vacuum_cost_delay` (lowered to `2ms` default in PG12+), `autovacuum_vacuum_cost_limit`. Symptoms of under-vacuuming: rising table/index size with flat row count, degrading index-only-scan ratio, `n_dead_tup` climbing in `pg_stat_user_tables`.

## Transaction-ID wraparound & freezing

XIDs are **32-bit** and wrap at ~4.2 billion. PostgreSQL treats the XID space as a circle: a tuple older than ~2.1B transactions would appear to be *in the future* and vanish. To prevent this, VACUUM **freezes** old tuples (marks them committed-and-visible-to-all via a frozen flag), removing them from XID-age accounting.

- Monitor `age(datfrozenxid)` per database and `age(relfrozenxid)` per table against `autovacuum_freeze_max_age` (default 200M).
- At ~200M unfrozen age, **anti-wraparound autovacuum** triggers automatically and **cannot be skipped** (it runs even if autovacuum is "off"); it takes stronger locks and can surprise you.
- At ~3M XIDs remaining the database enters protective shutdown to single-user mode. **This is an outage.** Catch it early via the `age()` queries; never disable autovacuum globally.
- PG18 and modern releases improve freezing (eager/opportunistic freezing, page-level freeze) reducing wraparound risk, but the model is unchanged.

## Bloat: diagnosis & remediation

**Bloat** = space occupied by dead tuples + free space that VACUUM made reusable but didn't return. Some steady-state bloat is *normal and healthy* (it gets reused). Pathological bloat means dead tuples accumulate faster than vacuum reclaims them.

Diagnose:
- `pg_stat_user_tables`: `n_dead_tup`, `n_live_tup`, `last_autovacuum`.
- Extension **`pgstattuple`** for exact dead-tuple %; community bloat-estimate queries (ioguix) for a fast approximation.
- Index bloat: compare index size to a freshly rebuilt copy; `REINDEX INDEX CONCURRENTLY` (PG12+) rebuilds online.

Common bloat causes & fixes:

| Cause | Fix |
|---|---|
| Long-running / idle-in-transaction sessions holding xmin | `idle_in_transaction_session_timeout`; kill offenders; check `pg_stat_activity.backend_xmin` |
| Stale/abandoned replication slot pinning xmin | Drop unused slots; alert on `pg_replication_slots` lag; cap with `max_slot_wal_keep_size` |
| `hot_standby_feedback=on` on a busy replica | Accept the trade-off or tune; it deliberately holds xmin on the primary |
| Autovacuum too lax for table size | Lower per-table `autovacuum_vacuum_scale_factor`; raise cost limit |
| Mass `UPDATE`/`DELETE` batch | Chunk it; vacuum between batches; or `pg_repack` after |
| Heavy non-HOT updates | Add `fillfactor` headroom; stop indexing hot columns |

> **Mental model:** PostgreSQL trades write-in-place simplicity for lock-free reads. The bill comes due as vacuum work. Healthy operation = autovacuum keeping pace with dead-tuple generation and freezing staying well ahead of wraparound.
