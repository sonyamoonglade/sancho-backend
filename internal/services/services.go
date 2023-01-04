package service

import storage "github.com/sonyamoonglade/sancho-backend/internal/storages"

type Services struct {
	Product Product
}

type Deps struct {
	Storages *storage.Storages
}

func NewServices(deps Deps) *Services {
	stg := deps.Storages
	return &Services{
		Product: NewProductService(stg.Product),
	}
}
