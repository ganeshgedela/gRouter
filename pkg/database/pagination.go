package database

// Pagination holds pagination and sorting parameters
type Pagination struct {
	Page     int
	PageSize int
	Sort     string                 // e.g., "created_at desc"
	Filters  map[string]interface{} // Dynamic filters, e.g., {"status": "active"}
}

// GetOffset computes the SQL offset
func (p Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10 // Default page size
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size constraint
func (p Pagination) GetLimit() int {
	if p.PageSize < 1 {
		return 10
	}
	return p.PageSize
}
