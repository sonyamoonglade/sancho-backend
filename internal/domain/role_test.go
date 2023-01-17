package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoleCustomJSONSerialize(t *testing.T) {
	role := RoleCustomer
	marshalled, err := json.Marshal(role)
	require.NoError(t, err)

	var decodedRole Role
	err = json.Unmarshal(marshalled, &decodedRole)
	require.NoError(t, err)
	require.Equal(t, RoleCustomer, decodedRole)
}
