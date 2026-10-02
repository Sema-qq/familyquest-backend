package uuidprovider

import "github.com/google/uuid"

type UuidProvider struct{}

func New() *UuidProvider {
	return &UuidProvider{}
}

func (g *UuidProvider) GenerateV4() uuid.UUID {
	return uuid.New()
}

func (g *UuidProvider) GenerateStringV4() string {
	return uuid.NewString()
}

func (g *UuidProvider) GenerateV7() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func (g *UuidProvider) GenerateStringV7() string {
	return uuid.Must(uuid.NewV7()).String()
}
