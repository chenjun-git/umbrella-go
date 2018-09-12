package errors

import (
	"fmt"
)

type Error interface {
	GetCode() int
	GetMessage() string
	SetMessage(message string)
	GetDescription() string
	SetDescription(description string)
	Error() string
}

type UmbrellaError struct {
	Code        int                `json:"code"`
	Message     string             `json:"message"`               // 用于显示前端错误提示
	Description string             `json:"description,omitempty"` // 用于内部显示错误信息
}

func (e *UmbrellaError) GetCode() int {
	return e.Code
}

func (e *UmbrellaError) GetMessage() string {
	return e.Message
}

func (e *UmbrellaError) SetMessage(message string) {
	e.Message = message
}

func (e *UmbrellaError) GetDescription() string {
	return e.Description
}

func (e *UmbrellaError) SetDescription(description string) {
	e.Description = description
}

func (e *UmbrellaError) Error() string {
	return fmt.Sprintf("%v, %v, %v", e.Code, e.Description, e.Validations)
}