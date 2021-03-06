package wrp

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type=MessageType

// MessageType indicates the kind of WRP message
type MessageType int64

const (
	SimpleRequestResponseMessageType MessageType = iota + 3
	SimpleEventMessageType
	CreateMessageType
	RetrieveMessageType
	UpdateMessageType
	DeleteMessageType
	ServiceRegistrationMessageType
	ServiceAliveMessageType
	UnknownMessageType
	lastMessageType
)

// SupportsTransaction tests if messages of this type are allowed to participate in transactions.
// If this method returns false, the TransactionUUID field should be ignored (but passed through
// where applicable).
func (mt MessageType) SupportsTransaction() bool {
	switch mt {
	case SimpleEventMessageType:
		return false
	case ServiceRegistrationMessageType:
		return false
	case ServiceAliveMessageType:
		return false
	case UnknownMessageType:
		return false
	default:
		return true
	}
}

// FriendlyName is just the String version of this type minus the "MessageType" suffix.
// This is used in most textual representations, such as HTTP headers.
func (mt MessageType) FriendlyName() string {
	return friendlyNames[mt]
}

var (
	// stringToMessageType is a simple map of allowed strings which uniquely indicate MessageType values.
	// Included in this map are integral string keys.  Keys are assumed to be case insensitive.
	stringToMessageType map[string]MessageType

	// friendlyNames are the string representations of each message type without the "MessageType" suffix
	friendlyNames map[MessageType]string
)

func init() {
	stringToMessageType = make(map[string]MessageType, lastMessageType-1)
	friendlyNames = make(map[MessageType]string, lastMessageType-1)
	suffixLength := len("MessageType")

	// for each MessageType, allow the following string representations:
	//
	// The integral value of the constant
	// The String() value
	// The String() value minus the MessageType suffix
	for v := SimpleRequestResponseMessageType; v < lastMessageType; v++ {
		stringToMessageType[strconv.Itoa(int(v))] = v

		vs := v.String()
		f := vs[0 : len(vs)-suffixLength]

		stringToMessageType[vs] = v
		stringToMessageType[f] = v
		friendlyNames[v] = f
	}

	stringToMessageType["event"] = SimpleEventMessageType
}

// StringToMessageType converts a string into an enumerated MessageType constant.
// If the value equals the friendly name of a type, e.g. "Auth" for AuthMessageType,
// that type is returned.  Otherwise, the value is converted to an integer and looked up,
// with an error being returned in the event the integer value is not valid.
func StringToMessageType(value string) (MessageType, error) {
	mt, ok := stringToMessageType[value]
	if !ok {
		return MessageType(-1), fmt.Errorf("Invalid message type: %s", value)
	}

	return mt, nil
}
