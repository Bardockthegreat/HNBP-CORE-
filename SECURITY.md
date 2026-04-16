# Security Notes — HNBP-CORE

## v1.0.0 Hash Construction

HNBP-CORE v1.0.0 uses string concatenation for the hash input:

```
SHA256(String(index) + String(timestamp) + canonical_json(data) + String(prev_hash))
```

**Known theoretical limitation:** String concatenation without explicit delimiters or length-prefixes is technically ambiguous — different combinations of field values could in theory produce the same input string.

**Why v1.0.0 is safe in practice:** In valid HNBP logs, `index` is always a small non-negative integer and `timestamp` is always a valid ISO 8601 string beginning with a 4-digit year (e.g. `"2026-..."`). No realistic combination of valid field values produces a collision. The spec enforces ISO 8601 timestamp validation, which closes the gap for all real-world usage.

**Resolution in v1.1.0:** HNBP-CORE v1.1.0 (draft) replaces string concatenation with canonical JSON of the full record object. This is formally unambiguous regardless of field values.

## Identity Layer

v1.0.0 does not bind actor labels to cryptographic keys. The `data` field may carry an actor identifier (e.g. `"actor": "alice"`) but this is a label, not a proof.

v1.1.0 introduces an optional Ed25519 signature over each record's hash, providing cryptographic actor binding.

## Reporting

To report a security issue, open a private GitHub advisory or contact the author directly via GitHub.
