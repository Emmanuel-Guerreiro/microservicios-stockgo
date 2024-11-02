package lib

import "encoding/json"

// Unauthorized el usuario no esta autorizado al recurso
var UnauthorizedError = NewRestError(401, "Unauthorized")

// NotFound cuando un registro no se encuentra en la db
var NotFoundError = NewRestError(404, "Document not found")

// AlreadyExist cuando no se puede ingresar un registro a la db
var AlreadyExistError = NewRestError(400, "Already exist")

// Internal esta aplicación no sabe como manejar el error
var InternalError = NewRestError(500, "Internal server error")

var InvalidError = NewRestError(400, "Invalid Document")

// - Creación de errors -
// NewRestError creates a new errCustom
func NewRestError(status int, message string) RestError {
	return &restError{
		status:  status,
		Message: message,
	}
}

//  - Algunas definiciones necesarias -

// RestError es una interfaz para definir errores custom
type RestError interface {
	Status() int
	Error() string
}

// restError es un error personalizado para http
type restError struct {
	status  int
	Message string `json:"error"`
}

func (e *restError) Error() string {
	return e.Message
}

// Status http status code
func (e *restError) Status() int {
	return e.status
}

// IValidationErr es una interfaz para definir errores custom
// IValidationErr es un error de validaciones de parameteros o de campos
type IValidationErr interface {
	Add(path string, message string) IValidationErr
	Error() string
}

func NewValidationError() IValidationErr {
	return &ValidationErr{
		Messages: []errField{},
	}
}

type ValidationErr struct {
	Messages []errField `json:"messages"`
}

func (e *ValidationErr) Error() string {
	body, err := json.Marshal(e)
	if err != nil {
		return "ErrValidation invalid."
	}
	return string(body)
}

// Add agrega errores a un validation error
func (e *ValidationErr) Add(path string, message string) IValidationErr {
	err := errField{
		Path:    path,
		Message: message,
	}
	e.Messages = append(e.Messages, err)
	return e
}

// errField define un campo inválido. path y mensaje de error
type errField struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}
