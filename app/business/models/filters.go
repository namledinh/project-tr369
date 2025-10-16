package models

import (
	"math"
	"usp-management-device-api/common/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type QueryOptions struct {
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
	FilterExpr []FilterExpr `json:"filter_expr"`
	OrderExpr  []OrderExpr  `json:"order_expr"`
}

type FilterExpr struct {
	Filter string `json:"filter"`
	Op     string `json:"op"`
	Value  string `json:"value"`
	Join   string `json:"join"` // "AND" | "OR"
}

type OrderExpr struct {
	Field     string `json:"field"`     // name column
	Direction string `json:"direction"` // "asc" | "desc"
}

func (f *Filters) ValidateFilters(v *validator.Validator) {
	// Check that the page and page_size parameters contain sensible values.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10000000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	// Check that the sort parameter matches a value in the safelist.
	// v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// Define a new Metadata struct for holding the pagination metadata.
type Metadata struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	PageSize     int   `json:"page_size,omitempty"`
	FirstPage    int   `json:"first_page,omitempty"`
	LastPage     int   `json:"last_page,omitempty"`
	TotalRecords int64 `json:"total_records,omitempty"`
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that the last page value is calculated using the math.Ceil() function, which rounds
// up a float to the nearest integer. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be math.Ceil(12/5) = 3.
func (f *Filters) CalculateMetadata(totalRecords int64) Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  f.Page,
		PageSize:     f.PageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(f.PageSize))),
		TotalRecords: totalRecords,
	}
}
