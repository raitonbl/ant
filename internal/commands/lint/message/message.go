package message

const (
	UNRESOLVABLE_FIELD                          = "field is unresolvable"
	BLANK_FIELD                                 = "field mustn't be blank"
	REPEATED_VALUE                              = "value must be unique"
	REQUIRED_FIELD                              = "field is required"
	FIELD_NOT_ALLOWED                           = "field is not allowed"
	FIELD_INDEX_GT_ZERO                         = "field must be greater than zero (0)"
	FIELD_WHEN_IN_ARGUMENTS                     = "field is required when in=arguments"
	FIELD_FORMAT_NOT_ALLOWED_IN_TYPE_STRING     = "specified format cannot be applied to type=string"
	FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING = "specified format can only be applied to type=string"
	FIELD_EXAMPLE_MUST_BE_PART_OF_ENUM          = "example must be contained in Enum"
	FIELD_MAX_LENGTH_GT_ZERO                    = "max-length cannot be lesser than zero (0)"
	FIELD_MIN_LENGTH_GT_ZERO                    = "min-length cannot be lesser than zero (0)"
	FIELD_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH  = "min-length mustn't be greater than max-length"
	FIELD_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER = "specified format can only be applied to type=string"
	FIELD_MIN_MUST_NOT_BE_GT_MAX                = "minimum mustn't be greater than maximum"
	FIELD_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS    = "min-items mustn't be greater than max-items"
	FIELD_MAX_ITEMS_GT_ZERO                     = "max-items cannot be lesser than zero (0)"
	FIELD_MIN_ITEMS_GT_ZERO                     = "min-items cannot be lesser than zero (0)"
	ARRAY_FIELD_TYPE_NOT_ALLOWED                = "type not allowed on array"
	NOT_AVAILABLE_IN_USE                        = "not available since the value has been defined"
)
