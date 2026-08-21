<!-- hub-reference-banner -->
> **Reference file — part of the `lang-js-ts` hub.** Formerly the standalone `nodejs-native-addons-napi` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

<!--
PROVENANCE: Authored by /dr (deep-research-and-build) on 2026-05-31.
HUB: programming-languages (reference spoke). NOT a standalone top-level skill.
SCOPE: Node.js native addons — the C ABI-stable Node-API (N-API) foundation, the
node-addon-api C++ wrapper layer on top of it, the node-gyp build toolchain
(binding.gyp / gyp-next), distributing prebuilt binaries (prebuildify + node-gyp-build,
prebuild/prebuild-install, node-pre-gyp), cmake-js as a node-gyp alternative, the
Rust alternatives (napi-rs, neon), loading .node files (require / process.dlopen /
--experimental-addon-modules), and debugging native addons. Node-API is the shared
foundational branch; node-addon-api and node-gyp build ON it and cross-reference it.
Cross-references sibling references: nodejs-concurrency-internals.md (libuv event loop
and thread pool — async work runs there), javascript-nodejs.md (Node runtime/module
system), nodejs-typescript-and-runtime-features.md (process model / SEA).
SOURCES: Node.js official docs — Node-API (nodejs.org/api/n-api.html, v26.2.0),
C++ Addons (nodejs.org/api/addons.html); nodejs/node-gyp GitHub repo + npm page;
nodejs/node-addon-api GitHub docs (object_wrap.md, threadsafe_function.md, async_worker);
The Node-API Resource / node-addon-examples (nodejs.github.io/node-addon-examples,
build-tools: prebuild, node-pre-gyp); prebuild/prebuildify + prebuild/node-gyp-build
GitHub repos; napi.rs (Announcing NAPI-RS v2) and napi-rs/napi-rs GitHub; neon-bindings/neon
GitHub; LogRocket "Solving common issues with node-gyp"; OpenReplay node-gyp troubleshooting.
-->

# Node.js Native Addons — Node-API (N-API), node-addon-api & node-gyp

A `programming-languages` hub reference for building, building-against, and shipping
**C/C++/Rust native addons** for Node.js. A Node.js addon is a *dynamically-linked
shared object loaded via `require()` like an ordinary module*; the compiled binary has
a **`.node`** extension ([Node.js C++ Addons](https://nodejs.org/api/addons.html)).

The mental model is a three-layer stack, foundation first:

1. **Node-API (N-API)** — the C ABI-stable foundation. Everything else sits on it.
2. **node-addon-api** — header-only **C++** sugar over Node-API (or napi-rs / neon for Rust).
3. **node-gyp** (or cmake-js) — the *build* toolchain that compiles the addon into a `.node`.

Plus a fourth concern orthogonal to all three: **distribution** of prebuilt binaries so
end users don't need a compiler. For the libuv event loop / thread pool that async addon
work runs on, defer to `nodejs-concurrency-internals.md`; for the module system and Node
runtime APIs, defer to `javascript-nodejs.md`.

---

## 1. Node-API (N-API) — the ABI-stable C foundation

Node-API is "an API for building native Addons... independent from the underlying
JavaScript runtime (for example, V8) and... maintained as part of Node.js itself. This
API will be **Application Binary Interface (ABI) stable** across versions of Node.js"
([n-api.html](https://nodejs.org/api/n-api.html)). It insulates addons from V8 changes
and lets a module **compiled for one major version run on later major versions without
recompilation** — the single biggest reason to choose it over raw V8.

**The ABI bargain:** to get that stability the addon must include only `<node_api.h>`
and must NOT include `<node.h>`, `<v8.h>`, `<uv.h>`, or `<node_buffer.h>`
([n-api.html](https://nodejs.org/api/n-api.html)). Touching V8/libuv directly re-couples
you to the engine and forfeits the guarantee.

### Versioning (the matrix and `NAPI_VERSION`)

- Current **maximum Node-API version is 10** (Node v22.14.0+, v23.6.0+); v9 landed in
  v18.17.0+/v20.3.0+; v8 is the broad baseline (v12.22+, v14.17+, v16+)
  ([n-api.html](https://nodejs.org/api/n-api.html)).
- **Default is version 8 when `NAPI_VERSION` is unset.** Define it to opt into more:
  ```c
  #define NAPI_VERSION 9
  #include <node_api.h>
  ```
  This "bakes in" the requested version used at runtime
  ([n-api.html](https://nodejs.org/api/n-api.html)).
- Pre-v9, versions were strictly additive (each a superset). **As of v9, an add-on
  written for v9 may require code updates for v10** — versioning is now independent, not
  purely additive. Runtimes supporting v9+ still support every version from 8 up to their
  max ([n-api.html](https://nodejs.org/api/n-api.html)).
- The **`node-api-headers`** package provides the headers for building against Node-API
  without the full Node source tree.

### Core opaque types

| Type | Meaning |
|---|---|
| `napi_env` | Context for VM-specific state, passed into every native function and back into every Node-API call. **Must not be cached for reuse or shared across `Worker` threads**; invalid once the addon instance unloads ([n-api.html](https://nodejs.org/api/n-api.html)). |
| `napi_value` | Opaque pointer representing a JS value ([n-api.html](https://nodejs.org/api/n-api.html)). |
| `napi_ref` | A reference that lets you manage a value's minimum lifetime explicitly. |
| `napi_handle_scope` / `napi_escapable_handle_scope` | Control the GC lifetime of `napi_value`s. |

### Handle scopes

A handle scope controls the lifetime of values created within it. A **default handle
scope exists when a native method is called from JS**, but **outside native-method
execution (e.g. inside a libuv callback) you must open a scope** before creating any JS
values ([n-api.html](https://nodejs.org/api/n-api.html)):

- `napi_open_handle_scope` / `napi_close_handle_scope`
- Escapable variant — `napi_open_escapable_handle_scope` / `napi_escape_handle` /
  `napi_close_escapable_handle_scope` — to promote one value into the parent scope.

Closing a scope tells the GC its values are no longer referenced from the current stack frame.

### References and reference counting

`napi_ref` keeps an object/function/symbol/external alive past its handle scope. Each
reference has a count ≥ 0; a count of 0 lets the value be collected
([n-api.html](https://nodejs.org/api/n-api.html)):

- `napi_create_reference` / `napi_delete_reference`
- `napi_reference_ref` (increment, pin) / `napi_reference_unref` (decrement, allow GC)
- `napi_get_reference_value`

### Error handling

Every Node-API function returns a `napi_status`. `napi_ok` = success; `napi_pending_exception`
= a JS exception is queued; anything else is an error
([n-api.html](https://nodejs.org/api/n-api.html)).

- After a non-`napi_ok`/non-`napi_pending_exception` status, call
  **`napi_is_exception_pending`** before continuing, rather than blindly returning.
- `napi_get_last_error_info` → `napi_extended_error_info` (valid only until the next
  Node-API call on the same `env`).
- Throw: `napi_throw`, `napi_throw_error`, `napi_throw_type_error`,
  `napi_throw_range_error`, `node_api_throw_syntax_error`.
- Inspect/clear: `napi_is_exception_pending`, `napi_get_and_clear_last_exception`,
  `napi_fatal_exception`.

### Finalizers

- `napi_add_finalizer(env, js_object, native_object, finalize_cb, finalize_hint, result)`
  — run cleanup when the JS object is GC'd.
- `node_api_basic_finalize` (experimental, v18.20+/v20.12+/v21.6+) — synchronous cleanup
  callable while the engine cannot run JS; only Node-APIs taking a `node_api_basic_env`
  may be used inside it. For external strings the `env` may be `null` during shutdown
  ([n-api.html](https://nodejs.org/api/n-api.html)).
- `node_api_post_finalizer` — schedule JS-touching work for after GC completes.

### Async work (libuv thread pool)

For CPU-bound work that must not block the event loop, use the async-work pattern. The
**execute callback runs on the libuv thread pool** (no JS calls allowed there); the
**complete callback runs back on the main thread** ([n-api.html](https://nodejs.org/api/n-api.html)).
See `nodejs-concurrency-internals.md` for the libuv pool sizing (`UV_THREADPOOL_SIZE`).

```c
napi_create_async_work(env, async_resource, async_resource_name,
                       execute /* libuv thread */, complete /* main thread */,
                       data, &work);
napi_queue_async_work(env, work);
// also: napi_cancel_async_work, napi_delete_async_work
```

### Thread-safe functions (calling JS from arbitrary threads)

`napi_threadsafe_function` (TSFN) is the *only* sanctioned way for a non-JS thread to
invoke a JS function ([n-api.html](https://nodejs.org/api/n-api.html)):

- `napi_create_threadsafe_function` (takes `max_queue_size`, `initial_thread_count`,
  a finalize cb, a context, and the `call_js_cb` that runs on the **main thread**).
- `napi_call_threadsafe_function(func, data, is_blocking)` —
  **`napi_tsfn_blocking`** waits when the queue is full; **`napi_tsfn_nonblocking`**
  returns `napi_queue_full` immediately.
- `napi_acquire_threadsafe_function` / `napi_release_threadsafe_function`
  (release mode `napi_tsfn_release` vs `napi_tsfn_abort`). The TSFN is destroyed once
  every using thread has released or received `napi_closing`.

### Registration

Addons register via the init macro. With Node-API:
```c
NAPI_MODULE_INIT(/* napi_env env, napi_value exports */) {
  return create_addon(env);
}
```
Context-aware C++ V8 addons use `NODE_MODULE_INIT()` / `NODE_MODULE_INITIALIZER`
(supplying `exports`, `module`, `context`). **Plain `NODE_MODULE()` addons cannot load
in multiple contexts/threads** — for `worker_threads`/Electron you need a Node-API addon
or a context-aware one, and `AddEnvironmentCleanupHook()` to release resources per thread
([addons.html](https://nodejs.org/api/addons.html)).

---

## 2. node-addon-api — the C++ wrapper

`node-addon-api` is a **header-only C++ wrapper** over the C Node-API. It hides raw
`napi_status` checks behind C++ exceptions/objects and RAII scopes. Key classes
([nodejs/node-addon-api docs](https://github.com/nodejs/node-addon-api)):

- `Napi::Env`, `Napi::Value`, `Napi::Object`, `Napi::Function`, `Napi::String`,
  `Napi::Number`, `Napi::Array`, `Napi::Buffer<T>`.
- **`Napi::ObjectWrap<T>`** — binds a C++ object's lifetime to a JS object; each JS
  instance creates a matching C++ instance
  ([object_wrap.md](https://github.com/nodejs/node-addon-api/blob/main/doc/object_wrap.md)).
- **`Napi::AsyncWorker`** — override `Execute()` (runs on a worker thread off the event
  loop) and `OnOK()` / `OnError()` (run on the main thread); `SetError()` routes to
  `OnError`. This is the C++ ergonomic form of `napi_create_async_work`.
- **`Napi::ThreadSafeFunction`** / `Napi::TypedThreadSafeFunction` — `New()` holds a
  persistent ref to a JS function callable from multiple threads asynchronously; threads
  call `BlockingCall()`/`NonBlockingCall()` and must `Release()` (or receive
  `napi_closing`) to tear it down ([threadsafe_function.md](https://github.com/nodejs/node-addon-api/blob/main/doc/threadsafe_function.md)).
- `Napi::HandleScope` / `Napi::EscapableHandleScope` — RAII over the C scopes.

**Exception-handling modes** — node-addon-api compiles in one of two modes, selected by
defining (or not) `NAPI_DISABLE_CPP_EXCEPTIONS` / `NAPI_CPP_EXCEPTIONS`:
- **C++ exceptions ON** — Node-API errors throw `Napi::Error`; idiomatic try/catch.
- **C++ exceptions OFF** — you must check `env.IsExceptionPending()` and return early
  manually (required where the toolchain forbids exceptions).
- **Known gotcha:** with a `ThreadSafeFunction`, a JS exception in the callback is thrown
  on the C++ side but does **not** surface as a normal JS throw the way a synchronous
  callback would ([node-addon-api issue #669](https://github.com/nodejs/node-addon-api/issues/669)).

---

## 3. node-gyp — the build toolchain

`node-gyp` is "a cross-platform command-line tool written in Node.js for compiling native
addon modules," vendoring a copy of **gyp-next** (the Chromium team's GYP fork) extended
for Node ([nodejs/node-gyp](https://github.com/nodejs/node-gyp)). npm shells out to it for
any package with a `binding.gyp`.

### `binding.gyp`

A JSON-like file describing build targets ([addons.html](https://nodejs.org/api/addons.html)):
```json
{
  "targets": [
    { "target_name": "addon", "sources": [ "hello.cc" ] }
  ]
}
```
Common keys: `target_name`, `sources`, `include_dirs`, `dependencies`, `defines`,
`cflags`/`cflags_cc`, `libraries`, and conditional `conditions` blocks per-OS. For
node-addon-api you add its include dir and (if header-only with napi version) set
`NAPI_VERSION` via `defines`.

### Commands

`configure` (generate Makefile/`.vcxproj`), `build` (run make/msbuild), `rebuild`
(clean+configure+build), `clean`, `install`/`list`/`remove` (Node header management).
Typical: `node-gyp configure build` → output lands in **`build/Release/addon.node`**,
loaded via `require('./build/Release/addon')` ([addons.html](https://nodejs.org/api/addons.html)).

### Prerequisites (the #1 source of pain)

| Platform | Needs |
|---|---|
| All | A supported **Python 3** (Python ≥ 3.12 requires node-gyp ≥ v10) ([node-gyp](https://github.com/nodejs/node-gyp)) |
| Linux/Unix | `make` + a C/C++ toolchain (GCC) |
| macOS | **Xcode Command Line Tools** (`xcode-select --install`) → clang/clang++/make |
| Windows | **Visual C++ Build Tools** (VS 2019+ "Desktop development with C++" workload); ARM64 needs the ARM64 VC++ compilers + ATL |

Config can be passed three ways: CLI flag (`--python /path/to/python`), `package.json`
`"config"` block, or `npm_package_config_*` env vars. Headers download to a devdir
(`node_gyp_devdir`); `--nodedir` points at a full source tree.

---

## 4. Distributing prebuilt binaries (so users skip the compiler)

End users without a C++ toolchain can't run node-gyp, so ship prebuilds
([The Node-API Resource](https://nodejs.github.io/node-addon-examples/build-tools/prebuild/)):

- **prebuildify + node-gyp-build (recommended).** `prebuildify` bundles all prebuilt
  binaries *inside the published npm package* (the `prebuilds/` folder) — no download step.
  `node-gyp-build` is the install-script/loader that picks the prebuild if present and
  only rebuilds as a fallback ([prebuild/node-gyp-build](https://github.com/prebuild/node-gyp-build)).
- **prebuild + prebuild-install.** Author runs `prebuild` (uploads to GitHub Releases by
  default); users' `prebuild-install` downloads the matching binary before falling back to
  source build. Supports both node-gyp and cmake-js backends. (The standalone `prebuild`
  CLI is now deprecated in favor of the prebuildify path.)
- **node-pre-gyp.** Older alternative that hosts prebuilt binaries (often on S3) and
  downloads them at install ([node-pre-gyp](https://nodejs.github.io/node-addon-examples/build-tools/node-pre-gyp/)).

Prebuilds are keyed by platform + arch + ABI; Node-API's ABI stability is what makes a
single prebuild work across Node majors (vs. one-per-Node-version with raw V8 addons).

---

## 5. Alternatives to node-gyp / C++

- **cmake-js** — drives a CMake-based build instead of GYP; preferred when the native lib
  already ships CMake or you want richer build logic than GYP expresses.
- **Rust via napi-rs** — a framework for compiled Node.js addons in Rust over Node-API.
  The `#[napi]` macro hides value casting and **auto-generates the `.js` bindings + `.d.ts`
  TypeScript definitions**; it targets the stable Node-API for cross-version compatibility
  ([napi.rs v2](https://napi.rs/blog/announce-v2)).
- **Rust via neon** — Rust bindings for native modules; historically a mix of N-API and V8
  backends, mid-transition to Node-API ([neon-bindings/neon](https://github.com/neon-bindings/neon)).
  Choose napi-rs for ABI-stable, TS-typed modules; neon for its simpler direct-integration model.
- **nan** — legacy "Native Abstractions for Node.js" for V8-era addons; new work should use
  Node-API instead ([addons.html](https://nodejs.org/api/addons.html)).

Rust buys language-level memory safety over hand-written C/C++ FFI; writing a raw addon
against the C Node-API from Rust without a helper means a lot of `unsafe`.

---

## 6. Loading & debugging addons

- **Load:** `require('./build/Release/addon')` (the `.node` extension is usually optional,
  but `.js` of the same basename wins). Under the hood this is `process.dlopen`. ESM import
  of a `.node` works under **`--experimental-addon-modules`** (v23.6.0/v22.20.0+)
  ([addons.html](https://nodejs.org/api/addons.html)).
- **Debug native code:** attach `gdb`/`lldb` to the `node` process; build with debug info
  (`node-gyp build --debug` → `build/Debug/`). `node --inspect` debugs the JS layer; for
  the C/C++ layer use a native debugger or a mixed-mode debugger. Symbol-load the `.node`
  shared object in the debugger to set breakpoints in addon source.
- **Common node-gyp failures** ([LogRocket](https://blog.logrocket.com/solving-common-issues-node-gyp/),
  [OpenReplay](https://blog.openreplay.com/node-gyp-troubleshooting-guide-fix-common-installation-build-errors/)):
  *"can't find Python"* (set `--python` or `PYTHON`), *"gyp failed with exit code 1"*
  (missing/incompatible VS Build Tools or Xcode CLT), `binding.gyp not found` (running in
  the wrong dir / package didn't ship it), and Node/Python version mismatches. Containerize
  the build (Docker) for reproducibility.

---

## Anti-patterns

- **Including `<v8.h>`/`<node.h>` in a Node-API addon** — silently forfeits ABI stability;
  the addon now breaks across Node majors.
- **Caching `napi_env`** or passing it between Worker threads — explicitly disallowed; it
  goes invalid on unload.
- **Calling JS-creating Node-API from a libuv/worker thread without a handle scope** — or
  worse, calling JS at all from the async *execute* callback (use a TSFN instead).
- **Calling a JS function directly from a native thread** instead of via
  `napi_threadsafe_function` — undefined behavior.
- **Shipping a source-only native package** with no prebuilds — every `npm install`
  forces a full toolchain build and breaks on machines without a compiler.
- **Choosing raw V8/nan for new addons** — start at Node-API; reach for V8 only for the
  rare API Node-API doesn't expose.

---

## Cross-references

- `nodejs-concurrency-internals.md` — libuv event loop phases + thread pool that async
  work / TSFNs dispatch onto (`UV_THREADPOOL_SIZE`), `worker_threads`.
- `javascript-nodejs.md` — module system, `require`/ESM, Node runtime APIs.
- `nodejs-typescript-and-runtime-features.md` — process model, SEA, permission model.
- For generic C++/Rust patterns and FFI safety, see the `software-engineering-patterns`
  hub and the Rust ecosystem references.
