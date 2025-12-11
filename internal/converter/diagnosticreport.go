package converter

import (
	"fmt"
	"time"

	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertToDiagnosticReports converts OBR segment to FHIR DiagnosticReport
func ConvertToDiagnosticReport(msg *hl7.Message, patientID string) (*fhir.DiagnosticReport, error) {
	obrSegment := msg.GetSegment("OBR")
	if obrSegment == nil {
		return nil, nil
	}

	report := &fhir.DiagnosticReport{
		ResourceType: "DiagnosticReport",
		ID:           getOBRID(obrSegment),
		Status:       mapOBRStatus(obrSegment),
		Code:         getOBRCode(obrSegment),
		Subject:      &fhir.Reference{Reference: "Patient/" + patientID},
		Issued:       getOBRDateTime(obrSegment),
	}

	return report, nil
}

// get OBRID extracts report ID from OBR-1
func getOBRID(seg *hl7.Segment) string {
	//OBR-2 Placer order number
	//OBR-2 Placer order number
	field := seg.GetField(2)
	if field != nil {
		comp := field.GetCompontent(1)
		if comp != "" {
			return "report-" + comp
		}
	}

	//Fallback to OBR-1
	field = seg.GetField(1)
	if field != nil {
		comp := field.GetCompontent(1)
		if comp != "" {
			return "report-" + comp
		}
	}

	return "report-unknown"
}

// mapOBRStatus maps OBR-25 to FHIR Status
func mapOBRStatus(seg *hl7.Segment) string {
	field := seg.GetField(25)
	if field == nil {
		return "unknown"
	}

	status := field.GetCompontent(1)

	switch status {
	case "O": // Order received
		return "registered"
	case "I": // In process
		return "partial"
	case "P": // Preliminary
		return "preliminary"
	case "F": // Final
		return "final"
	case "C": // Corrected
		return "corrected"
	case "X": // Cancelled
		return "cancelled"
	default:
		return "unknown"
	}
}

// getOBRCode extracts test/procedure code from OBR-4
func getOBRCode(seg *hl7.Segment) *fhir.CodeableConcept {
	field := seg.GetField(4)
	if field == nil {
		return nil
	}

	identifier := field.GetCompontent(1)
	text := field.GetCompontent(2)
	system := field.GetCompontent(3)

	if identifier == "" {
		return nil
	}

	//Map coding system
	codingSystem := "http://loinc.org"
	if system == "L" || system == "LN" {
		codingSystem = "http://loinc.org"
	}

	return &fhir.CodeableConcept{
		Coding: []fhir.Coding{
			{
				System:  codingSystem,
				Code:    identifier,
				Display: text,
			},
		},
		Text: text,
	}
}

// getOBRDateTime extracts observation date/time from OBR-7
func getOBRDateTime(seg *hl7.Segment) string {
	field := seg.GetField(7)
	if field == nil {
		return time.Now().Format(time.RFC3339)
	}

	dateTime := field.GetCompontent(1)
	if dateTime == "" {
		return time.Now().Format(time.RFC3339)
	}

	//Parse HL7 datetime
	if len(dateTime) >= 8 {
		// Simple conversion: YYYYMMDD -> YYYY-MM-DD
		year := dateTime[0:4]
		month := dateTime[4:6]
		day := dateTime[6:8]

		if len(dateTime) >= 14 {
			hour := dateTime[8:10]
			minute := dateTime[10:12]
			second := dateTime[12:14]
			return fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", year, month, day, hour, minute, second)
		}

		return fmt.Sprintf("%s-%s-%sT00:00:00Z", year, month, day)
	}

	return time.Now().Format(time.RFC3339)
}
