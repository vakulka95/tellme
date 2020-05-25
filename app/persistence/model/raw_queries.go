package model

const (
	StatementGetExpert                = "GetExpert"
	StatementGetExpertByEmail         = "GetExpertByEmail"
	StatementGetExpertByPhone         = "GetExpertByPhone"
	StatementGetExpertList            = "GetExpertList"
	StatementGetExpertCount           = "GetCountExpert"
	StatementCreateExpert             = "CreateExpert"
	StatementDeleteExpert             = "DeleteExpert"
	StatementUpdateExpertStatus       = "UpdateExpertStatus"
	StatementUpdateExpertPassword     = "UpdateExpertPassword"
	StatementUpdateExpertDocumentURLs = "UpdateExpertDocumentURLs"

	StatementCreateRequisition       = "CreateRequisition"
	StatementGetRequisitionList      = "GetRequisitionList"
	StatementGetRequisitionCount     = "GetCountRequisition"
	StatementUpdateRequisitionStatus = "UpdateRequisitionStatus"

	StatementGetAdminByLogin = "GetAdminByLogin"
)

func GetRawQueries() map[string]string {
	return queries
}

var queries = map[string]string{
	StatementGetExpert:                `SELECT id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at FROM experts WHERE id=$1`,
	StatementGetExpertByEmail:         `SELECT id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at FROM experts WHERE email=$1`,
	StatementGetExpertByPhone:         `SELECT id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at FROM experts WHERE phone=$1`,
	StatementGetExpertList:            `SELECT id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at, processing_count, completed_count FROM v$experts ORDER BY created_at desc LIMIT $1 OFFSET $2`,
	StatementGetExpertCount:           `SELECT COUNT(id) FROM experts`,
	StatementCreateExpert:             `INSERT INTO experts (id, username, gender, phone, email, password, specializations, education, document_urls, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
	StatementDeleteExpert:             `DELETE FROM experts WHERE id=$1`,
	StatementUpdateExpertStatus:       `UPDATE experts SET status=$2 WHERE id=$1`,
	StatementUpdateExpertPassword:     `UPDATE experts SET password=$2 WHERE id=$1`,
	StatementUpdateExpertDocumentURLs: `UPDATE experts SET document_urls=$2 WHERE id=$1`,
	StatementCreateRequisition:        `INSERT INTO requisitions (id, expert_id, username, gender, phone, diagnosis, diagnosis_description, expert_gender, feedback_type, feedback_contact, feedback_time, feedback_week_day, is_adult, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
	StatementGetRequisitionList:       `SELECT id, expert_id, username, gender, phone, diagnosis, diagnosis_description, expert_gender, feedback_type, feedback_contact, feedback_time, feedback_week_day, is_adult, status, created_at, updated_at FROM requisitions ORDER BY created_at desc LIMIT $1 OFFSET $2`,
	StatementGetRequisitionCount:      `SELECT COUNT(id) FROM requisitions`,
	StatementUpdateRequisitionStatus:  `UPDATE requisitions SET status=$2, expert_id=$3 WHERE id=$1`,
	StatementGetAdminByLogin:          `SELECT id, username, password, status, updated_at, created_at FROM admins WHERE username=$1`,
}
