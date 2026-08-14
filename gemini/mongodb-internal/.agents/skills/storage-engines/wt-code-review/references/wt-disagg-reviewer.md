# WiredTiger disaggregated-storage review

Correctness defects specific to disagg mode — bugs classic WT would not exhibit, or bugs from disagg's constraints (leader/follower asymmetry, remote page storage, layered tables, precise checkpoints, deltas, materialization frontier, shared metadata). Assumes the reader knows what disagg storage is; focuses on what a reviewer must *check*.

## Scope

Engage on any of these diff signals:

- **Paths**: `src/block_disagg/`, `src/conn/conn_layered*`, `src/conn/conn_layered_page_log.c`, `src/conn/conn_layered_ingest.c`, `src/conn/conn_layered_table_manager.c`.
- **Symbols / structs**: `WT_DISAGGREGATED_STORAGE`, `WT_PAGE_LOG`, `WT_PAGE_LOG_HANDLE`, `WT_PAGE_BLOCK_META`, `WT_BLOCK_DISAGG`, `WT_BLOCK_DISAGG_ADDRESS_COOKIE`, `WT_DISAGG_CHECKPOINT_META`, `page->disagg_info`, `layered_table_manager`, `layered_drain_data`, `WT_LAYERED_TABLE*`.
- **Naming**: `__wt_disagg_*`, `__wti_disagg_*`, `__wti_block_disagg_*`, `__wti_layered_*`, `__clayered_*`, `WT_DISAGG_*`.
- **Flags / branches**: `WT_BTREE_DISAGGREGATED`, `WT_CONN_PRECISE_CHECKPOINT`, `WT_CONN_PRESERVE_PREPARED`, `WT_DHANDLE_DISAGG_META`, `WT_REF_FLAG_LEAF` checks in eviction/discard.
- **Config keys**: `disaggregated=(...)`, `block_manager=disagg`, `type=layered`, `precise_checkpoint=true`, `preserve_prepared=true`, `page_delta=(...)`, `page_log=...`, `role=leader|follower`.
- **Helpers**: `__wt_conn_is_disagg`, `layered_table_manager.leader`, `__wt_disagg_enqueue_metadata_operation`, `__wt_materialization_check`, `__wt_btree_can_discard`.
- **Verbose / stats**: `WT_VERB_DISAGGREGATED_STORAGE`, `disagg_*` stats.

If none appear and no behaviour is conditional on disagg mode, the diff is out of scope — say so and stop. If the diff also affects classic mode, note it — disagg and classic coexist (mixed-mode), and divergence between them is a frequent bug source.

## Reference

Verify findings against current source in `src/include/` (`disaggregated.h`, `block.h`, `block_disagg.h`, `btree.h`, `layered.h`, `meta.h`). Starting greps:
- Leader guards: `grep -rnE 'layered_table_manager\.leader|WT_BTREE_DISAGGREGATED' src/`
- PALI callbacks: `grep -rnE 'plh_(put|get|discard|get_page_ids)' src/`
- Address cookie pack/unpack: `grep -rnE '__wti?_block_disagg_addr' src/`
- Metadata enqueue: `grep -rnE '__wt_disagg_enqueue_metadata_operation' src/`
- Frontier / discard: `grep -rnE '(materialization_check|btree_can_discard|last_materialized_lsn)' src/`
- Precise checkpoint flag: `grep -rnE 'WT_CONN_PRECISE_CHECKPOINT|WT_CONN_PRESERVE_PREPARED' src/`

## What to flag

### 1. Leader / follower asymmetry
- **Leader-only writes.** Any new mutation, reconciliation, or `plh_put` on a disagg-eligible btree needs a leader guard; compare sibling writes in `src/cursor/cur_layered.c`.
- **Removed guards.** A refactor that deletes, weakens (`&&` → `||`), or reorders past a side-effectful call is almost always wrong.
- **Checkpoint advance.** Only the leader advances checkpoints and writes shared metadata; follower-reachable callers of the checkpoint or metadata-queue path are bugs.
- **Role transitions.** Step-up / step-down: checkpoint lock held, eviction/checkpoint cannot race the flip, L1 drains into L0 before metadata reopen on step-up, btrees readonly before step-down.
- **Cursor role awareness.** Layered cursors must reopen the stable handle on checkpoint advance or role change (stale-read bug otherwise). Stable-upgrade is unsafe mid-iter without `read_timestamp` set; `WT_CLAYERED_RANDOM` cursors must never upgrade — internal sampling state is lost on reopen.
- **Error asymmetry.** Stable-table `ENOENT` on a follower is transient — remap to success + verbose log; the same `ENOENT` on a leader is an error. Followers can legally have `file:T.wt_ingest` without `file:T.wt_stable`; leaders must have both — "incomplete layered table" assertions must encode this asymmetry.
- **Layered modify finding an ingest tombstone returns `WT_NOTFOUND`** — do not fall back to the stable table just because ingest "didn't have a usable value".
- **On the stable side treat shared metadata as read-only at the checkpoint** — opening the live handle marks pages dirty, blocking disagg eviction.
- **Role plumbing.** Extract `__wt_disagg_config_get_role()` rather than parsing role in multiple subsystems; `layered_table_manager.leader` is unset during recovery — plumb role into a recovery-local field early.

### 2. Block manager and the Page-And-Log interface (PALI)
- **Address cookie invariants.** Cookie carries `page_id` (non-zero), `lsn >= base_lsn` (underflow on unpack is a bug), `size > 0`; a new packer/unpacker must round-trip with the existing one. Checkpoint cookie = root address cookie — code that synthesizes/rewrites/drops fields must preserve this.
- **PALI discard arguments.** `plh_discard` takes `backlink_lsn` / `base_lsn` from the unpacked address cookie, not `block_meta`: full pages send `base_lsn = 0`, deltas send `base_lsn = cookie.base_lsn`.
- **Delta vs base reconstruction.** Block manager verifies checksums and the `previous_checksum` chain; a mismatch must surface as `WT_CONN_DATA_CORRUPTION`, not be swallowed.
- **`plh_put` before discard.** A page cannot be `plh_discard`-ed before its `plh_put` is durable; async paths that let discard race put are corruption bugs.
- **Optional vs required PALI methods.** Missing required ones must error; missing optional ones (e.g. optional `plh_discard`) require a NULL-pointer skip with `__wt_verbose_warning` and `return 0` — returning a hard error when the page service hasn't shipped the op is wrong. Trim / SLS-discard must run *before* removing the corresponding metadata entries (reversed order orphans data on crash); skip on followers.
- **Magic numbers / version.** Adding a header field without bumping the compatible version (or vice versa) breaks old readers.

### 3. Layered tables (L0 / L1)
- **Read precedence.** Reads merge L1 and L0 with L1 taking precedence; bypassing L1, or reading a stale L0 checkpoint when the leader has the live stable open, is data loss.
- **L1 is in-memory only.** Never written to shared storage.
- **Drain on step-up.** L1 must drain into L0 (prepared transactions fixed up) *before* the leader exposes itself as writer.
- **Tombstone encoding.** Ingest tables use `__clayered_deleted_encode`; a plain `cursor->remove()` on an ingest table corrupts the format.
- **Snapshot isolation.** Layered modify requires snapshot isolation — reject read-committed / read-uncommitted with `ENOTSUP` at entry.
- **Buffer copies.** When copying key/value buffers between constituent cursors use `__wt_buf_set` deep copies, not `WT_ITEM_SET` pointer aliasing — shared buffers cause use-after-free.
- **Truncate entries outliving a dhandle must hold the URI string**, not a `WT_LAYERED_TABLE *` — the dhandle can be swept before commit.
- **Ingest-vs-stable timestamp rule** is "latest version in ingest >= latest version in stable", not "any version newer"; drain must also filter by `oldest_timestamp`.
- **Layered table manager.** Add/remove must be safe with respect to the manager's lock and per-table flags.

### 4. Precise checkpoints and `preserve_prepared`
- **No RTS on disagg restart.** Code relying on "RTS will fix it up later" is wrong for disagg.
- **Stable-only snapshot.** Reconciliation under precise checkpoint must not write unstable updates.
- **Prepared transactions.** With `preserve_prepared=true`, prepared updates must carry a non-`WT_PREPARED_ID_NONE` id across checkpoint and recovery. HS-loop termination must check `WT_CONN_PRESERVE_PREPARED` *and* distinguish `WT_PREPARE_INPROGRESS` vs `WT_PREPARE_LOCKED` before treating aborted prepared updates as "skip me".
- **Threshold adjustments.** Any new error path must restore thresholds; otherwise eviction stays in checkpoint-mode pressure forever.
- **Durability severities.** Missing ingest updates after a crash = data loss; missing stable = corruption — different severities, different recovery paths.
- **Mixed mode.** A single connection can host both disagg (precise) and classic (fuzzy) btrees; mode choice must be per-btree.

### 5. Eviction and the materialization frontier
- **`__wt_btree_can_discard` / `__wt_materialization_check`.** Any new discard path on a disagg btree must consult the frontier; a bypass — even a "fast path for clean pages" — is corruption. `__wt_page_materialization_check` must guard clean evictions too, using `old_rec_lsn_max` (not the current LSN) for the dirty path; scrub eviction legitimately bypasses the frontier.
- **Dirty internal pages.** Evicting a dirty internal page on a disagg btree is disallowed.
- **Synchronous put before discard.** Reconciliation must complete `plh_put` before eviction releases the page.
- **Frontier monotonicity.** `last_materialized_lsn` moves forward only.
- **Victim/server-side cache.** Only clean leaf pages belong there; stashing dirty/internal pages or serving without revalidation is a bug.

### 6. Page deltas
- **Chain bound.** `block_meta->delta_count` stays below `page_delta.max_consecutive_delta`; at the bound the next write must be a full page. Static arrays sized to `WT_DELTA_LIMIT` must stay consistent with the runtime config bound.
- **Base / backlink immutability.** Within a chain `base_lsn` is fixed; each delta's `backlink_lsn` points at the previous delta or base.
- **Delta policy per page type.** `WT_INTERNAL_PAGE_DELTA` / `WT_LEAF_PAGE_DELTA` gate separately.
- **Obsoleting a chain.** When a full page replaces a chain, the previous chain must be marked obsolete for GC.
- **Merge invariants.** Compute time-aggregates as cells are emitted from the merge, not in a post-merge scan — a post-pass can't catch build-logic bugs. Overflow cells (`WT_CELL_KEY_OVFL`, `WT_CELL_VALUE_OVFL`) must never appear in disagg B-tree merges (assert absence; data lives in the object store), and `WT_CELL_ADDR_DEL` (fast-truncated entries) must come only from the base image, never a delta (assert in page-merge code).

### 7. Shared metadata and garbage collection
- **No local turtle on disagg.** Metadata changes flow through `__wt_disagg_enqueue_metadata_operation` under the schema lock.
- **Dhandle / metadata flags.** Tables holding shared metadata are marked `WT_DHANDLE_DISAGG_META`; classic metadata rules (sweep, drop, free-block list) do not apply.
- **No local free-block lists.** Reclamation happens via discard records emitted to the phylog plus checkpoint-completion markers. DSC "checkpoint complete" means the SLS completion marker was sent, not that internal reconciliation finished.
- **Schema / drop paths.** A table drop must emit the disagg discard machinery, not just remove a local file.

### 8. History store
- **Read source per role.** Followers serve historical reads from the shared HS at the appropriate checkpoint; falling back to the local HS file returns potentially non-durable data. HS-on-disagg iteration must loop over multiple HS btrees via `__wt_curhs_next_hs_id` — assuming a single HS handle undercounts on layered tables.
- **HS checkpoint binding.** Changes to how `btree->hs_checkpoint_name` is propagated/cached/refreshed are high-risk.
- **HS lifecycle and sweep.** HS dhandles must remain available across operations that depend on them.

### 9. Encryption and the key provider
- **Key version match.** A payload encrypted with key version N must be decrypted with key version N.
- **Header validation.** The crypto header carries version, compatible-version, checksum; any path that loads keys must validate these.
- **Block-layer encryption flags.** Must stay in sync with what the phylog actually did.

### 10. Unsupported / restricted in disagg
- **Unsupported operations**: classic compaction, salvage on layered tables, local (unreplicated) user tables, OS-cache / mmap optimizations, in-place block rewrites, classic free-block reuse.
- **Stub block-manager methods.** New callers on a disagg btree must be guarded out or error explicitly; a path that "works" via a silent no-op is the most insidious failure.
- **Config validation.** Disagg config can be incompatible with `in_memory`, `log`, or local-only flags; watch for new unvalidated combinations.

### 11. Mixed mode (disagg + classic in one connection)
- **Per-btree branches.** Behavior must branch per btree (`WT_BTREE_DISAGGREGATED`) and per connection state (`__wt_conn_is_disagg`).
- **Shared subsystems.** Eviction, sweep, schema, statistics, prefetch, recovery must keep the two paths separate where they diverge and equivalent where they don't.
