package uspstore

import (
	"context"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
)

func (s *store) DeleteDevice(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.Device{}.TableName()).Where("id = ?", id).Delete(&models.Device{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteGroup(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&models.Group{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteModel(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.Model{}.TableName()).Where("id = ?", id).Delete(&models.Model{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteFirmware(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.Firmware{}.TableName()).Where("id = ?", id).Delete(&models.Firmware{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteProfile(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.Profile{}.TableName()).Where("id = ?", id).Delete(&models.Profile{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteParameter(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.Parameter{}.TableName()).Where("id = ?", id).Delete(&models.Parameter{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteProfileParameter(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.ProfileParameter{}.TableName()).Where("id = ?", id).Delete(&models.ProfileParameter{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) DeleteModelCustomField(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(models.ModelCustomField{}.TableName()).Where("id = ?", id).Delete(&models.ModelCustomField{}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// change status from ENABLE to DISABLE
func (s *store) ChangeStatusParameterToDelete(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.Parameter{}.
			TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "DELETE",
			"updated_by": updatedBy,
		}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// change status profile_name from ENABLE to DISABLE
func (s *store) ChangeStatusProfileToDelete(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.Profile{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "DELETE",
			"updated_by": updatedBy,
		}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) ChangeStatusModelsToDelete(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.Model{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "DELETE",
			"updated_by": updatedBy,
		}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) ChangeStatusFirmwareToDelete(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.Firmware{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "DELETE",
			"updated_by": updatedBy,
		}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) ChangeStatusGroupToDelete(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.Group{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "DELETE",
			"updated_by": updatedBy,
		}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

func (s *store) ChangeStatusDeviceToDelete(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(models.Device{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "DELETE",
			"updated_by": updatedBy,
		}).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}
