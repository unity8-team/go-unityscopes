package scopes

// RadioButtonsFilter is a filter that displays mutually exclusive list of options
type RadioButtonsFilter struct {
	filterWithOptions
	Label string
}

// NewRadioButtonsFilter creates a new radio button filter.
func NewRadioButtonsFilter(id, label string) *RadioButtonsFilter {
	return &RadioButtonsFilter{
		filterWithOptions: filterWithOptions{
			filterBase: filterBase{
				Id:           id,
				DisplayHints: FilterDisplayDefault,
				FilterType:   "radio_buttons",
			},
		},
		Label: label,
	}
}

// UpdateState updates the value of a particular option in the filter state.
func (f *RadioButtonsFilter) UpdateState(state FilterState, optionId string, active bool) {
	if !f.isValidOption(optionId) {
		panic("invalid option ID")
	}
	// If the state isn't in a form we expect, treat it as empty
	selected, _ := state[f.Id].([]interface{})

	if active {
		if len(selected) == 0 {
			// just add the optionId
			selected = append(selected, optionId)
		} else if len(selected) > 0 && selected[0] != optionId {
			// we have another option selected, just select the current one
			selected[0] = optionId
		}
	} else {
		if len(selected) > 0 && selected[0] == optionId {
			// we have 1 option selected and it's the current one.
			// clear the state
			selected = make([]interface{}, 0)
		}
	}
	state[f.Id] = selected
}

func (f *RadioButtonsFilter) serializeFilter() map[string]interface{} {
	v := f.filterBase.serializeFilter()
	v["label"] = f.Label
	v["options"] = f.Options
	return v
}
