package entity

import "github.com/google/uuid"

type FamilyID uuid.UUID

func (id FamilyID) String() string {
	return id.UUID().String()
}

func (id FamilyID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type FamilyMemberID uuid.UUID

func (id FamilyMemberID) String() string {
	return id.UUID().String()
}

func (id FamilyMemberID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type FamilyRole string

const (
	FamilyRoleParent FamilyRole = "parent"
	FamilyRoleChild  FamilyRole = "child"
)

type Family struct {
	ID       FamilyID
	Name     string
	Timezone string
	OwnerID  UserID
}

type FamilyCreate struct {
	ID       FamilyID
	Name     string
	Timezone string
	OwnerID  UserID
}

type FamilyCreateRequest struct {
	Name     string
	Timezone string
	OwnerID  UserID
}

type FamilyWithRole struct {
	Family Family
	Role   FamilyRole
}

type FamilyMember struct {
	ID       FamilyMemberID
	FamilyID FamilyID
	UserID   UserID
	Role     FamilyRole
}

type FamilyMemberProfile struct {
	ID          FamilyMemberID
	FamilyID    FamilyID
	UserID      UserID
	Login       string
	DisplayName string
	Role        FamilyRole
}

type FamilyMemberCreate struct {
	ID       FamilyMemberID
	FamilyID FamilyID
	UserID   UserID
	Role     FamilyRole
}

type FamilyMemberCreateRequest struct {
	Login       string
	Password    string
	DisplayName string
	ParentID    UserID
}
