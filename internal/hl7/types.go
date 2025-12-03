package hl7

// Message represents a HL7 message
type Message struct {
	Segments []Segment
}

// Segments respresents a single line like MSH, PID, etc.
type Segment struct {
	Name   string
	Fields []Field
}

//Fields can have multiple repititions

type Field struct {
	Repetitions []Repetition
}

//Repetition contains component

type Repetition struct {
	Components []Component
}

// Component represents a single component
type Component struct {
	Subcomponents []string
}
