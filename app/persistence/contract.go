package persistence

import "gitlab.com/tellmecomua/tellme.api/app/persistence/model"

type Repository interface {
	Connect() error
	Disconnect() error

	// Experts
	GetExpert(id string) (*model.Expert, error)
	GetExpertByEmail(email string) (*model.Expert, error)
	GetExpertByPhone(phone string) (*model.Expert, error)
	GetExpertList(q *model.QueryExpertList) (*model.ExpertList, error)
	CreateExpert(e *model.Expert) (*model.Expert, error)
	DeleteExpert(e *model.Expert) error
	UpdateExpertStatus(e *model.Expert) (*model.Expert, error)
	UpdateExpertPassword(e *model.Expert) (*model.Expert, error)
	UpdateExpertDocumentURLs(e *model.Expert) (*model.Expert, error)

	// Requisitions
	CreateRequisition(req *model.Requisition) (*model.Requisition, error)
	GetRequisitionList(q *model.QueryRequisitionList) (*model.RequisitionList, error)
	UpdateRequisitionStatus(q *model.Requisition) (*model.Requisition, error)

	// Admins
	GetAdminByLogin(login string) (*model.Admin, error)

	// Service
	Name() string
	Stat() map[string]string
}
