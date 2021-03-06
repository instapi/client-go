package schema

// Schema represents a schema.
type Schema struct {
	Name       string    `json:"name"`
	PrimaryKey []string  `json:"primaryKey,omitempty"`
	Indexed    []string  `json:"indexed,omitempty"`
	Count      int       `json:"count,omitempty"`
	Settings   *Settings `json:"settings,omitempty"`
	Fields     []Field   `json:"fields"`
}

// Settings represents schema settings.
type Settings struct {
	Delimiter   *string `json:"delimiter,omitempty"`
	QuoteValues *string `json:"quoteValues,omitempty"`
}

// Field represents a schema field.
type Field struct {
	Name     string  `json:"name"`
	Mapping  string  `json:"mapping,omitempty"`
	Type     string  `json:"type"`
	Format   *string `json:"format,omitempty"`
	Required bool    `json:"required,omitempty"`
}
