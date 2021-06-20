package schema

// Schema represents a schema.
type Schema struct {
	Name       string    `json:"name"`
	PrimaryKey []string  `json:"primaryKey,omitempty"`
	Indexed    []string  `json:"indexed,omitempty"`
	Count      int       `json:"count,omitempty"`
	Settings   *Settings `json:"settings,omitempty"`
	Fields     []*Field  `json:"fields"`
}

// Settings represents schema settings.
type Settings struct {
	ExternalID     string `json:"externalId,omitempty"`
	Sheet          string `json:"sheet,omitempty"`
	Delimiter      string `json:"delimiter,omitempty"`
	QuoteValues    string `json:"quoteValues,omitempty"`
	InputFunction  string `json:"inputFunction,omitempty"`
	OutputFunction string `json:"outputFunction,omitempty"`
}

// Field represents a schema field.
type Field struct {
	Name     string  `json:"name"`
	Mapping  string  `json:"mapping,omitempty"`
	Type     string  `json:"type"`
	Format   *string `json:"format,omitempty"`
	Required bool    `json:"required,omitempty"`
}

// Import represents a schema import.
type Import struct {
	Count  int     `json:"count"`
	Schema *Schema `json:"schema"`
}
