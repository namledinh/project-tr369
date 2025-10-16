package uspstore

import (
	"fmt"
	"strings"
	"usp-management-device-api/business/models"

	"gorm.io/gorm"
)

// Specification interface defined how to apply conditions to GORM queries
type Specification interface {
	Apply(db *gorm.DB) *gorm.DB
}

// FilterSpecification implementation for FilterExpr
type FilterSpecification struct {
	filters []models.FilterExpr
}

// NewFilterSpecification create FilterSpecification
func NewFilterSpecification(filters []models.FilterExpr) *FilterSpecification {
	return &FilterSpecification{filters: filters}
}

// Apply applies filters to GORM query
func (fs *FilterSpecification) Apply(db *gorm.DB) *gorm.DB {
	return applyFilterExprsV2(db, fs.filters)
}

// OrderSpecification implementation for OrderExpr
type OrderSpecification struct {
	orders []models.OrderExpr
}

// NewOrderSpecification create OrderSpecification
func NewOrderSpecification(orders []models.OrderExpr) *OrderSpecification {
	return &OrderSpecification{orders: orders}
}

// Apply applies orders to GORM query
func (os *OrderSpecification) Apply(db *gorm.DB) *gorm.DB {
	for _, o := range os.orders {
		direction := "ASC"
		if strings.ToUpper(o.Direction) == "DESC" {
			direction = "DESC"
		}
		db = db.Order(fmt.Sprintf("%s %s", o.Field, direction))
	}
	return db
}

// PaginationSpecification implementation for pagination
type PaginationSpecification struct {
	limit  int
	offset int
}

// NewPaginationSpecification create PaginationSpecification
func NewPaginationSpecification(limit, offset int) *PaginationSpecification {
	return &PaginationSpecification{limit: limit, offset: offset}
}

// Apply applies pagination to GORM query
func (ps *PaginationSpecification) Apply(db *gorm.DB) *gorm.DB {
	return db.Limit(ps.limit).Offset(ps.offset)
}

// CompositeSpecification implementation for combining multiple specifications
type CompositeSpecification struct {
	specs []Specification
}

// NewCompositeSpecification create CompositeSpecification
func NewCompositeSpecification(specs ...Specification) *CompositeSpecification {
	return &CompositeSpecification{specs: specs}
}

// Apply applies all specifications
func (cs *CompositeSpecification) Apply(db *gorm.DB) *gorm.DB {
	for _, spec := range cs.specs {
		db = spec.Apply(db)
	}
	return db
}

func applyFilterExprsV2(db *gorm.DB, filters []models.FilterExpr) *gorm.DB {
	if len(filters) == 0 {
		return db
	}

	// Group filters by join type
	groups := groupFiltersByJoin(filters)

	for i, group := range groups {
		query := buildGroupQuery(group.filters)

		if i == 0 {
			db = db.Where(query.condition, query.args...)
		} else {
			if strings.ToUpper(group.joinType) == "OR" {
				db = db.Or(query.condition, query.args...)
			} else {
				db = db.Where(query.condition, query.args...)
			}
		}
	}

	return db
}

// FilterGroup groups filters theo join type
type FilterGroup struct {
	filters  []models.FilterExpr
	joinType string
}

// QueryResult result of build query
type QueryResult struct {
	condition string
	args      []any
}

// groupFiltersByJoin groups filters by join type
func groupFiltersByJoin(filters []models.FilterExpr) []FilterGroup {
	if len(filters) == 0 {
		return nil
	}

	var groups []FilterGroup
	var currentGroup FilterGroup

	for i, filter := range filters {
		if i == 0 {
			// First filter starts new group
			currentGroup = FilterGroup{
				filters:  []models.FilterExpr{filter},
				joinType: "AND", // Default join type
			}
		} else {
			if strings.ToUpper(filter.Join) == "OR" {
				// Continue current group if OR
				currentGroup.filters = append(currentGroup.filters, filter)
			} else {
				// Start new group if AND
				groups = append(groups, currentGroup)
				currentGroup = FilterGroup{
					filters:  []models.FilterExpr{filter},
					joinType: filter.Join,
				}
			}
		}
	}

	// Add last group
	if len(currentGroup.filters) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// buildGroupQuery builds query for a group of filters
func buildGroupQuery(filters []models.FilterExpr) QueryResult {
	if len(filters) == 0 {
		return QueryResult{}
	}

	var conditions []string
	var args []any

	for i, filter := range filters {
		condition, filterArgs := buildCondition(filter)
		conditions = append(conditions, condition)
		args = append(args, filterArgs...)

		// Add join operator if not the last filter
		if i < len(filters)-1 {
			// All filters in a group should be joined by OR within parentheses
		}
	}

	// Join conditions with OR and wrap in parentheses
	finalCondition := "(" + strings.Join(conditions, " OR ") + ")"

	return QueryResult{
		condition: finalCondition,
		args:      args,
	}
}

// buildCondition xây dựng điều kiện cho một filter
func buildCondition(filter models.FilterExpr) (string, []any) {
	op := getSQLOperator(filter.Op)
	condition := fmt.Sprintf("%s %s ?", filter.Filter, op)

	value := filter.Value
	if strings.ToLower(filter.Op) == "like" {
		value = "%" + filter.Value + "%"
	}

	return condition, []any{value}
}

// getSQLOperator chuyển đổi operator thành SQL operator
func getSQLOperator(op string) string {
	switch strings.ToLower(op) {
	case "eq":
		return "="
	case "ne":
		return "!="
	case "lt":
		return "<"
	case "gt":
		return ">"
	case "lte":
		return "<="
	case "gte":
		return ">="
	case "like":
		return "LIKE"
	default:
		return "="
	}
}
