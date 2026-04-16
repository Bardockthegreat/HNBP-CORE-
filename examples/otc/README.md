# Example: OTC Trade Log

## Scenario

Alice and Bob execute an over-the-counter Bitcoin trade. Every step — proposal, acceptance, settlement — is logged to a shared HNBP chain. This creates a binding, verifiable record of the deal that neither party can alter.

## How HNBP-CORE Is Used

1. The deal proposal is the genesis record (terms, amounts, parties).
2. Acceptance is appended by the seller.
3. Settlement confirmation (Bitcoin txid + fiat reference) is the final record.
4. Both parties sign off by anchoring the HEAD hash.

## Why This Matters

OTC trades are typically done over chat with no verifiable record. An HNBP log gives both parties a cryptographic receipt of the entire deal flow that can be used in any dispute.
