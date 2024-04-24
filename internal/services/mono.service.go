package services

import (
	"sigo/internal/lib"
)

type MonoService struct {
	DB struct {
		idManager *lib.IdentifierManager
		Rooms     map[int64]*Room
	}
}
