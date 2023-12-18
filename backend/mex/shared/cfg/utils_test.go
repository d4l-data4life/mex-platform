package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EmptyString(t *testing.T) {
	require.True(t, StringIsEmpty(""))
	require.True(t, StringIsEmpty(" "))
	require.True(t, StringIsEmpty("  "))
	require.True(t, StringIsEmpty("   "))
	require.True(t, StringIsEmpty("∅"))
	require.True(t, StringIsEmpty("∅∅"))
	require.True(t, StringIsEmpty(" ∅ "))
	require.True(t, StringIsEmpty(" ∅ ∅ ∅ "))

	require.False(t, StringIsEmpty("\t"))
	require.False(t, StringIsEmpty("\n"))
	require.False(t, StringIsEmpty("∅|∅"))
}

func Test_BytesAreString(t *testing.T) {
	require.True(t, BytesAreEmpty(nil))
	require.True(t, BytesAreEmpty([]byte{}))

	require.True(t, BytesAreEmpty([]byte("")))
	require.True(t, BytesAreEmpty([]byte("∅")))
	require.True(t, BytesAreEmpty([]byte("\u2205"))) // Codepoint for ∅
	require.True(t, BytesAreEmpty([]byte("    ")))
	require.True(t, BytesAreEmpty([]byte(" ∅  ∅ ")))

	require.False(t, BytesAreEmpty([]byte{0}))
	require.False(t, BytesAreEmpty([]byte("\t")))
	require.False(t, BytesAreEmpty([]byte("∅\n∅")))
}
