package models

import (
	"time"

	"github.com/google/uuid"
)

/*

CREATE TABLE public.devices (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	mac_address VARCHAR(12) NOT NULL,
	endpoint_id VARCHAR NOT NULL,
	model_id UUID NOT NULL,
	group_id UUID NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	CONSTRAINT devices_pkey PRIMARY KEY (id ASC),
	CONSTRAINT devices_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id),
	CONSTRAINT devices_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(id),
	UNIQUE INDEX devices_mac_address_idx (mac_address ASC) STORING (endpoint_id, model_id, group_id, created_at, updated_at),
	INDEX devices_model_id_idx (model_id ASC) STORING (mac_address, endpoint_id, group_id, created_at, updated_at),
	INDEX devices_group_id_idx (group_id ASC) STORING (mac_address, endpoint_id, model_id, created_at, updated_at),
    INDEX devices_endpoint_id_idx (endpoint_id ASC) STORING (mac_address, model_id, group_id, created_at, updated_at)
);
COMMENT ON COLUMN public.devices.mac_address IS 'WAN MAC Address of device, 12 chars, not contains `:`';
COMMENT ON COLUMN public.devices.endpoint_id IS 'endpoint of device in TR369 System, example: `os::4485DA-4485DA68A1E7`';
COMMENT ON COLUMN public.devices.model_id IS 'ref to models table, define device model.';
*/

const USPDeviceTableName = "devices"
const USPDeviceEntityName = "Device"

type Device struct {
	Id          *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	MacAddress  string     `gorm:"column:mac_address;type:varchar(12);not null" json:"mac_address"`
	EndpointId  string     `gorm:"column:endpoint_id;type:varchar;not null" json:"endpoint_id"`
	Status      string     `gorm:"column:status;type:varchar;default:'ENABLE'" json:"status"`
	UpdatedBy   string     `gorm:"column:updated_by;type:varchar;not null" json:"updated_by"`
	ModelId     *uuid.UUID `gorm:"column:model_id;type:uuid;not null" json:"-"`
	GroupId     *uuid.UUID `gorm:"column:group_id;type:uuid;default:null" json:"-"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
	Description string     `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`

	Model *Model `gorm:"foreignKey:ModelId;references:Id;" json:"model,omitempty"`
	Group *Group `gorm:"foreignKey:GroupId;references:Id;" json:"group,omitempty"`
}

func (Device) TableName() string     { return USPDeviceTableName }
func (Device) GetEntityName() string { return USPDeviceEntityName }

func (Device) GetIdColumnName() string          { return "id" }
func (Device) GetMacAddressColumnName() string  { return "mac_address" }
func (Device) GetEndpointIdColumnName() string  { return "endpoint_id" }
func (Device) GetModelIdColumnName() string     { return "model_id" }
func (Device) GetGroupIdColumnName() string     { return "group_id" }
func (Device) GetCreatedAtColumnName() string   { return "created_at" }
func (Device) GetUpdatedAtColumnName() string   { return "updated_at" }
func (Device) GetUpdatedByColumnName() string   { return "updated_by" }
func (Device) GetStatusColumnName() string      { return "status" }
func (Device) GetDescriptionColumnName() string { return "description" }

type DeviceUpdate struct {
	GroupId     *uuid.UUID `gorm:"column:group_id;type:uuid;default:null" json:"group_id,omitempty"`
	Status      *string    `gorm:"column:status;type:varchar;default:'ENABLE'" json:"status,omitempty"`
	UpdatedBy   *string    `gorm:"column:updated_by;type:varchar;default:null" json:"updated_by,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
	Description *string    `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`
}

func (DeviceUpdate) TableName() string     { return USPDeviceTableName }
func (DeviceUpdate) GetEntityName() string { return USPDeviceEntityName }

func NewDeviceUpdate(
	groupId *uuid.UUID,
	status *string,
	updatedBy *string,
	description *string,
) *DeviceUpdate {
	now := time.Now()
	return &DeviceUpdate{
		GroupId:     groupId,
		Status:      status,
		UpdatedBy:   updatedBy,
		Description: description,
		UpdatedAt:   &now,
	}
}
