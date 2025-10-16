package models

import (
	"time"

	"github.com/google/uuid"
)

/*
CREATE TABLE public.firmwares (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	model_id UUID NOT NULL,
	name VARCHAR(255) NOT NULL,
	file_path VARCHAR(255) NULL,
	status public.firmware_status NOT NULL DEFAULT 'ENABLE':::public.firmware_status,
	description VARCHAR(255) NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	CONSTRAINT firmwares_pkey PRIMARY KEY (id ASC),
	CONSTRAINT firmwares_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id),
	UNIQUE INDEX firmwares_name_idx (name ASC) STORING (model_id, file_path, status, description, created_at, updated_at),
	INDEX firmwares_model_id_idx (model_id ASC) STORING (name, file_path, status, description, created_at, updated_at)
);
*/

const USPFirmwareTableName = "firmwares"
const USPFirmwareEntityName = "Firmware"

type Firmware struct {
	Id        *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	ModelId   *uuid.UUID `gorm:"column:model_id;type:uuid;not null" json:"-"`
	CreatedAt *time.Time `json:"-" gorm:"created_at"`
	UpdatedAt *time.Time `json:"-" gorm:"updated_at"`

	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	FilePath    string `json:"file_path,omitempty" gorm:"column:file_path;type:varchar(255);default:null"`
	Status      string `json:"status" gorm:"column:status;type:public.firmware_status;default:'ENABLE':::public.firmware_status"`
	Description string `json:"description,omitempty" gorm:"column:description;type:varchar(255);default:null"`
	UpdatedBy   string `json:"updated_by,omitempty" gorm:"column:updated_by;type:varchar(255);default:null"`

	Model *Model `gorm:"foreignKey:ModelId;references:Id;" json:"model,omitempty"`
}

func (Firmware) TableName() string     { return USPFirmwareTableName }
func (Firmware) GetEntityName() string { return USPFirmwareEntityName }

func (Firmware) GetIdColumnName() string          { return "id" }
func (Firmware) GetModelIdColumnName() string     { return "model_id" }
func (Firmware) GetNameColumnName() string        { return "name" }
func (Firmware) GetFilePathColumnName() string    { return "file_path" }
func (Firmware) GetStatusColumnName() string      { return "status" }
func (Firmware) GetDescriptionColumnName() string { return "description" }
func (Firmware) GetCreatedAtColumnName() string   { return "created_at" }
func (Firmware) GetUpdatedAtColumnName() string   { return "updated_at" }
func (Firmware) GetUpdatedByColumnName() string   { return "updated_by" }

type FirmwareUpdate struct {
	Name        *string    `gorm:"column:name;type:varchar(255);not null" json:"name,omitempty"`
	FilePath    *string    `gorm:"column:file_path;type:varchar(255);default:null" json:"file_path,omitempty"`
	Status      *string    `gorm:"column:status;type:public.firmware_status;default:'ENABLE':::public.firmware_status" json:"status"`
	Description *string    `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`
	UpdatedBy   *string    `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (FirmwareUpdate) TableName() string     { return USPFirmwareTableName }
func (FirmwareUpdate) GetEntityName() string { return USPFirmwareEntityName }

func NewFirmwareUpdate(
	filePath *string,
	status *string,
	description *string,
	updatedBy *string,
) *FirmwareUpdate {
	var firmwareUpdate = FirmwareUpdate{}
	if filePath != nil {
		firmwareUpdate.FilePath = filePath
	}

	if status != nil {
		firmwareUpdate.Status = status
	}

	if description != nil {
		firmwareUpdate.Description = description
	}

	if updatedBy != nil {
		firmwareUpdate.UpdatedBy = updatedBy
	}
	now := time.Now()
	firmwareUpdate.UpdatedAt = &now

	return &firmwareUpdate
}
