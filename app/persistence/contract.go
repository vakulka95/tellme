package persistence

import (
	"time"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

type Repository interface {
	Connect() error
	Disconnect() error

	// Experts
	GetExpert(id string) (*model.Expert, error)
	GetExpertView(id string) (*model.Expert, error)
	GetExpertByEmail(email string) (*model.Expert, error)
	GetExpertByPhone(phone string) (*model.Expert, error)
	GetExpertList(q *model.QueryExpertList) (*model.ExpertList, error)
	CreateExpert(e *model.Expert) (*model.Expert, error)
	DeleteExpert(e *model.Expert) error
	UpdateExpert(e *model.Expert) (*model.Expert, error)
	UpdateExpertStatus(e *model.Expert) (*model.Expert, error)
	UpdateExpertPassword(e *model.Expert) (*model.Expert, error)
	UpdateExpertDocumentURLs(e *model.Expert) (*model.Expert, error)
	UpdateExpertSpecializations(e *model.Expert) (*model.Expert, error)
	GetExpertRating(q *model.QueryExpertRatingList) (*model.ExpertList, error)
	GetExpertRatingTable(q *model.QueryExpertRatingList) ([]*model.Expert, error)

	// Requisitions
	GetRequisition(id string) (*model.Requisition, error)
	CreateRequisition(req *model.Requisition) (*model.Requisition, error)
	GetRequisitionList(q *model.QueryRequisitionList) (*model.RequisitionList, error)
	UpdateRequisitionStatus(q *model.Requisition) (*model.Requisition, error)
	UpdateRequisitionSMSReplyCount(q *model.Requisition) (*model.Requisition, error)
	DeleteRequisition(q *model.Requisition) error
	GetNotReviewedRequisition() ([]*model.Requisition, error)
	GetNotProcessedRequisition(lifetimeTo *time.Time) ([]*model.Requisition, error)

	// Admins
	GetAdmin(id string) (*model.Admin, error)
	GetAdminByLogin(login string) (*model.Admin, error)

	// Reviews
	GetReview(id string) (*model.Review, error)
	CreateReview(v *model.Review) (*model.Review, error)
	GetReviewByToken(token string) (*model.Review, error)
	GetReviewList(q *model.QueryReviewList) (*model.ReviewList, error)
	UpdateReviewBodyStatus(v *model.Review) (*model.Review, error)

	// Sessions
	GetSession(id string) (*model.Session, error)
	CreateSession(session *model.Session) (*model.Session, error)
	GetRequisitionSessionList(requisitionID string) ([]*model.Session, error)
	DeleteRequisitionSessions(requisitionID string) error
	DeleteLastRequisitionSession(requisitionID string) error

	// Comments
	GetComment(id string) (*model.Comment, error)
	GetCommentList(q *model.QueryCommentList) (*model.CommentList, error)
	GetExpertCommentList(expertID string) ([]*model.Comment, error)
	CreateComment(c *model.Comment) (*model.Comment, error)
	DeleteComment(c *model.Comment) error

	// Service
	Name() string
	Stat() map[string]string
}
