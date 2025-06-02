package eventPaymentFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/eventPayment"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"time"
)

var (
	EventPaymentFetchResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []eventPayment.Response{
			{
				FetchAllRow: eventPayment.FetchAllRow{
					ID:           uuid.New(),
					EventID:      uuid.New(),
					EventName:    "event-name",
					UniqueNumber: 123,
					PaymentLink:  "https://google.com",
					Status:       db.EventReceiptStatusApproved,
					ClubID:       uuid.New(),
					CreatedAt:    time.Now(),
				},
				ClassEvents: []db.ClassEventFetchByPaymentIDRow{
					{
						ID:    uuid.New(),
						Price: 1000,
						Name:  "name",
					},
				},
			},
		},
	}

	EventPaymentCartResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.EventRegistrationFetchCartRow{
			{
				ID:        uuid.New(),
				Name:      "name",
				Thumbnail: "https://google.com",
				Total:     1,
			},
		},
	}

	EventPaymentCartDetailResponse = eventPayment.CartDetailResponse{
		Results: []eventPayment.EventRegistrationFetchCartDetailsRow{
			{
				ID:           uuid.New(),
				ClassName:    "class name",
				Price:        1000,
				Participants: []string{"participant-1", "participant-2"},
			},
		},
		Event: db.EventFetchForCartRow{
			EventOwner: "eventOwner",
			EventName:  "eventName",
			SportName:  "sport name",
			Thumbnail:  "https://google.com",
			Deadline:   time.Now().Add(1 * time.Hour),
		},
		UniqueNumber: eventPayment.UniqueNumberData{
			Number: "123",
			Time:   time.Now().Add(1 * time.Hour).String(),
		},
	}

	DetailPaymentResponse = eventPayment.DetailPaymentResponse{
		Detail: db.EventPaymentFetchOneForAdminRow{
			ID:           uuid.New(),
			Status:       db.EventReceiptStatusApproved,
			Total:        1,
			UniqueNumber: 123,
			CreatedAt:    time.Now(),
			PaymentLink:  "https://google.com",
			EventID:      uuid.New(),
			ClubName:     "club name",
			Owner:        "owner",
			EventName:    "event name",
		},
		ClassEvents: []eventPayment.EventRegistrationFetchCartDetailsRow{
			{
				ID:           uuid.New(),
				ClassName:    "class name",
				Price:        1000,
				Participants: []string{"participant-1"},
			},
		},
	}

	SummaryResponse = eventPayment.SummaryResponse{
		TotalApproved: 1,
		TotalWaiting:  1,
		TotalRefund:   1,
	}
)
