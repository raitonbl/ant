package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_is_blank_where_value_is_empty(t *testing.T) {
	assert.Equal(t, true, IsBlank(""))
}

func Test_is_blank_where_value_is_blank(t *testing.T) {
	assert.Equal(t, true, IsBlank(" "))
}

func Test_is_blank_where_value_isnt_blank(t *testing.T) {
	assert.Equal(t, false, IsBlank("A"))
}
