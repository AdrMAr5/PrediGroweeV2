package models

import (
	"encoding/json"
	"io"
)

type Case struct {
	ID              int              `json:"id"`
	Code            string           `json:"code"`
	Gender          string           `json:"gender"`
	Age1            int              `json:"age1"`
	Age2            int              `json:"age2"`
	Image1          string           `json:"image1"`
	Image2          string           `json:"image2"`
	Parameters      []Parameter      `json:"parameters"`
	ParameterValues []ParameterValue `json:"parameters_values"`
}

func (c *Case) ToJSON(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(c)
}
func (c *Case) FromJSON(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(c)
}
