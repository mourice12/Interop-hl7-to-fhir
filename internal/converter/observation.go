package converter

import (
	"strconv"

	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertToObservations converts OBX segments to FHIR Observations
func ConvertToObservations(msg *hl7.Message, patientID string) ([]*fhir.Oberservation, error) {
	var observations []*fhir.Oberservation

	obxSegments := msg.GetSegments("OBX")

	for _, obx := range obxSegments {
		obs := &fhir.Oberservation{
			ResourceType: "Observation",
			ID:           "observation-" + obx.GetField(1).GetCompontent(1),
			Subject: &fhir.Reference{
				Reference: "Patient/" + patientID,
			},
		}

		// OBX-3 observation ID
		obs.Code = buildObservationCode(obx)

		//OBX-5 value
		obs.ValueQuantity = buildValueQuantity(obx)

		//OBX-7 Reference Range
		refRange := obx.GetField(7).GetCompontent(1)
		if refRange != "" {
			obs.ReferenceRange = []fhir.ReferenceRange{{Text: refRange}}
		}

		//OBX-11 Status
		obs.Status = mapObservationStatus(obx.GetField(11).GetCompontent(1))

		//OBX-14 DateTime
		obsDateTime := obx.GetField(14).GetCompontent(1)
		if obsDateTime != "" {
			obs.EffectiveDateTime = formatDateTime(obsDateTime)
		}

		observations = append(observations, obs)
	}

	return observations, nil
}

// buildObservationCode extracts LOINC code
func buildObservationCode(obx *hl7.Segment) *fhir.CodeableConcept {
	codeField := obx.GetField(3)

	if codeField == nil {
		return nil
	}

	code := codeField.GetCompontent(1)
	display := codeField.GetCompontent(2)
	system := codeField.GetCompontent(3)

	//map common coding systems

	systemURL := ""

	if system == "LN" {
		systemURL = "http://loinc.org"
	}

	return &fhir.CodeableConcept{
		Coding: []fhir.Coding{{
			System:  systemURL,
			Code:    code,
			Display: display,
		}},
		Text: display,
	}
}

// buildValueQuantity extracts value
func buildValueQuantity(obx *hl7.Segment) *fhir.Quantity {
	valueStr := obx.GetField(5).GetCompontent(1)
	unit := obx.GetField(6).GetCompontent(1)

	if valueStr == "" {
		return nil
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil
	}

	return &fhir.Quantity{
		Value: value,
		Unit:  unit,
	}
}

// mapObservationsStatus converts
func mapObservationStatus(status string) string {
	switch status {
	case "F":
		return "final"
	case "P":
		return "preliminary"
	case "C":
		return "corrected"
	default:
		return "unknown"
	}
}
