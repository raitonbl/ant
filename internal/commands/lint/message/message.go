package message

const (
	BLANK_FIELD_MESSAGE = "field mustn't be blank"
	REPEATED_VALUE_MESSAGE                          = "value must be unique"
	REQUIRED_FIELD_MESSAGE                          = "field is required"
	FIELD_NOT_ALLOWED                               = "field is not allowed"
	PARAMETER_INDEX_GT_ZERO_MESSAGE                 = "field must be greater than zero (0)"
	PARAMETER_FIELD_NOT_ALLOWED_IN_FLAGS            = "field is not allowed when in=flags"
	PARAMETER_FIELD_NOT_ALLOWED_IN_ARGUMENTS        = "field is not allowed when in=arguments"
	REQUIRED_PARAMETER_FIELD_WHEN_IN_ARGUMENTS      = "field is required when in=arguments"
	PAREMETER_FORMAT_NOT_ALLOWED_IN_TYPE_STRING     = "specified format cannot be applied to type=string"
	PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_STRING = "specified format can only be applied to type=string"
	PARAMETER_EXAMPLE_MUST_BE_PART_OF_ENUM          = "example must be contained in Enum"
	PARAMETER_MAX_LENGTH_GT_ZERO                    = "max-length cannot be lesser than zero (0)"
	PARAMETER_MIN_LENGTH_GT_ZERO                    = "min-length cannot be lesser than zero (0)"
	PARAMETER_MIN_LENGTH_MUST_NOT_BE_GT_MAX_LENGTH  = "min-length mustn't be greater than max-length"
	PAREMETER_FORMAT_IS_ONLY_ALLOWED_IN_TYPE_NUMBER = "specified format can only be applied to type=string"
	PARAMETER_MIN_MUST_NOT_BE_GT_MAX                = "minimum mustn't be greater than maximum"
	PARAMETER_MIN_ITEMS_MUST_NOT_BE_GT_MAX_ITEMS    = "min-items mustn't be greater than max-items"
	PARAMETER_MAX_ITEMS_GT_ZERO                     = "max-items cannot be lesser than zero (0)"
	PARAMETER_MIN_ITEMS_GT_ZERO                     = "min-items cannot be lesser than zero (0)"
	PARAMETER_TYPE_NOT_ALLOWED                      = "type not allowed on array"
)
