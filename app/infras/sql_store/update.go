package uspstore

import (
	"context"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
)

// UpdateDevice updates an existing device in the database.
func (s *store) UpdateDevice(
	ctx context.Context,
	id string,
	device *models.DeviceUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(device.TableName()).
		Where("id = ?", id).
		Updates(device).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// UpdateGroup updates an existing group in the database.
func (s *store) UpdateGroup(
	ctx context.Context,
	id string,
	group *models.GroupUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(group.TableName()).
		Where("id = ?", id).
		Updates(group).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// UpdateModel updates an existing model in the database.
func (s *store) UpdateModel(
	ctx context.Context,
	id string,
	model *models.ModelUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(model.TableName()).
		Where("id = ?", id).
		Updates(model).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// UpdateFirmware updates an existing firmware in the database.
func (s *store) UpdateFirmware(
	ctx context.Context,
	id string,
	firmware *models.FirmwareUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(firmware.TableName()).
		Where("id = ?", id).
		Updates(firmware).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// UpdateProfile updates an existing profile in the database.
func (s *store) UpdateProfile(
	ctx context.Context,
	id string,
	profile *models.ProfileUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(profile.TableName()).
		Where("id = ?", id).
		Updates(profile).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// UpdateParameter updates an existing parameter in the database.
func (s *store) UpdateParameter(
	ctx context.Context,
	id string,
	parameter *models.ParameterUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).
		Table(parameter.TableName()).
		Where("id = ?", id).
		Updates(parameter).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}

// UpdateProfileParameter updates an existing profile parameter in the database.
func (s *store) UpdateProfileParameter(
	ctx context.Context,
	id string,
	profileParameter *models.ProfileParameterUpdate,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)
	if err := db.WithContext(ctx).Table(profileParameter.TableName()).
		Where("id = ?", id).
		Updates(profileParameter).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}
	return nil
}
