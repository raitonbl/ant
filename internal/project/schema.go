package project

type Schema struct {
	// applies to number
	MultipleOf *int `yaml:"multiple-of" json:"multiple-of,omitempty"`
	Maximum    *int `yaml:"maximum" json:"maximum,omitempty"`
	Minimum    *int `yaml:"minimum" json:"minimum,omitempty"`
	// applies to string
	MaxLength *int `yaml:"max-length" json:"max-length,omitempty"`
	MinLength *int `yaml:"min-length" json:"min-length,omitempty"`
	// applies to number
	ExclusiveMinimum *bool `yaml:"exclusive-minimum" json:"exclusive-minimum,omitempty"`
	ExclusiveMaximum *bool `yaml:"exclusive-maximum" json:"exclusive-maximum,omitempty"`
	// applies to array
	MaxItems    *int    `yaml:"max-items" json:"max-items,omitempty"`
	MinItems    *int    `yaml:"min-items" json:"min-items,omitempty"`
	UniqueItems *bool   `yaml:"unique-items" json:"unique-items,omitempty"`
	Items       *Schema `yaml:"items" json:"items,omitempty"`
	// applies to everything
	Enum []string `yaml:"enum" json:"enum,omitempty"`
	// object
	Examples []string      `yaml:"examples" json:"examples,omitempty"`
	TypeOf   *SchemaType   `yaml:"type" json:"type,omitempty"`
	Format   *SchemaFormat `yaml:"format" json:"format,omitempty"`
	Pattern  *string       `yaml:"pattern" json:"pattern,omitempty"`
	RefersTo *string       `yaml:"refers-to" json:"refers-to,omitempty"`
}

type SchemaType string

const (
	String SchemaType = "string"
	Number SchemaType = "number"
	Bool   SchemaType = "boolean"
	Array  SchemaType = "array"
)

type SchemaFormat string

const (
	Byte   SchemaFormat = "byte"
	Int64  SchemaFormat = "int64"
	Int32  SchemaFormat = "int32"
	Float  SchemaFormat = "float"
	Double SchemaFormat = "double"

	Date     SchemaFormat = "date"
	Binary   SchemaFormat = "binary"
	DateTime SchemaFormat = "datetime"
)
