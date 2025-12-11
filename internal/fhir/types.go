package fhir

//Patient represents a FHIR R4 patient resource

type Patient struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id,omitempty"`
	Identifier   []Identifier   `json:"identifier,omitempty"`
	Name         []HumanName    `json:"name,omitempty"`
	Telecom      []ContactPoint `json:"telecom,omitempty"`
	Gender       string         `json:"gender,omitempty"`
	BirthDate    string         `json:"birthDate,omitempty"`
	Address      []Address      `json:"address,omitempty"`
}

//Identifier represents a FHIR Identifier

type Identifier struct {
	Use    string           `json:"use,omitempty"`
	Type   *CodeableConcept `json:"type,omitempty"`
	System string           `json:"system,omitempty"`
	Value  string           `json:"value,omitempty"`
}

//CodeableConcept represents a FHIR CodeableConcept

type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

//Coding represents a FHIR coding

type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

// HumanName respresents a persons name in FHIR
type HumanName struct {
	Family string   `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
}

// ContactPoint represents phone/email
type ContactPoint struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
	Use    string `json:"use,omitempty"`
}

// Address represents a FHIR address
type Address struct {
	Use        string   `json:"use,omitempty"`  // home, work, temp
	Type       string   `json:"type,omitempty"` // postal, physical
	Line       []string `json:"line,omitempty"` // Street address lines
	City       string   `json:"city,omitempty"`
	State      string   `json:"state,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
}

//Bundle represents a FHIR Bundle

type Bundle struct {
	ResourceType string        `json:"resourceType"`
	Type         string        `json:"type"`
	Entry        []BundleEntry `json:"entry,omitempty"`
}

type BundleEntry struct {
	FullURL  string      `json:"fullUrl,omitempty"`
	Resource interface{} `json:"resource"`
}

//Encounter represents a FHIR encounter resource

type Encounter struct {
	ResourceType string              `json:"resourceType"`
	ID           string              `json:"id,omitempty"`
	Status       string              `json:"status"` // planned, arrived, in-progress, finished
	Class        *Coding             `json:"class,omitempty"`
	Type         []CodeableConcept   `json:"type,omitempty"`
	Subject      *Reference          `json:"subject,omitempty"` // Reference to Patient
	Participant  []Participant       `json:"participant,omitempty"`
	Period       *Period             `json:"period,omitempty"`
	Location     []EncounterLocation `json:"location,omitempty"`
}

// Reference is a FHIR reference to another resource
type Reference struct {
	Reference string `json:"reference,omitempty"`
	Display   string `json:"display,omitempty"`
}

// Period represents a time period
type Period struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// Participant represents someone involved in the encounter
type Participant struct {
	Type       []CodeableConcept `json:"type,omitempty"`
	Individual *Reference        `json:"individual,omitempty"`
}

type EncounterLocation struct {
	Location *Reference `json:"location,omitempty"`
	Status   string     `json:"status,omitempty"`
}

// Condition represents a FHIR condition
type Condition struct {
	ResourceType   string           `json:"resourceType"`
	ID             string           `json:"id,omitempty"`
	ClinicalStatus *CodeableConcept `json:"clinicalStatus,omitempty"`
	Code           *CodeableConcept `json:"code,omitempty"`
	Subject        *Reference       `json:"subject,omitempty"`
	RecordedDate   string           `json:"recordedDate,omitempty"`
}

//AllergyIntolerance represents a FHIR AllergyIntolerance Resource

type AllergyIntolerance struct {
	ResourceType   string            `json:"resourceType"`
	ID             string            `json:"id,omitempty"`
	ClinicalStatus *CodeableConcept  `json:"clinicalStatus,omitempty"`
	Type           string            `json:"type,omitempty"`     // allergy | intolerance
	Category       []string          `json:"category,omitempty"` // food | medication | environment
	Code           *CodeableConcept  `json:"code,omitempty"`
	Patient        *Reference        `json:"patient,omitempty"`
	RecordedDate   string            `json:"recordedDate,omitempty"`
	Reaction       []AllergyReaction `json:"reaction,omitempty"`
}

// AllergyReaction represents a reaction to an allergen
type AllergyReaction struct {
	Manifestation []CodeableConcept `json:"manifestation,omitempty"`
}

// Oberservation represents a FHIR Observation resource
type Oberservation struct {
	ResourceType      string           `json:"resourceType"`
	ID                string           `json:"id,omitempty"`
	Status            string           `json:"status"` // final, preliminary
	Code              *CodeableConcept `json:"code,omitempty"`
	Subject           *Reference       `json:"subject,omitempty"`
	EffectiveDateTime string           `json:"effectiveDateTime,omitempty"`
	ValueQuantity     *Quantity        `json:"valueQuantity,omitempty"`
	ReferenceRange    []ReferenceRange `json:"referenceRange,omitempty"`
}

// Quanity represents a FHIR Quantity
type Quantity struct {
	Value  float64 `json:"value,omitempty"`
	Unit   string  `json:"unit,omitempty"`
	System string  `json:"system,omitempty"`
	Code   string  `json:"code,omitempty"`
}

// ReferenceRange represents normal ranges
type ReferenceRange struct {
	Text string `json:"text,omitempty"`
}

// DiagnosticReport represents a diagnostic report
type DiagnosticReport struct {
	ResourceType      string           `json:"resourceType"`
	ID                string           `json:"id,omitempty"`
	Status            string           `json:"status"` // final, preliminary
	Code              *CodeableConcept `json:"code,omitempty"`
	Subject           *Reference       `json:"subject,omitempty"`
	EffectiveDateTime string           `json:"effectiveDateTime,omitempty"`
	Issued            string           `json:"issued,omitempty"`
	Performer         []Reference      `json:"performer,omitempty"`
	Result            []Reference      `json:"result,omitempty"`
}
