# HNBP-CORE

**A minimal, tamper-evident, append-only witness log for any actor (human, AI, or system), with optional Bitcoin anchoring.**

[![License: PolyForm Noncommercial 1.0.0](https://img.shields.io/badge/license-PolyForm%20Noncommercial-blue)](LICENSE)

---

## What Is This?

Modern systems lack verifiable memory. HNBP-CORE provides a cryptographic witness log — a hash chain where every record commits to all previous records. Any modification is instantly detectable. With Bitcoin anchoring, the log's existence at a specific moment in time becomes globally provable.

**Bitcoin timestamps the log's existence. HNBP-CORE proves it hasn't changed.**

---

## Features

- Tamper-evident append-only log (SHA-256 hash chain)
- Actor-agnostic: human, AI, system, or hybrid actors
- Optional Bitcoin anchoring via `OP_RETURN`
- Offline-verifiable — no server, no trust required
- Zero-dependency reference implementations in Python, JS, Go, and Rust
- Single-file web verifier (drag & drop, WebCrypto)
- Frozen v1.0.0 spec — logs written today are verifiable forever
- PolyForm Noncommercial license; commercial license available

---

## Protocol Overview

### Record Fields

| Field | Type | Description |
|---|---|---|
| `index` | integer | Position starting at 0, strictly sequential |
| `timestamp` | string | ISO 8601 |
| `data` | any | Any JSON-serializable value |
| `prev_hash` | string/null | Hash of previous record; `null` for record 0 |
| `hash` | string | SHA-256 hex digest of this record |

### Hashing Rule

```
hash = SHA256(
  String(index) +
  String(timestamp) +
  JSON.stringify(data, sorted_keys) +
  String(prev_hash)   // "null" for index 0
)
```

### Validity Conditions

1. Non-empty JSON array
2. Indices sequential: `0..n`
3. `records[0].prev_hash === null`
4. `records[i].prev_hash === records[i-1].hash` for `i > 0`
5. Recomputed hash matches stored hash for every record
6. All timestamps are valid ISO 8601

---

## Bitcoin Anchoring

**Anchor:**
1. Compute `HEAD = records[last].hash`
2. Create a Bitcoin transaction with `OP_RETURN <HEAD>` (32 bytes)
3. Broadcast via any Bitcoin wallet
4. Record the `txid`

**Verify:**
1. Validate the log (chain integrity)
2. Recompute HEAD
3. Extract `OP_RETURN` from the Bitcoin transaction
4. If `OP_RETURN == HEAD` → log is provably unmodified since that block

---

## Universal Actor Model

| Actor Type | Examples |
|---|---|
| Human Node | Person using a device, wallet, or companion AI |
| AI Node | LLM agent, workflow agent, autonomous system |
| System Node | Server, IoT device, pipeline, script |
| Hybrid Node | Human + AI pair acting together |

HNBP-CORE is actor-agnostic. The `data` field carries actor-specific semantics.

---

## Quickstart

### Run the Web Verifier

Open `verifiers/web/index.html` in any browser. No build step. No server.

### Run Tests

```bash
# Python
python3 tests/python_test.py

# JavaScript
node tests/js_test.mjs

# Go
cd tests && go test .

# Rust
rustc tests/rust_test.rs -o hnbp_test && ./hnbp_test
```

### Create and Verify a Log (Python)

```python
from reference.python.hnbp import append, validate_log, head_hash, save_log

log = []
append(log, {"event": "genesis", "actor": "alice"})
append(log, {"event": "signed", "doc": "contract-001"})
validate_log(log)
print("HEAD:", head_hash(log))
save_log(log, "mylog.morus")
```

### Create and Verify a Log (JavaScript)

```javascript
import { append, validateLog, headHash, saveLog } from './reference/javascript/hnbp.js';

const log = [];
append(log, { event: 'genesis', actor: 'alice' });
append(log, { event: 'signed', doc: 'contract-001' });
validateLog(log);
console.log('HEAD:', headHash(log));
saveLog(log, 'mylog.morus');
```

---

## Repo Structure

```
/spec/              Frozen protocol specification
/reference/         Zero-dependency implementations (Python, JS, Go, Rust)
/tests/             Valid/invalid/bitcoin test fixtures + runners
/verifiers/web/     Single-file HTML verifier
/examples/          Real-world scenarios (human, ai, system, hybrid, multisig, otc, custody, governance)
/docs/              Extended documentation
LICENSE             PolyForm Noncommercial 1.0.0
GOVERNANCE.md       Protocol freeze policy and licensing terms
```

---

## Examples

| Scenario | Description |
|---|---|
| [Human](examples/human/) | Document signing session |
| [AI](examples/ai/) | Agent audit trail |
| [System](examples/system/) | Payment processing log |
| [Hybrid](examples/hybrid/) | Human + AI decision log |
| [Multisig](examples/multisig/) | Key ceremony witness |
| [OTC](examples/otc/) | Bitcoin trade deal flow |
| [Custody](examples/custody/) | Proof of reserves chain |
| [Governance](examples/governance/) | Voting and proposal log |

---

## Reference Implementations

All implementations are zero-dependency and ≤100 lines:

- `reference/python/hnbp.py`
- `reference/javascript/hnbp.js`
- `reference/go/hnbp.go`
- `reference/rust/lib.rs`

---

## Governance

HNBP-CORE v1.0.0 is **frozen**. The format does not change. See [GOVERNANCE.md](GOVERNANCE.md).

---

## License

**PolyForm Noncommercial License 1.0.0** — noncommercial use is allowed.  
Commercial use requires a separate commercial license from the author.  
© 2026 BobbiéB (github.com/Bardockthegreat). All commercial rights reserved.

See [LICENSE](LICENSE) and [GOVERNANCE.md](GOVERNANCE.md) for details.

---

## Verified Actor Badge

```
✔ HNBP-Verified Actor
This actor maintains a tamper-evident witness log anchored to Bitcoin.
```

Use this badge only when you maintain a valid, verifiable HNBP log. See GOVERNANCE.md for rules.
