package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commonProto "umbrella-go/umbrella-common/proto"
)

func TestFromProtoError(t *testing.T) {
	assert := assert.New(t)

	protoErr := &commonProto.Error{
		Code:        int32(1),
		Description: "test description",
		Message:     "test message",
	}

	umbrellaErr := FromProtoError(protoErr)

	assert.Equal(int(protoErr.Code), umbrellaErr.GetCode())
	assert.Equal(protoErr.Description, umbrellaErr.GetDescription())
	assert.Equal(protoErr.Message, umbrellaErr.GetMessage())
}

func TestToProtoError(t *testing.T) {
	assert := assert.New(t)

	umbrellaErr := &UmbrellaError{
		Code:        int(1),
		Description: "test description",
		Message:     "test message",
	}

	protoErr := ToProtoError(umbrellaErr)

	assert.Equal(umbrellaErr.Code, int(protoErr.Code))
	assert.Equal(umbrellaErr.Description, protoErr.Description)
	assert.Equal(umbrellaErr.Message, protoErr.Message)
}
