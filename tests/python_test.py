"""
HNBP-CORE v1.0.0 — Python Test Runner
Runs all valid/invalid/bitcoin test cases against the reference implementation.
"""

import json
import os
import sys

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'reference', 'python'))
from hnbp import validate_log, head_hash

TESTS_DIR = os.path.dirname(__file__)
passed = failed = 0

def check(label, fn):
    global passed, failed
    try:
        fn()
        print(f"  PASS  {label}")
        passed += 1
    except Exception as e:
        print(f"  FAIL  {label}: {e}")
        failed += 1

print("\n=== Valid Cases ===")
valid_dir = os.path.join(TESTS_DIR, 'valid')
for fname in sorted(os.listdir(valid_dir)):
    if not fname.endswith('.json'):
        continue
    path = os.path.join(valid_dir, fname)
    with open(path) as f:
        records = json.load(f)
    check(fname, lambda r=records: validate_log(r))

print("\n=== Invalid Cases ===")
invalid_dir = os.path.join(TESTS_DIR, 'invalid')
for fname in sorted(os.listdir(invalid_dir)):
    if not fname.endswith('.json'):
        continue
    path = os.path.join(invalid_dir, fname)
    with open(path) as f:
        records = json.load(f)
    def must_fail(r=records, n=fname):
        try:
            validate_log(r)
            raise AssertionError(f"{n} should have failed but passed")
        except (ValueError, KeyError, AssertionError) as e:
            if "should have failed" in str(e):
                raise
    check(fname, must_fail)

print("\n=== Bitcoin Cases ===")
btc_dir = os.path.join(TESTS_DIR, 'bitcoin')
for fname in sorted(os.listdir(btc_dir)):
    if not fname.endswith('.json'):
        continue
    path = os.path.join(btc_dir, fname)
    with open(path) as f:
        fixture = json.load(f)
    records = fixture['log']
    op_return = fixture['anchor']['op_return']
    if fname == 'anchored_valid.json':
        def btc_valid(r=records, op=op_return):
            validate_log(r)
            assert head_hash(r) == op, "HEAD mismatch"
        check(fname, btc_valid)
    elif fname == 'anchored_mismatch.json':
        def btc_mismatch(r=records, op=op_return):
            validate_log(r)
            assert head_hash(r) != op, "Should mismatch"
        check(fname, btc_mismatch)
    elif fname == 'anchored_malformed_opreturn.json':
        def btc_malformed(r=records, op=op_return):
            validate_log(r)
            import re
            assert not re.fullmatch(r'[0-9a-f]{64}', op), "Should be malformed"
        check(fname, btc_malformed)

print(f"\nResults: {passed} passed, {failed} failed")
sys.exit(1 if failed else 0)
