package test

import "github.com/project-ippl-dev/tanding-api/internal/db"

type Usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return Usecase{repository: repository}
}
