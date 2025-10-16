package managementuc

import (
	"fmt"
	"strings"
	"usp-management-device-api/business/models"
	"usp-management-device-api/common/app_errors"
)

type QueryBuilder struct {
	conditions map[string]any
	filters    []models.FilterExpr
	orders     []models.OrderExpr
	limit      int
	offset     int
}

// NewQueryBuilder create query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		conditions: make(map[string]any),
		filters:    []models.FilterExpr{},
		orders:     []models.OrderExpr{},
	}
}

// AddCondition append base condition
func (qb *QueryBuilder) AddCondition(key string, value any) *QueryBuilder {
	qb.conditions[key] = value
	return qb
}

// AddFilter append filter expression with validation
func (qb *QueryBuilder) AddFilter(field, op, value, join string) error {
	// Validate operator
	validOps := []string{"eq", "ne", "lt", "gt", "lte", "gte", "like"}
	if !contains(validOps, strings.ToLower(op)) {
		return apperrors.NewInvalidRequestError(nil, fmt.Sprintf("invalid operator: %s", op), "invalid_operator")
	}

	// Validate join
	validJoins := []string{"and", "or"}
	if !contains(validJoins, strings.ToLower(join)) {
		return apperrors.NewInvalidRequestError(nil, fmt.Sprintf("invalid join: %s", join), "invalid_join")
	}

	qb.filters = append(qb.filters, models.FilterExpr{
		Filter: field,
		Op:     op,
		Value:  value,
		Join:   strings.ToUpper(join),
	})
	return nil
}

// AddOrder append order expression with validation
func (qb *QueryBuilder) AddOrder(field, direction string) error {
	// Validate direction
	validDirections := []string{"asc", "desc"}
	if !contains(validDirections, strings.ToLower(direction)) {
		return apperrors.NewInvalidRequestError(nil, fmt.Sprintf("invalid direction: %s", direction), "invalid_direction")
	}

	qb.orders = append(qb.orders, models.OrderExpr{
		Field:     field,
		Direction: strings.ToUpper(direction),
	})
	return nil
}

// SetPagination
func (qb *QueryBuilder) SetPagination(limit, offset int) error {
	if limit <= 0 || limit > 100 {
		return apperrors.NewInvalidRequestError(nil, "limit must be between 1 and 100", "invalid_limit")
	}
	if offset < 0 {
		return apperrors.NewInvalidRequestError(nil, "offset must be greater than or equal to 0", "invalid_offset")
	}

	qb.limit = limit
	qb.offset = offset
	return nil
}

// ValidateFields
func (qb *QueryBuilder) ValidateFields(validColumns map[string]bool) error {
	// Validate filter fields
	for _, filter := range qb.filters {
		if _, ok := validColumns[strings.ToLower(filter.Filter)]; !ok {
			return apperrors.NewInvalidRequestError(
				nil,
				fmt.Sprintf("invalid filter field: %s", filter.Filter),
				"invalid_filter_field",
			)
		}
	}

	// Validate order fields
	for _, order := range qb.orders {
		if _, ok := validColumns[strings.ToLower(order.Field)]; !ok {
			return apperrors.NewInvalidRequestError(
				nil,
				fmt.Sprintf("invalid order field: %s", order.Field),
				"invalid_order_field",
			)
		}
	}

	return nil
}

// ApplyStatusFilter default
func (qb *QueryBuilder) ApplyStatusFilter() error {
	hasStatusFilter := false
	for _, expr := range qb.filters {
		if strings.EqualFold(expr.Filter, "status") {
			hasStatusFilter = true
			if strings.EqualFold(expr.Op, "eq") && strings.EqualFold(expr.Value, "DELETE") {
				return apperrors.NewInvalidRequestError(
					nil,
					"cannot list items with DELETE status",
					"invalid_status_filter",
				)
			}
		}
	}

	if !hasStatusFilter {
		qb.conditions["status"] = []string{"ENABLE", "DISABLE"}
	}

	return nil
}

// SetDefaultOrder default order by updated_at desc
func (qb *QueryBuilder) SetDefaultOrder() {
	if len(qb.orders) == 0 {
		qb.orders = []models.OrderExpr{
			{
				Field:     "updated_at",
				Direction: "DESC",
			},
		}
	}
}

// Build QueryOptions and conditions
func (qb *QueryBuilder) Build() (map[string]any, models.QueryOptions) {
	return qb.conditions, models.QueryOptions{
		Limit:      qb.limit,
		Offset:     qb.offset,
		FilterExpr: qb.filters,
		OrderExpr:  qb.orders,
	}
}

func (qb *QueryBuilder) BuildSafeQuery(validColumns map[string]bool) (map[string]any, models.QueryOptions, error) {
	// Validate fields
	if err := qb.ValidateFields(validColumns); err != nil {
		return nil, models.QueryOptions{}, err
	}

	// Apply status filter
	if err := qb.ApplyStatusFilter(); err != nil {
		return nil, models.QueryOptions{}, err
	}

	// Set default order
	qb.SetDefaultOrder()

	// Build result
	conditions, opts := qb.Build()
	return conditions, opts, nil
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

/**/
type ProfileQueryBuilder struct {
	*QueryBuilder
	validColumns map[string]bool
}

// NewProfileQueryBuilder create query builder for Profile
func NewProfileQueryBuilder() *ProfileQueryBuilder {
	return &ProfileQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
		validColumns: map[string]bool{
			"id":                     true,
			"profile_name":           true,
			"msg_type":               true,
			"return_commands":        true,
			"return_events":          true,
			"return_params":          true,
			"return_unique_key_sets": true,
			"allow_partial":          true,
			"send_resp":              true,
			"first_level_only":       true,
			"max_depth":              true,
			"tags":                   true,
			"created_at":             true,
			"updated_at":             true,
			"updated_by":             true,
			"status":                 true,
		},
	}
}
func (pqb *ProfileQueryBuilder) BuildProfile() (map[string]any, models.QueryOptions, error) {
	return pqb.BuildSafeQuery(pqb.validColumns)
}

/**/
type ParameterQueryBuilder struct {
	*QueryBuilder
	validColumns map[string]bool
}

func NewParameterQueryBuilder() *ParameterQueryBuilder {
	return &ParameterQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
		validColumns: map[string]bool{
			"id":          true,
			"path":        true,
			"data_type":   true,
			"description": true,
			"created_at":  true,
			"updated_at":  true,
			"updated_by":  true,
			"status":      true,
		},
	}
}
func (pqb *ParameterQueryBuilder) BuildParameter() (map[string]any, models.QueryOptions, error) {
	return pqb.BuildSafeQuery(pqb.validColumns)
}

/**/
type ModelQueryBuilder struct {
	*QueryBuilder
	validColumns map[string]bool
}

// NewModelQueryBuilder create query builder for Model
func NewModelQueryBuilder() *ModelQueryBuilder {
	return &ModelQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
		validColumns: map[string]bool{
			"id":           true,
			"name":         true,
			"vendor_name":  true,
			"manufacturer": true,
			"status":       true,
			"description":  true,
			"created_at":   true,
			"updated_at":   true,
			"updated_by":   true,
		},
	}
}
func (mqb *ModelQueryBuilder) BuildModel() (map[string]any, models.QueryOptions, error) {
	return mqb.BuildSafeQuery(mqb.validColumns)
}

/**/
type FirmwareQueryBuilder struct {
	*QueryBuilder
	validColumns map[string]bool
}

func NewFirmwareQueryBuilder() *FirmwareQueryBuilder {
	return &FirmwareQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
		validColumns: map[string]bool{
			"id":          true,
			"model_id":    true,
			"name":        true,
			"file_path":   true,
			"status":      true,
			"description": true,
			"created_at":  true,
			"updated_at":  true,
			"updated_by":  true,
		},
	}
}
func (fqb *FirmwareQueryBuilder) BuildFirmware() (map[string]any, models.QueryOptions, error) {
	return fqb.BuildSafeQuery(fqb.validColumns)
}

/**/
type GroupQueryBuilder struct {
	*QueryBuilder
	validColumns map[string]bool
}

// NewGroupQueryBuilder create query builder for Group
func NewGroupQueryBuilder() *GroupQueryBuilder {
	return &GroupQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
		validColumns: map[string]bool{
			"id":          true,
			"model_id":    true,
			"firmware_id": true,
			"name":        true,
			"description": true,
			"status":      true,
			"created_at":  true,
			"updated_at":  true,
			"updated_by":  true,
		},
	}
}
func (gqb *GroupQueryBuilder) BuildGroup() (map[string]any, models.QueryOptions, error) {
	return gqb.BuildSafeQuery(gqb.validColumns)
}

/**/
type DeviceQueryBuilder struct {
	*QueryBuilder
	validColumns map[string]bool
}

// NewDeviceQueryBuilder create query builder for Device
func NewDeviceQueryBuilder() *DeviceQueryBuilder {
	return &DeviceQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
		validColumns: map[string]bool{
			"id":          true,
			"mac_address": true,
			"endpoint_id": true,
			"model_id":    true,
			"group_id":    true,
			"created_at":  true,
			"updated_at":  true,
			"updated_by":  true,
			"status":      true,
			"description": true,
		},
	}
}
func (dqb *DeviceQueryBuilder) BuildDevice() (map[string]any, models.QueryOptions, error) {
	return dqb.BuildSafeQuery(dqb.validColumns)
}
