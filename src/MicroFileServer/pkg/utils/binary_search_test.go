package utils_test

import (
	"sort"
	"testing"

	"github.com/MicroFileServer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestFunc_BinarySearch(t *testing.T) {
	t.Run("don't find", func(t *testing.T) {
		slice := []string{
			"a", "b", "c", "d", "f","F",
		}
		strSlice := sort.StringSlice(slice)
		strSlice.Sort()

		_, find := utils.FindString(
			strSlice,
			"G",
		)
		require.Equal(t, false, find)

		_, find = utils.FindString(
			strSlice,
			"zqnfjqfnwelf",
		)
		require.Equal(t, false, find)

		_, find = utils.FindString(
			strSlice,
			"-",
		)
		require.Equal(t, false, find)

		_, find = utils.FindString(
			strSlice,
			"192839124809",
		)
		require.Equal(t, false, find)

		_, find = utils.FindString(
			strSlice,
			"z",
		)
		require.Equal(t, false, find)
	})

	t.Run("find", func(t *testing.T) {
		slice := []string{
			"a", "b", "c", "d", "f","F",
		}
		strSlice := sort.StringSlice(slice)
		strSlice.Sort()

		s, find := utils.FindString(
			strSlice,
			"a",
		)
		require.Equal(t, true, find)
		require.Equal(t, "a", s)

		s, find = utils.FindString(
			strSlice,
			"b",
		)
		require.Equal(t, true, find)
		require.Equal(t, "b", s)

		s, find = utils.FindString(
			strSlice,
			"c",
		)
		require.Equal(t, true, find)
		require.Equal(t, "c", s)

		s, find = utils.FindString(
			strSlice,
			"d",
		)
		require.Equal(t, true, find)
		require.Equal(t, "d", s)

		s, find = utils.FindString(
			strSlice,
			"f",
		)
		require.Equal(t, true, find)
		require.Equal(t, "f", s)

		s, find = utils.FindString(
			strSlice,
			"F",
		)
		require.Equal(t, true, find)
		require.Equal(t, "F", s)
	})
}