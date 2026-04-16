# HNBP-CORE Integration Guide

## Human-Facing App

### Create a log

```python
from hnbp import append, validate_log, save_log

log = []
append(log, {"actor": "alice", "event": "session-start", "device": "mobile"})
save_log(log, "alice.morus")
```

### Append an entry

```python
append(log, {"actor": "alice", "event": "document-signed", "doc_id": "contract-001"})
save_log(log, "alice.morus")
```

### Validate

```python
validate_log(log)  # raises ValueError if invalid
print("Log valid. HEAD:", head_hash(log))
```

---

## AI Agent Integration

```javascript
import { append, validateLog, headHash, saveLog } from './reference/javascript/hnbp.js';

const log = [];
append(log, { actor: 'my-agent', event: 'init', task: 'summarize-docs' });

// On each action:
append(log, { actor: 'my-agent', event: 'tool-call', tool: 'read_file', path: 'doc.pdf' });
append(log, { actor: 'my-agent', event: 'output', output_hash: 'sha256:abc123' });

validateLog(log);
console.log('HEAD:', headHash(log));
saveLog(log, 'agent.morus');
```

---

## Backend System

```python
import json
from hnbp import append, validate_log, save_log, load_log

# Load existing log or start fresh
try:
    log = load_log("system.morus")
except FileNotFoundError:
    log = []
    append(log, {"actor": "payment-service", "event": "start", "version": "3.2.1"})

# On each event:
def log_event(event_data):
    append(log, event_data)
    save_log(log, "system.morus")

log_event({"event": "tx-processed", "tx_id": "TXN-0042", "amount": 500})
```

---

## Compute HEAD Hash

```python
from hnbp import head_hash
print(head_hash(log))  # 64-character hex string
```

```javascript
import { headHash } from './hnbp.js';
console.log(headHash(log));
```

---

## Anchor to Bitcoin (PSBT with OP_RETURN)

1. Compute HEAD hash:

```python
head = head_hash(log)  # e.g. "4424b820..."
```

2. Create a Bitcoin transaction with `OP_RETURN <head_hash>` using any wallet or library that supports PSBTs. The OP_RETURN payload is the HEAD hash as raw bytes (32 bytes, decoded from hex).

3. Broadcast and record the resulting `txid`.

4. Append an anchor record to the log:

```python
append(log, {
    "event": "bitcoin-anchor",
    "head_hash": head,
    "txid": "a1b2c3...",
    "block_height": 882000
})
save_log(log, "mylog.morus")
```

---

## Verify a Bitcoin Anchor

```python
from hnbp import validate_log, head_hash, verify_bitcoin_anchor

log = load_log("mylog.morus")
validate_log(log)

# op_return_hex comes from parsing the Bitcoin transaction
op_return_hex = "4424b82044a8c3a264cdb3bcb8fb968326772dc132c3ac82c7ed9b8b62962cf0"
if verify_bitcoin_anchor(log, op_return_hex):
    print("Bitcoin anchor verified — log is immutably timestamped.")
else:
    print("Anchor mismatch — log may have been modified after anchoring.")
```

---

## Web Verifier

Open `verifiers/web/index.html` in any browser. Drag and drop a `.morus` file or paste JSON. No server required.
