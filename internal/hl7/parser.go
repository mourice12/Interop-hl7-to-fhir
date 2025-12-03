package hl7

import (
	"errors"
	"strings"
)

// Delimiters hold the HL7 separator characters
type Delimiters struct {
	Field        string // |
	Component    string // ^
	Repetition   string // ~
	Escape       string // \
	Subcomponent string // &
}

//DefaultDelimiters returns standard HL7 delimiters

func DefaultDelimiters() Delimiters {
	return Delimiters{
		Field:        "|",
		Component:    "^",
		Repetition:   "~",
		Escape:       "\\",
		Subcomponent: "&",
	}
}

// Parse takes a raw HL7 message string and returns a message struct
func Parse(raw string) (*Message, error) {
	//normaline line endings
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")

	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		return nil, errors.New("empty message")
	}
	if !strings.HasPrefix(lines[0], "MSH") {
		return nil, errors.New("message must start with MSH segment")
	}

	delimiters := DefaultDelimiters()
	message := &Message{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		segment, err := parseSegment(line, delimiters)
		if err != nil {
			return nil, err
		}
		message.Segments = append(message.Segments, segment)
	}
	return message, nil
}

// parseSegment parses a single line
func parseSegment(line string, delim Delimiters) (Segment, error) {
	parts := strings.Split(line, delim.Field)

	segmentName := parts[0]
	segment := Segment{
		Name:   segmentName,
		Fields: []Field{},
	}

	//MSH is special: MSH-1 is the field separator itself
	if segmentName == "MSH" {
		//MSH-1 = field separator "|"
		segment.Fields = append(segment.Fields, Field{
			Repetitions: []Repetition{{
				Components: []Component{{
					Subcomponents: []string{delim.Field},
				}},
			}},
		})
		//MSH-2 onwards = rest of fields starting at parts[1]
		for i := 1; i < len(parts); i++ {
			field := parseField(parts[i], delim)
			segment.Fields = append(segment.Fields, field)
		}
	} else {
		for i := 1; i < len(parts); i++ {
			field := parseField(parts[i], delim)
			segment.Fields = append(segment.Fields, field)
		}

	}
	return segment, nil
}

func parseField(fieldStr string, delim Delimiters) Field {
	field := Field{
		Repetitions: []Repetition{},
	}
	//Splits repeition delimiter
	repParts := strings.Split(fieldStr, delim.Repetition)

	for _, repStr := range repParts {
		rep := Repetition{
			Components: []Component{},
		}

		//Split the component delimiter
		compParts := strings.Split(repStr, delim.Component)

		for _, compStr := range compParts {
			//split by subcomponent delim
			subParts := strings.Split(compStr, delim.Subcomponent)

			comp := Component{
				Subcomponents: subParts,
			}
			rep.Components = append(rep.Components, comp)
		}
		field.Repetitions = append(field.Repetitions, rep)
	}

	return field
}
