package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfrr/go-cqrsify/domain/aggregate"
)

func Test_SearchCriteria_AggregateIDs(t *testing.T) {
	sut := aggregate.SearchCriteria().WithSearchAggregateIDs("1", "2", "3")
	require.ElementsMatch(t, sut.AggregateIDs(), []string{"1", "2", "3"})
}

func Test_SearchCriteria_AggregateNames(t *testing.T) {
	sut := aggregate.SearchCriteria().WithSearchAggregateNames("test-1", "test-2", "test-3")
	require.ElementsMatch(t, sut.AggregateNames(), []string{"test-1", "test-2", "test-3"})
}

func Test_SearchCriteria_AggregateVersions(t *testing.T) {
	sut := aggregate.SearchCriteria().WithSearchAggregateVersions(1, 2, 3)
	require.ElementsMatch(t, sut.AggregateVersions(), []int{1, 2, 3})
}
