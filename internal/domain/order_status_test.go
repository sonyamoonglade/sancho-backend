package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderStatusCustomJSONSerialize(t *testing.T) {
	orderStatus := StatusWaitingForVerification
	marshalled, err := json.Marshal(orderStatus)
	require.NoError(t, err)

	var decodedOrderStatus OrderStatus
	err = json.Unmarshal(marshalled, &decodedOrderStatus)
	require.NoError(t, err)
	require.Equal(t, StatusWaitingForVerification, decodedOrderStatus)
}
