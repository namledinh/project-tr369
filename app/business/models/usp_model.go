package models

import (
	"github.com/google/uuid"
	"time"
)

// CREATE TABLE public.models (
// 	id UUID NOT NULL DEFAULT gen_random_uuid(),
// 	name VARCHAR NOT NULL,
// 	vendor_name VARCHAR NOT NULL,
// 	manufacturer VARCHAR NOT NULL,
// 	status public.model_status NOT NULL DEFAULT 'ENABLE':::public.model_status,
// 	description VARCHAR(255) NULL,
// 	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
// 	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
// 	CONSTRAINT models_pkey PRIMARY KEY (id ASC),
// 	UNIQUE INDEX models_name_idx (name ASC) STORING (vendor_name, manufacturer, status, description, created_at, updated_at),
//     INDEX models_vendor_name_idx (vendor_name ASC) STORING (name, manufacturer, status, description, created_at, updated_at)
// );
// COMMENT ON COLUMN public.models.name IS 'factory model name: Model Wifi 6';
// COMMENT ON COLUMN public.models.vendor_name IS 'factory model name: AX3000S';
// COMMENT ON COLUMN public.models.manufacturer IS 'ODM name, such as: CIG, Skyworth';

const USPModelTableName = "models"
const USPModedlEntityName = "Model"

type Model struct {
	Id        *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	CreatedAt *time.Time `json:"-" gorm:"created_at"`
	UpdatedAt *time.Time `json:"-" gorm:"updated_at"`

	Name         string `json:"name" gorm:"column:name;type:varchar;not null"`
	VendorName   string `json:"vendor_name" gorm:"column:vendor_name;type:varchar;not null"`
	Manufacturer string `json:"manufacturer" gorm:"column:manufacturer;type:varchar;not null"`
	Status       string `json:"status" gorm:"column:status;type:public.model_status;default:'ENABLE':::public.model_status"`
	Description  string `json:"description,omitempty" gorm:"column:description;type:varchar(255);default:null"`
	UpdatedBy    string `json:"updated_by,omitempty" gorm:"column:updated_by;type:varchar(255);default:null"`
	Image        string `json:"image,omitempty" gorm:"column:image;type:text;default:null"`

	Groups            []Group            `gorm:"foreignKey:ModelId;references:Id;" json:"groups,omitempty"`
	ModelCustomFields []ModelCustomField `gorm:"foreignKey:ModelId;references:Id;" json:"model_custom_fields,omitempty"`
	Firmwares         []Firmware         `gorm:"foreignKey:ModelId;references:Id;" json:"firmwares,omitempty"`
}

func (Model) TableName() string     { return USPModelTableName }
func (Model) GetEntityName() string { return USPModedlEntityName }

func (Model) GetIdColumnName() string           { return "id" }
func (Model) GetNameColumnName() string         { return "name" }
func (Model) GetVendorNameColumnName() string   { return "vendor_name" }
func (Model) GetManufacturerColumnName() string { return "manufacturer" }
func (Model) GetStatusColumnName() string       { return "status" }
func (Model) GetDescriptionColumnName() string  { return "description" }
func (Model) GetCreatedAtColumnName() string    { return "created_at" }
func (Model) GetUpdatedAtColumnName() string    { return "updated_at" }

type ModelUpdate struct {
	Name         *string    `gorm:"column:name;type:varchar;not null" json:"name"`
	VendorName   *string    `gorm:"column:vendor_name;type:varchar;not null" json:"vendor_name"`
	Manufacturer *string    `gorm:"column:manufacturer;type:varchar;not null" json:"manufacturer"`
	Status       *string    `gorm:"column:status;type:public.model_status;default:'ENABLE':::public.model_status" json:"status"`
	Description  *string    `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`
	UpdatedBy    *string    `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"column:updated_at"`
	Image        *string    `json:"image,omitempty" gorm:"column:image;type:text;default:null"`
}

func (ModelUpdate) TableName() string     { return USPModelTableName }
func (ModelUpdate) GetEntityName() string { return USPModedlEntityName }

func NewModelUpdate(
	name *string,
	vendorName *string,
	manufacturer *string,
	status *string,
	description *string,
	updatedBy *string,
	image *string,
) *ModelUpdate {
	var modelUpdate = ModelUpdate{}
	if name != nil {
		modelUpdate.Name = name
	}
	if vendorName != nil {
		modelUpdate.VendorName = vendorName
	}
	if manufacturer != nil {
		modelUpdate.Manufacturer = manufacturer
	}
	if status != nil {
		modelUpdate.Status = status
	}
	if description != nil {
		modelUpdate.Description = description
	}
	if updatedBy != nil {
		modelUpdate.UpdatedBy = updatedBy
	}
	if image != nil {
		modelUpdate.Image = image
	}
	now := time.Now()
	modelUpdate.UpdatedAt = &now
	return &modelUpdate
}
