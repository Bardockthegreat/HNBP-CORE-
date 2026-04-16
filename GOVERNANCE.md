# HNBP-CORE Governance

## Protocol Status

**HNBP-CORE v1.0.0 is frozen.**

The core format MUST NOT be broken. No changes to record structure, hashing rules, or validation rules are permitted once published. The spec exists to be a permanent anchor — not a living document.

## Rules for Implementations

- Implementations may evolve; the format may not.
- Any implementation claiming "HNBP-CORE compatible" MUST support v1.0.0 exactly.
- Reference implementations (Python, JS, Go, Rust) serve as canonical correctness anchors.
- Deviations in hashing, serialization, or validation make an implementation non-conformant.

## Versioning

- The frozen spec lives at `/spec/HNBP-CORE-v1.0.0.md`.
- Future spec versions (if any) will be additive and will NOT break v1.0.0 logs.
- All versions are independently verifiable.

## Licensing

- This project is licensed under **PolyForm Noncommercial License 1.0.0**.
- Noncommercial use is allowed under that license.
- **Commercial use requires a separate commercial license from the author.**
- The author retains all commercial rights.

## HNBP-Verified Actor Badge

The badge:

```
✔ HNBP-Verified Actor
This actor maintains a tamper-evident witness log anchored to Bitcoin.
```

MUST only be used when:
- The actor maintains a valid HNBP log conforming to v1.0.0.
- The log can be independently verified against the spec.
- (Optionally) the HEAD hash is anchored to a Bitcoin OP_RETURN transaction.

Misuse of the badge is a protocol violation.

## Permanence

Mirrors on GitHub, GitLab, Codeberg, IPFS, and Archive.org are encouraged. The protocol must survive any single point of failure.

## Contact

For commercial licensing inquiries, contact the author via the GitHub repository.
