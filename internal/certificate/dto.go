package certificate

import "github.com/project-ippl-dev/tanding-api/internal/db"

type Response struct {
	Certificate db.CertificateFetchOneRow `json:"certificate"`
	Event       EventDetail               `json:"event"`
	Recipient   string                    `json:"recipient"`
}

type EventDetail struct {
	db.EventFetchOneInfiniteByIDRow
	Participants int64 `json:"participants"`
}
