package managementuc

import (
	"context"
	"io"
	"time"
	"usp-management-device-api/business/models"

	"github.com/google/uuid"
)

type service struct {
	store      iUSPStoreRepository
	minioStore iUSPMinioRepository
}

type IManagementUsecase interface {
	IProfileUsecase
	IParameterUsecase
	IModelUsecase
	IFirmwareUsecase
	IGroupUsecase
	IDeviceUsecase
}

type IProfileUsecase interface {
	UpsertProfileWithParameters(
		ctx context.Context,
		profile *models.Profile,
		parameters []models.Parameter,
	) error

	CreateProfileWithNewParameters(
		ctx context.Context,
		profile *models.Profile,
		parameters []models.Parameter,

	) error

	CreateProfileWithParameterId(
		ctx context.Context,
		profile *models.Profile,
		parameterIds []string,
	) (string, error)

	DeleteProfileWithProfileId(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	ListTotalProfiles(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Profile, error)

	CountProfilesByStatusUC(
		ctx context.Context,
	) (int64, error)

	CreateProfileWithBatch(
		ctx context.Context,
		file io.Reader,
		updatedBy string,
	) error

	UpdateProfileWithParameterId(
		ctx context.Context,
		id string,
		profileParameter *models.ProfileUpdate,
	) error

	ExportProfilesCSV(
		ctx context.Context,
		condition map[string]any,
		oppts models.QueryOptions,
	) ([]byte, error)
}

type IParameterUsecase interface {
	CreateNewParameter(
		ctx context.Context,
		parameters *models.Parameter,
	) (string, error)

	UpdateParameterWithId(
		ctx context.Context,
		parameterId string,
		parameter *models.ParameterUpdate,
	) error

	DeleteParametersWithParameterId(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	ListTotalParameters(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Parameter, error)

	CountParametersByStatus(
		ctx context.Context,
	) (int64, error)

	CreateParametersFromCSVFile(
		ctx context.Context,
		file io.Reader,
		updatedBy string,
		batchSize int,
	) ([]string, error)

	ExportParametersCSV(
		ctx context.Context,
		condition map[string]any,
		oppts models.QueryOptions,
	) ([]byte, error)

	GetParameters(
		ctx context.Context,
		condition map[string]any,
	) ([]models.Parameter, error)
}

type IModelUsecase interface {
	CreateModels(
		ctx context.Context,
		input *models.Model,
	) (string, error)

	UpdateModelWithModelId(
		ctx context.Context,
		modelID string,
		model *models.ModelUpdate,
	) error

	DeleteModelWithId(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	ListTotalModels(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Model, error)

	CountModelsByStatus(
		ctx context.Context,
	) (int64, error)

	FindModelWithModelId(
		ctx context.Context,
		modelId string,
	) (*models.Model, error)
}

type IFirmwareUsecase interface {
	CreateFirmware(
		ctx context.Context,
		file io.Reader,
		fileSize int64,
		firmware *models.Firmware,
	) (string, error)

	UpdateFirmwareWithId(
		ctx context.Context,
		condition map[string]any,
		file io.Reader,
		fileSize int64,
		firmware *models.FirmwareUpdate,
	) error

	DeleteFirmwareWithId(
		ctx context.Context,
		modelId string,
		id string,
		updatedBy string,
	) error

	ListFirmwares(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Firmware, error)

	CountFirmwaresByStatus(
		ctx context.Context,
		modelId string,
	) (int64, error)

	GetFirmwareWithId(
		ctx context.Context,
		condition map[string]any,
	) (*models.Firmware, error)

	GetFirmwares(
		ctx context.Context,
		condition map[string]any,
		modelId string,
	) ([]models.Firmware, error)
}

type IGroupUsecase interface {
	CreateGroup(
		ctx context.Context,
		group *models.Group,
	) (string, error)

	UpdateGroupWithId(
		ctx context.Context,
		modelId string,
		id string,
		group *models.GroupUpdate,
	) error

	DeleteGroupWithId(
		ctx context.Context,
		modelId string,
		id string,
		updatedBy string,
	) error

	ListTotalGroups(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
		modelId string,
	) ([]models.Group, error)

	CountGroupsByStatus(
		ctx context.Context,
		modelId string,
	) (int64, error)

	GetGroupWithId(
		ctx context.Context,
		condition map[string]any,
	) (*models.Group, error)

	TotalGroupsWithModelId(
		ctx context.Context,
		condition map[string]any,
		modelId string,
	) ([]models.Group, error)
}

type IDeviceUsecase interface {
	CreateDevice(
		ctx context.Context,
		device *models.Device,
	) (string, error)

	UpdateDeviceWithId(
		ctx context.Context,
		modelId string,
		id string,
		device *models.DeviceUpdate,
	) error

	DeleteDeviceWithId(
		ctx context.Context,
		modelId string,
		id string,
		updatedBy string,
	) error

	ListTotalDevices(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Device, error)

	CountDevicesByStatus(
		ctx context.Context,
		modelId string,
	) (int64, error)

	CreateDevicesWithBatch(
		ctx context.Context,
		file io.Reader,
		updatedBy string,
		batchSize int,
		modelID uuid.UUID,
		groupID uuid.UUID,
	) ([]string, error)

	GetDeviceWithId(
		ctx context.Context,
		condition map[string]any,
	) (*models.Device, error)

	TotalDevicesWithGroupId(
		ctx context.Context,
		modelId string,
		groupId string,
	) (int64, error)
}

func NewManagementUsecase(
	store iUSPStoreRepository,
	minioStore iUSPMinioRepository,
) IManagementUsecase {
	return &service{
		store:      store,
		minioStore: minioStore,
	}
}

type iUSPMinioRepository interface {
	UploadFile(bucketName, objectName string, reader io.Reader, size int64) error
	DownloadFile(bucketName, objectName, localPath string) error
	GetFileURL(bucketName, objectName string, expiry time.Duration) (string, error)
	MoveFileBetweenBuckets(srcBucket, srcPath, dstBucket, dstPath string) error
}

type iUSPStoreRepository interface {
	// BeginTx starts a transaction and stores it in the context
	BeginTx(ctx context.Context) (context.Context, error)

	// CommitTx commits a transaction stored in the context
	CommitTx(ctx context.Context) error

	// RollbackTx rolls back a transaction stored in the context
	RollbackTx(ctx context.Context) error

	// FindDevice retrieves a device by condition and optionally preloads related device entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindDevice(
		ctx context.Context,
		condition map[string]any,
		moreKeys ...string,
	) (*models.Device, error)

	// FindGroup retrieves a group by condition and optionally preloads related group entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindGroup(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.Group, error)

	// FindModel retrieves a model by condition and optionally preloads related model entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindModel(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.Model, error)

	// FindFirmware retrieves a firmware by condition and optionally preloads related firmware entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindFirmware(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.Firmware, error)

	// FindProfile retrieves a profile by condition and optionally preloads related profile entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindProfile(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.Profile, error)

	// FindParameter retrieves a parameter by condition and optionally preloads related parameter entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindParameter(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.Parameter, error)

	// FindProfileParameter retrieves a profile parameter by condition and optionally preloads related profile parameter entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindProfileParameter(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.ProfileParameter, error)

	// FindModelCustomField retrieves a model custom field by condition and optionally preloads related model custom field entity.
	// If record not found, returns an error indicating the entity does not exist.
	FindModelCustomField(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*models.ModelCustomField, error)

	ListProfileParameter(
		ctx context.Context,
		condition map[string]any,
	) ([]models.ProfileParameter, error)

	// InsertDevice inserts a new device into the database.
	InsertDevice(
		ctx context.Context,
		device *models.Device,
	) error

	// InsertGroup inserts a new group into the database.
	InsertGroup(
		ctx context.Context,
		group *models.Group,
	) error

	// InsertModel inserts a new model into the database.
	InsertModel(
		ctx context.Context,
		model *models.Model,
	) error

	// InsertFirmware inserts a new firmware into the database.
	InsertFirmware(
		ctx context.Context,
		firmware *models.Firmware,
	) error

	// InsertProfile inserts a new profile into the database.
	InsertProfile(
		ctx context.Context,
		profile *models.Profile,
	) error

	// InsertParameter inserts a new parameter into the database.
	InsertParameter(
		ctx context.Context,
		parameter *models.Parameter,
	) error

	// InsertProfileParameter inserts a new profile parameter into the database.
	InsertProfileParameter(
		ctx context.Context,
		profileParameter *models.ProfileParameter,
	) error

	// InsertModelCustomField inserts a new model custom field into the database.
	InsertModelCustomField(
		ctx context.Context,
		modelCustomField *models.ModelCustomField,
	) error

	// UpdateDevice updates an existing device in the database.
	UpdateDevice(
		ctx context.Context,
		id string,
		device *models.DeviceUpdate,
	) error

	// UpdateGroup updates an existing group in the database.
	UpdateGroup(
		ctx context.Context,
		id string,
		group *models.GroupUpdate,
	) error

	// UpdateModel updates an existing model in the database.
	UpdateModel(
		ctx context.Context,
		id string,
		model *models.ModelUpdate,
	) error

	// UpdateFirmware updates an existing firmware in the database.
	UpdateFirmware(
		ctx context.Context,
		id string,
		firmware *models.FirmwareUpdate,
	) error

	// UpdateProfile updates an existing profile in the database.
	UpdateProfile(
		ctx context.Context,
		id string,
		profile *models.ProfileUpdate,
	) error

	// UpdateParameter updates an existing parameter in the database.
	UpdateParameter(
		ctx context.Context,
		id string,
		parameter *models.ParameterUpdate,
	) error

	// UpdateProfileParameter updates an existing profile parameter in the database.
	UpdateProfileParameter(
		ctx context.Context,
		id string,
		profileParameter *models.ProfileParameterUpdate,
	) error

	// DeleteDevice deletes a device by ID from the database.
	DeleteDevice(
		ctx context.Context,
		id string,
	) error

	// DeleteGroup deletes a group by ID from the database.
	DeleteGroup(
		ctx context.Context,
		id string,
	) error

	// DeleteModel deletes a model by ID from the database.
	DeleteModel(
		ctx context.Context,
		id string,
	) error

	// DeleteFirmware deletes a firmware by ID from the database.
	DeleteFirmware(
		ctx context.Context,
		id string,
	) error

	// DeleteProfile deletes a profile by ID from the database.
	DeleteProfile(
		ctx context.Context,
		id string,
	) error

	// DeleteParameter deletes a parameter by ID from the database.
	DeleteParameter(
		ctx context.Context,
		id string,
	) error

	// DeleteProfileParameter deletes a profile parameter by ID from the database.
	DeleteProfileParameter(
		ctx context.Context,
		id string,
	) error

	// DeleteModelCustomField deletes a model custom field by ID from the database.
	DeleteModelCustomField(
		ctx context.Context,
		id string,
	) error

	// ListProfiles retrieves a list of profiles based on certain conditions
	ListProfiles(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Profile, error)

	ListParameters(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Parameter, error)

	ChangeStatusParameterToDelete(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	ChangeStatusProfileToDelete(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	CountProfilesByStatus(
		ctx context.Context,
	) (int64, error)

	CountParametersByStatus(
		ctx context.Context,
	) (int64, error)

	ListModels(
		ctx context.Context,
		condition map[string]any,
		oppts models.QueryOptions,
	) ([]models.Model, error)

	ChangeStatusModelsToDelete(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	CountModelByStatus(
		ctx context.Context,
	) (int64, error)

	ListFirmware(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Firmware, error)

	CountFirmwareByStatus(
		ctx context.Context,
		modelId string,
	) (int64, error)

	ChangeStatusFirmwareToDelete(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	ListGroups(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Group, error)

	ChangeStatusGroupToDelete(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	CountGroupByStatus(
		ctx context.Context,
		modelId string,
	) (int64, error)

	ListDevices(
		ctx context.Context,
		condition map[string]any,
		opts models.QueryOptions,
	) ([]models.Device, error)

	CountDeviceByStatus(
		ctx context.Context,
		modelId string,
	) (int64, error)

	ChangeStatusDeviceToDelete(
		ctx context.Context,
		id string,
		updatedBy string,
	) error

	InsertParametersBatch(
		ctx context.Context,
		parameters []*models.Parameter,
	) error

	InsertDevicesBatch(
		ctx context.Context,
		devices []*models.Device,
	) error

	InsertProfileWithBatch(
		ctx context.Context,
		profiles []*models.Profile,
	) error

	ListTotalParameters(
		ctx context.Context,
		condition map[string]any,
	) ([]models.Parameter, error)

	ListTotalFirmwares(
		ctx context.Context,
		condition map[string]any,
		modelId string,
	) ([]models.Firmware, error)

	CountDeviceByGroupId(
		ctx context.Context,
		modelId string,
		groupId string,
	) (int64, error)

	ListTotalGroups(
		ctx context.Context,
		condition map[string]any,
		modelId string,
	) ([]models.Group, error)
}
