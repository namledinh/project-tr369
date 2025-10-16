package uspstore

import (
	"context"
	"errors"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"

	"gorm.io/gorm"
)

func (s *store) ListProfileParameter(
	ctx context.Context,
	condition map[string]any,
) ([]models.ProfileParameter, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var profileParameters []models.ProfileParameter

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.ProfileParameter{}.TableName())

	query = queryConditionBuilder(query, condition)
	query = query.Order("updated_at DESC")

	if err := query.Find(&profileParameters).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.ProfileParameter{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return profileParameters, nil
}

func (s *store) ListProfiles(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Profile, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var profiles []models.Profile

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.Profile{}.TableName())

	query = queryConditionBuilder(query, condition)

	// Use Specification Pattern instead of applyFilterExprs directly
	var specs []Specification

	// Add filter specification
	if len(oppts.FilterExpr) > 0 {
		filterSpec := NewFilterSpecification(oppts.FilterExpr)
		specs = append(specs, filterSpec)
	}

	// Add order specification
	if len(oppts.OrderExpr) > 0 {
		orderSpec := NewOrderSpecification(oppts.OrderExpr)
		specs = append(specs, orderSpec)
	}

	// Add pagination specification
	if oppts.Limit > 0 {
		paginationSpec := NewPaginationSpecification(oppts.Limit, oppts.Offset)
		specs = append(specs, paginationSpec)
	}

	// Apply all specifications
	if len(specs) > 0 {
		compositeSpec := NewCompositeSpecification(specs...)
		query = compositeSpec.Apply(query)
	}

	if err := query.Preload("ProfileParameters.Parameter").Find(&profiles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Profile{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return profiles, nil
}

func (s *store) ListParameters(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Parameter, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var parameters []models.Parameter
	query := s.getDBFromContext(ctx).
		Table(models.Parameter{}.TableName()).
		WithContext(ctx)

	query = queryConditionBuilder(query, condition)

	var specs []Specification

	// Add filter specification
	if len(oppts.FilterExpr) > 0 {
		filterSpec := NewFilterSpecification(oppts.FilterExpr)
		specs = append(specs, filterSpec)
	}

	// Add order specification
	if len(oppts.OrderExpr) > 0 {
		orderSpec := NewOrderSpecification(oppts.OrderExpr)
		specs = append(specs, orderSpec)
	}

	// Add pagination specification
	if oppts.Limit > 0 {
		paginationSpec := NewPaginationSpecification(oppts.Limit, oppts.Offset)
		specs = append(specs, paginationSpec)
	}

	// Apply all specifications
	if len(specs) > 0 {
		compositeSpec := NewCompositeSpecification(specs...)
		query = compositeSpec.Apply(query)
	}

	if err := query.Find(&parameters).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Parameter{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return parameters, nil
}

func (s *store) ListModels(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var modelsList []models.Model

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.Model{}.TableName())

	query = queryConditionBuilder(query, condition)

	// Sử dụng Specification Pattern thay vì applyFilterExprs trực tiếp
	var specs []Specification

	// Add filter specification
	if len(oppts.FilterExpr) > 0 {
		filterSpec := NewFilterSpecification(oppts.FilterExpr)
		specs = append(specs, filterSpec)
	}

	// Add order specification
	if len(oppts.OrderExpr) > 0 {
		orderSpec := NewOrderSpecification(oppts.OrderExpr)
		specs = append(specs, orderSpec)
	}

	// Add pagination specification
	if oppts.Limit > 0 {
		paginationSpec := NewPaginationSpecification(oppts.Limit, oppts.Offset)
		specs = append(specs, paginationSpec)
	}

	// Apply all specifications
	if len(specs) > 0 {
		compositeSpec := NewCompositeSpecification(specs...)
		query = compositeSpec.Apply(query)
	}

	if err := query.Find(&modelsList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Model{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return modelsList, nil
}

func (s *store) ListFirmware(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Firmware, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var firmwareList []models.Firmware

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).Debug().
		Table(models.Firmware{}.TableName())

	query = queryConditionBuilder(query, condition)

	// Sử dụng Specification Pattern thay vì applyFilterExprs trực tiếp
	var specs []Specification

	// Add filter specification
	if len(oppts.FilterExpr) > 0 {
		filterSpec := NewFilterSpecification(oppts.FilterExpr)
		specs = append(specs, filterSpec)
	}

	// Add order specification
	if len(oppts.OrderExpr) > 0 {
		orderSpec := NewOrderSpecification(oppts.OrderExpr)
		specs = append(specs, orderSpec)
	}

	// Add pagination specification
	if oppts.Limit > 0 {
		paginationSpec := NewPaginationSpecification(oppts.Limit, oppts.Offset)
		specs = append(specs, paginationSpec)
	}

	// Apply all specifications
	if len(specs) > 0 {
		compositeSpec := NewCompositeSpecification(specs...)
		query = compositeSpec.Apply(query)
	}

	if err := query.Find(&firmwareList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Firmware{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return firmwareList, nil
}

func (s *store) ListGroups(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var groups []models.Group
	query := s.getDBFromContext(ctx).
		Table(models.Group{}.TableName()).
		Preload("Firmware").
		WithContext(ctx)

	query = queryConditionBuilder(query, condition)

	// Sử dụng Specification Pattern thay vì applyFilterExprs trực tiếp
	var specs []Specification

	// Add filter specification
	if len(oppts.FilterExpr) > 0 {
		filterSpec := NewFilterSpecification(oppts.FilterExpr)
		specs = append(specs, filterSpec)
	}

	// Add order specification
	if len(oppts.OrderExpr) > 0 {
		orderSpec := NewOrderSpecification(oppts.OrderExpr)
		specs = append(specs, orderSpec)
	}

	// Add pagination specification
	if oppts.Limit > 0 {
		paginationSpec := NewPaginationSpecification(oppts.Limit, oppts.Offset)
		specs = append(specs, paginationSpec)
	}

	// Apply all specifications
	if len(specs) > 0 {
		compositeSpec := NewCompositeSpecification(specs...)
		query = compositeSpec.Apply(query)
	}

	if err := query.Find(&groups).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Group{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return groups, nil
}

func (s *store) ListDevices(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Device, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var devices []models.Device
	query := s.getDBFromContext(ctx).
		Table(models.Device{}.TableName()).
		Preload("Group").
		WithContext(ctx)

	query = queryConditionBuilder(query, condition)

	// Use Specification Pattern thay vì applyFilterExprs trực tiếp
	var specs []Specification

	// Add filter specification
	if len(oppts.FilterExpr) > 0 {
		filterSpec := NewFilterSpecification(oppts.FilterExpr)
		specs = append(specs, filterSpec)
	}

	// Add order specification
	if len(oppts.OrderExpr) > 0 {
		orderSpec := NewOrderSpecification(oppts.OrderExpr)
		specs = append(specs, orderSpec)
	}

	// Add pagination specification
	if oppts.Limit > 0 {
		paginationSpec := NewPaginationSpecification(oppts.Limit, oppts.Offset)
		specs = append(specs, paginationSpec)
	}

	// Apply all specifications
	if len(specs) > 0 {
		compositeSpec := NewCompositeSpecification(specs...)
		query = compositeSpec.Apply(query)
	}

	if err := query.Find(&devices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Device{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return devices, nil
}

func (s *store) ListTotalParameters(
	ctx context.Context,
	condition map[string]any,
) ([]models.Parameter, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var parameters []models.Parameter

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.Parameter{}.TableName())

	query = queryConditionBuilder(query, condition)

	if err := query.Find(&parameters).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Parameter{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return parameters, nil
}

func (s *store) ListTotalFirmwares(
	ctx context.Context,
	condition map[string]any,
	modelId string,
) ([]models.Firmware, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var firmwares []models.Firmware

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Where("model_id = ?", modelId).
		Table(models.Firmware{}.TableName())

	query = queryConditionBuilder(query, condition)

	if err := query.Find(&firmwares).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Firmware{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return firmwares, nil
}

func (s *store) ListTotalGroups(
	ctx context.Context,
	condition map[string]any,
	modelId string,
) ([]models.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var groups []models.Group

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Where("model_id = ?", modelId).
		Table(models.Group{}.TableName())

	query = queryConditionBuilder(query, condition)

	if err := query.Find(&groups).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(models.Group{}.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return groups, nil
}
