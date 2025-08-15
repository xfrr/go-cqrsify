package message

import "reflect"

// NameOf returns the name of the message.
func NameOf[M any]() string {
	return reflect.TypeOf((*M)(nil)).Elem().Name()
}
