// Package status Error and Validation message response
package status

import (
	"encoding/json"
	"main/core/messages"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ErrInternal represents internal server error.
var ErrInternal = ErrServiceStatus{
	ServiceStatus{Code: http.StatusInternalServerError, Message: messages.InternalServerError},
}

// ErrNotFound represents an error when a domain artifact was not found.
var ErrNotFound = ErrServiceStatus{
	ServiceStatus{Code: http.StatusNotFound, Message: messages.NotFound},
}

// ServiceStatus captures basic information about a status construct.
type ServiceStatus struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"msg"`
	Details []*Dtl `json:"details,omitempty"`
}

// Dtl captures basic information about a status construct.
type Dtl struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

// ErrServiceStatus captures basic information about an error.
type ErrServiceStatus struct {
	ServiceStatus
}

// ErrStatusConflict represents conflict because of inconsistent or duplicated info.
var ErrStatusConflict = ErrServiceStatus{
	ServiceStatus{Code: http.StatusConflict, Message: messages.ConflictError},
}

// ErrStatusUnprocessableEntity represents conflict because of inconsistent or duplicated info.
var ErrStatusUnprocessableEntity = ErrServiceStatus{
	ServiceStatus{Code: http.StatusUnprocessableEntity, Message: messages.UnprocessableEntityError},
}

// WithMessage returns an error status with given message.
func (e ErrServiceStatus) WithMessage(msg string) ErrServiceStatus {
	return ErrServiceStatus{ServiceStatus{Code: e.Code, Message: msg}}
}

// WithError returns an error status with given err.Error().
func (e ErrServiceStatus) WithError(err error) ErrServiceStatus {
	return ErrServiceStatus{ServiceStatus{Code: e.Code, Message: err.Error()}}
}

// Error returns the error object
func (e ErrServiceStatus) Error() string {
	if errB, err := json.Marshal(&e); err == nil {
		return string(errB)
	}
	return `{"code":500, "msg": "error marshal failed"}`
}

// WithValidationError returns an error status with given err.Error().
func (e ErrServiceStatus) WithValidationError(err validation.Errors) ErrServiceStatus {
	errSvc := ErrServiceStatus{ServiceStatus{Code: e.Code, Message: e.Message, Details: nil}}
	for key, msg := range err {
		errSvc.AddDtl(key, msg.Error())
	}
	return errSvc
}

// AddDtl returns an error status with given message.
func (e *ErrServiceStatus) AddDtl(key, msg string) {
	if e.Details == nil {
		e.Details = []*Dtl{}
	}
	d := &Dtl{Key: key, Message: msg}
	e.Details = append(e.Details, d)
}
