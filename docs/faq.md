# HNBP-CORE FAQ

## Is this a blockchain?

No. HNBP-CORE is a hash chain — a simple, append-only linked list of records where each record commits to all previous records via a SHA-256 hash. It does not have distributed consensus, mining, or nodes. It is a single-actor (or multi-actor) witness log. Bitcoin provides the global anchor if you want external immutability.

## Do I need to run a Bitcoin node?

No. Anchoring is optional. You can use HNBP-CORE entirely offline, for local tamper-evidence, without any Bitcoin interaction. If you want a Bitcoin anchor, you only need a wallet that supports OP_RETURN outputs (any modern Bitcoin wallet). Verification of an anchor requires looking up a txid — any public block explorer or node will do.

## Can I use this without Bitcoin?

Yes. Bitcoin anchoring is an optional extension. A HNBP log is independently valuable as a tamper-evident, append-only log for any actor. The chain integrity is self-verifying — no Bitcoin required.

## How is this different from normal logs?

A normal log file (syslog, JSON log, database table) can be modified, deleted, or backdated without detection. An HNBP log cannot — any change to any record breaks the chain, and verification instantly detects the tampering. With a Bitcoin anchor, even the existence of the log at a specific time is globally provable.

## What does "noncommercial" mean under PolyForm?

Noncommercial means personal use, research, education, open source projects, and non-profit use. It excludes:
- Using HNBP-CORE as part of a paid product or service.
- Internal use at a for-profit company (even if not sold directly).
- Any use "directed toward commercial advantage or monetary compensation."

When in doubt, assume commercial and get a license.

## How do I get a commercial license?

Contact the author via the GitHub repository. Commercial licenses are available for companies and products that want to build on HNBP-CORE.

## What is the `.morus` file extension?

`.morus` is the canonical file extension for HNBP log files. The content is standard JSON (a JSON array of records). The extension signals that the file is an HNBP witness log rather than arbitrary JSON.

## Can multiple actors share one log?

Yes. The `data` field can include an `actor` identifier. Multiple actors can append to the same log sequentially. For simultaneous multi-actor signing, each actor maintains their own log and cross-references other actors' HEAD hashes in their records.

## Is the spec going to change?

No. HNBP-CORE v1.0.0 is frozen. The format will not change. Implementations may be updated; the core hashing and validation rules are permanent. This is a deliberate design choice — logs written today must be verifiable forever.

## Can I fork this?

Yes, under PolyForm Noncommercial 1.0.0 for noncommercial use. For commercial forks or products, contact the author for a commercial license.
