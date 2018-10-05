package errors

import (
	"fmt"

	proto "umbrella-go/umbrella-common/proto"
)

type Error interface {
	GetCode() int
	GetMessage() string
	SetMessage(message string)
	GetDescription() string
	SetDescription(description string)
	Error() string
}

func NewError(code int, description string) Error {
	return &UmbrellaError{
		Code:        code,
		Description: description,
	}
}

type UmbrellaError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`               // 用于显示前端错误提示
	Description string `json:"description,omitempty"` // 用于内部显示错误信息
}

func (ue *UmbrellaError) GetCode() int {
	return ue.Code
}

func (ue *UmbrellaError) GetMessage() string {
	return ue.Message
}

func (ue *UmbrellaError) SetMessage(message string) {
	ue.Message = message
}

func (ue *UmbrellaError) GetDescription() string {
	return ue.Description
}

func (ue *UmbrellaError) SetDescription(description string) {
	ue.Description = description
}

func (ue *UmbrellaError) Error() string {
	return fmt.Sprintf("%v, %v", ue.Code, ue.Description)
}

func FromProtoError(err *proto.Error) Error {
	r := NewError(int(err.Code), err.Description)
	r.SetMessage(err.Message)

	return r
}

func ToProtoError(err Error) *proto.Error {
	r := proto.Error{
		Code:        int32(err.GetCode()),
		Message:     err.GetMessage(),
		Description: err.GetDescription(),
	}

	return &r
}
