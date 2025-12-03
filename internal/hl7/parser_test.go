package hl7

import (
	"testing"
)

func TestParse_BasicMessage(t *testing.T) {
	raw := `MSH|^~\&|EPIC|FAC1|CERNER|FAC2|20231115||ADT^A01|MSG001|P|2.5
PID|1||12345^^^MRN||Doe^John||19800115|M`

	msg, err := Parse(raw)

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)

	}

	if len(msg.Segments) != 2 {
		t.Errorf("Expected 2 segments, got %d", len(msg.Segments))
	}

	if msg.Segments[0].Name != "MSH" {
		t.Errorf("Expected first segment MSH, got %s", msg.Segments[0].Name)
	}

	if msg.Segments[1].Name != "PID" {
		t.Errorf("Expected second segment PID, got %s", msg.Segments[1].Name)
	}

}

func TestParse_Repetitions(t *testing.T) {
	raw := `MSH|^~\&|APP|||FAC|20231115||ADT^A01|1|P|2.5
PID|1||123||Doe^John||19800115|M|||123 Main||555-1234~555-5678`

	msg, err := Parse(raw)

	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	pid := msg.GetSegment("PID")
	phoneField := pid.GetField(13)

	if len(phoneField.Repetitions) != 2 {
		t.Errorf("Expected 2 phone repetitions, got %d", len(phoneField.Repetitions))
	}
}

func TestParse_EmptyMessage(t *testing.T) {
	_, err := Parse("")

	if err == nil {
		t.Error("Expected error for empty message, got nil")
	}
}

func TestParse_InvalidStart(t *testing.T) {
	_, err := Parse("PID|1||12345")

	if err == nil {
		t.Error(("Expected error for message not starting with MSH"))
	}
}
