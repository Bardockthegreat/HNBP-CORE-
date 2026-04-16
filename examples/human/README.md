# Example: Human Node

## Scenario

Alice is a human actor who needs a tamper-evident log of her document-signing session. She uses HNBP-CORE to create a witness log that proves she signed `contract-2026-Q1-001` at a specific time.

## How HNBP-CORE Is Used

1. Alice's app creates a new log with a `session-start` event as the genesis record.
2. Each significant action (signing a document, confirming a transaction) is appended as a new record.
3. The log can be shared with any counterparty who can verify it offline using the web verifier or any reference implementation.
4. Optionally, Alice anchors the HEAD hash to Bitcoin via an `OP_RETURN` transaction for a globally immutable timestamp.

## Verification

```bash
python3 reference/python/hnbp.py
# or open verifiers/web/index.html and drag in example_log.json
```

## Why This Matters

Without HNBP-CORE, Alice has no way to prove to a third party that she performed these actions at these exact times without trusting a central server. With HNBP-CORE, the chain is self-verifying.
