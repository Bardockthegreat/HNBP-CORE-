/**
 * HNBP-CORE v1.0.0 — Reference Implementation (JavaScript / Node.js)
 * Zero dependencies. Deterministic. Spec-compliant.
 */

import { createHash } from 'crypto';
import { readFileSync, writeFileSync } from 'fs';

function canonicalJSON(data) {
  if (data === null || typeof data !== 'object' || Array.isArray(data)) {
    return JSON.stringify(data);
  }
  const sorted = Object.keys(data).sort().reduce((acc, k) => {
    acc[k] = data[k];
    return acc;
  }, {});
  return '{' + Object.keys(sorted).map(k =>
    JSON.stringify(k) + ':' + canonicalJSON(sorted[k])
  ).join(',') + '}';
}

function computeHash(index, timestamp, data, prevHash) {
  const prevStr = prevHash === null ? 'null' : String(prevHash);
  const raw = String(index) + String(timestamp) + canonicalJSON(data) + prevStr;
  return createHash('sha256').update(raw, 'utf8').digest('hex');
}

function validateLog(records) {
  if (!Array.isArray(records) || records.length === 0)
    throw new Error('Log must be a non-empty JSON array');
  for (let i = 0; i < records.length; i++) {
    const rec = records[i];
    for (const f of ['index', 'timestamp', 'data', 'prev_hash', 'hash']) {
      if (!(f in rec)) throw new Error(`Record ${i}: missing field '${f}'`);
    }
    if (rec.index !== i) throw new Error(`Record ${i}: index mismatch`);
    if (isNaN(Date.parse(rec.timestamp))) throw new Error(`Record ${i}: invalid timestamp`);
    if (i === 0) {
      if (rec.prev_hash !== null) throw new Error('Record 0: prev_hash must be null');
    } else {
      if (rec.prev_hash !== records[i - 1].hash)
        throw new Error(`Record ${i}: prev_hash mismatch`);
    }
    const expected = computeHash(rec.index, rec.timestamp, rec.data, rec.prev_hash);
    if (rec.hash !== expected) throw new Error(`Record ${i}: hash mismatch`);
  }
  return true;
}

function headHash(records) {
  return records[records.length - 1].hash;
}

function append(records, data, timestamp = null) {
  if (!timestamp) timestamp = new Date().toISOString().replace(/\.\d+Z$/, 'Z');
  const index = records.length;
  const prevHash = records.length > 0 ? records[records.length - 1].hash : null;
  const hash = computeHash(index, timestamp, data, prevHash);
  records.push({ index, timestamp, data, prev_hash: prevHash, hash });
  return records;
}

function loadLog(path) {
  return JSON.parse(readFileSync(path, 'utf8'));
}

function saveLog(records, path) {
  writeFileSync(path, JSON.stringify(records, null, 2), 'utf8');
}

function verifyBitcoinAnchor(records, opReturnHex) {
  return headHash(records) === opReturnHex.toLowerCase();
}

// Demo
const log = [];
append(log, { event: 'genesis', actor: 'human' });
append(log, { event: 'signed', doc: 'contract-001' });
append(log, { event: 'confirmed', by: 'counterparty' });
validateLog(log);
console.log('Log valid. HEAD:', headHash(log));
saveLog(log, 'example.morus');
console.log('Saved to example.morus');

export { computeHash, validateLog, headHash, append, loadLog, saveLog, verifyBitcoinAnchor };
