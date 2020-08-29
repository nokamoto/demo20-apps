package test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

type protoMatcher struct {
	x proto.Message
}

func (p protoMatcher) Matches(x interface{}) bool {
	m, ok := x.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(p.x, m)
}

func (p protoMatcher) String() string {
	return fmt.Sprintf("is equal to %v", p.x)
}

// ProtoEq returns Matcher for the proto message.
func ProtoEq(x proto.Message) gomock.Matcher {
	return protoMatcher{x: x}
}
