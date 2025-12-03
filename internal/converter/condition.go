package converter

import (
	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertToConditions converts DG1 segments to FHIR Conditions
func ConvertToConditions(msg *hl7.Message, patientID string) ([]*fhir.Condition, error) {
	var conditions []*fhir.Condition

	dg1Segments := msg.GetSegments("DG1")
	for _, dg1 := range dg1Segments {
		condition := &fhir.Condition{
			ResourceType: "Condition",
			ID:           "condition-" + dg1.GetField(1).GetCompontent(1),
			Subject: &fhir.Reference{
				Reference: "Patient/" + patientID,
			},
		}

		//DG1-3: Diagnosis Code
		condition.Code = buildDiagnosisCode(dg1)

		//DG1-5 Diagnosis Date Time
		diagDate := dg1.GetField(5).GetCompontent(1)
		if diagDate != "" {
			condition.RecordedDate = formatDateTime(diagDate)
		}

		//DG1-6 Diagnosis Type
		condition.ClinicalStatus = mapDiagnosisType(dg1.GetField(6).GetCompontent(1))

		conditions = append(conditions, condition)
	}
	return conditions, nil
}

// buildDiagnosisCode extracts diagnosis code
func buildDiagnosisCode(dg1 *hl7.Segment) *fhir.CodeableConcept {
	codeField := dg1.GetField(3)
	if codeField == nil {
		return nil
	}

	code := codeField.GetCompontent(1)
	display := codeField.GetCompontent(2)

	if code == "" {
		return nil
	}

	//DG1-2 tells us the coding system
	codingMethod := dg1.GetField(2).GetCompontent(1)
	system := mapCodingSystem(codingMethod)

	return &fhir.CodeableConcept{
		Coding: []fhir.Coding{{
			System:  system,
			Code:    code,
			Display: display,
		}},
		Text: display,
	}
}

// mapCodingSystem converts hl7 coding method to FHIR system URK
func mapCodingSystem(method string) string {
	switch method {
	case "ICD10":
		return "http://hl7.org/fhir/sid/icd-10"
	case "ICD9":
		return "http://hl7.org/fhir/sid/icd-9-cm"
	default:
		return ""
	}
}

// mapDiagnosisType converts DG1-6 to Clinical Status
func mapDiagnosisType(diagType string) *fhir.CodeableConcept {
	return &fhir.CodeableConcept{
		Coding: []fhir.Coding{{
			System: "http://terminology.hl7.org/CodeSystem/condition-clinical",
			Code:   "active",
		}},
	}
}
