package structure

type Schema struct {
	MultipleOf       *int          `yaml:"multiple-of"`
	Maximum          *int          `yaml:"maximum"`
	Minimum          *int          `yaml:"minimum"`
	MaxLength        *int          `yaml:"max-length"`
	MinLength        *int          `yaml:"min-length"`
	ExclusiveMinimum *bool         `yaml:"exclusive-minimum"`
	ExclusiveMaximum *bool         `yaml:"exclusive-maximum"`
	MaxItems         *int          `yaml:"max-items"`
	MinItems         *int          `yaml:"min-items"`
	UniqueItems      *bool         `yaml:"unique-items"`
	Enum             []string      `yaml:"enum"`
	Pattern          *string       `yaml:"pattern"`
	Examples         *string       `yaml:"examples"`
	TypeOf           *SchemaType   `yaml:"type"`
	Format           *SchemaFormat `yaml:"format"`
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
