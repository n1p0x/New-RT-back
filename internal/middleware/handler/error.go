package handler

import (
	"encoding/json"
	"fmt"
)

type ErrorCode string

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"detail": e.Detail,
	}
	return json.Marshal(&m)
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("ErrorResponse{Detail:%s}", e.Detail)
}
