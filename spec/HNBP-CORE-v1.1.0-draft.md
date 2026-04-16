# HNBP-CORE v1.1.0 — Draft Specification

**Status:** DRAFT — not yet frozen.  
**Date:** 2026-04-16  
**Supersedes:** HNBP-CORE v1.0.0 (frozen, unchanged)  
**Format:** `.morus` (JSON array)

---

## What Changed From v1.0.0

| Issue | v1.0.0 | v1.1.0 |
|---|---|---|
| Hash construction | String concatenation (technically ambiguous) | Canonical JSON of full record (unambiguous) |
| Identity | Actor label in `data` (unbound) | Optional Ed25519 signature over each record |
| Consensus language | Removed in README | Not applicable — this is a witness log, not a consensus mechanism |

v1.0.0 logs remain valid and verifiable forever. v1.1.0 is additive.

---

## 1. Scope

HNBP-CORE v1.1.0 defines a tamper-evident, append-only witness log for any actor — human, AI, or system — with optional Bitcoin anchoring and optional Ed25519 record signatures.

This specification covers:
- Record structure (extended from v1.0.0)
- Hashing rule (fixed)
- Identity and signature layer (new)
- Log validity conditions (extended)
- Bitcoin anchoring procedure (unchanged)
- Determinism requirements
- Versioning policy

---

## 2. Record Structure

Each record contains the v1.0.0 fields plus an optional `sig` field:

| Field | Type | Required | Description |
|---|---|---|---|
| `index` | integer | yes | Position in the log, starting at 0, strictly sequential |
| `timestamp` | string | yes | ISO 8601 datetime string |
| `data` | any | yes | Any JSON-serializable value |
| `prev_hash` | string or null | yes | Hash of the previous record; `null` only for `index 0` |
| `hash` | string | yes | SHA-256 hex digest per Section 3 |
| `sig` | object | no | Ed25519 signature object per Section 4 |

---

## 3. Hashing Rule (Fixed)

The hash for each record is computed by serializing the full record (minus `hash` and `sig`) as canonical JSON, then taking its SHA-256.

**Canonical input object:**

```json
{
  "index": <integer>,
  "timestamp": "<ISO 8601 string>",
  "data": <any JSON value>,
  "prev_hash": <string or null>
}
```

**Hash computation:**

```
canonical = JSON.stringify(input, sorted_keys, no_whitespace)
hash = SHA256(UTF8(canonical))
```

**Rules:**
- Keys sorted alphabetically at every level (recursive).
- No whitespace (no spaces or newlines between tokens).
- `prev_hash` is JSON `null` for record 0, the previous hash string for all others.
- SHA-256 output is lowercase hex.
- This replaces v1.0.0's string concatenation approach.

**Why this is unambiguous:** The canonical JSON encoding is fully typed. An integer `12` and a string `"12"` produce different encodings. Field boundaries are defined by JSON structure, not string length. There is no way to produce the same canonical input from different field values.

---

## 4. Identity and Signature Layer (New)

### 4.1 Purpose

v1.0.0 allows any actor to write `"actor": "alice"` in the `data` field, but nothing cryptographically binds that label to a real key. v1.1.0 adds an optional signature that proves the record was produced by the holder of a specific private key.

### 4.2 Key Format

- Algorithm: **Ed25519**
- Public key encoding: **base64url** (no padding)
- A keypair is generated once per actor and registered in a key registry (out of scope for this spec — any PKI, DID, or manual exchange works).

### 4.3 Signature Object

If present, the `sig` field is:

```json
{
  "alg": "ed25519",
  "kid": "<key identifier string>",
  "pub": "<base64url-encoded public key>",
  "sig": "<base64url-encoded signature over the record hash>"
}
```

### 4.4 What Is Signed

The signature is computed over the **record hash** (the SHA-256 hex string from Section 3), encoded as UTF-8 bytes:

```
signature = Ed25519Sign(private_key, UTF8(record.hash))
```

This means:
- The signature commits to the entire record (via the hash).
- Verifying the hash verifies the record. Verifying the signature verifies the actor.

### 4.5 Verification

To verify a signed record:
1. Verify the record hash per Section 3.
2. Obtain the actor's public key by `kid`.
3. Verify: `Ed25519Verify(pub_key, UTF8(record.hash), sig)`.

### 4.6 Unsigned Records

The `sig` field is optional. Unsigned records are valid HNBP logs. The signature layer is an extension for use cases that require actor binding.

---

## 5. Log Validity

A v1.1.0 log is valid if and only if ALL of the following hold:

1. **Non-empty**: The log is a JSON array with at least one record.
2. **Sequential indices**: `records[i].index == i` for all `i`.
3. **Genesis prev_hash**: `records[0].prev_hash === null`.
4. **Chain linkage**: For all `i > 0`, `records[i].prev_hash === records[i-1].hash`.
5. **Hash integrity**: For every record, the stored `hash` matches the recomputed hash per Section 3.
6. **Valid timestamps**: Every `timestamp` is a valid ISO 8601 string.
7. **Signature validity** (if `sig` present): Signature verifies per Section 4.5.

---

## 6. Bitcoin Anchoring

Unchanged from v1.0.0. See Section 5 of HNBP-CORE-v1.0.0.md.

---

## 7. Compatibility

- v1.0.0 logs are **not** v1.1.0 hash-compatible (different hash construction).
- Implementations MUST declare which version they implement.
- A v1.1.0 verifier SHOULD be able to verify v1.0.0 logs using the v1.0.0 hash rule.
- Log files MAY include a `"hnbp_version"` field in the first record's `data` to declare version.

---

## 8. Determinism

All implementations MUST be deterministic. Same inputs → same hash, always, across all languages and platforms.

---

## 9. Versioning

- HNBP-CORE v1.0.0 is frozen and will never change.
- HNBP-CORE v1.1.0 will be frozen upon completion of the reference implementation and test suite.
- This document is a draft and may be revised before freezing.

---

*End of HNBP-CORE v1.1.0 Draft Specification.*
