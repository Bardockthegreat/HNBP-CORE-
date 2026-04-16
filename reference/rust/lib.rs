// HNBP-CORE v1.0.0 — Reference Implementation (Rust)
// Zero external dependencies (uses only std). Deterministic. Spec-compliant.
// Note: SHA-256 is implemented inline to avoid external crates.

use std::collections::BTreeMap;
use std::fs;

// Minimal SHA-256 (FIPS 180-4 compliant, no external deps)
fn sha256(input: &[u8]) -> [u8; 32] {
    let k: [u32; 64] = [
        0x428a2f98,0x71374491,0xb5c0fbcf,0xe9b5dba5,0x3956c25b,0x59f111f1,0x923f82a4,0xab1c5ed5,
        0xd807aa98,0x12835b01,0x243185be,0x550c7dc3,0x72be5d74,0x80deb1fe,0x9bdc06a7,0xc19bf174,
        0xe49b69c1,0xefbe4786,0x0fc19dc6,0x240ca1cc,0x2de92c6f,0x4a7484aa,0x5cb0a9dc,0x76f988da,
        0x983e5152,0xa831c66d,0xb00327c8,0xbf597fc7,0xc6e00bf3,0xd5a79147,0x06ca6351,0x14292967,
        0x27b70a85,0x2e1b2138,0x4d2c6dfc,0x53380d13,0x650a7354,0x766a0abb,0x81c2c92e,0x92722c85,
        0xa2bfe8a1,0xa81a664b,0xc24b8b70,0xc76c51a3,0xd192e819,0xd6990624,0xf40e3585,0x106aa070,
        0x19a4c116,0x1e376c08,0x2748774c,0x34b0bcb5,0x391c0cb3,0x4ed8aa4a,0x5b9cca4f,0x682e6ff3,
        0x748f82ee,0x78a5636f,0x84c87814,0x8cc70208,0x90befffa,0xa4506ceb,0xbef9a3f7,0xc67178f2,
    ];
    let mut h: [u32; 8] = [
        0x6a09e667,0xbb67ae85,0x3c6ef372,0xa54ff53a,
        0x510e527f,0x9b05688c,0x1f83d9ab,0x5be0cd19,
    ];
    let mut msg = input.to_vec();
    let bit_len = (input.len() as u64) * 8;
    msg.push(0x80);
    while msg.len() % 64 != 56 { msg.push(0); }
    msg.extend_from_slice(&bit_len.to_be_bytes());
    for chunk in msg.chunks(64) {
        let mut w = [0u32; 64];
        for i in 0..16 { w[i] = u32::from_be_bytes(chunk[i*4..i*4+4].try_into().unwrap()); }
        for i in 16..64 {
            let s0 = w[i-15].rotate_right(7) ^ w[i-15].rotate_right(18) ^ (w[i-15] >> 3);
            let s1 = w[i-2].rotate_right(17) ^ w[i-2].rotate_right(19) ^ (w[i-2] >> 10);
            w[i] = w[i-16].wrapping_add(s0).wrapping_add(w[i-7]).wrapping_add(s1);
        }
        let (mut a,mut b,mut c,mut d,mut e,mut f,mut g,mut hh) = (h[0],h[1],h[2],h[3],h[4],h[5],h[6],h[7]);
        for i in 0..64 {
            let s1 = e.rotate_right(6) ^ e.rotate_right(11) ^ e.rotate_right(25);
            let ch = (e & f) ^ ((!e) & g);
            let t1 = hh.wrapping_add(s1).wrapping_add(ch).wrapping_add(k[i]).wrapping_add(w[i]);
            let s0 = a.rotate_right(2) ^ a.rotate_right(13) ^ a.rotate_right(22);
            let maj = (a & b) ^ (a & c) ^ (b & c);
            let t2 = s0.wrapping_add(maj);
            hh=g; g=f; f=e; e=d.wrapping_add(t1); d=c; c=b; b=a; a=t1.wrapping_add(t2);
        }
        h[0]=h[0].wrapping_add(a); h[1]=h[1].wrapping_add(b); h[2]=h[2].wrapping_add(c);
        h[3]=h[3].wrapping_add(d); h[4]=h[4].wrapping_add(e); h[5]=h[5].wrapping_add(f);
        h[6]=h[6].wrapping_add(g); h[7]=h[7].wrapping_add(hh);
    }
    let mut out = [0u8; 32];
    for (i, &v) in h.iter().enumerate() { out[i*4..i*4+4].copy_from_slice(&v.to_be_bytes()); }
    out
}

fn to_hex(bytes: &[u8]) -> String {
    bytes.iter().map(|b| format!("{:02x}", b)).collect()
}

// Minimal JSON value type for canonical serialization
#[derive(Clone, Debug)]
enum JVal {
    Null,
    Bool(bool),
    Num(f64),
    Str(String),
    Arr(Vec<JVal>),
    Obj(BTreeMap<String, JVal>), // BTreeMap = sorted keys
}

impl JVal {
    fn canonical(&self) -> String {
        match self {
            JVal::Null => "null".into(),
            JVal::Bool(b) => if *b { "true" } else { "false" }.into(),
            JVal::Num(n) => {
                if n.fract() == 0.0 { format!("{}", *n as i64) } else { format!("{}", n) }
            }
            JVal::Str(s) => format!("\"{}\"", s.replace('\\', "\\\\").replace('"', "\\\"")),
            JVal::Arr(a) => format!("[{}]", a.iter().map(|v| v.canonical()).collect::<Vec<_>>().join(",")),
            JVal::Obj(m) => {
                let pairs: Vec<String> = m.iter().map(|(k,v)| {
                    format!("\"{}\":{}", k.replace('\\', "\\\\").replace('"', "\\\""), v.canonical())
                }).collect();
                format!("{{{}}}", pairs.join(","))
            }
        }
    }
}

#[derive(Clone, Debug)]
pub struct Record {
    pub index: usize,
    pub timestamp: String,
    pub data: JVal,
    pub prev_hash: Option<String>,
    pub hash: String,
}

fn compute_hash(index: usize, timestamp: &str, data: &JVal, prev_hash: Option<&str>) -> String {
    let prev = prev_hash.unwrap_or("null");
    let raw = format!("{}{}{}{}", index, timestamp, data.canonical(), prev);
    to_hex(&sha256(raw.as_bytes()))
}

pub fn validate_log(records: &[Record]) -> Result<(), String> {
    if records.is_empty() { return Err("Log must be non-empty".into()); }
    for (i, rec) in records.iter().enumerate() {
        if rec.index != i { return Err(format!("Record {i}: index mismatch")); }
        if i == 0 && rec.prev_hash.is_some() { return Err("Record 0: prev_hash must be null".into()); }
        if i > 0 {
            let expected_prev = &records[i-1].hash;
            if rec.prev_hash.as_deref() != Some(expected_prev.as_str()) {
                return Err(format!("Record {i}: prev_hash mismatch"));
            }
        }
        let expected = compute_hash(rec.index, &rec.timestamp, &rec.data, rec.prev_hash.as_deref());
        if rec.hash != expected { return Err(format!("Record {i}: hash mismatch")); }
    }
    Ok(())
}

pub fn head_hash(records: &[Record]) -> &str {
    &records[records.len()-1].hash
}

pub fn append(records: &mut Vec<Record>, data: JVal) {
    let index = records.len();
    let timestamp = "2026-04-15T00:00:00Z".to_string(); // use real time in production
    let prev_hash = records.last().map(|r| r.hash.clone());
    let hash = compute_hash(index, &timestamp, &data, prev_hash.as_deref());
    records.push(Record { index, timestamp, data, prev_hash, hash });
}

fn main() {
    let mut log: Vec<Record> = vec![];
    let mut genesis = BTreeMap::new();
    genesis.insert("actor".to_string(), JVal::Str("human".to_string()));
    genesis.insert("event".to_string(), JVal::Str("genesis".to_string()));
    append(&mut log, JVal::Obj(genesis));
    let mut signed = BTreeMap::new();
    signed.insert("doc".to_string(), JVal::Str("contract-001".to_string()));
    signed.insert("event".to_string(), JVal::Str("signed".to_string()));
    append(&mut log, JVal::Obj(signed));
    match validate_log(&log) {
        Ok(()) => println!("Log valid. HEAD: {}", head_hash(&log)),
        Err(e) => println!("Invalid: {}", e),
    }
}
