package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RawActivity map[string]interface{}

type ActivityEvent struct {
	EventID string `json:"event_id"`

	UserWallet  string `json:"user_wallet"`
	ProxyWallet string `json:"proxy_wallet"`

	TimestampMs int64     `json:"timestamp_ms"`
	EventTime   time.Time `json:"event_time"`

	ConditionID     string `json:"condition_id"`
	ActivityType    string `json:"activity_type"`
	Size            string `json:"size"`
	USDCSize        string `json:"usdc_size"`
	TransactionHash string `json:"transaction_hash"`
	Price           string `json:"price"`
	Asset           string `json:"asset"`
	Side            string `json:"side"`
	OutcomeIndex    *int   `json:"outcome_index"`

	Title     string `json:"title"`
	Slug      string `json:"slug"`
	EventSlug string `json:"event_slug"`
	Outcome   string `json:"outcome"`

	RawJSON string `json:"raw_json"`

	SourceOffset int       `json:"source_offset"`
	SourceIndex  int       `json:"source_index"`
	PulledAt     time.Time `json:"pulled_at"`
}

func parseActivityEvent(raw RawActivity, userWallet string, sourceOffset, sourceIndex int, pulledAt time.Time) (ActivityEvent, error) {
	ts, _ := extractInt64(raw["timestamp"])
	eventTime := time.Time{}
	if ts > 0 {
		eventTime = time.UnixMilli(ts).UTC()
	}
	outcomeIndex, _ := extractOptionalInt(raw["outcomeIndex"])
	rawJSON := canonicalJSON(raw)

	e := ActivityEvent{
		EventID:          buildEventID(userWallet, rawJSON),
		UserWallet:       strings.ToLower(strings.TrimSpace(userWallet)),
		ProxyWallet:      strings.ToLower(strings.TrimSpace(extractString(raw["proxyWallet"]))),
		TimestampMs:      ts,
		EventTime:        eventTime,
		ConditionID:      strings.TrimSpace(extractString(raw["conditionId"])),
		ActivityType:     strings.ToUpper(strings.TrimSpace(extractString(raw["type"]))),
		Size:             normalizeNumberString(raw["size"]),
		USDCSize:         normalizeNumberString(raw["usdcSize"]),
		TransactionHash:  strings.ToLower(strings.TrimSpace(extractString(raw["transactionHash"]))),
		Price:            normalizeNumberString(raw["price"]),
		Asset:            strings.ToLower(strings.TrimSpace(extractString(raw["asset"]))),
		Side:             strings.ToLower(strings.TrimSpace(extractString(raw["side"]))),
		OutcomeIndex:     outcomeIndex,
		Title:            strings.TrimSpace(extractString(raw["title"])),
		Slug:             strings.TrimSpace(extractString(raw["slug"])),
		EventSlug:        strings.TrimSpace(extractString(raw["eventSlug"])),
		Outcome:          strings.TrimSpace(extractString(raw["outcome"])),
		RawJSON:          rawJSON,
		SourceOffset:     sourceOffset,
		SourceIndex:      sourceIndex,
		PulledAt:         pulledAt.UTC(),
	}
	if e.EventID == "" {
		return ActivityEvent{}, fmt.Errorf("empty event id")
	}
	return e, nil
}

func buildEventID(userWallet, canonical string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(userWallet)) + "|" + canonical))
	return hex.EncodeToString(sum[:])
}

func extractString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", t)
	}
}

func extractInt64(v interface{}) (int64, bool) {
	switch t := v.(type) {
	case nil:
		return 0, false
	case json.Number:
		n, err := t.Int64()
		if err == nil {
			return n, true
		}
		f, err := t.Float64()
		if err != nil {
			return 0, false
		}
		return int64(f), true
	case float64:
		return int64(t), true
	case int:
		return int64(t), true
	case int64:
		return t, true
	case string:
		if t == "" {
			return 0, false
		}
		n, err := strconv.ParseInt(t, 10, 64)
		if err == nil {
			return n, true
		}
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, false
		}
		return int64(f), true
	default:
		return 0, false
	}
}

func extractOptionalInt(v interface{}) (*int, bool) {
	switch t := v.(type) {
	case nil:
		return nil, false
	case json.Number:
		n, err := t.Int64()
		if err != nil {
			f, err2 := t.Float64()
			if err2 != nil {
				return nil, false
			}
			tmp := int(f)
			return &tmp, true
		}
		tmp := int(n)
		return &tmp, true
	case float64:
		tmp := int(t)
		return &tmp, true
	case int:
		tmp := t
		return &tmp, true
	case int64:
		tmp := int(t)
		return &tmp, true
	case string:
		t = strings.TrimSpace(t)
		if t == "" {
			return nil, false
		}
		n, err := strconv.Atoi(t)
		if err != nil {
			return nil, false
		}
		tmp := n
		return &tmp, true
	default:
		return nil, false
	}
}

func normalizeNumberString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case json.Number:
		return strings.TrimSpace(t.String())
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", t))
	}
}

func canonicalJSON(v interface{}) string {
	var b strings.Builder
	writeCanonical(&b, v)
	return b.String()
}

func writeCanonical(b *strings.Builder, v interface{}) {
	switch t := v.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if t {
			b.WriteString("true")
			return
		}
		b.WriteString("false")
	case string:
		enc, _ := json.Marshal(t)
		b.Write(enc)
	case json.Number:
		b.WriteString(strings.TrimSpace(t.String()))
	case float64:
		b.WriteString(strconv.FormatFloat(t, 'f', -1, 64))
	case int:
		b.WriteString(strconv.Itoa(t))
	case int64:
		b.WriteString(strconv.FormatInt(t, 10))
	case []interface{}:
		b.WriteByte('[')
		for i := range t {
			if i > 0 {
				b.WriteByte(',')
			}
			writeCanonical(b, t[i])
		}
		b.WriteByte(']')
	case map[string]interface{}:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			keyEnc, _ := json.Marshal(k)
			b.Write(keyEnc)
			b.WriteByte(':')
			writeCanonical(b, t[k])
		}
		b.WriteByte('}')
	default:
		enc, _ := json.Marshal(t)
		b.Write(enc)
	}
}
