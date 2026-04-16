# Example: Multisig Key Ceremony

## Scenario

Five participants perform a 3-of-5 multisig Bitcoin key generation ceremony. Every step of the ceremony is logged to an HNBP chain so it can be independently verified that the ceremony was conducted correctly.

## How HNBP-CORE Is Used

1. Ceremony start is the genesis record.
2. Each participant's key attestation is appended as they complete their step.
3. The ceremony completion and resulting multisig address are the final record.
4. The HEAD hash is anchored to Bitcoin — proof that this ceremony happened at this block.

## Verification

Any auditor can take the `.morus` log and verify that:
- All 5 participants attested.
- The ceremony happened in the correct order.
- No records were tampered with after the fact.
- The Bitcoin anchor proves the log existed before that block.
