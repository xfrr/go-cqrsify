package valueobject_test

import (
	"testing"

	domainkit "github.com/xfrr/go-cqrsify/domain/value-object"
)

func BenchmarkMoneyAddition(b *testing.B) {
	money1, _ := domainkit.NewMoney(10.0, "USD")
	money2, _ := domainkit.NewMoney(5.0, "USD")

	for b.Loop() {
		money1.Add(money2)
	}
}
