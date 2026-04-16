# Example: Custody Log

## Scenario

A Bitcoin custodian maintains an HNBP log of all custody events for a fund client. Monthly proof-of-reserves attestations are appended to the same chain, creating a continuous, verifiable record of asset custody.

## How HNBP-CORE Is Used

1. Custody establishment is the genesis record.
2. Monthly proof-of-reserves entries are appended.
3. Any transfer or withdrawal gets its own record.
4. The HEAD hash is anchored to Bitcoin monthly.

## Verification

A client (or regulator) can verify:
- The custody was established on a specific date.
- Proof-of-reserves attestations happened at specific intervals.
- No records were tampered with between attestations.
- The Bitcoin anchors prove each attestation existed before the corresponding block.
