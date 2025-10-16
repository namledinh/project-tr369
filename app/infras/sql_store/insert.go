package uspstore

import (
	"context"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
)

// InsertDevice inserts a new device into the database.
func (s *store) InsertDevice(
	ctx context.Context,
	device *models.Device,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(device.TableName()).Create(device).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertGroup inserts a new group into the database.
func (s *store) InsertGroup(
	ctx context.Context,
	group *models.Group,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(group.TableName()).Create(group).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertModel inserts a new model into the database.
func (s *store) InsertModel(
	ctx context.Context,
	model *models.Model,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(model.TableName()).Create(model).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertFirmware inserts a new firmware into the database.
func (s *store) InsertFirmware(
	ctx context.Context,
	firmware *models.Firmware,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(firmware.TableName()).Create(firmware).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertProfile inserts a new profile into the database.
func (s *store) InsertProfile(
	ctx context.Context,
	profile *models.Profile,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).
		Table(profile.TableName()).
		Create(profile).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertParameter inserts a new parameter into the database.
func (s *store) InsertParameter(
	ctx context.Context,
	parameter *models.Parameter,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(parameter.TableName()).Create(parameter).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertProfileParameter inserts a new profile parameter into the database.
func (s *store) InsertProfileParameter(
	ctx context.Context,
	profileParameter *models.ProfileParameter,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(profileParameter.TableName()).Create(profileParameter).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertModelCustomField inserts a new model custom field into the database.
func (s *store) InsertModelCustomField(
	ctx context.Context,
	modelCustomField *models.ModelCustomField,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).Table(modelCustomField.TableName()).Create(modelCustomField).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertParameters inserts multiple parameters into the database within a transaction.
func (s *store) InsertParametersBatch(
	ctx context.Context,
	parameters []*models.Parameter,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	if len(parameters) == 0 {
		return nil
	}

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).
		Table(parameters[0].TableName()).
		Create(&parameters).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

// InsertDevices inserts multiple devices into the database within a transaction.
func (s *store) InsertDevicesBatch(
	ctx context.Context,
	devices []*models.Device,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).
		Table(devices[0].TableName()).
		Create(&devices).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}

func (s *store) InsertProfileWithBatch(
	ctx context.Context,
	profiles []*models.Profile,
) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	if len(profiles) == 0 {
		return nil
	}

	db := s.getDBFromContext(ctx)

	if err := db.WithContext(ctx).
		Table(profiles[0].TableName()).
		Create(&profiles).Error; err != nil {
		return apperrors.NewDBError(err, s.GetDBName())
	}

	return nil
}
