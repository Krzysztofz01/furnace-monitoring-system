package protocol

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

func GetHostIdFromEventPayload(ep string) (uuid.UUID, error) {
	epParts := strings.Split(ep, PayloadSeparatorSymbolString)
	if len(epParts) < 2 {
		return uuid.UUID{}, errors.New("protocol: invalid event payload format")
	}

	hostId, err := uuid.Parse(epParts[1])
	if err != nil {
		return uuid.UUID{}, errors.New("protocol: failed to parse the provided host id")
	}

	return hostId, nil
}
