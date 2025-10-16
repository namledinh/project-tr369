package models

import (
	"time"

	"github.com/google/uuid"
)

/*
CREATE TABLE public.parameters (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	path VARCHAR NOT NULL,
	data_type VARCHAR(32) NOT NULL,
	description STRING NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	status public.status NOT NULL DEFAULT 'ENABLE':::public.status,
	updated_by VARCHAR(255) NULL,
	CONSTRAINT parameters_pkey PRIMARY KEY (id ASC),
	UNIQUE INDEX parameters_path_idx (path ASC) STORING (data_type, description, created_at, updated_at)
);
COMMENT ON COLUMN public.parameters.path IS 'TR369 key, example: Device.DeviceInfo.NetworkProperties';
COMMENT ON COLUMN public.parameters.data_type IS 'The datatype of TR369 key, example: `string`, ref: https://usp-data-models.broadband-forum.org/tr-181-2-19-1-usp.html';
COMMENT ON COLUMN public.parameters.description IS 'just some description for human';

*/

const USPParameterTableName = "parameters"
const USPParameterEntityName = "Parameter"

type Parameter struct {
	Id          *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	Path        string     `gorm:"column:path;type:varchar;not null" json:"path"`
	DataType    string     `gorm:"column:data_type;type:varchar(32);not null" json:"data_type"`
	Description string     `gorm:"column:description;type:string;default:null" json:"description,omitempty"`
	CreatedAt   *time.Time `json:"-" gorm:"created_at"`
	UpdatedAt   *time.Time `json:"-" gorm:"updated_at"`
	Status      string     `gorm:"column:status;type:public.status;default:'ENABLE':::public.status" json:"status"`
	UpdatedBy   string     `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`

	// Metadata, not ORM field
	DefaultValue string `json:"default_value,omitempty" gorm:"-"`
	Required     bool   `json:"required,omitempty" gorm:"-"` // Indicates if this parameter is required
}

func (Parameter) TableName() string                { return USPParameterTableName }
func (Parameter) GetEntityName() string            { return USPParameterEntityName }
func (Parameter) GetIdColumnName() string          { return "id" }
func (Parameter) GetPathColumnName() string        { return "path" }
func (Parameter) GetDataTypeColumnName() string    { return "data_type" }
func (Parameter) GetDescriptionColumnName() string { return "description" }
func (Parameter) GetStatusColumnName() string      { return "status" }
func (Parameter) GetCreatedAtColumnName() string   { return "created_at" }
func (Parameter) GetUpdatedAtColumnName() string   { return "updated_at" }
func (Parameter) GetUpdatedByColumnName() string   { return "updated_by" }

type ParameterUpdate struct {
	Path         *string    `gorm:"column:path;type:varchar;not null" json:"path"`
	DataType     *string    `gorm:"column:data_type;type:varchar(32);not null" json:"data_type"`
	Description  *string    `gorm:"column:description;type:string;default:null" json:"description,omitempty"`
	Status       *string    `gorm:"column:status;type:public.status;default:'ENABLE':::public.status" json:"status"`
	DefaultValue *string    `gorm:"column:default_value;type:string;default:null" json:"default_value,omitempty"`
	Required     *bool      `gorm:"column:required;type:boolean;default:false" json:"required,omitempty"`
	CreatedAt    *time.Time `json:"-" gorm:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy    *string    `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
}

func (ParameterUpdate) TableName() string     { return USPParameterTableName }
func (ParameterUpdate) GetEntityName() string { return USPParameterEntityName }

func NewParameterUpdate(
	path *string, dataType *string, description *string, status *string, updatedBy *string,
) *ParameterUpdate {
	var p = ParameterUpdate{}
	if path != nil {
		p.Path = path
	}
	if dataType != nil {
		p.DataType = dataType
	}
	if description != nil {
		p.Description = description
	}
	if status != nil {
		p.Status = status
	}
	if updatedBy != nil {
		p.UpdatedBy = updatedBy
	}
	p.UpdatedAt = new(time.Time)
	*p.UpdatedAt = time.Now()
	return &p
}
