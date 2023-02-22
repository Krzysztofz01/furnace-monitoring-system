package protocol

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func CreateDashboardConnectedEventPayload(hostId uuid.UUID) ([]byte, error) {
	hostIdString := hostId.String()
	if len(hostIdString) == 0 {
		return nil, errors.New("protocol: invalid host id provided in create dashboard connected event payload")
	}

	payloadBuilder := strings.Builder{}
	payloadBuilder.WriteString(EventTypeToString(DashboardConnectedEvent))
	payloadBuilder.WriteRune(PayloadSeparatorSymbolRune)
	payloadBuilder.WriteString(hostIdString)

	return []byte(payloadBuilder.String()), nil
}

func CreateDashboardDisconnectedEventPayload(hostId uuid.UUID) ([]byte, error) {
	hostIdString := hostId.String()
	if len(hostIdString) == 0 {
		return nil, errors.New("protocol: invalid host id provided in create dashboard connected event payload")
	}

	payloadBuilder := strings.Builder{}
	payloadBuilder.WriteString(EventTypeToString(DashboardDisconnectedEvent))
	payloadBuilder.WriteRune(PayloadSeparatorSymbolRune)
	payloadBuilder.WriteString(hostIdString)

	return []byte(payloadBuilder.String()), nil
}
