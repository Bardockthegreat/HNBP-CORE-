// HNBP-CORE v1.0.0 — Reference Implementation (Go)
// Zero dependencies (stdlib only). Deterministic. Spec-compliant.

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Record struct {
	Index     int         `json:"index"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
	PrevHash  interface{} `json:"prev_hash"`
	Hash      string      `json:"hash"`
}

// canonicalJSON produces sorted-key, no-whitespace JSON for any value.
func canonicalJSON(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			kj, _ := json.Marshal(k)
			parts = append(parts, string(kj)+":"+canonicalJSON(val[k]))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case []interface{}:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = canonicalJSON(item)
		}
		return "[" + strings.Join(parts, ",") + "]"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func computeHash(index int, timestamp string, data interface{}, prevHash interface{}) string {
	prev := "null"
	if prevHash != nil {
		prev = fmt.Sprintf("%v", prevHash)
	}
	raw := fmt.Sprintf("%d%s%s%s", index, timestamp, canonicalJSON(data), prev)
	h := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", h)
}

func validateLog(records []Record) error {
	if len(records) == 0 {
		return fmt.Errorf("log must be a non-empty array")
	}
	for i, rec := range records {
		if rec.Index != i {
			return fmt.Errorf("record %d: index mismatch", i)
		}
		if _, err := time.Parse(time.RFC3339, rec.Timestamp); err != nil {
			return fmt.Errorf("record %d: invalid timestamp", i)
		}
		if i == 0 {
			if rec.PrevHash != nil {
				return fmt.Errorf("record 0: prev_hash must be null")
			}
		} else {
			if fmt.Sprintf("%v", rec.PrevHash) != records[i-1].Hash {
				return fmt.Errorf("record %d: prev_hash mismatch", i)
			}
		}
		expected := computeHash(rec.Index, rec.Timestamp, rec.Data, rec.PrevHash)
		if rec.Hash != expected {
			return fmt.Errorf("record %d: hash mismatch", i)
		}
	}
	return nil
}

func headHash(records []Record) string {
	return records[len(records)-1].Hash
}

func appendRecord(records []Record, data interface{}) []Record {
	ts := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	index := len(records)
	var prevHash interface{} = nil
	if index > 0 {
		prevHash = records[index-1].Hash
	}
	h := computeHash(index, ts, data, prevHash)
	return append(records, Record{
		Index: index, Timestamp: ts, Data: data, PrevHash: prevHash, Hash: h,
	})
}

func loadLog(path string) ([]Record, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records []Record
	return records, json.Unmarshal(b, &records)
}

func saveLog(records []Record, path string) error {
	b, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func main() {
	var log []Record
	log = appendRecord(log, map[string]interface{}{"event": "genesis", "actor": "human"})
	log = appendRecord(log, map[string]interface{}{"event": "signed", "doc": "contract-001"})
	log = appendRecord(log, map[string]interface{}{"event": "confirmed", "by": "counterparty"})
	if err := validateLog(log); err != nil {
		fmt.Println("Invalid:", err)
		os.Exit(1)
	}
	fmt.Println("Log valid. HEAD:", headHash(log))
	if err := saveLog(log, "example.morus"); err != nil {
		fmt.Println("Save error:", err)
	}
	fmt.Println("Saved to example.morus")
}
