package converter

import (
	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertToEncounter converts PV1 segment to FHIR encounter
func ConvertToEncounter(msg *hl7.Message, patientID string) (*fhir.Encounter, error) {
	pv1 := msg.GetSegment("PV1")
	if pv1 == nil {
		return nil, nil
	}

	encounter := &fhir.Encounter{
		ResourceType: "Encounter",
		ID:           pv1.GetField(19).GetCompontent(1),
		Status:       "finished",
		Subject: &fhir.Reference{
			Reference: "Patient/" + patientID,
		},
	}

	//PV1-2 Patient Class
	patientClass := pv1.GetField(2).GetCompontent(1)
	encounter.Class = mapPatientClass(patientClass)

	//PV1-3 Assigned Location
	location := buildLocation(pv1)
	if location != nil {
		encounter.Location = []fhir.EncounterLocation{*location}
	}

	//PV1-7 Attending Doctor

	attendingDoc := buildAttendingDoctor(pv1)
	if attendingDoc != nil {
		encounter.Participant = []fhir.Participant{*attendingDoc}
	}

	//PV1-44 Admit DateTime
	admitField := pv1.GetField(44)
	if admitField != nil {
		admitDate := admitField.GetCompontent(1)
		if admitDate != "" {
			encounter.Period = &fhir.Period{
				Start: formatDateTime(admitDate),
			}
		}
	}

	return encounter, nil
}

//mapPatientClass converts HL7 Patient Class to FHIR

func mapPatientClass(hl7class string) *fhir.Coding {
	classMap := map[string]struct {
		code    string
		display string
	}{
		"I": {"IMP", "inpatient encounter"},
		"O": {"AMB", "ambulatory"},
		"E": {"EMER", "emergency"},
		"P": {"PRENC", "pre-admission"},
	}

	if mapped, ok := classMap[hl7class]; ok {
		return &fhir.Coding{
			System:  "http://terminology.hl7.org/CodeSystem/v3-ActCode",
			Code:    mapped.code,
			Display: mapped.display,
		}
	}

	return nil
}

// buildLocation extracts location from PV1-3
func buildLocation(pv1 *hl7.Segment) *fhir.EncounterLocation {
	locField := pv1.GetField(3)
	if locField == nil {
		return nil
	}

	//PV1-3 format
	unit := locField.GetCompontent(1)
	room := locField.GetCompontent(2)
	bed := locField.GetCompontent(3)

	if unit == "" && room == "" {
		return nil
	}

	locationDisplay := unit
	if room != "" {
		locationDisplay += " Room " + room
	}
	if bed != "" {
		locationDisplay += " Bed " + bed
	}

	return &fhir.EncounterLocation{
		Location: &fhir.Reference{
			Display: locationDisplay,
		},
		Status: "active",
	}
}

// buildAttendingDoctor
func buildAttendingDoctor(pv1 *hl7.Segment) *fhir.Participant {
	docField := pv1.GetField(7)
	if docField == nil {
		return nil
	}

	//Format ID^LastName^FirstName
	lastName := docField.GetCompontent(2)
	firstName := docField.GetCompontent(3)

	if lastName == "" {
		return nil
	}

	displayName := lastName
	if firstName != "" {
		displayName = firstName + " " + lastName
	}

	return &fhir.Participant{
		Type: []fhir.CodeableConcept{{
			Coding: []fhir.Coding{{
				System:  "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
				Code:    "ATND",
				Display: "attender",
			}},
		}},
		Individual: &fhir.Reference{
			Display: "Dr. " + displayName,
		},
	}
}

// formatDateTime converts HL7 datetime to FHIR Format
func formatDateTime(hl7DateTime string) string {
	if len(hl7DateTime) < 8 {
		return ""
	}

	//Basic YYYY-MM-DD
	result := hl7DateTime[0:4] + "-" + hl7DateTime[4:6] + "-" + hl7DateTime[6:8]

	//add time if present
	if len(hl7DateTime) >= 12 {
		result += "T" + hl7DateTime[8:10] + ":" + hl7DateTime[10:12] + ":00"
	}

	return result
}
