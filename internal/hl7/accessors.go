package hl7

// GetSegment returns the first segment with the given name
func (m *Message) GetSegment(name string) *Segment {
	for i := range m.Segments {
		if m.Segments[i].Name == name {
			return &m.Segments[i]
		}
	}
	return nil
}

// GetField returns the field at the given index

func (s *Segment) GetField(index int) *Field {
	actualIndex := index - 1
	if actualIndex < 0 || actualIndex >= len(s.Fields) {
		return nil
	}
	return &s.Fields[actualIndex]
}

//GetRepetition returns a specific repetition

func (f *Field) GetRepetition(index int) *Repetition {
	actualIndex := index - 1
	if actualIndex < 0 || actualIndex >= len(f.Repetitions) {
		return nil
	}
	return &f.Repetitions[actualIndex]
}

// GetCompontent returns the component at the given index
func (f *Field) GetCompontent(index int) string {
	if len(f.Repetitions) == 0 {
		return ""
	}
	return f.Repetitions[0].GetCompontent(index)

}

func (r *Repetition) GetCompontent(index int) string {
	actualIndex := index - 1
	if actualIndex < 0 || actualIndex >= len(r.Components) {
		return ""
	}
	if len(r.Components[actualIndex].Subcomponents) > 0 {
		return r.Components[actualIndex].Subcomponents[0]
	}
	return ""
}

func (c *Component) GetCompontent(index int) string {
	actualIndex := index - 1
	if actualIndex < 0 || actualIndex >= len(c.Subcomponents) {
		return ""
	}
	return c.Subcomponents[actualIndex]
}

// GetSegments returns all segments with a given name
func (m *Message) GetSegments(name string) []*Segment {
	var segments []*Segment
	for i := range m.Segments {
		if m.Segments[i].Name == name {
			segments = append(segments, &m.Segments[i])
		}
	}
	return segments
}
