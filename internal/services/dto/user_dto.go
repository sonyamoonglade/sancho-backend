package dto

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

type SaveSessionDTO struct {
	UserID  string
	Role    domain.Role
	Session domain.Session
}

type LoginAdminDTO struct {
	Login    string
	Password string
}

type RegisterAdminDTO struct {
	Login    string
	Password string
}
