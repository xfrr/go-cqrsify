package aggregate_test

import (
	"reflect"
	"testing"

	"github.com/xfrr/go-cqrsify/aggregate"
)

func Test_SearchCriteria_AggregateIDs(t *testing.T) {
	sut := aggregate.SearchCriteria().WithSearchAggregateIDs("1", "2", "3")
	if !reflect.DeepEqual(sut.AggregateIDs(), []string{"1", "2", "3"}) {
		t.Fatalf("expected aggregate ids to be [1, 2, 3], got %v", sut.AggregateIDs())
	}
}

func Test_SearchCriteria_AggregateNames(t *testing.T) {
	sut := aggregate.SearchCriteria().WithSearchAggregateNames("test-1", "test-2", "test-3")
	if !reflect.DeepEqual(sut.AggregateNames(), []string{"test-1", "test-2", "test-3"}) {
		t.Fatalf("expected aggregate names to be [test-1, test-2, test-3], got %v", sut.AggregateNames())
	}
}

func Test_SearchCriteria_AggregateVersions(t *testing.T) {
	sut := aggregate.SearchCriteria().WithSearchAggregateVersions(1, 2, 3)
	if !reflect.DeepEqual(sut.AggregateVersions(), []int{1, 2, 3}) {
		t.Fatalf("expected aggregate versions to be [1, 2, 3], got %v", sut.AggregateVersions())
	}
}
