package document

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type request struct {
	BirthCertificate      string `json:"birth_certificate"`
	FamilyCard            string `json:"family_card"`
	UserIdentity          string `json:"user_identity"`
	BeltCertificate       string `json:"belt_certificate"`
	ElementaryCertificate string `json:"elementary_certificate"`
	MiddleCertificate     string `json:"middle_certificate"`
	HighCertificate       string `json:"high_certificate"`
	BachelorCertificate   string `json:"bachelor_certificate"`
	MasterCertificate     string `json:"master_certificate"`
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BirthCertificate, is.URL),
		validation.Field(&r.FamilyCard, is.URL),
		validation.Field(&r.UserIdentity, is.URL),
		validation.Field(&r.BeltCertificate, is.URL),
		validation.Field(&r.ElementaryCertificate, is.URL),
		validation.Field(&r.MiddleCertificate, is.URL),
		validation.Field(&r.HighCertificate, is.URL),
		validation.Field(&r.BachelorCertificate, is.URL),
		validation.Field(&r.MasterCertificate, is.URL),
	)
}
