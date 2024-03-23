package aggregate

import (
	"reflect"
	"testing"
)

func Test_searchCriteria_AggregateIDs(t *testing.T) {
	sut := NewSearchCriteria(WithSearchAggregateIDs("1", "2", "3"))
	if !reflect.DeepEqual(sut.AggregateIDs(), []string{"1", "2", "3"}) {
		t.Fatalf("expected aggregate ids to be [1, 2, 3], got %v", sut.AggregateIDs())
	}
}

func Test_searchCriteria_AggregateNames(t *testing.T) {
	sut := NewSearchCriteria(WithSearchAggregateNames("test-1", "test-2", "test-3"))
	if !reflect.DeepEqual(sut.AggregateNames(), []string{"test-1", "test-2", "test-3"}) {
		t.Fatalf("expected aggregate names to be [test-1, test-2, test-3], got %v", sut.AggregateNames())
	}
}

func Test_searchCriteria_AggregateVersions(t *testing.T) {
	sut := NewSearchCriteria(WithSearchAggregateVersions(1, 2, 3))
	if !reflect.DeepEqual(sut.AggregateVersions(), []int{1, 2, 3}) {
		t.Fatalf("expected aggregate versions to be [1, 2, 3], got %v", sut.AggregateVersions())
	}
}

func Test_searchCriteria_Apply(t *testing.T) {
	sut := &searchCriteria{}
	sut.Apply(WithSearchAggregateIDs("1", "2", "3"))
	if !reflect.DeepEqual(sut.AggregateIDs(), []string{"1", "2", "3"}) {
		t.Fatalf("expected aggregate ids to be [1, 2, 3], got %v", sut.AggregateIDs())
	}
}
