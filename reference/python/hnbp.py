"""
HNBP-CORE v1.0.0 — Reference Implementation (Python)
Zero dependencies. Deterministic. Spec-compliant.
"""

import hashlib
import json
import sys
from datetime import datetime, timezone


def canonical_json(data):
    return json.dumps(data, sort_keys=True, separators=(',', ':'))


def compute_hash(index, timestamp, data, prev_hash):
    prev_str = 'null' if prev_hash is None else str(prev_hash)
    raw = str(index) + str(timestamp) + canonical_json(data) + prev_str
    return hashlib.sha256(raw.encode('utf-8')).hexdigest()


def validate_log(records):
    if not isinstance(records, list) or len(records) == 0:
        raise ValueError("Log must be a non-empty JSON array")
    for i, rec in enumerate(records):
        for field in ('index', 'timestamp', 'data', 'prev_hash', 'hash'):
            if field not in rec:
                raise ValueError(f"Record {i}: missing field '{field}'")
        if rec['index'] != i:
            raise ValueError(f"Record {i}: index mismatch (got {rec['index']})")
        try:
            datetime.fromisoformat(rec['timestamp'].replace('Z', '+00:00'))
        except ValueError:
            raise ValueError(f"Record {i}: invalid ISO 8601 timestamp")
        if i == 0:
            if rec['prev_hash'] is not None:
                raise ValueError("Record 0: prev_hash must be null")
        else:
            if rec['prev_hash'] != records[i - 1]['hash']:
                raise ValueError(f"Record {i}: prev_hash does not match record {i-1} hash")
        expected = compute_hash(rec['index'], rec['timestamp'], rec['data'], rec['prev_hash'])
        if rec['hash'] != expected:
            raise ValueError(f"Record {i}: hash mismatch")
    return True


def head_hash(records):
    return records[-1]['hash']


def append(records, data, timestamp=None):
    if timestamp is None:
        timestamp = datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')
    index = len(records)
    prev = records[-1]['hash'] if records else None
    h = compute_hash(index, timestamp, data, prev)
    records.append({
        'index': index,
        'timestamp': timestamp,
        'data': data,
        'prev_hash': prev,
        'hash': h,
    })
    return records


def load_log(path):
    with open(path, 'r', encoding='utf-8') as f:
        return json.load(f)


def save_log(records, path):
    with open(path, 'w', encoding='utf-8') as f:
        json.dump(records, f, indent=2)


def verify_bitcoin_anchor(records, op_return_hex):
    return head_hash(records) == op_return_hex.lower()


def export_log(records):
    return json.dumps(records, indent=2)


if __name__ == '__main__':
    log = []
    append(log, {'event': 'genesis', 'actor': 'human'})
    append(log, {'event': 'signed', 'doc': 'contract-001'})
    append(log, {'event': 'confirmed', 'by': 'counterparty'})
    validate_log(log)
    print("Log valid. HEAD:", head_hash(log))
    save_log(log, 'example.morus')
    print("Saved to example.morus")
