package managementuc

import (
	"context"
	"fmt"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
)

func (s *service) ListTotalProfiles(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Profile, error) {
	// Sử dụng ProfileQueryBuilder để xử lý logic ở tầng usecase
	qb := NewProfileQueryBuilder()

	// Set pagination
	if err := qb.SetPagination(oppts.Limit, oppts.Offset); err != nil {
		return nil, err
	}

	// Add base conditions
	for key, value := range condition {
		qb.AddCondition(key, value)
	}

	// Add filters with validation
	for _, filter := range oppts.FilterExpr {
		if err := qb.AddFilter(filter.Filter, filter.Op, filter.Value, filter.Join); err != nil {
			return nil, err
		}
	}

	// Add orders with validation
	for _, order := range oppts.OrderExpr {
		if err := qb.AddOrder(order.Field, order.Direction); err != nil {
			return nil, err
		}
	}

	// Build safe query với tất cả validations
	finalCondition, finalOpts, err := qb.BuildProfile()
	if err != nil {
		return nil, err
	}

	profiles, err := s.store.ListProfiles(ctx, finalCondition, finalOpts)
	if err != nil {
		return nil, err
	}
	profileNames := append([]models.Profile(nil), profiles...)
	return profileNames, nil
}

func (s *service) ListTotalParameters(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Parameter, error) {
	// Sử dụng ParameterQueryBuilder để xử lý logic ở tầng usecase
	qb := NewParameterQueryBuilder()

	// Set pagination
	if err := qb.SetPagination(oppts.Limit, oppts.Offset); err != nil {
		return nil, err
	}

	// Add base conditions
	for key, value := range condition {
		qb.AddCondition(key, value)
	}

	// Add filters with validation
	for _, filter := range oppts.FilterExpr {
		if err := qb.AddFilter(filter.Filter, filter.Op, filter.Value, filter.Join); err != nil {
			return nil, err
		}
	}

	// Add orders with validation
	for _, order := range oppts.OrderExpr {
		if err := qb.AddOrder(order.Field, order.Direction); err != nil {
			return nil, err
		}
	}

	// Build safe query với tất cả validations
	finalCondition, finalOpts, err := qb.BuildParameter()
	if err != nil {
		return nil, err
	}

	parameters, err := s.store.ListParameters(
		ctx,
		finalCondition,
		finalOpts,
	)
	if err != nil {
		return nil, err
	}
	parameterNames := append([]models.Parameter(nil), parameters...)
	return parameterNames, nil
}

func (s *service) ListTotalModels(
	ctx context.Context,
	condition map[string]any,
	opts models.QueryOptions,
) ([]models.Model, error) {
	// Sử dụng ModelQueryBuilder để xử lý logic ở tầng usecase
	qb := NewModelQueryBuilder()

	// Set pagination
	if err := qb.SetPagination(opts.Limit, opts.Offset); err != nil {
		return nil, err
	}

	// Add base conditions
	for key, value := range condition {
		qb.AddCondition(key, value)
	}

	// Add filters with validation
	for _, filter := range opts.FilterExpr {
		if err := qb.AddFilter(filter.Filter, filter.Op, filter.Value, filter.Join); err != nil {
			return nil, err
		}
	}

	// Add orders with validation
	for _, order := range opts.OrderExpr {
		if err := qb.AddOrder(order.Field, order.Direction); err != nil {
			return nil, err
		}
	}

	// Build safe query với tất cả validations
	finalCondition, finalOpts, err := qb.BuildModel()
	if err != nil {
		return nil, err
	}

	listmodels, err := s.store.ListModels(ctx, finalCondition, finalOpts)
	if err != nil {
		return nil, err
	}

	modelNames := append([]models.Model(nil), listmodels...)
	return modelNames, nil
}

func (s *service) ListFirmwares(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Firmware, error) {
	// Validate model exists first
	if modelId := condition["model_id"]; modelId != nil {
		_, err := s.store.FindModel(ctx, map[string]any{
			"id": modelId,
		})
		if err != nil {
			return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+fmt.Sprint(modelId), "find_model_error")
		}
	}

	// Sử dụng FirmwareQueryBuilder để xử lý logic ở tầng usecase
	qb := NewFirmwareQueryBuilder()

	// Set pagination
	if err := qb.SetPagination(oppts.Limit, oppts.Offset); err != nil {
		return nil, err
	}

	// Add base conditions
	for key, value := range condition {
		qb.AddCondition(key, value)
	}

	// Add filters with validation
	for _, filter := range oppts.FilterExpr {
		if err := qb.AddFilter(filter.Filter, filter.Op, filter.Value, filter.Join); err != nil {
			return nil, err
		}
	}

	// Add orders with validation
	for _, order := range oppts.OrderExpr {
		if err := qb.AddOrder(order.Field, order.Direction); err != nil {
			return nil, err
		}
	}

	// Build safe query với tất cả validations
	finalCondition, finalOpts, err := qb.BuildFirmware()
	if err != nil {
		return nil, err
	}

	firmwareList, err := s.store.ListFirmware(ctx, finalCondition, finalOpts)
	if err != nil {
		return nil, err
	}
	firmwareNames := append([]models.Firmware(nil), firmwareList...)

	return firmwareNames, nil
}

func (s *service) ListTotalGroups(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
	modelId string,
) ([]models.Group, error) {
	// Validate model exists first
	_, err := s.store.FindModel(ctx, map[string]any{
		"id": modelId,
	})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+modelId, "find_model_error")
	}

	// Sử dụng GroupQueryBuilder để xử lý logic ở tầng usecase
	qb := NewGroupQueryBuilder()

	// Set pagination
	if err := qb.SetPagination(oppts.Limit, oppts.Offset); err != nil {
		return nil, err
	}

	// Add base conditions
	for key, value := range condition {
		qb.AddCondition(key, value)
	}
	// Add modelId condition
	qb.AddCondition("model_id", modelId)

	// Add filters with validation
	for _, filter := range oppts.FilterExpr {
		if err := qb.AddFilter(filter.Filter, filter.Op, filter.Value, filter.Join); err != nil {
			return nil, err
		}
	}

	// Add orders with validation
	for _, order := range oppts.OrderExpr {
		if err := qb.AddOrder(order.Field, order.Direction); err != nil {
			return nil, err
		}
	}

	// Build safe query với tất cả validations
	finalCondition, finalOpts, err := qb.BuildGroup()
	if err != nil {
		return nil, err
	}
	groups, err := s.store.ListGroups(ctx, finalCondition, finalOpts)
	if err != nil {
		return nil, err
	}
	groupNames := append([]models.Group(nil), groups...)
	return groupNames, nil
}

func (s *service) ListTotalDevices(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]models.Device, error) {
	// Validate model exists first
	if modelId := condition["model_id"]; modelId != nil {
		_, err := s.store.FindModel(ctx, map[string]any{
			"id": modelId,
		})
		if err != nil {
			return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+fmt.Sprint(modelId), "find_model_error")
		}
	}

	// Sử dụng DeviceQueryBuilder để xử lý logic ở tầng usecase
	qb := NewDeviceQueryBuilder()

	// Set pagination
	if err := qb.SetPagination(oppts.Limit, oppts.Offset); err != nil {
		return nil, err
	}

	// Add base conditions
	for key, value := range condition {
		qb.AddCondition(key, value)
	}

	// Add filters with validation
	for _, filter := range oppts.FilterExpr {
		if err := qb.AddFilter(filter.Filter, filter.Op, filter.Value, filter.Join); err != nil {
			return nil, err
		}
	}

	// Add orders with validation
	for _, order := range oppts.OrderExpr {
		if err := qb.AddOrder(order.Field, order.Direction); err != nil {
			return nil, err
		}
	}

	// Build safe query với tất cả validations
	finalCondition, finalOpts, err := qb.BuildDevice()
	if err != nil {
		return nil, err
	}

	devices, err := s.store.ListDevices(ctx, finalCondition, finalOpts)
	if err != nil {
		return nil, err
	}
	deviceNames := append([]models.Device(nil), devices...)
	return deviceNames, nil
}

func (s *service) FindModelWithModelId(
	ctx context.Context,
	modelId string,
) (*models.Model, error) {
	condition := map[string]any{
		"id": modelId,
	}
	model, err := s.store.FindModel(ctx, condition)
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+modelId, "find_model_error")
	}
	return model, nil
}

func (s *service) GetFirmwareWithId(
	ctx context.Context,
	condition map[string]any,
) (*models.Firmware, error) {
	firmwareID := condition["id"]
	modelID := condition["model_id"]

	// check model exists
	_, err := s.store.FindModel(ctx, map[string]any{"id": modelID})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+fmt.Sprint(modelID), "db_error")
	}

	// find firmware
	firmware, err := s.store.FindFirmware(ctx, map[string]any{
		"id":       firmwareID,
		"model_id": modelID,
	})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "firmware not found with id: "+fmt.Sprint(firmwareID), "db_error")
	}

	return firmware, nil
}

func (s *service) GetGroupWithId(
	ctx context.Context,
	condition map[string]any,
) (*models.Group, error) {
	groupID := condition["id"]
	modelID := condition["model_id"]

	// check model exists
	_, err := s.store.FindModel(ctx, map[string]any{"id": modelID})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+fmt.Sprint(modelID), "db_error")
	}

	// find group
	group, err := s.store.FindGroup(ctx, map[string]any{
		"id":       groupID,
		"model_id": modelID,
	}, "Firmware")
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "group not found with id: "+fmt.Sprint(groupID), "db_error")
	}

	return group, nil
}

func (s *service) GetDeviceWithId(
	ctx context.Context,
	condition map[string]any,
) (*models.Device, error) {
	deviceID := condition["id"]
	modelID := condition["model_id"]

	// check model exists
	_, err := s.store.FindModel(ctx, map[string]any{"id": modelID})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+fmt.Sprint(modelID), "db_error")
	}

	// find device
	device, err := s.store.FindDevice(ctx, map[string]any{
		"id":       deviceID,
		"model_id": modelID,
	}, "Group", "Model")
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "device not found with id: "+fmt.Sprint(deviceID), "db_error")
	}

	return device, nil
}
func (s *service) CountProfilesByStatusUC(
	ctx context.Context,
) (int64, error) {
	return s.store.CountProfilesByStatus(ctx)
}

func (s *service) CountParametersByStatus(
	ctx context.Context,
) (int64, error) {
	return s.store.CountParametersByStatus(ctx)
}

func (s *service) CountGroupsByStatus(
	ctx context.Context,
	modelId string,
) (int64, error) {
	//Find model with modelId
	_, err := s.store.FindModel(ctx, map[string]any{
		"id": modelId,
	})
	if err != nil {
		return 0, apperrors.NewInvalidRequestError(err, "model not found with id: "+modelId, "find_model_error")
	}
	return s.store.CountGroupByStatus(ctx, modelId)
}

func (s *service) CountModelsByStatus(
	ctx context.Context,
) (int64, error) {
	return s.store.CountModelByStatus(ctx)
}

func (s *service) CountFirmwaresByStatus(
	ctx context.Context,
	modelId string,
) (int64, error) {
	//Find model with modelId
	_, err := s.store.FindModel(ctx, map[string]any{
		"id": modelId,
	})
	if err != nil {
		return 0, apperrors.NewInvalidRequestError(err, "model not found with id: "+modelId, "find_model_error")
	}
	return s.store.CountFirmwareByStatus(ctx, modelId)
}

func (s *service) CountDevicesByStatus(
	ctx context.Context,
	modelId string,
) (int64, error) {
	//Find model with modelId
	_, err := s.store.FindModel(ctx, map[string]any{
		"id": modelId,
	})
	if err != nil {
		return 0, apperrors.NewInvalidRequestError(err, "model not found with id: "+modelId, "find_model_error")
	}
	return s.store.CountDeviceByStatus(ctx, modelId)
}

func (s *service) GetParameters(
	ctx context.Context,
	condition map[string]any,
) ([]models.Parameter, error) {

	if condition == nil {
		condition = make(map[string]any)
	}
	condition["status"] = []string{"ENABLE", "DISABLE"}

	parameters, err := s.store.ListTotalParameters(ctx, condition)
	if err != nil {
		return nil, err
	}
	return parameters, nil
}

func (s *service) GetFirmwares(
	ctx context.Context,
	condition map[string]any,
	modelId string,
) ([]models.Firmware, error) {
	//Find model with modelId
	_, err := s.store.FindModel(ctx, map[string]any{
		"id": modelId,
	})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+modelId, "find_model_error")
	}

	if condition == nil {
		condition = make(map[string]any)
	}
	condition["status"] = []string{"ENABLE", "DISABLE"}

	firmwares, err := s.store.ListTotalFirmwares(ctx, condition, modelId)
	if err != nil {
		return nil, err
	}
	return firmwares, nil
}

func (s *service) TotalDevicesWithGroupId(
	ctx context.Context,
	modelId string,
	groupId string,
) (int64, error) {
	// check group exists
	_, err := s.store.FindGroup(ctx, map[string]any{"id": groupId})
	if err != nil {
		return 0, apperrors.NewInvalidRequestError(err, "group not found with id: "+fmt.Sprint(groupId), "db_error")
	}

	// count devices in group
	count, err := s.store.CountDeviceByGroupId(ctx, modelId, groupId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *service) TotalGroupsWithModelId(
	ctx context.Context,
	condition map[string]any,
	modelId string,
) ([]models.Group, error) {

	if condition == nil {
		condition = make(map[string]any)
	}
	condition["status"] = []string{"ENABLE", "DISABLE"}

	// check model exists
	_, err := s.store.FindModel(ctx, map[string]any{"id": modelId})
	if err != nil {
		return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+fmt.Sprint(modelId), "db_error")
	}

	// count groups in model
	groups, err := s.store.ListTotalGroups(ctx, condition, modelId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
