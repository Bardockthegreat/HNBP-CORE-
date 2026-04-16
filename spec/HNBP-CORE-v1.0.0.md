# HNBP-CORE v1.0.0 — Frozen Specification

**Status:** FROZEN — this document does not change.  
**Date:** 2026-04-15  
**Format:** `.morus` (JSON array)

---

## 1. Scope

HNBP-CORE defines a universal, tamper-evident, append-only witness log for any actor — human, AI, or system — with optional anchoring to the Bitcoin blockchain for global immutability.

This specification covers:
- Record structure
- Hashing rule
- Log validity conditions
- Bitcoin anchoring procedure
- Actor model
- Determinism requirements
- Versioning policy
- Checksums

---

## 2. Record Structure

Each record in a HNBP log is a JSON object with exactly these fields:

| Field | Type | Description |
|---|---|---|
| `index` | integer | Position in the log, starting at 0, strictly sequential |
| `timestamp` | string | ISO 8601 datetime string (e.g. `"2026-04-15T07:06:00Z"`) |
| `data` | any | Any JSON-serializable value (object, string, number, array, etc.) |
| `prev_hash` | string or null | Hash of the previous record; `null` only for `index 0` |
| `hash` | string | SHA-256 hex digest of this record's canonical input |

A log file is a JSON array of records:

```json
[
  {
    "index": 0,
    "timestamp": "2026-04-15T07:06:00Z",
    "data": { "event": "genesis" },
    "prev_hash": null,
    "hash": "<sha256-hex>"
  },
  {
    "index": 1,
    "timestamp": "2026-04-15T08:00:00Z",
    "data": { "event": "signed", "doc": "contract-001" },
    "prev_hash": "<hash-of-record-0>",
    "hash": "<sha256-hex>"
  }
]
```

Log files use the `.morus` extension (e.g., `mylog.morus`). The content is valid JSON.

---

## 3. Hashing

The hash for each record is computed as:

```
hash = SHA256(
    String(index) +
    String(timestamp) +
    JSON.stringify(data, canonical) +
    String(prev_hash)
)
```

Rules:
- All components are converted to strings before concatenation.
- `index` → `"0"`, `"1"`, etc.
- `timestamp` → used as-is (the stored string value).
- `data` → **canonical JSON**: keys sorted alphabetically, no extra whitespace. Specifically: `JSON.stringify` with sorted keys and no spaces.
- `prev_hash` → the stored hex string, **or the literal string `"null"` for `index 0`** (not the JSON null value — the four-character string `null`).
- The SHA-256 is computed over the UTF-8 encoding of the concatenated string.
- Output is lowercase hex.

### Canonical JSON for `data`

For determinism across implementations, `data` MUST be serialized with:
- Keys sorted alphabetically at every level (recursive).
- No whitespace (no spaces or newlines).
- Equivalent to Python's `json.dumps(data, sort_keys=True, separators=(',', ':'))`.

---

## 4. Log Validity

A log is valid if and only if ALL of the following hold:

1. **Non-empty**: The log is a JSON array with at least one record.
2. **Sequential indices**: `records[i].index == i` for all `i`.
3. **Genesis prev_hash**: `records[0].prev_hash === null` (JSON null).
4. **Chain linkage**: For all `i > 0`, `records[i].prev_hash === records[i-1].hash`.
5. **Hash integrity**: For every record, the stored `hash` matches the recomputed hash per Section 3.
6. **Valid timestamps**: Every `timestamp` is a valid ISO 8601 string.

Any violation of any rule means the log is invalid.

---

## 5. Bitcoin Anchoring

### 5.1 Purpose

Anchoring the HEAD hash to Bitcoin provides:
- A globally verifiable timestamp (block height)
- Immutability backed by Bitcoin's proof-of-work
- Censorship resistance
- Independent auditability without trusting any server

### 5.2 Anchoring Procedure

1. Compute the current HEAD hash: the `hash` field of the last record in the log.
2. Create a Bitcoin transaction with an `OP_RETURN` output containing the HEAD hash as raw bytes (32 bytes, decoded from hex).
3. Broadcast the transaction via any Bitcoin wallet or node.
4. Record the resulting `txid` for later verification.

### 5.3 Verification Procedure

Given a `.morus` log and a Bitcoin `txid`:

1. Validate the log per Section 4.
2. Compute HEAD hash = `records[last].hash`.
3. Retrieve the transaction by `txid`.
4. Extract the `OP_RETURN` payload (hex).
5. Compare: if `OP_RETURN payload == HEAD hash` → the log is provably unmodified since that Bitcoin block.

---

## 6. Actor Model

A **Node** is any actor that produces actions, decisions, or state transitions requiring a verifiable record.

| Actor Type | Description |
|---|---|
| **Human Node** | A person using a device, wallet, or companion AI |
| **AI Node** | Any autonomous or semi-autonomous agent (LLM, workflow agent, swarm member, etc.) |
| **System Node** | Any non-AI system: server, IoT device, daemon, pipeline, script |
| **Hybrid Node** | A human + AI pair acting together |

HNBP-CORE is actor-agnostic. The `data` field carries actor-specific semantics. The protocol does not care who or what is writing records — only that the chain is intact.

---

## 7. Determinism

All implementations MUST be deterministic:

- The same record inputs MUST always produce the same hash.
- Canonical JSON serialization (sorted keys, no whitespace) is mandatory.
- No implementation-specific floating point, locale, or encoding differences are permitted.
- SHA-256 is the only permitted hash function.

---

## 8. Versioning

- This spec is **HNBP-CORE v1.0.0** and is frozen.
- The format does not change.
- Implementations may be updated; the format may not.
- Future spec versions (if any) will be assigned new version numbers and will not alter v1.0.0 behavior.
- Any implementation claiming "HNBP-CORE compatible" MUST pass the v1.0.0 test suite.

---

## 9. Checksums

SHA-256 checksums of this specification file are stored in `/spec/spec_hashes.txt` for tamper-evidence of the spec itself.

---

*End of HNBP-CORE v1.0.0 Frozen Specification.*
