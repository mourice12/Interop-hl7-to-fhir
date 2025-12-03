package converter

import (
	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertToAllergies converts AL1 segments to FHIR AllergyIntolerance
func ConvertToAllergies(msg *hl7.Message, patientID string) ([]*fhir.AllergyIntolerance, error) {
	var allergies []*fhir.AllergyIntolerance

	allSegments := msg.GetSegments("AL1")

	for _, al1 := range allSegments {
		allergy := &fhir.AllergyIntolerance{
			ResourceType: "AllergyIntolerance",
			ID:           "allergy-" + al1.GetField(1).GetCompontent(1),
			Type:         "allergy",
			Patient: &fhir.Reference{
				Reference: "Patient/" + patientID,
			},
			ClinicalStatus: &fhir.CodeableConcept{
				Coding: []fhir.Coding{{
					System: "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical",
					Code:   "active",
				}},
			},
		}

		//AL1-2 Allergy Type
		allergy.Category = mapAllergyCategory(al1.GetField(2).GetCompontent(1))

		//AL1-3 Allergen Code
		allergy.Code = buildAllergenCode(al1)

		//AL5 Reactions
		allergy.Reaction = buildReactions(al1)

		//AL1-6 Identification date

		identDate := al1.GetField(6).GetCompontent(1)
		if identDate != "" {
			allergy.RecordedDate = formatDateTime(identDate)
		}

		allergies = append(allergies, allergy)
	}

	return allergies, nil

}

// mapAllergyCategory converts AL1-2 to FHIR
func mapAllergyCategory(al1Type string) []string {
	switch al1Type {
	case "DA":
		return []string{"medication"}
	case "FA":
		return []string{"food"}
	case "EA":
		return []string{"environment"}
	default:
		return nil
	}
}

// buildAllergenCode extracts allergen from AL1-3
func buildAllergenCode(al1 *hl7.Segment) *fhir.CodeableConcept {
	allergenField := al1.GetField(3)
	if allergenField == nil {
		return nil
	}

	//Format code^description

	code := allergenField.GetCompontent(1)
	description := allergenField.GetCompontent(2)

	//use description if code is empty
	displayText := description
	if displayText == "" {
		displayText = code
	}

	if displayText == "" {
		return nil
	}

	return &fhir.CodeableConcept{
		Text: displayText,
	}
}

// buildReactions extracts reactions for AL1-5
func buildReactions(al1 *hl7.Segment) []fhir.AllergyReaction {
	reactionField := al1.GetField(5)
	if reactionField == nil {
		return nil
	}

	var manifestations []fhir.CodeableConcept

	//al1-5 can have multiple reactions via repetition
	for _, rep := range reactionField.Repetitions {
		reactionText := ""
		if len(rep.Components) > 0 && len(rep.Components[0].Subcomponents) > 0 {
			reactionText = rep.Components[0].Subcomponents[0]
		}
		if reactionText != "" {
			manifestations = append(manifestations, fhir.CodeableConcept{
				Text: reactionText,
			})
		}
	}

	if len(manifestations) == 0 {
		return nil
	}

	return []fhir.AllergyReaction{{
		Manifestation: manifestations,
	}}
}
