package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPayCustomJSONSerialize(t *testing.T) {
	pay := PayOnPickup
	marshalled, err := json.Marshal(pay)
	require.NoError(t, err)
	expected := `{"pay":"on pickup"}`
	require.Equal(t, expected, string(marshalled))

	var decodedPay Pay
	err = json.Unmarshal(marshalled, &decodedPay)
	require.NoError(t, err)
	require.Equal(t, PayOnPickup, decodedPay)
}
