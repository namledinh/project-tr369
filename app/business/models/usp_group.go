package models

import (
	"time"

	"github.com/google/uuid"
)

/*
CREATE TABLE public.groups (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	model_id UUID NOT NULL,
	firmware_id UUID NULL,
	name VARCHAR NOT NULL,
	status public.group_status NOT NULL DEFAULT 'ENABLE':::public.group_status,
	description VARCHAR(255) NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	CONSTRAINT groups_pkey PRIMARY KEY (id ASC),
	CONSTRAINT groups_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id),
	CONSTRAINT groups_firmware_id_fkey FOREIGN KEY (firmware_id) REFERENCES public.firmwares(id),
	UNIQUE INDEX groups_name_idx (name ASC) STORING (model_id, firmware_id, status, description, created_at, updated_at),
	INDEX groups_model_id_idx (model_id ASC) STORING (firmware_id, name, status, description, created_at, updated_at),
	INDEX groups_firmware_id_idx (firmware_id ASC) STORING (model_id, name, status, description, created_at, updated_at)
);
*/

const USPGroupTableName = "groups"
const USPGroupEntityName = "Group"

type Group struct {
	Id             *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	ModelId        *uuid.UUID `gorm:"column:model_id;type:uuid;not null" json:"-"`
	FirmwareId     *uuid.UUID `gorm:"column:firmware_id;type:uuid;default:null" json:"-"`
	CreatedAt      *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"column:updated_at"`
	Name           string     `gorm:"column:name;type:varchar;not null" json:"name"`
	Status         string     `gorm:"column:status;type:public.group_status;default:'ENABLE':::public.group_status" json:"status"`
	Description    string     `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`
	UpdatedBy      string     `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
	DownloadPeriod string     `gorm:"column:download_period;type:varchar(11);default:null" json:"download_period,omitempty"`

	Model    *Model    `gorm:"foreignKey:ModelId;references:Id;" json:"model,omitempty"`
	Firmware *Firmware `gorm:"foreignKey:FirmwareId;references:Id;" json:"firmware,omitempty"`
}

func (Group) TableName() string     { return USPGroupTableName }
func (Group) GetEntityName() string { return USPGroupEntityName }

func (Group) GetIdColumnName() string             { return "id" }
func (Group) GetModelIdColumnName() string        { return "model_id" }
func (Group) GetFirmwareIdColumnName() string     { return "firmware_id" }
func (Group) GetNameColumnName() string           { return "name" }
func (Group) GetStatusColumnName() string         { return "status" }
func (Group) GetDescriptionColumnName() string    { return "description" }
func (Group) GetCreatedAtColumnName() string      { return "created_at" }
func (Group) GetUpdatedAtColumnName() string      { return "updated_at" }
func (Group) GetUpdatedByColumnName() string      { return "updated_by" }
func (Group) GetDownloadPeriodColumnName() string { return "download_period" }

type GroupUpdate struct {
	FirmwareId     *uuid.UUID `gorm:"column:firmware_id;type:uuid;default:null" json:"-"`
	Name           *string    `gorm:"column:name;type:varchar;not null" json:"name"`
	Status         *string    `gorm:"column:status;type:public.group_status;default:'ENABLE':::public.group_status" json:"status"`
	Description    *string    `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`
	UpdatedBy      *string    `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
	DownloadPeriod *string    `gorm:"column:download_period;type:varchar(11);default:null" json:"download_period,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (GroupUpdate) TableName() string     { return USPGroupTableName }
func (GroupUpdate) GetEntityName() string { return USPGroupEntityName }

func NewGroupUpdate(
	firmwareId *uuid.UUID,
	name *string,
	status *string,
	description *string,
	updatedBy *string,
	DownloadPeriod *string,
) *GroupUpdate {
	var groupUpdate = GroupUpdate{}
	if firmwareId != nil {
		groupUpdate.FirmwareId = firmwareId
	}
	if name != nil {
		groupUpdate.Name = name
	}

	if status != nil {
		groupUpdate.Status = status
	}

	if description != nil {
		groupUpdate.Description = description
	}

	if updatedBy != nil {
		groupUpdate.UpdatedBy = updatedBy
	}
	if DownloadPeriod != nil {
		groupUpdate.DownloadPeriod = DownloadPeriod
	}
	now := time.Now()
	groupUpdate.UpdatedAt = &now

	return &groupUpdate
}
