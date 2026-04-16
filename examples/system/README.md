# Example: System Node

## Scenario

A payment processing service logs every transaction to an HNBP chain. This creates a tamper-evident, append-only audit log that can be verified by regulators, auditors, or customers without trusting the payment provider's database.

## How HNBP-CORE Is Used

1. The system writes a `service-start` genesis record on boot.
2. Every processed transaction is appended with status and relevant metadata.
3. The log is periodically anchored to Bitcoin for external immutability.
4. In a dispute, any party can independently verify the log and the Bitcoin anchor.

## Why Not Just Use a Database?

A database can be modified without detection. An HNBP log cannot — any modification breaks the chain, and any Bitcoin-anchored HEAD proves the log was exactly this way at that block height.
