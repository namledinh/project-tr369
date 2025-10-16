package managementuc

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	utils "usp-management-device-api/common/utils"
)

func (s *service) ExportParametersCSV(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]byte, error) {
	if len(oppts.OrderExpr) == 0 {
		oppts.OrderExpr = []models.OrderExpr{
			{
				Field:     "updated_at",
				Direction: "DESC",
			},
		}
	}

	validColumns := map[string]bool{
		"path":        true,
		"data_type":   true,
		"description": true,
		"created_at":  true,
		"updated_at":  true,
		"updated_by":  true,
		"status":      true,
	}
	for _, expr := range oppts.OrderExpr {
		if _, ok := validColumns[strings.ToLower(expr.Field)]; !ok {
			return nil, apperrors.NewInvalidRequestError(
				nil,
				fmt.Sprintf("invalid order field: %s", expr.Field),
				"invalid order field",
			)
		}
	}

	hasStatusFilter := false
	for _, expr := range oppts.FilterExpr {
		if strings.EqualFold(expr.Filter, "status") {
			hasStatusFilter = true
			if strings.EqualFold(expr.Op, "eq") &&
				strings.EqualFold(fmt.Sprint(expr.Value), "DELETE") {
				return nil, apperrors.NewInvalidRequestError(nil, "cannot export models with DELETE status", "invalid status filter")
			}
		}
	}
	if !hasStatusFilter {
		condition["status"] = []string{"ENABLE", "DISABLE"}
	}

	parameters, err := s.store.ListParameters(ctx, condition, oppts)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// header
	headers := []string{"path", "data_type", "description", "created_at", "updated_at", "updated_by", "status"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// rows
	for _, p := range parameters {
		row := []string{
			p.Path,
			p.DataType,
			p.Description,
			utils.FormatTimeGMT7(p.CreatedAt, time.RFC3339),
			utils.FormatTimeGMT7(p.UpdatedAt, time.RFC3339),
			p.UpdatedBy,
			p.Status,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *service) ExportProfilesCSV(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]byte, error) {
	if len(oppts.OrderExpr) == 0 {
		oppts.OrderExpr = []models.OrderExpr{
			{
				Field:     "updated_at",
				Direction: "DESC",
			},
		}
	}

	validColumns := map[string]bool{
		"name":                   true,
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
		"status":                 true,
		"updated_by":             true,
		"description":            true,
	}
	for _, expr := range oppts.OrderExpr {
		if _, ok := validColumns[strings.ToLower(expr.Field)]; !ok {
			return nil, apperrors.NewInvalidRequestError(
				nil,
				fmt.Sprintf("invalid order field: %s", expr.Field),
				"invalid order field",
			)
		}
	}

	hasStatusFilter := false
	for _, expr := range oppts.FilterExpr {
		if strings.EqualFold(expr.Filter, "status") {
			hasStatusFilter = true
			if strings.EqualFold(expr.Op, "eq") &&
				strings.EqualFold(fmt.Sprint(expr.Value), "DELETE") {
				return nil, apperrors.NewInvalidRequestError(nil, "cannot export models with DELETE status", "invalid status filter")
			}
		}
	}
	if !hasStatusFilter {
		condition["status"] = []string{"ENABLE", "DISABLE"}
	}

	profiles, err := s.store.ListProfiles(ctx, condition, oppts)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// header
	headers := []string{"name",
		"msg_type",
		"return_commands",
		"return_events",
		"return_params",
		"return_unique_key_sets",
		"allow_partial", "send_resp",
		"first_level_only",
		"max_depth",
		"tags",
		"created_at",
		"updated_at",
		"updated_by",
		"status",
		"description",
		"parameter_paths",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// rows
	for _, p := range profiles {
		var paramPaths []string
		for _, pp := range p.ProfileParameters {
			if pp.Parameter != nil {
				paramPaths = append(paramPaths, pp.Parameter.Path)
			}
		}
		tagsStr := strings.Join(p.Tags, ";")
		paramPathsStr := strings.Join(paramPaths, ";")
		row := []string{
			p.Name,
			fmt.Sprintf("%v", p.MsgType),
			fmt.Sprintf("%v", p.ReturnCommands),
			fmt.Sprintf("%v", p.ReturnEvents),
			fmt.Sprintf("%v", p.ReturnParams),
			fmt.Sprintf("%v", p.ReturnUniqueKeySets),
			fmt.Sprintf("%v", p.AllowPartial),
			fmt.Sprintf("%v", p.SendResp),
			fmt.Sprintf("%v", p.FirstLevelOnly),
			fmt.Sprintf("%d", p.MaxDepth),
			tagsStr,
			utils.FormatTimeGMT7(p.CreatedAt, time.RFC3339),
			utils.FormatTimeGMT7(p.UpdatedAt, time.RFC3339),
			p.UpdatedBy,
			p.Status,
			p.Description,
			paramPathsStr,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *service) ExportDevicesCSV(
	ctx context.Context,
	condition map[string]any,
	oppts models.QueryOptions,
) ([]byte, error) {
	if len(oppts.OrderExpr) == 0 {
		oppts.OrderExpr = []models.OrderExpr{
			{
				Field:     "updated_at",
				Direction: "DESC",
			},
		}
	}

	validColumns := map[string]bool{
		"mac_address": true,
		"endpoint_id": true,
		"model_id":    true,
		"group_id":    true,
		"created_at":  true,
		"updated_at":  true,
		"updated_by":  true,
		"status":      true,
		"description": true,
	}
	for _, expr := range oppts.OrderExpr {
		if _, ok := validColumns[strings.ToLower(expr.Field)]; !ok {
			return nil, apperrors.NewInvalidRequestError(
				nil,
				fmt.Sprintf("invalid order field: %s", expr.Field),
				"invalid order field",
			)
		}
	}

	hasStatusFilter := false
	for _, expr := range oppts.FilterExpr {
		if strings.EqualFold(expr.Filter, "status") {
			hasStatusFilter = true
			if strings.EqualFold(expr.Op, "eq") &&
				strings.EqualFold(fmt.Sprint(expr.Value), "DELETE") {
				return nil, apperrors.NewInvalidRequestError(nil, "cannot export models with DELETE status", "invalid status filter")
			}
		}
	}
	if !hasStatusFilter {
		condition["status"] = []string{"ENABLE", "DISABLE"}
	}

	devices, err := s.store.ListDevices(ctx, condition, oppts)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// header
	headers := []string{
		"mac_address",
		"endpoint_id",
		"model_id",
		"group_id",
		"created_at",
		"updated_at",
		"updated_by",
		"status",
		"description",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// rows
	for _, d := range devices {
		var modelId, groupId string
		if d.ModelId != nil {
			modelId = d.ModelId.String()
		}
		if d.GroupId != nil {
			groupId = d.GroupId.String()
		}

		row := []string{
			d.MacAddress,
			d.EndpointId,
			modelId,
			groupId,
			utils.FormatTimeGMT7(d.CreatedAt, time.RFC3339),
			utils.FormatTimeGMT7(d.UpdatedAt, time.RFC3339),
			d.UpdatedBy,
			d.Status,
			d.Description,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *service) ExportFirmwaresCSV(
	ctx context.Context,
	condition map[string]any,
	modelId string,
	oppts models.QueryOptions,
) ([]byte, error) {
	if len(oppts.OrderExpr) == 0 {
		oppts.OrderExpr = []models.OrderExpr{
			{
				Field:     "updated_at",
				Direction: "DESC",
			},
		}
	}

	firmwares, err := s.store.ListTotalFirmwares(ctx, condition, modelId)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// header
	headers := []string{"name", "file_path", "description", "created_at", "updated_at", "updated_by", "status"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// rows
	for _, f := range firmwares {
		row := []string{
			f.Name,
			f.FilePath,
			f.Description,
			utils.FormatTimeGMT7(f.CreatedAt, time.RFC3339),
			utils.FormatTimeGMT7(f.UpdatedAt, time.RFC3339),
			f.UpdatedBy,
			f.Status,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *service) ExportGroupsCSV(
	ctx context.Context,
	condition map[string]any,
	modelId string,
	oppts models.QueryOptions,
) ([]byte, error) {
	// Always order by updated_at DESC
	if len(oppts.OrderExpr) == 0 {
		oppts.OrderExpr = []models.OrderExpr{
			{
				Field:     "updated_at",
				Direction: "DESC",
			},
		}
	}

	groups, err := s.store.ListTotalGroups(ctx, condition, modelId)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// header
	headers := []string{"name", "description", "created_at", "updated_at", "updated_by", "status"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// rows
	for _, g := range groups {
		row := []string{
			g.Name,
			g.Description,
			utils.FormatTimeGMT7(g.CreatedAt, time.RFC3339),
			utils.FormatTimeGMT7(g.UpdatedAt, time.RFC3339),
			g.UpdatedBy,
			g.Status,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
