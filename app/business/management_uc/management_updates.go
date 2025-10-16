package managementuc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"regexp"
	"strings"
	"time"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	"usp-management-device-api/common/logging"
)

// UpsertProfileWithParameters creates or updates a profile and its parameters.
// Both profile and parameters will be upserted (created or updated as needed).
func (s *service) UpsertProfileWithParameters(
	ctx context.Context,
	profile *models.Profile,
	parameters []models.Parameter,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	// First, upsert all parameters
	for index, param := range parameters {
		existingParam, err := s.store.FindParameter(
			txCtx,
			map[string]any{
				models.Parameter{}.GetPathColumnName(): param.Path,
			},
		)

		if err == nil {
			// Parameter exists - update it
			param.Id = existingParam.Id
			paramUpdate := models.NewParameterUpdate(&param.Path, &param.DataType, &param.Description, &param.Status, &param.UpdatedBy)
			if err := s.store.UpdateParameter(txCtx, existingParam.Id.String(), paramUpdate); err != nil {
				logging.Errorf("failed to update parameter: %v", err)
				return err
			}
			parameters[index].Id = param.Id // Update the slice with the ID
		} else {
			if appErr, ok := err.(*apperrors.AppError); ok && appErr.ErrorKey() == apperrors.ErrEntityNotExist {
				// Parameter doesn't exist - create it
				if err := s.store.InsertParameter(txCtx, &param); err != nil {
					logging.Errorf("failed to insert parameter: %v", err)
					return err
				}
				parameters[index].Id = param.Id // Update the slice with the generated ID
			} else {
				logging.Errorf("failed to find parameter: %v", err)
				return err
			}
		}
	}

	// Then, upsert the profile
	existingProfile, err := s.store.FindProfile(
		txCtx,
		map[string]any{
			models.Profile{}.GetProfileNameColumnName(): profile.Name,
		},
	)

	if err == nil {
		// Profile exists - update it
		profile.Id = existingProfile.Id
		profileUpdate := models.NewProfileUpdate()
		profileUpdate.Name = &profile.Name
		profileUpdate.MsgType = &profile.MsgType
		profileUpdate.ReturnCommands = &profile.ReturnCommands
		profileUpdate.ReturnEvents = &profile.ReturnEvents
		profileUpdate.ReturnParams = &profile.ReturnParams
		profileUpdate.ReturnUniqueKeySets = &profile.ReturnUniqueKeySets
		profileUpdate.AllowPartial = &profile.AllowPartial
		profileUpdate.SendResp = &profile.SendResp
		profileUpdate.FirstLevelOnly = &profile.FirstLevelOnly
		profileUpdate.MaxDepth = &profile.MaxDepth
		profileUpdate.Tags = profile.Tags
		profileUpdate.Status = &profile.Status
		if err := s.store.UpdateProfile(txCtx, existingProfile.Id.String(), profileUpdate); err != nil {
			logging.Errorf("failed to update profile: %v", err)
			return err
		}
	} else {
		if appErr, ok := err.(*apperrors.AppError); ok && appErr.ErrorKey() == apperrors.ErrEntityNotExist {
			// Profile doesn't exist - create it
			if err := s.store.InsertProfile(txCtx, profile); err != nil {
				logging.Errorf("failed to insert profile: %v", err)
				return err
			}
		} else {
			logging.Errorf("failed to find profile: %v", err)
			return err
		}
	}

	// Delete existing profile parameters
	profileParameters, err := s.store.ListProfileParameter(
		txCtx,
		map[string]any{
			models.ProfileParameter{}.GetProfileIdColumnName(): profile.Id,
		},
	)
	if err != nil {
		logging.Errorf("failed to list profile parameters: %v", err)
		return err
	}

	for _, pp := range profileParameters {
		if err := s.store.DeleteProfileParameter(txCtx, pp.Id.String()); err != nil {
			logging.Errorf("failed to delete profile parameter: %v", err)
			return err
		}
	}

	// Create new associations
	for _, param := range parameters {
		profileParam := &models.ProfileParameter{
			ProfileId:    profile.Id,
			ParameterId:  param.Id,
			DefaultValue: param.DefaultValue,
			Required:     param.Required,
		}

		if err := s.store.InsertProfileParameter(txCtx, profileParam); err != nil {
			logging.Errorf("failed to insert profile parameter: %v", err)
			return err
		}
	}

	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Profile and parameters upserted successfully: %s", profile.Name)
	return nil
}

// Update Parameter with parameterId
func (s *service) UpdateParameterWithId(
	ctx context.Context,
	parameterId string,
	parameter *models.ParameterUpdate,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find parameter with Id
	existingParam, err := s.store.FindParameter(txCtx, map[string]any{
		models.Parameter{}.GetIdColumnName(): parameterId,
	})
	if err != nil {
		logging.Errorf("parameter not found with id=%s: %v", parameterId, err)
		return apperrors.NewInvalidRequestError(err, "parameter not found with id="+parameterId, "id")
	}
	if existingParam == nil || existingParam.Id == nil {
		logging.Errorf("parameter not found or id is nil for id=%s", parameterId)
		return apperrors.NewInvalidRequestError(err, "parameter not found with id="+parameterId, "id")
	}
	//Find Parameter path exists
	if parameter.Path != nil && *parameter.Path != "" {
		existingParamPath, err := s.store.FindParameter(txCtx, map[string]any{
			models.Parameter{}.GetPathColumnName(): parameter.Path,
		})
		if err == nil && existingParamPath != nil && existingParamPath.Id.String() != parameterId {
			logging.Errorf("parameter path already exists with path: %s", *parameter.Path)
			return apperrors.NewInvalidRequestError(err, "parameter path already exists with path: "+*parameter.Path, "path")
		}
	}
	// Update parameter
	if err := s.store.UpdateParameter(txCtx, parameterId, parameter); err != nil {
		logging.Errorf("failed to update parameter: %v", err)
		return err
	}

	// Commit transaction nếu thành công
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Parameter updated successfully with ID: %s", parameterId)
	return nil
}

func (s *service) UpdateModelWithModelId(
	ctx context.Context,
	modelID string,
	model *models.ModelUpdate,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find model with Id
	existingModel, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): modelID,
	})
	if err != nil {
		logging.Errorf("model not found with id=%s: %v", modelID, err)
		return apperrors.NewInvalidRequestError(err, "model not found with id="+modelID, "id")
	}
	if existingModel == nil || existingModel.Id == nil {
		logging.Errorf("model not found or id is nil for id=%s", modelID)
		return apperrors.NewInvalidRequestError(err, "model not found with id="+modelID, "id")
	}

	// validate image size if image is not nil
	if model.Image != nil {
		if len(*model.Image) > 150*1024 {
			logging.Errorf("image size exceeds 150KB")
			return apperrors.NewInvalidRequestError(nil, "image size exceeds 150KB", "image")
		}
	}
	//Find Model name exists
	if model.Name != nil && *model.Name != "" {
		existingModelName, err := s.store.FindModel(txCtx, map[string]any{
			models.Model{}.GetNameColumnName(): *model.Name,
		})
		if err == nil && existingModelName != nil && existingModelName.Id.String() != modelID {
			logging.Errorf("model already exists with name: %s", *model.Name)
			return apperrors.NewInvalidRequestError(err, "model already exists with name: "+*model.Name, "name")
		}
	}
	// Update model
	if err := s.store.UpdateModel(txCtx, modelID, model); err != nil {
		logging.Errorf("failed to update model: %v", err)
		return err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Model updated successfully with ID: %s", modelID)
	return nil
}

func (s *service) UpdateFirmwareWithId(
	ctx context.Context,
	condition map[string]any,
	file io.Reader,
	fileSize int64,
	firmware *models.FirmwareUpdate,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find model with modelId
	model, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): condition["model_id"],
	})
	if err != nil {
		logging.Errorf("model not found with id=%v: %v", condition["model_id"], err)
		return apperrors.NewInvalidRequestError(err, "model not found with id="+condition["model_id"].(string), "id")
	}
	// Find firmware with firmwareId and modelId
	existingFirmware, err := s.store.FindFirmware(txCtx, map[string]any{
		models.Firmware{}.GetIdColumnName():      condition["id"],
		models.Firmware{}.GetModelIdColumnName(): condition["model_id"],
	})
	if err != nil {
		logging.Errorf("firmware not found with id=%v and model_id=%v: %v", condition["id"], condition["model_id"], err)
		return apperrors.NewInvalidRequestError(err, "firmware not found with id="+condition["id"].(string)+" and model_id="+condition["model_id"].(string), "id")
	}
	if existingFirmware == nil || existingFirmware.Id == nil {
		logging.Errorf("firmware not found or id is nil for id=%v and model_id=%v", condition["id"], condition["model_id"])
		return apperrors.NewInvalidRequestError(err, "firmware not found with id="+condition["id"].(string)+" and model_id="+condition["model_id"].(string), "id")
	}

	// Check if new name conflicts with existing firmware (excluding current firmware)
	if firmware.Name != nil && *firmware.Name != "" {
		existingFirmwareByName, err := s.store.FindFirmware(txCtx, map[string]any{
			models.Firmware{}.GetNameColumnName(): *firmware.Name,
		})
		if err == nil && existingFirmwareByName != nil && existingFirmwareByName.Id != nil {
			// Found a firmware with the same name
			if existingFirmwareByName.Id.String() != condition["id"].(string) {
				// It's a different firmware, not the one we're updating
				logging.Errorf("firmware name '%s' already exists", *firmware.Name)
				return apperrors.NewInvalidRequestError(nil, "firmware name '"+*firmware.Name+"' already exists", "name")
			}
		}
	}
	//Find Firmware name exists
	if firmware.Name != nil && *firmware.Name != "" {
		existingFirmwareName, err := s.store.FindFirmware(txCtx, map[string]any{
			models.Firmware{}.GetNameColumnName(): *firmware.Name,
		})
		if err == nil && existingFirmwareName != nil && existingFirmwareName.Id.String() != condition["id"].(string) {
			logging.Errorf("firmware already exists with name: %s", *firmware.Name)
			return apperrors.NewInvalidRequestError(err, "firmware already exists with name: "+*firmware.Name, "name")
		}
	}
	//Upload firmware file to MInIO
	if file == nil {
		// Update firmware
		if err := s.store.UpdateFirmware(txCtx, condition["id"].(string), firmware); err != nil {
			logging.Errorf("failed to update firmware: %v", err)
			return err
		}
	} else {
		objectName := model.Name + "/" + *firmware.Name
		if err := s.minioStore.UploadFile("firmware", objectName, file, fileSize); err != nil {
			logging.Errorf("failed to upload firmware to MinIO: %v", err)
			return err
		}
		// Get file URL
		url, err := s.minioStore.GetFileURL("firmware", objectName, 24*time.Hour)
		if err != nil {
			logging.Errorf("failed to get file URL: %v", err)
			return err
		}
		firmware.FilePath = &url
		// Update firmware
		if err := s.store.UpdateFirmware(txCtx, condition["id"].(string), firmware); err != nil {
			logging.Errorf("failed to update firmware: %v", err)
			return err
		}
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Firmware updated successfully with ID: %s", condition["id"].(string))
	return nil
}

func (s *service) UpdateGroupWithId(
	ctx context.Context,
	modelId string,
	id string,
	group *models.GroupUpdate,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
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
	existingGroup, err := s.store.FindGroup(txCtx, map[string]any{
		models.Group{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("group not found with id=%s: %v", id, err)
		return apperrors.NewInvalidRequestError(err, "group not found with id="+id, "id")
	}
	if existingGroup == nil || existingGroup.Id == nil {
		logging.Errorf("group not found or id is nil for id=%s", id)
		return apperrors.NewInvalidRequestError(err, "group not found with id="+id, "id")
	}
	//Find Group name exists
	if group.Name != nil && *group.Name != "" {
		existingGroupName, err := s.store.FindGroup(txCtx, map[string]any{
			models.Group{}.GetNameColumnName(): *group.Name,
		})
		if err == nil && existingGroupName != nil && existingGroupName.Id.String() != id {
			logging.Errorf("group already exists with name: %s", *group.Name)
			return apperrors.NewInvalidRequestError(err, "group already exists with name: "+*group.Name, "name")
		}
	}
	//Find firmware with firmwareId
	if group.FirmwareId != nil {
		_, err = s.store.FindFirmware(txCtx, map[string]any{
			models.Firmware{}.GetIdColumnName():      *group.FirmwareId,
			models.Firmware{}.GetModelIdColumnName(): modelId,
		})
		if err != nil {
			logging.Errorf("firmware_id=%s does not belong to model: %v", group.FirmwareId.String(), err)
			return apperrors.NewInvalidRequestError(err, "firmware id="+group.FirmwareId.String()+" does not belong to model", "firmware_id")
		}
	}

	//validate download period
	if group.DownloadPeriod != nil {
		var downloadPeriodRegex = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)~([01]\d|2[0-3]):([0-5]\d)$`)
		downloadPeriod := *group.DownloadPeriod
		if !downloadPeriodRegex.MatchString(downloadPeriod) {
			return apperrors.NewInvalidRequestError(nil, "invalid format, expected format: 'HH:MM~HH:MM'", "download_period")
		}
		layout := "15:04"
		parts := strings.Split(downloadPeriod, "~")
		if len(parts) != 2 {
			return apperrors.NewInvalidRequestError(nil, "invalid format, expected format: 'start~end'", "download_period")
		}
		from, err := time.Parse(layout, parts[0])
		if err != nil {
			return apperrors.NewInvalidRequestError(err, "invalid start time: "+parts[0], "download_period")
		}
		logging.Infof("Parsed start time: %s", from.Format(layout))
		to, err := time.Parse(layout, parts[1])
		if err != nil {
			return apperrors.NewInvalidRequestError(err, "invalid end time: "+parts[1], "download_period")
		}
		logging.Infof("Parsed end time: %s", to.Format(layout))
		// if from.Equal(to) {
		// 	return apperrors.NewInvalidRequestError(nil, "start time and end time cannot be the same", "download_period")
		// }
	}
	// Update group
	if err := s.store.UpdateGroup(txCtx, id, group); err != nil {
		logging.Errorf("failed to update group: %v", err)
		return err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Group updated successfully with ID: %s", id)
	return nil
}

func (s *service) UpdateDeviceWithId(
	ctx context.Context,
	modelId string,
	id string,
	device *models.DeviceUpdate,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
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
	existingDevice, err := s.store.FindDevice(txCtx, map[string]any{
		models.Device{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("device not found with id=%s: %v", id, err)
		return apperrors.NewInvalidRequestError(err, "device not found with id="+id, "id")
	}
	if existingDevice == nil || existingDevice.Id == nil {
		logging.Errorf("device not found or id is nil for id=%s", id)
		return apperrors.NewInvalidRequestError(err, "device not found with id="+id, "id")
	}
	//Find Group with groupId
	if device.GroupId != nil {
		_, err = s.store.FindGroup(txCtx, map[string]any{
			models.Group{}.GetIdColumnName():      *device.GroupId,
			models.Group{}.GetModelIdColumnName(): modelId,
		})
		if err != nil {
			logging.Errorf("group not found with id=%v for model_id=%s: %v", *device.GroupId, modelId, err)
			return apperrors.NewInvalidRequestError(err, "group not found with id="+(*device.GroupId).String()+" for model_id="+modelId, "group_id")
		}
	}
	// Update device
	if err := s.store.UpdateDevice(txCtx, id, device); err != nil {
		logging.Errorf("failed to update device: %v", err)
		return err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Device updated successfully with ID: %s", id)
	return nil
}

func (s *service) UpdateProfileWithParameterId(
	ctx context.Context,
	id string,
	profileParameter *models.ProfileUpdate,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Find profile with Id
	existingProfile, err := s.store.FindProfile(txCtx, map[string]any{
		models.Profile{}.GetIdColumnName(): id,
	})
	if err != nil {
		logging.Errorf("profile not found with id=%s: %v", id, err)
		return apperrors.NewInvalidRequestError(err, "profile not found with id="+id, "id")
	}
	if existingProfile == nil || existingProfile.Id == nil {
		logging.Errorf("profile not found or id is nil for id=%s", id)
		return apperrors.NewInvalidRequestError(err, "profile not found with id="+id, "id")
	}
	//Find parameters with id
	if profileParameter.Parameters != nil {
		for index, param := range profileParameter.Parameters {
			pid, err := uuid.Parse(param.Id)
			if err != nil {
				logging.Errorf("invalid parameter id at index %d: id=%s, error=%v", index, param.Id, err)
				return apperrors.NewInvalidRequestError(err, "invalid parameter id at index "+fmt.Sprintf("%d", index)+": "+param.Id, "parameters")
			}

			_, err = s.store.FindParameter(txCtx, map[string]any{
				models.Parameter{}.GetIdColumnName(): pid,
			})
			if err != nil {
				logging.Errorf("parameter not found at index %d: id=%s, error=%v", index, pid.String(), err)
				return apperrors.NewInvalidRequestError(err, "parameter not found at index "+fmt.Sprintf("%d", index)+": "+pid.String(), "parameters")
			}
		}
	}
	// Update profile
	if err := s.store.UpdateProfile(txCtx, id, profileParameter); err != nil {
		logging.Errorf("failed to update profile: %v", err)
		return err
	}
	//just update profile parameters if Parameters is not nil
	if profileParameter.Parameters != nil {
		// Delete existing profile parameters
		profileParameters, err := s.store.ListProfileParameter(txCtx, map[string]any{
			models.ProfileParameter{}.GetProfileIdColumnName(): existingProfile.Id,
		})
		if err != nil {
			logging.Errorf("failed to list profile parameters: %v", err)
			return err
		}

		for _, pp := range profileParameters {
			if err := s.store.DeleteProfileParameter(txCtx, pp.Id.String()); err != nil {
				logging.Errorf("failed to delete profile parameter: %v", err)
				return err
			}
		}
		// Create new associations
		for index, param := range profileParameter.Parameters {
			pid, err := uuid.Parse(param.Id)
			if err != nil {
				logging.Errorf("invalid parameter id at index %d during insertion: id=%s, error=%v", index, param.Id, err)
				return apperrors.NewInvalidRequestError(err, "invalid parameter id at index "+fmt.Sprintf("%d", index)+" during insertion: "+param.Id, "parameters")
			}

			profileParam := &models.ProfileParameter{
				ProfileId:   existingProfile.Id,
				ParameterId: &pid,
			}

			if err := s.store.InsertProfileParameter(txCtx, profileParam); err != nil {
				logging.Errorf("failed to insert profile parameter at index %d: parameter_id=%s, error=%v", index, pid.String(), err)
				return apperrors.NewInvalidRequestError(err, "failed to insert profile parameter at index "+fmt.Sprintf("%d", index)+": "+pid.String(), "parameters")
			}
		}
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Profile updated successfully with ID: %s", id)
	return nil
}
