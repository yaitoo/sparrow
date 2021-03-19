package db

type Search struct {
	// db               *DB
	whereConditions   []map[string]interface{}
	orConditions      []map[string]interface{}
	notConditions     []map[string]interface{}
	havingConditions  []map[string]interface{}
	joinConditions    []map[string]interface{}
	inConditions      []map[string]interface{}
	initAttrs         []interface{}
	assignAttrs       []interface{}
	selects           map[string]interface{}
	omits             []string
	orders            []interface{}
	orderByDescending bool
	// preload          []searchPreload
	offset           interface{}
	limit            interface{}
	group            string
	tableName        string
	raw              bool
	Unscoped         bool
	ignoreOrderQuery bool
}

func (s *Search) Where(query interface{}, values ...interface{}) *Search {
	s.whereConditions = append(s.whereConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *Search) OR(query interface{}, values ...interface{}) *Search {
	s.orConditions = append(s.orConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *Search) IN(query interface{}, values ...interface{}) *Search {
	s.inConditions = append(s.inConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *Search) OrderBy(values ...interface{}) *Search {
	s.orders = append(s.orders, values...)
	return s
}

func (s *Search) OrderByDescending(values ...interface{}) *Search {
	s.orders = append(s.orders, values...)
	s.orderByDescending = true
	return s
}

func (s *Search) Limit(limit interface{}) *Search {
	s.limit = limit
	return s
}

func (s *Search) Offset(offset interface{}) *Search {
	s.offset = offset
	return s
}
