package httpcontroller

import (
	// "context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	managementuc "usp-management-device-api/business/management_uc"
	httphelper "usp-management-device-api/common/http_helper"
	middleware "usp-management-device-api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(
		validator.WithRequiredStructEnabled(),
	)

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type invalidField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func generateInvalidFieldError(v validator.FieldError) invalidField {
	return invalidField{
		Field:   v.Field(),
		Message: fmt.Sprintf("value: %v does not satisfy the validation rule %s: %v", v.Value(), v.Tag(), v.Param()),
	}
}

type httpController struct {
	usecase managementuc.IManagementUsecase
}

func (h *httpController) ping() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(map[string]string{"result": "pong"}, nil, nil))
	}
}

func NewHTTPController(usecase managementuc.IManagementUsecase) *httpController {
	return &httpController{
		usecase: usecase,
	}
}

func (h *httpController) SetupRoute(
	router *gin.Engine,
	buildTime string,
	version string,
) {
	router.RedirectTrailingSlash = true
	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = true
	router.Use(middleware.CORSMiddleware())
	router.GET("/ping", h.ping())

	profiles := router.Group("/profiles")
	{
		profiles.GET("", h.listTotalProfiles())
		profiles.GET("/count", h.countProfilesByStatus())
		profiles.POST("/", h.createProfileWithParameterId())
		profiles.POST("/import-csv", h.createProfilesWithBatch())
		profiles.GET("/export-csv", h.exportProfilesCSV())
		profiles.PUT("/:profile_id", h.updateProfileWithParametersId())
		profiles.DELETE("/:profile_id", h.deleteProfileWithId())
	}

	parameters := router.Group("/parameters")
	{
		parameters.GET("", h.listTotalParameters())
		parameters.GET("/combobox", h.listParameters())
		parameters.GET("/count", h.countParametersByStatus())
		parameters.GET("/export-csv", h.exportParameterCSV())
		parameters.POST("/", h.createNewParameter())
		parameters.POST("/import-csv", h.createParametersFromCSVFile())
		parameters.PUT("/:parameter_id", h.updateParameterWithId())
		parameters.DELETE("/:parameter_id", h.deleteParametersWithId())
	}

	models := router.Group("/models")
	{
		// ----- Models -----
		models.GET("", h.listTotalModels())
		models.GET("/:model_id", h.getModelWithModelId())
		models.POST("", h.createModels())
		models.GET("/count", h.countModelsByStatus())
		models.PUT("/:model_id", h.updateModelWithModelId())
		models.DELETE("/:model_id", h.deleteModelWithId())

		// ----- Groups -----
		models.GET("/:model_id/groups", h.listTotalGroups())
		models.GET("/:model_id/groups/combobox", h.listGroups())
		models.GET("/:model_id/groups/:group_id", h.getGroupWithGroupId())
		models.POST("/:model_id/groups", h.createGroup())
		models.PUT("/:model_id/groups/:group_id", h.updateGroupWithId())
		models.DELETE("/:model_id/groups/:group_id", h.deleteGroupWithId())
		models.GET("/:model_id/groups/count", h.countGroupsByStatus())

		// ----- Firmwares -----
		models.GET("/:model_id/firmwares", h.listFirmwares())
		models.GET("/:model_id/firmwares/combobox", h.listTotalFirmwares())
		models.GET("/:model_id/firmwares/:firmware_id", h.getFirmwareWithId())
		models.POST("/:model_id/firmwares", h.createFirmware())
		models.PUT("/:model_id/firmwares/:firmware_id", h.updateFirmwareWithId())
		models.DELETE("/:model_id/firmwares/:firmware_id", h.deleteFirmwareWithId())
		models.GET("/:model_id/firmwares/count", h.countFirmwaresByStatus())

		// ----- Devices -----
		models.GET("/:model_id/devices", h.listTotalDevices())
		models.GET("/:model_id/groups/:group_id/devices/count", h.totalDevicesWithGroupId())
		models.GET("/:model_id/devices/:device_id", h.getDeviceWithId())
		models.GET("/:model_id/devices/count", h.countDevicesByStatus())
		models.POST("/:model_id/groups/:group_id/devices", h.createDevice())
		models.POST("/:model_id/groups/:group_id/devices/import-csv", h.createDevicesWithBatch())
		models.PUT("/:model_id/devices/:device_id", h.updateDeviceWithId())
		models.DELETE("/:model_id/devices/:device_id", h.deleteDeviceWithId())
	}
}
