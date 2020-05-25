package representation

type GatewayError struct {
	Code string `json:"code"`
}

// Codes
var (
	ErrInternal               = &GatewayError{Code: "Internal"}
	ErrConflict               = &GatewayError{Code: "Conflict"}
	ErrNotFound               = &GatewayError{Code: "NotFound"}
	ErrInvalidRequest         = &GatewayError{Code: "InvalidRequest"}
	ErrExpertAlreadyProcessed = &GatewayError{Code: "ExpertAlreadyProcessed"}
)
