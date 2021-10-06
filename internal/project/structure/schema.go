package structure

type Schema struct {
	MultipleOf       *int          `yaml:"multiple-of" json:"multiple-of,omitempty"`
	Maximum          *int          `yaml:"maximum" json:"maximum,omitempty"`
	Minimum          *int          `yaml:"minimum" json:"minimum,omitempty"`
	MaxLength        *int          `yaml:"max-length" json:"max-length,omitempty"`
	MinLength        *int          `yaml:"min-length" json:"min-length,omitempty"`
	ExclusiveMinimum *bool         `yaml:"exclusive-minimum" json:"exclusive-minimum,omitempty"`
	ExclusiveMaximum *bool         `yaml:"exclusive-maximum" json:"exclusive-maximum,omitempty"`
	MaxItems         *int          `yaml:"max-items" json:"max-items,omitempty"`
	MinItems         *int          `yaml:"min-items" json:"min-items,omitempty"`
	UniqueItems      *bool         `yaml:"unique-items" json:"unique-items,omitempty"`
	Enum             []string      `yaml:"enum" json:"enum,omitempty"`
	Pattern          *string       `yaml:"pattern" json:"pattern,omitempty"`
	Examples         []string      `yaml:"examples" json:"examples,omitempty"`
	TypeOf           *SchemaType   `yaml:"type" json:"type,omitempty"`
	Format           *SchemaFormat `yaml:"format" json:"format,omitempty"`
}

type SchemaType string

const (
	String SchemaType = "string"
	Number SchemaType = "number"
	Bool   SchemaType = "boolean"
)

type SchemaFormat string

const (
	None     SchemaFormat = ""
	Byte     SchemaFormat = "byte"
	Date     SchemaFormat = "date"
	Int64    SchemaFormat = "int64"
	Int32    SchemaFormat = "int32"
	Float    SchemaFormat = "float"
	Double   SchemaFormat = "double"
	Binary   SchemaFormat = "binary"
	DateTime SchemaFormat = "datetime"
	Pattern  SchemaFormat = "pattern"
)
