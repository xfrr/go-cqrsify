package valueobject_test

import (
	"testing"

	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func BenchmarkEmailCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		valueobject.NewEmail("test@example.com")
	}
}

func BenchmarkEmailEquality(b *testing.B) {
	email1, _ := valueobject.NewEmail("test@example.com")
	email2, _ := valueobject.NewEmail("test@example.com")

	for b.Loop() {
		email1.Equals(email2)
	}
}
