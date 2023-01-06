package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomJSONSerialize(t *testing.T) {
	u := User{
		UserID: "abcd",
		Role:   RoleCustomer,
	}

	marshalled, err := json.Marshal(u)
	require.NoError(t, err)

	var user User
	err = json.Unmarshal(marshalled, &user)
	require.NoError(t, err)

	require.Equal(t, RoleCustomer.String(), user.Role.String())
	require.Equal(t, u.UserID, user.UserID)
}
