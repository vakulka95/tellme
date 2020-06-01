package generator

import (
	"github.com/rs/xid"
)

func NewReviewToken() string {
	return xid.New().String()
}
