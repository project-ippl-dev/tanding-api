package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase struct {
	repository *db.Queries
}

func NewUsecase(repository *db.Queries) Usecase {
	return Usecase{repository: repository}
}

func (u Usecase) store(ctx context.Context, req request, userID string) error {
	return u.repository.DocumentCreate(ctx, db.DocumentCreateParams{
		UserID:                uuid.MustParse(userID),
		BirthCertificate:      req.BirthCertificate,
		FamilyCard:            req.FamilyCard,
		UserIdentity:          req.UserIdentity,
		BeltCertificate:       req.BeltCertificate,
		ElementaryCertificate: req.ElementaryCertificate,
		MiddleCertificate:     req.MiddleCertificate,
		HighCertificate:       req.HighCertificate,
		BachelorCertificate:   req.BachelorCertificate,
		MasterCertificate:     req.MasterCertificate,
	})
}

func (u Usecase) fetchAll(ctx context.Context, page int32, pageSize int32, userID string) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)
	documents, err := u.repository.DocumentFetchAll(ctx, db.DocumentFetchAllParams{
		UserID: uuid.MustParse(userID),
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch documents : %s", err.Error())
	}
	count, err := u.repository.DocumentCountAll(ctx, uuid.MustParse(userID))
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count documents : %s", err.Error())
	}
	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      documents,
	}, nil
}

func (u Usecase) update(ctx context.Context, req request, userID string, documentID int64) error {
	docID, err := u.repository.DocumentCheckOne(ctx, db.DocumentCheckOneParams{
		ID:     documentID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check document : %s", err.Error())
	}
	return u.repository.DocumentUpdate(ctx, db.DocumentUpdateParams{
		BirthCertificate:      req.BirthCertificate,
		FamilyCard:            req.FamilyCard,
		UserIdentity:          req.UserIdentity,
		BeltCertificate:       req.BeltCertificate,
		ElementaryCertificate: req.ElementaryCertificate,
		MiddleCertificate:     req.MiddleCertificate,
		HighCertificate:       req.HighCertificate,
		BachelorCertificate:   req.BachelorCertificate,
		MasterCertificate:     req.MasterCertificate,
		ID:                    docID,
	})
}

func (u Usecase) delete(ctx context.Context, userID string, documentID int64) error {
	docID, err := u.repository.DocumentCheckOne(ctx, db.DocumentCheckOneParams{
		ID:     documentID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check document : %s", err.Error())
	}
	return u.repository.DocumentDelete(ctx, docID)
}
