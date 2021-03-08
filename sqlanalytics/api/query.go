package api

import (
	"encoding/json"
)

// Query ...
type Query struct {
	ID             string            `json:"id,omitempty"`
	DataSourceID   string            `json:"data_source_id"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	Query          string            `json:"query"`
	Schedule       *QuerySchedule    `json:"schedule,omitempty"`
	Options        *QueryOptions     `json:"options,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	Visualizations []json.RawMessage `json:"visualizations,omitempty"`
}

// QuerySchedule ...
type QuerySchedule struct {
	Interval int `json:"interval"`
}

// QueryOptions ...
type QueryOptions struct {
	Parameters    []interface{}     `json:"-"`
	RawParameters []json.RawMessage `json:"parameters,omitempty"`
}

// MarshalJSON ...
func (o *QueryOptions) MarshalJSON() ([]byte, error) {
	if o.Parameters != nil {
		o.RawParameters = []json.RawMessage{}
		for _, p := range o.Parameters {
			b, err := json.Marshal(p)
			if err != nil {
				return nil, err
			}
			o.RawParameters = append(o.RawParameters, b)
		}
	}

	type localQueryOptions QueryOptions
	return json.Marshal((*localQueryOptions)(o))
}

// UnmarshalJSON ...
func (o *QueryOptions) UnmarshalJSON(b []byte) error {
	type localQueryOptions QueryOptions
	err := json.Unmarshal(b, (*localQueryOptions)(o))
	if err != nil {
		return err
	}

	o.Parameters = []interface{}{}
	for _, rp := range o.RawParameters {
		var qp QueryParameter

		// Unmarshal into base parameter type to figure out the right type.
		err = json.Unmarshal(rp, &qp)
		if err != nil {
			return err
		}

		// Acquire pointer to the correct parameter type.
		var i interface{}
		switch qp.Type {
		case "text":
			i = &QueryParameterText{}
		case "number":
			i = &QueryParameterNumber{}
		case "enum":
			i = &QueryParameterEnum{}
		case "query":
			i = &QueryParameterQuery{}
		default:
			panic("don't know what to do...")
		}

		// Unmarshal into correct parameter type.
		err = json.Unmarshal(rp, &i)
		if err != nil {
			return err
		}

		o.Parameters = append(o.Parameters, i)
	}

	return nil
}

// QueryParameter ...
type QueryParameter struct {
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// QueryParameterText ...
type QueryParameterText struct {
	QueryParameter
}

// MarshalJSON sets the type before marshaling.
func (p QueryParameterText) MarshalJSON() ([]byte, error) {
	p.QueryParameter.Type = "text"
	type localQueryParameter QueryParameterText
	return json.Marshal((localQueryParameter)(p))
}

// QueryParameterNumber ...
type QueryParameterNumber struct {
	QueryParameter
}

// MarshalJSON sets the type before marshaling.
func (p QueryParameterNumber) MarshalJSON() ([]byte, error) {
	p.QueryParameter.Type = "number"
	type localQueryParameter QueryParameterNumber
	return json.Marshal((localQueryParameter)(p))
}

// QueryParameterMultipleValuesOptions ...
type QueryParameterMultipleValuesOptions struct {
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Separator string `json:"separator"`
}

// QueryParameterEnum ...
type QueryParameterEnum struct {
	QueryParameter
	Options string                               `json:"enumOptions"`
	Multi   *QueryParameterMultipleValuesOptions `json:"multiValuesOptions,omitempty"`
}

// MarshalJSON sets the type before marshaling.
func (p QueryParameterEnum) MarshalJSON() ([]byte, error) {
	p.QueryParameter.Type = "enum"
	type localQueryParameter QueryParameterEnum
	return json.Marshal((localQueryParameter)(p))
}

// QueryParameterQuery ...
type QueryParameterQuery struct {
	QueryParameter
	QueryID string                               `json:"queryId"`
	Multi   *QueryParameterMultipleValuesOptions `json:"multiValuesOptions,omitempty"`
}

// MarshalJSON sets the type before marshaling.
func (p QueryParameterQuery) MarshalJSON() ([]byte, error) {
	p.QueryParameter.Type = "query"
	type localQueryParameter QueryParameterQuery
	return json.Marshal((localQueryParameter)(p))
}
