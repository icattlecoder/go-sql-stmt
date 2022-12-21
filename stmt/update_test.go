package stmt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	setMap := map[Column]interface{}{
		"name": "bob",
		"age":  20,
		"foo":  "foo",
		"baz":  1,
	}
	query := Update(Column("users")).Set(setMap)

	require.Equal(t, "UPDATE users SET age = %s, baz = %s, foo = %s, name = %s", query.SqlString())

	require.Equal(t, query.Values()[0], 20)
	require.Equal(t, query.Values()[1], 1)
	require.Equal(t, query.Values()[2], "foo")
	require.Equal(t, query.Values()[3], "bob")
}
