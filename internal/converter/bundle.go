package converter

import (
	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertToBundle converts HL7 message to FHIR Bundle
func ConvertToBundle(msg *hl7.Message) (*fhir.Bundle, error) {
	bundle := fhir.NewBundle()

	//Convert Patient
	patient, err := ConvertToPatient(msg)
	if err != nil {
		return nil, err
	}
	if patient != nil {
		bundle.AddEntry("Patient", patient.ID, patient)
	}
	//Convert Encounter
	if patient != nil {
		encounter, err := ConvertToEncounter(msg, patient.ID)
		if err != nil {
			return nil, err
		}
		if encounter != nil {
			bundle.AddEntry("Encounter", encounter.ID, encounter)
		}
	}

	//Convert Conditions
	if patient != nil {
		conditions, err := ConvertToConditions(msg, patient.ID)
		if err != nil {
			return nil, err
		}

		for _, condition := range conditions {
			bundle.AddEntry("Condition", condition.ID, condition)
		}
	}

	//Convert Allergies
	if patient != nil {
		allergies, err := ConvertToAllergies(msg, patient.ID)
		if err != nil {
			return nil, err
		}

		for _, allergy := range allergies {
			bundle.AddEntry("AllergyIntolerance", allergy.ID, allergy)
		}
	}

	//Convert Observations
	if patient != nil {
		observations, err := ConvertToObservations(msg, patient.ID)
		if err != nil {
			return nil, err
		}

		for _, obs := range observations {
			bundle.AddEntry("Observation", obs.ID, obs)
		}
	}

	return bundle, nil
}
