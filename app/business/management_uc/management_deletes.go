package managementuc

import (
	"context"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	"usp-management-device-api/common/logging"
)

func (s *service) DeleteParametersWithParameterId(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	// start transaction
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	// Find parameter with Id
	param, err := s.store.FindParameter(txCtx, map[string]any{
		models.Parameter{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("parameter not found with id=%s: %v", id, err)
		return err
	}
	if param == nil || param.Id == nil {
		logging.Errorf("parameter not found or id is nil for id=%s", id)
		return apperrors.NewDBError(err, "parameterID not found")
	}

	// Change status of parameter to DELETE
	if err := s.store.ChangeStatusParameterToDelete(txCtx, param.Id.String(), updatedBy); err != nil {
		logging.Errorf("failed to change status of parameter id=%s: %v", id, err)
		return err
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit disable parameter: %v", err)
		return err
	}

	success = true
	logging.Infof("Parameter disabled successfully with ID: %s", id)
	return nil
}

func (s *service) DeleteProfileWithProfileId(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	// start transaction
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find profile with profileId
	profile, err := s.store.FindProfile(txCtx, map[string]any{
		models.Profile{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("profile not found with id=%s: %v", id, err)
		return err
	}
	if profile == nil || profile.Id == nil {
		logging.Errorf("profile not found or id is nil for id=%s", id)
		return apperrors.NewDBError(err, "profileID not found")
	}

	// Change status of profile to DELETE
	if err := s.store.ChangeStatusProfileToDelete(txCtx, profile.Id.String(), updatedBy); err != nil {
		logging.Errorf("failed to change status of profile id=%s: %v", id, err)
		return err
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit disable profile: %v", err)
		return err
	}

	success = true
	logging.Infof("Profile disabled successfully with ID: %s", id)
	return nil
}

func (s *service) DeleteModelWithId(
	ctx context.Context,
	id string,
	updatedBy string,
) error {
	// start transaction
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	// Find model with Id
	model, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("model not found with id=%s: %v", id, err)
		return err
	}
	if model == nil || model.Id == nil {
		logging.Errorf("model not found or id is nil for id=%s", id)
		return apperrors.NewDBError(err, "modelID not found")
	}

	// Disable model
	if err := s.store.ChangeStatusModelsToDelete(txCtx, model.Id.String(), updatedBy); err != nil {
		logging.Errorf("failed to disable model id=%s: %v", id, err)
		return err
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit disable model: %v", err)
		return err
	}

	success = true
	logging.Infof("Model disabled successfully with ID: %s", id)
	return nil
}

func (s *service) DeleteFirmwareWithId(
	ctx context.Context,
	modelId string,
	id string,
	updatedBy string,
) error {
	// start transaction
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	//Find model with modelId
	_, err = s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): modelId,
	})
	if err != nil {
		logging.Errorf("model not found with id=%s: %v", modelId, err)
		return apperrors.NewInvalidRequestError(err, "model not found with id="+modelId, "id")
	}
	// Find firmware with Id
	firmware, err := s.store.FindFirmware(txCtx, map[string]any{
		models.Firmware{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("firmware not found with id=%s: %v", id, err)
		return err
	}
	if firmware == nil || firmware.Id == nil {
		logging.Errorf("firmware not found or id is nil for id=%s", id)
		return apperrors.NewDBError(err, "firmwareID not found")
	}

	// Find model with modelId
	model, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): firmware.ModelId,
	})
	if err != nil {
		logging.Errorf("model not found with id=%v: %v", firmware.ModelId, err)
		return apperrors.NewDBError(err, "modelID not found")
	}
	if model == nil || model.Id == nil {
		logging.Errorf("model not found or id is nil for id=%v", firmware.ModelId)
		return apperrors.NewDBError(err, "modelID not found")
	}
	// Change status of firmware to DELETE
	if err := s.store.ChangeStatusFirmwareToDelete(txCtx, firmware.Id.String(), updatedBy); err != nil {
		logging.Errorf("failed to change status of firmware id=%s: %v", id, err)
		return err
	}

	//Move firmware file from firmware bucket to trash bucket in MinIO
	firmwarePath := model.Name + "/" + firmware.Name

	// For now, move to trash bucket with same path structure
	if err := s.minioStore.MoveFileBetweenBuckets("firmware", firmwarePath, "trash", firmwarePath); err != nil {
		logging.Errorf("failed to move file from firmware bucket to trash bucket in MinIO: %v", err)
		return err
	}
	logging.Infof("Moved firmware file from firmware bucket to trash bucket: %s", firmwarePath)
	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit disable firmware: %v", err)
		return err
	}

	success = true
	logging.Infof("Firmware disabled successfully with ID: %s", id)
	return nil
}

func (s *service) DeleteGroupWithId(
	ctx context.Context,
	modelId string,
	id string,
	updatedBy string,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find model with modelId
	_, err = s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): modelId,
	})
	if err != nil {
		logging.Errorf("model not found with id=%s: %v", modelId, err)
		return apperrors.NewInvalidRequestError(err, "model not found with id="+modelId, "id")
	}
	// Find group with Id
	group, err := s.store.FindGroup(txCtx, map[string]any{
		models.Group{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("group not found with id=%s: %v", id, err)
		return err
	}
	if group == nil || group.Id == nil {
		logging.Errorf("group not found or id is nil for id=%s", id)
		return apperrors.NewDBError(err, "groupID not found")
	}

	// Change status of group to DELETE
	if err := s.store.ChangeStatusGroupToDelete(txCtx, group.Id.String(), updatedBy); err != nil {
		logging.Errorf("failed to change status of group id=%s: %v", id, err)
		return err
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit disable group: %v", err)
		return err
	}

	success = true
	logging.Infof("Group disabled successfully with ID: %s", id)
	return nil
}

func (s *service) DeleteDeviceWithId(
	ctx context.Context,
	modelId string,
	id string,
	updatedBy string,
) error {
	// start transaction
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find model with modelId
	_, err = s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): modelId,
	})
	if err != nil {
		logging.Errorf("model not found with id=%s: %v", modelId, err)
		return apperrors.NewInvalidRequestError(err, "model not found with id="+modelId, "id")
	}
	// Find device with Id
	device, err := s.store.FindDevice(txCtx, map[string]any{
		models.Device{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("device not found with id=%s: %v", id, err)
		return err
	}
	if device == nil || device.Id == nil {
		logging.Errorf("device not found or id is nil for id=%s", id)
		return apperrors.NewDBError(err, "deviceID not found")
	}

	// Change status of device to DELETE
	if err := s.store.ChangeStatusDeviceToDelete(txCtx, device.Id.String(), updatedBy); err != nil {
		logging.Errorf("failed to change status of device id=%s: %v", id, err)
		return err
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit disable device: %v", err)
		return err
	}

	success = true
	logging.Infof("Device disabled successfully with ID: %s", id)
	return nil
}
