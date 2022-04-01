package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUniqueStrinng(t *testing.T) {

	t.Run("should return a string of the given length", func(t *testing.T) {
		str := NewUniqueString(10)
		assert.Len(t, str, 10)
	})

	t.Run("should return unique string each time", func(t *testing.T) {
		str1 := NewUniqueString(10)
		str2 := NewUniqueString(10)
		assert.NotEqual(t, str1, str2)
	})

}

func TestCopyNonEmptyField(t *testing.T) {
	t.Run("should be combined", func(t *testing.T) {
		name := "test"
		uid := "id"
		src := &AugmontUserInfo{
			Name: name,
		}

		dest := &AugmontUserInfo{
			UniqueID: uid,
		}
		CopyNonEmptyFiled(src, dest)
		assert.Equal(t, uid, dest.UniqueID)
		assert.Equal(t, name, dest.Name)
	})

	t.Run("should be overwrite old values", func(t *testing.T) {
		src := &AugmontUserInfo{
			UniqueID:   "src_id",
			Name:       "src_name",
			DOB:        "src_dob",
			NomineeDOB: "src_dob",
		}

		dest := &AugmontUserInfo{
			Name:       "dest_name",
			DOB:        "dest_dob",
			NomineeDOB: "dest_dob",
		}
		CopyNonEmptyFiled(src, dest)
		assert.Equal(t, "src_id", dest.UniqueID)
		assert.Equal(t, "dest_name", dest.Name)
		assert.Equal(t, "dest_dob", dest.DOB)
		assert.Equal(t, "dest_dob", dest.NomineeDOB)
	})
}

func TestGetNonEmptyFields(t *testing.T) {
	t.Run("should return Dict with non nil fields", func(t *testing.T) {
		name := "name"
		src := &AugmontUserInfo{
			Name: name,
		}
		expected := Dict{
			"userName": name,
		}
		got := GetNonEmptyFields(src)
		assert.Equal(t, expected, got)
	})

	t.Run("should return empty Dict", func(t *testing.T) {
		src := &AugmontUserInfo{}
		got := GetNonEmptyFields(src)
		assert.Equal(t, Dict{}, got)
	})
}
