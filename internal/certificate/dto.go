package certificate

import "github.com/project-ippl-dev/tanding-api/internal/db"

type response struct {
	Certificate db.CertificateFetchOneRow `json:"certificate"`
	Event       eventDetail               `json:"event"`
	Recipient   string                    `json:"recipient"`
}

type eventDetail struct {
	db.EventFetchOneInfiniteByIDRow
	Participants int64 `json:"participants"`
}
