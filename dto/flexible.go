package dto

import (
	"encoding/json"
	"strings"
)

// FlexibleStringSlice accepts assignees as a JSON array or legacy string values ("null", etc.).
type FlexibleStringSlice []string

func (f *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*f = FlexibleStringSlice{}
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*f = FlexibleStringSlice(normalizeAssignees(arr))
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(s)
		if s == "" || s == "null" {
			*f = FlexibleStringSlice{}
			return nil
		}
		if strings.HasPrefix(s, "[") {
			_ = json.Unmarshal([]byte(s), &arr)
			*f = FlexibleStringSlice(normalizeAssignees(arr))
			return nil
		}
		*f = FlexibleStringSlice{s}
		return nil
	}

	*f = FlexibleStringSlice{}
	return nil
}

func normalizeAssignees(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" || s == "null" || s == "[" || s == "]" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func (f FlexibleStringSlice) Slice() []string {
	if f == nil {
		return []string{}
	}
	return []string(f)
}
