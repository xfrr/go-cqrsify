package aggregate

import "slices"

type SearchCriteriaOptions struct {
	aggregateIDs      []string
	aggregateNames    []string
	aggregateVersions []int
}

func (sc *SearchCriteriaOptions) AggregateIDs() []string {
	return sc.aggregateIDs
}

func (sc *SearchCriteriaOptions) AggregateNames() []string {
	return sc.aggregateNames
}

func (sc *SearchCriteriaOptions) AggregateVersions() []int {
	return sc.aggregateVersions
}

func (sc *SearchCriteriaOptions) IsEmpty() bool {
	return len(sc.aggregateIDs) == 0 &&
		len(sc.aggregateNames) == 0 &&
		len(sc.aggregateVersions) == 0
}

// WithSearchAggregateIDs returns a search option that sets the aggregate ids to the search criteria.
func (sc *SearchCriteriaOptions) WithSearchAggregateIDs(ids ...string) *SearchCriteriaOptions {
	sc.aggregateIDs = ids
	return sc
}

// WithSearchAggregateNames returns a search option that sets the aggregate names to the search criteria.
func (sc *SearchCriteriaOptions) WithSearchAggregateNames(names ...string) *SearchCriteriaOptions {
	sc.aggregateNames = names
	return sc
}

// WithSearchAggregateVersions returns a search option that sets the aggregate versions to the search criteria.
func (sc *SearchCriteriaOptions) WithSearchAggregateVersions(versions ...int) *SearchCriteriaOptions {
	sc.aggregateVersions = versions
	return sc
}

// Matches checks if the given aggregate matches the search criteria.
func (sc *SearchCriteriaOptions) Matches(agg Aggregate[string]) bool {
	if sc.IsEmpty() {
		return true
	}

	if len(sc.aggregateIDs) > 0 && !contains(sc.aggregateIDs, agg.AggregateID()) {
		return false
	}

	if len(sc.aggregateNames) > 0 && !contains(sc.aggregateNames, agg.AggregateName()) {
		return false
	}

	if versionedAggregate, ok := agg.(VersionedAggregate[string]); ok {
		if len(sc.aggregateVersions) > 0 && !contains(sc.aggregateVersions, int(versionedAggregate.AggregateVersion())) {
			return false
		}
	}

	return true
}

func contains[T comparable](slice []T, item T) bool {
	return slices.Contains(slice, item)
}

func SearchCriteria() *SearchCriteriaOptions {
	return &SearchCriteriaOptions{
		aggregateIDs:      make([]string, 0),
		aggregateNames:    make([]string, 0),
		aggregateVersions: make([]int, 0),
	}
}
