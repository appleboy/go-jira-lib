package onpremise

import (
	"encoding/json"
	"testing"

	"github.com/fatih/structs"
)

func TestIssueIDWireFormat(t *testing.T) {
	tests := map[string]any{
		"worklog": WorklogRecord{IssueID: "10002"},
		"request": Request{IssueID: "10002"},
	}
	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]any
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			for encoding, encoded := range map[string]map[string]any{
				"json":    fields,
				"structs": structs.Map(value),
			} {
				if encoded["issueId"] != "10002" {
					t.Errorf("%s: expected Jira field issueId=10002, got %v", encoding, encoded)
				}
				if _, exists := encoded["issueID"]; exists {
					t.Errorf("%s: unexpected field issueID; Jira expects issueId", encoding)
				}
			}
		})
	}
}

// Empty nested structs were always encoded by encoding/json despite omitempty.
func TestEmptyCommentStructFieldsRemainEncoded(t *testing.T) {
	data, err := json.Marshal(Comment{})
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"author", "updateAuthor", "visibility"} {
		if _, exists := fields[key]; !exists {
			t.Errorf("missing existing wire field %q in %s", key, data)
		}
	}
}
