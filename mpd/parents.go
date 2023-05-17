package mpd

// SetParents sets the parent pointers for all children.
// Call this after xml.Unmarshal of MPD to set the parent pointers.
func (m *MPD) SetParents() {
	for _, p := range m.Periods {
		p.SetParent(m)
		for _, a := range p.AdaptationSets {
			a.SetParent(p)
			for _, r := range a.Representations {
				r.SetParent(a)
				for _, sr := range r.SubRepresentations {
					sr.SetParent(r)
				}
			}
		}
	}
}

// AppendPeriod appends a Period to the MPD and sets parent pointer.
func (m *MPD) AppendPeriod(p *Period) {
	m.Periods = append(m.Periods, p)
	p.SetParent(m)
}

// SetParents sets the parent pointers for all children of p.
func (p *Period) SetParents() {
	for _, a := range p.AdaptationSets {
		a.SetParent(p)
		for _, r := range a.Representations {
			r.SetParent(a)
			for _, sr := range r.SubRepresentations {
				sr.SetParent(r)
			}
		}
	}
}

// AppendAdaptationSet appends an AdaptationSet to the Period and sets parent pointer.
func (p *Period) AppendAdaptationSet(a *AdaptationSetType) {
	p.AdaptationSets = append(p.AdaptationSets, a)
	a.SetParent(p)
}

// SetParents sets the parent pointers for all children of a.
func (a *AdaptationSetType) SetParents() {
	for _, r := range a.Representations {
		r.SetParent(a)
		for _, sr := range r.SubRepresentations {
			sr.SetParent(r)
		}
	}
}

// AppendRepresentation appends a Representation to the AdaptationSet and sets parent pointer.
func (a *AdaptationSetType) AppendRepresentation(r *RepresentationType) {
	a.Representations = append(a.Representations, r)
	r.SetParent(a)
}

// SetParents sets the parent pointers for all children of r.
func (r *RepresentationType) SetParents() {
	for _, sr := range r.SubRepresentations {
		sr.SetParent(r)
	}
}

// AppendSubRepresentation appends a SubRepresentation to the Representation and sets parent pointer.
func (r *RepresentationType) AppendSubRepresentation(sr *SubRepresentationType) {
	r.SubRepresentations = append(r.SubRepresentations, sr)
	sr.SetParent(r)
}
