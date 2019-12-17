package schema

// Schema represents a schema.
type Schema struct {
	Name       string    `json:"name"`
	PrimaryKey []string  `json:"primaryKey,omitempty"`
	Count      int       `json:"count,omitempty"`
	Settings   *Settings `json:"settings,omitempty"`
	Fields     []Field   `json:"fields"`
}

// Settings represents schema settings.
type Settings struct {
	QuoteValues *string `json:"quoteValues,omitempty"`
}

// Field represents a schema field.
type Field struct {
	Name     string `json:"name"`
	Mapping  string `json:"mapping,omitempty"`
	Type     string `json:"type"`
	Required bool   `json:"required,omitempty"`
}
