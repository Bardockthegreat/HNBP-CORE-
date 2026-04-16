/**
 * HNBP-CORE v1.0.0 — JavaScript Test Runner (ESM)
 * Run: node tests/js_test.mjs
 */

import { readFileSync, readdirSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { createHash } from 'crypto';

const __dir = dirname(fileURLToPath(import.meta.url));

// Inline reference logic (mirrors reference/javascript/hnbp.js)
function canonicalJSON(data) {
  if (data === null || typeof data !== 'object' || Array.isArray(data)) return JSON.stringify(data);
  const sorted = Object.keys(data).sort().reduce((acc, k) => { acc[k] = data[k]; return acc; }, {});
  return '{' + Object.keys(sorted).map(k => JSON.stringify(k) + ':' + canonicalJSON(sorted[k])).join(',') + '}';
}
function computeHash(index, timestamp, data, prevHash) {
  const prevStr = prevHash === null ? 'null' : String(prevHash);
  const raw = String(index) + String(timestamp) + canonicalJSON(data) + prevStr;
  return createHash('sha256').update(raw, 'utf8').digest('hex');
}
function validateLog(records) {
  if (!Array.isArray(records) || records.length === 0) throw new Error('Non-empty array required');
  for (let i = 0; i < records.length; i++) {
    const rec = records[i];
    for (const f of ['index', 'timestamp', 'data', 'prev_hash', 'hash'])
      if (!(f in rec)) throw new Error(`Record ${i}: missing '${f}'`);
    if (rec.index !== i) throw new Error(`Record ${i}: index mismatch`);
    if (isNaN(Date.parse(rec.timestamp))) throw new Error(`Record ${i}: invalid timestamp`);
    if (i === 0 && rec.prev_hash !== null) throw new Error('Record 0: prev_hash must be null');
    if (i > 0 && rec.prev_hash !== records[i-1].hash) throw new Error(`Record ${i}: prev_hash mismatch`);
    if (rec.hash !== computeHash(rec.index, rec.timestamp, rec.data, rec.prev_hash))
      throw new Error(`Record ${i}: hash mismatch`);
  }
}
const headHash = records => records[records.length - 1].hash;

let passed = 0, failed = 0;
function check(label, fn) {
  try { fn(); console.log(`  PASS  ${label}`); passed++; }
  catch (e) { console.log(`  FAIL  ${label}: ${e.message}`); failed++; }
}

const load = path => JSON.parse(readFileSync(path, 'utf8'));
const jsonFiles = dir => readdirSync(dir).filter(f => f.endsWith('.json')).sort();

console.log('\n=== Valid Cases ===');
const validDir = join(__dir, 'valid');
for (const f of jsonFiles(validDir)) {
  const records = load(join(validDir, f));
  check(f, () => validateLog(records));
}

console.log('\n=== Invalid Cases ===');
const invalidDir = join(__dir, 'invalid');
for (const f of jsonFiles(invalidDir)) {
  const records = load(join(invalidDir, f));
  check(f, () => {
    let threw = false;
    try { validateLog(records); } catch { threw = true; }
    if (!threw) throw new Error(`${f} should have failed`);
  });
}

console.log('\n=== Bitcoin Cases ===');
const btcDir = join(__dir, 'bitcoin');
for (const f of jsonFiles(btcDir)) {
  const fixture = load(join(btcDir, f));
  const { log, anchor } = fixture;
  if (f === 'anchored_valid.json') {
    check(f, () => { validateLog(log); if (headHash(log) !== anchor.op_return) throw new Error('HEAD mismatch'); });
  } else if (f === 'anchored_mismatch.json') {
    check(f, () => { validateLog(log); if (headHash(log) === anchor.op_return) throw new Error('Should mismatch'); });
  } else if (f === 'anchored_malformed_opreturn.json') {
    check(f, () => { validateLog(log); if (/^[0-9a-f]{64}$/.test(anchor.op_return)) throw new Error('Should be malformed'); });
  }
}

console.log(`\nResults: ${passed} passed, ${failed} failed`);
process.exit(failed ? 1 : 0);
