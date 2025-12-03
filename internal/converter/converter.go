package converter

import (
	"github.com/mourice12/hl7-to-fhir/internal/fhir"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

// ConvertPatuebt converts an HL7 message to a FHIR patient
func ConvertToPatient(msg *hl7.Message) (*fhir.Patient, error) {
	pid := msg.GetSegment("PID")
	if pid == nil {
		return nil, nil
	}

	patient := &fhir.Patient{
		ResourceType: "Patient",
		ID:           pid.GetField(3).GetCompontent(1),
		Gender:       mapGender(pid.GetField(8).GetCompontent(1)),
		BirthDate:    formatDate(pid.GetField(7).GetCompontent(1)),
	}

	//build Identifiers

	patient.Identifier = buildIdentifiers(pid)

	//Build name
	patient.Name = buildNames(pid)

	//Build telecom
	patient.Telecom = buildTelecom(pid)

	//Build address
	patient.Address = buildAddresses(pid)

	return patient, nil
}

// mapGender converts HL7 gender codes to fhir
func mapGender(hl7Gender string) string {
	switch hl7Gender {
	case "M":
		return "male"
	case "F":
		return "female"
	case "O":
		return "other"
	default:
		return "unknown"
	}
}

// formatDate converts HL7 date (YYYYMMDD) to FHIR date (YYYY-MM-DD)
func formatDate(hl7Date string) string {
	if len(hl7Date) < 8 {
		return ""
	}
	return hl7Date[0:4] + "-" + hl7Date[4:6] + "-" + hl7Date[6:8]

}

//buildNames extracts names from PID

func buildNames(pid *hl7.Segment) []fhir.HumanName {
	names := []fhir.HumanName{}

	nameField := pid.GetField(5)
	if nameField == nil {
		return names
	}

	for _, rep := range nameField.Repetitions {
		name := fhir.HumanName{
			Family: getComponentValue(rep, 1),
			Given:  []string{},
		}

		firstName := getComponentValue(rep, 2)
		middleName := getComponentValue(rep, 3)

		if firstName != "" {
			name.Given = append(name.Given, firstName)
		}

		if middleName != "" {
			name.Given = append(name.Given, middleName)
		}

		if name.Family != "" || len(name.Given) > 0 {
			names = append(names, name)
		}
	}

	return names
}

// Build telecome extracts phone numbers
func buildTelecom(pid *hl7.Segment) []fhir.ContactPoint {
	telecoms := []fhir.ContactPoint{}

	homeField := pid.GetField((13))

	if homeField != nil {
		for _, rep := range homeField.Repetitions {
			phone := getComponentValue(rep, 1)
			if phone != "" {
				telecoms = append(telecoms, fhir.ContactPoint{
					System: "phone",
					Value:  phone,
					Use:    "home",
				})
			}
		}
	}

	workField := pid.GetField(14)
	if workField != nil {
		for _, rep := range workField.Repetitions {
			phone := getComponentValue(rep, 1)
			if phone != "" {
				telecoms = append(telecoms, fhir.ContactPoint{
					System: "phone",
					Value:  phone,
					Use:    "work",
				})
			}
		}
	}
	return telecoms
}

// buildAddresses extracts addresses
func buildAddresses(pid *hl7.Segment) []fhir.Address {
	addresses := []fhir.Address{}

	addrField := pid.GetField(11)
	if addrField == nil {
		return addresses
	}

	for _, rep := range addrField.Repetitions {
		addr := fhir.Address{
			Use:  "home",
			Type: "physical",
			Line: []string{},
		}
		street := getComponentValue(rep, 1)
		if street != "" {
			addr.Line = append(addr.Line, street)
		}

		other := getComponentValue(rep, 2)
		if other != "" {
			addr.Line = append(addr.Line, other)
		}

		addr.City = getComponentValue(rep, 3)
		addr.State = getComponentValue(rep, 4)
		addr.PostalCode = getComponentValue(rep, 5)
		addr.Country = getComponentValue(rep, 6)

		if len(addr.Line) > 0 || addr.City != "" {
			addresses = append(addresses, addr)
		}
	}

	return addresses
}

// getComponentValue safely gets a component value
func getComponentValue(rep hl7.Repetition, index int) string {
	actualIndex := index - 1
	if actualIndex < 0 || actualIndex >= len(rep.Components) {
		return ""
	}

	if len(rep.Components[actualIndex].Subcomponents) > 0 {
		return rep.Components[actualIndex].Subcomponents[0]
	}

	return ""
}

// buildIdentifiers extracts patient identifiers from PID-3
func buildIdentifiers(pid *hl7.Segment) []fhir.Identifier {
	identifiers := []fhir.Identifier{}

	idField := pid.GetField(3)
	if idField == nil {
		return identifiers
	}

	// PID-3 can have multiple repetitions (multiple IDs)
	for _, rep := range idField.Repetitions {
		id := fhir.Identifier{
			Value: getComponentValue(rep, 1), // ID value
			Use:   "usual",
		}

		// PID-3.4 is the assigning authority
		authority := getComponentValue(rep, 4)
		if authority != "" {
			id.System = "urn:oid:" + authority
		}

		// PID-3.5 is the identifier type code
		typeCode := getComponentValue(rep, 5)
		if typeCode != "" {
			id.Type = &fhir.CodeableConcept{
				Coding: []fhir.Coding{{
					System: "http://terminology.hl7.org/CodeSystem/v2-0203",
					Code:   typeCode,
				}},
			}
		}

		if id.Value != "" {
			identifiers = append(identifiers, id)
		}
	}

	return identifiers
}
