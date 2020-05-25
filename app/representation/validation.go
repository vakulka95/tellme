package representation

import (
	"fmt"
	"strings"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/tellmecomua/tellme.api/util/validation"
)

func ExpertAPItoPersistence(e *CreateExpertRequest) (*model.Expert, error) {
	id := uuid.NewV4()

	// Name
	if strings.TrimSpace(e.Name) == "" {
		return nil, fmt.Errorf("name mustn't be empty")
	}

	// Password
	if err := validation.ValidatePassword(e.Password); err != nil {
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %v", e.Password)
	}

	// Gender
	if ok := allowedGender[e.Gender]; !ok {
		return nil, fmt.Errorf("invalid gender %s", e.Gender)
	}

	// Email
	if err := validation.ValidateEmail(e.Email); err != nil {
		return nil, err
	}
	e.Email = strings.ToLower(e.Email)

	// Phone
	e.Phone = strings.TrimPrefix(e.Phone, "+")
	if err := validation.ValidatePhoneNumber(e.Phone); err != nil {
		return nil, err
	}

	// Education
	if strings.TrimSpace(e.Education) == "" {
		return nil, fmt.Errorf("education mustn't be empty")
	}

	// Specializations
	for _, v := range e.Specializations {
		if _, ok := Diagnoses[v]; !ok {
			return nil, fmt.Errorf("diagnosis %s isn't supported", v)
		}
	}

	if len(e.Specializations) == 0 {
		return nil, fmt.Errorf("specialization mustn't be empty")
	}

	return &model.Expert{
		ID:              id.String(),
		Username:        e.Name,
		Gender:          e.Gender,
		Phone:           e.Phone,
		Email:           e.Email,
		Password:        string(hashedPass),
		Specializations: e.Specializations,
		Education:       e.Education,
		DocumentURLs:    make([]string, 0),
		Status:          model.StatusOnReview,
	}, nil
}

func RequisitionAPItoPersistence(r *CreateRequisitionRequest) (*model.Requisition, error) {
	id := uuid.NewV4()

	// Name
	if strings.TrimSpace(r.Name) == "" {
		return nil, fmt.Errorf("name mustn't be empty")
	}

	// Gender
	if ok := allowedGender[r.Gender]; !ok {
		return nil, fmt.Errorf("invalid gender %s", r.Gender)
	}
	if ok := allowedGender[r.ExpertGender]; !ok {
		return nil, fmt.Errorf("invalid expert gender %s", r.ExpertGender)
	}

	// Phone
	if err := validation.ValidatePhoneNumber(strings.TrimPrefix(r.Phone, "+")); err != nil {
		return nil, err
	}

	// Diagnosis
	if _, ok := Diagnoses[r.Diagnosis]; !ok {
		return nil, fmt.Errorf("diagnosis %s isn't supported", r.Diagnosis)
	}

	// Feedback
	if ok := allowedFeedbackType[r.FeedbackType]; !ok {
		return nil, fmt.Errorf("invalid feedback type %s", r.FeedbackType)
	}

	if r.FeedbackType != FeedbackTypePhone && strings.TrimSpace(r.FeedbackContact) == "" {
		return nil, fmt.Errorf("feedback type mustn't be empty")
	}

	// Is adult
	if !r.IsAdult {
		return nil, fmt.Errorf("user must be adult")
	}

	return &model.Requisition{
		ID:                   id.String(),
		Username:             r.Name,
		Gender:               r.Gender,
		Phone:                r.Phone,
		Diagnosis:            r.Diagnosis,
		DiagnosisDescription: r.DiagnosisDescription,
		ExpertGender:         r.ExpertGender,
		FeedbackType:         r.FeedbackType,
		FeedbackContact:      r.FeedbackContact,
		FeedbackTime:         r.FeedbackTime,
		FeedbackWeekDay:      r.FeedbackWeekDay,
		IsAdult:              r.IsAdult,
		Status:               model.StatusCreated,
	}, nil
}
