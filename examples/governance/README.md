# Example: Governance Log

## Scenario

A decentralized project uses HNBP-CORE to log its governance process — proposals, votes, and outcomes. This creates a permanent, verifiable record of how decisions were made.

## How HNBP-CORE Is Used

1. Proposal creation is the genesis record (body hash, not full text — for brevity).
2. Each vote is appended as a record.
3. The final outcome is appended with vote tallies.
4. The HEAD hash is anchored to Bitcoin at the close of voting.

## Why This Matters

Without a tamper-evident log, governance results can be disputed or altered after the fact. An HNBP chain + Bitcoin anchor makes the vote count permanent and independently verifiable by anyone with a copy of the log.
