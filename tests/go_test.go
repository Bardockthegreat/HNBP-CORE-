// HNBP-CORE v1.0.0 — Go Test Runner
// Run: go test ./tests/

package tests

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

type Record struct {
	Index     int         `json:"index"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
	PrevHash  interface{} `json:"prev_hash"`
	Hash      string      `json:"hash"`
}

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
		return fmt.Errorf("log must be non-empty")
	}
	for i, rec := range records {
		if rec.Index != i {
			return fmt.Errorf("record %d: index mismatch", i)
		}
		if _, err := time.Parse(time.RFC3339, rec.Timestamp); err != nil {
			return fmt.Errorf("record %d: invalid timestamp", i)
		}
		if i == 0 && rec.PrevHash != nil {
			return fmt.Errorf("record 0: prev_hash must be null")
		}
		if i > 0 && fmt.Sprintf("%v", rec.PrevHash) != records[i-1].Hash {
			return fmt.Errorf("record %d: prev_hash mismatch", i)
		}
		expected := computeHash(rec.Index, rec.Timestamp, rec.Data, rec.PrevHash)
		if rec.Hash != expected {
			return fmt.Errorf("record %d: hash mismatch", i)
		}
	}
	return nil
}

func loadLog(path string) ([]Record, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records []Record
	return records, json.Unmarshal(b, &records)
}

func TestValidCases(t *testing.T) {
	files, _ := filepath.Glob("valid/*.json")
	sort.Strings(files)
	for _, f := range files {
		records, err := loadLog(f)
		if err != nil {
			t.Errorf("%s: parse error: %v", f, err)
			continue
		}
		if err := validateLog(records); err != nil {
			t.Errorf("%s: expected valid, got: %v", f, err)
		} else {
			t.Logf("PASS %s", f)
		}
	}
}

func TestInvalidCases(t *testing.T) {
	files, _ := filepath.Glob("invalid/*.json")
	sort.Strings(files)
	for _, f := range files {
		records, err := loadLog(f)
		if err != nil {
			t.Logf("PASS %s (parse error as expected)", f)
			continue
		}
		if err := validateLog(records); err == nil {
			t.Errorf("%s: expected failure but passed", f)
		} else {
			t.Logf("PASS %s (failed as expected: %v)", f, err)
		}
	}
}

type btcFixture struct {
	Log    []Record `json:"log"`
	Anchor struct {
		TXID      string `json:"txid"`
		OPReturn  string `json:"op_return"`
		BlockHeight int  `json:"block_height"`
	} `json:"anchor"`
}

func TestBitcoinCases(t *testing.T) {
	validHex := regexp.MustCompile(`^[0-9a-f]{64}$`)
	cases := []struct {
		file    string
		wantMatch bool
		wantMalformed bool
	}{
		{"bitcoin/anchored_valid.json", true, false},
		{"bitcoin/anchored_mismatch.json", false, false},
		{"bitcoin/anchored_malformed_opreturn.json", false, true},
	}
	for _, tc := range cases {
		b, err := os.ReadFile(tc.file)
		if err != nil {
			t.Errorf("%s: read error: %v", tc.file, err)
			continue
		}
		var fix btcFixture
		if err := json.Unmarshal(b, &fix); err != nil {
			t.Errorf("%s: parse error: %v", tc.file, err)
			continue
		}
		if err := validateLog(fix.Log); err != nil {
			t.Errorf("%s: log invalid: %v", tc.file, err)
			continue
		}
		head := fix.Log[len(fix.Log)-1].Hash
		op := fix.Anchor.OPReturn
		if tc.wantMalformed {
			if validHex.MatchString(op) {
				t.Errorf("%s: expected malformed OP_RETURN", tc.file)
			} else {
				t.Logf("PASS %s (malformed OP_RETURN)", tc.file)
			}
		} else if tc.wantMatch {
			if head != op {
				t.Errorf("%s: HEAD mismatch", tc.file)
			} else {
				t.Logf("PASS %s (HEAD matches)", tc.file)
			}
		} else {
			if head == op {
				t.Errorf("%s: expected mismatch but matched", tc.file)
			} else {
				t.Logf("PASS %s (mismatch confirmed)", tc.file)
			}
		}
	}
}
