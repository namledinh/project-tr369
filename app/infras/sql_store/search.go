package uspstore

import (
	"context"
	"errors"
	"reflect"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"

	"gorm.io/gorm"
)

// queryConditionBuilder applies conditions to a GORM query
// Supports both single values (uses =) and slices (uses IN)
func queryConditionBuilder(query *gorm.DB, condition map[string]any) *gorm.DB {
	for key, value := range condition {
		if value == nil {
			continue
		}

		// Check if the value is a slice using reflection
		valueType := reflect.TypeOf(value)
		if valueType.Kind() == reflect.Slice {
			// Use IN query for slice values
			query = query.Where(key+" IN ?", value)
		} else {
			// Use equality query for single values
			query = query.Where(key+" = ?", value)
		}
	}
	return query
}

// FindDevice retrieves a device by condition and optionally preloads related device entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindDevice(
	ctx context.Context,
	condition map[string]any,
	moreKeys ...string,
) (*models.Device, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	// check if condtion value:
	// + has slice: using 'IN' query
	// + has normal value: using '=' query
	var device = models.Device{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPDeviceTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(device.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &device, nil
}

// FindGroup retrieves a group by condition and optionally preloads related group entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindGroup(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var group = models.Group{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPGroupTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(group.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &group, nil
}

// FindModel retrieves a model by condition and optionally preloads related model entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindModel(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var model = models.Model{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPModelTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(model.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &model, nil
}

// FindFirmware retrieves a firmware by condition and optionally preloads related firmware entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindFirmware(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.Firmware, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var firmware = models.Firmware{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).Debug().
		Table(models.USPFirmwareTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&firmware).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(firmware.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &firmware, nil
}

// FindProfile retrieves a profile by condition and optionally preloads related profile entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindProfile(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var profile = models.Profile{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPProfileTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(profile.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &profile, nil
}

// FindParameter retrieves a parameter by condition and optionally preloads related parameter entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindParameter(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.Parameter, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var parameter = models.Parameter{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPParameterTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&parameter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(parameter.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &parameter, nil
}

// FindProfileParameter retrieves a profile parameter by condition and optionally preloads related profile parameter entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindProfileParameter(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.ProfileParameter, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var profileParameter = models.ProfileParameter{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPProfileParameterTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&profileParameter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(profileParameter.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &profileParameter, nil
}

// FindModelCustomField retrieves a model custom field by condition and optionally preloads related model custom field entity.
// If record not found, returns an error indicating the entity does not exist.
func (s *store) FindModelCustomField(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*models.ModelCustomField, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var modelCustomField = models.ModelCustomField{}

	db := s.getDBFromContext(ctx)
	query := db.WithContext(ctx).
		Table(models.USPModelCustomFieldTableName)

	for _, key := range moreKeys {
		query = query.Preload(key)
	}

	query = queryConditionBuilder(query, condition)

	if err := query.First(&modelCustomField).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrEntityNotExist(modelCustomField.GetEntityName())
		}
		return nil, apperrors.NewDBError(err, s.GetDBName())
	}

	return &modelCustomField, nil
}

// count total profile with status ENABLE/DISABLE
func (s *store) CountProfilesByStatus(
	ctx context.Context,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPProfileTableName).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}

// count total parameter with status ENABLE/DISABLE
func (s *store) CountParametersByStatus(
	ctx context.Context,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPParameterTableName).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}

func (s *store) CountModelByStatus(
	ctx context.Context,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPModelTableName).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}

func (s *store) CountFirmwareByStatus(
	ctx context.Context,
	modelId string,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPFirmwareTableName).
		Where("model_id = ?", modelId).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}

func (s *store) CountGroupByStatus(
	ctx context.Context,
	modelId string,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPGroupTableName).
		Where("model_id = ?", modelId).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}
func (s *store) CountDeviceByStatus(
	ctx context.Context,
	modelId string,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPDeviceTableName).
		Where("model_id = ?", modelId).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}

func (s *store) CountDeviceByGroupId(
	ctx context.Context,
	modelId string,
	groupId string,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	var count int64
	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.USPDeviceTableName).
		Where("model_id = ?", modelId).
		Where("group_id = ?", groupId).
		Where("status IN ?", []string{"ENABLE", "DISABLE"}).
		Count(&count).Error; err != nil {
		return 0, apperrors.NewDBError(err, s.GetDBName())
	}

	return count, nil
}
