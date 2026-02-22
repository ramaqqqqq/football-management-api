package provider

import "github.com/google/uuid"

type UUIDProvider interface {
	NewUUID() uuid.UUID
}

type GoogleUUID struct {
}

func (g *GoogleUUID) NewUUID() uuid.UUID {
	return uuid.New()
}
