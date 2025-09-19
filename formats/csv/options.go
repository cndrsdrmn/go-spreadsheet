package csv

import "encoding/json"

type Options struct {
	BatchSize        int  `json:"batch_size"`
	Comma            rune `json:"delimiter"`
	Comment          rune `json:"comment"`
	LazyQuotes       bool `json:"lazy_quotes"`
	TrimLeadingSpace bool `json:"trim_leading_space"`
}

func (o *Options) Merge(other Options) {
	if other.Comma != 0 {
		o.Comma = other.Comma
	}
	if other.Comment != 0 {
		o.Comment = other.Comment
	}
	if other.BatchSize > 0 {
		o.BatchSize = other.BatchSize
	}
	o.LazyQuotes = other.LazyQuotes
	o.TrimLeadingSpace = other.TrimLeadingSpace
}

func (o *Options) UnmarshalJSON(data []byte) error {
	type Alias Options
	aux := &struct {
		Comma   string `json:"delimiter"`
		Comment string `json:"comment"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Comma != "" {
		o.Comma = []rune(aux.Comma)[0]
	}
	if aux.Comment != "" {
		o.Comment = []rune(aux.Comment)[0]
	}

	return nil
}
