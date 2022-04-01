package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictGet(t *testing.T) {
	t.Run("test Get", func(t *testing.T) {
		d := Dict{
			"a": Dict{
				"b": Dict{
					"c": "d",
				},
			},
		}
		got := d.Get("a", "b", "c")
		assert.Equal(t, "d", got)
	})

	t.Run("test Get with nil", func(t *testing.T) {
		d := Dict{
			"a": Dict{
				"b": Dict{
					"c": "d",
				},
			},
		}
		got := d.Get("b")
		assert.Equal(t, nil, got)

	})
}

func TestDictToString(t *testing.T) {
	t.Run("test ToString", func(t *testing.T) {
		d := Dict{
			"a": Dict{
				"b": Dict{
					"c": "d",
				},
			},
		}
		got := d.ToString()
		assert.Equal(t, `{"a":{"b":{"c":"d"}}}`, got)
	})
}
