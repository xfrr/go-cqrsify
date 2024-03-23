package aggregate

// SearchOption is the type that represents the options for criteriaing aggregates.
type SearchOption func(SearchCriteria)

// SearchCriteria is the interface that wraps the basic methods for criteriaing aggregates.
type SearchCriteria interface {
	// Apply applies the given options to the search criteria.
	Apply(...SearchOption) error
}

type searchCriteria struct {
	aggregateIDs      []string
	aggregateNames    []string
	aggregateVersions []int
}

func (s *searchCriteria) AggregateIDs() []string {
	return s.aggregateIDs
}

func (s *searchCriteria) AggregateNames() []string {
	return s.aggregateNames
}

func (s *searchCriteria) AggregateVersions() []int {
	return s.aggregateVersions
}

func (s *searchCriteria) Apply(options ...SearchOption) error {
	for _, option := range options {
		option(s)
	}
	return nil
}

// NewSearchCriteria returns a new search criteria.
func NewSearchCriteria(options ...SearchOption) *searchCriteria {
	criteria := &searchCriteria{
		aggregateIDs:      make([]string, 0),
		aggregateNames:    make([]string, 0),
		aggregateVersions: make([]int, 0),
	}
	criteria.Apply(options...)
	return criteria
}

// WithSearchAggregateIDs returns a search option that sets the aggregate ids to the search criteria.
func WithSearchAggregateIDs(ids ...string) SearchOption {
	return func(c SearchCriteria) {
		c.(*searchCriteria).aggregateIDs = ids
	}
}

// WithSearchAggregateNames returns a search option that sets the aggregate names to the search criteria.
func WithSearchAggregateNames(names ...string) SearchOption {
	return func(c SearchCriteria) {
		c.(*searchCriteria).aggregateNames = names
	}
}

// WithSearchAggregateVersions returns a search option that sets the aggregate versions to the search criteria.
func WithSearchAggregateVersions(versions ...int) SearchOption {
	return func(c SearchCriteria) {
		c.(*searchCriteria).aggregateVersions = versions
	}
}
