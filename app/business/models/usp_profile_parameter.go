package models

import (
	"time"

	"github.com/google/uuid"
)

/*
CREATE TABLE public.profile_parameters (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	profile_id UUID NOT NULL,
	parameter_id UUID NOT NULL,
	default_value STRING NULL,
	required BOOL NULL,
	created_at TIMESTAMP NOT NULL DEFAULT 'now()':::STRING::TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT 'now()':::STRING::TIMESTAMP,
	CONSTRAINT profile_parameters_pkey PRIMARY KEY (id ASC),
	CONSTRAINT profile_parameters_profile_id_fkey FOREIGN KEY (profile_id) REFERENCES public.profiles(id),
	CONSTRAINT profile_parameters_parameter_id_fkey FOREIGN KEY (parameter_id) REFERENCES public.parameters(id),
	UNIQUE INDEX profile_parameters_profile_id_parameter_id_idx (profile_id ASC, parameter_id ASC) STORING (default_value, created_at, updated_at),
	INDEX profile_parameters_profile_id_idx (profile_id ASC) STORING (parameter_id, default_value, created_at, updated_at),
	INDEX profile_parameters_parameter_id_idx (parameter_id ASC) STORING (profile_id, default_value, created_at, updated_at)
);
COMMENT ON COLUMN public.profile_parameters.default_value IS 'just default value if theres any empty return';
*/

const USPProfileParameterTableName = "profile_parameters"
const USPProfileParameterEntityName = "ProfileParameter"

type ProfileParameter struct {
	Id           *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	ProfileId    *uuid.UUID `gorm:"column:profile_id;type:uuid;not null" json:"-"`
	ParameterId  *uuid.UUID `gorm:"column:parameter_id;type:uuid;not null" json:"-"`
	DefaultValue string     `gorm:"column:default_value;type:string;default:null" json:"default_value,omitempty"`
	Required     bool       `gorm:"column:required;type:bool;default:false" json:"required,omitempty"`
	CreatedAt    *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy    string     `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`

	// Foreign keys
	// ProfileId  -> Profile.Id
	// ParameterId -> Parameter.Id

	// Relationships - Preload
	Profile   *Profile   `gorm:"foreignKey:ProfileId;references:Id" json:"profile"`
	Parameter *Parameter `gorm:"foreignKey:ParameterId;references:Id" json:"parameter"`
}

func (ProfileParameter) TableName() string                 { return USPProfileParameterTableName }
func (ProfileParameter) GetEntityName() string             { return USPProfileParameterEntityName }
func (ProfileParameter) GetIdColumnName() string           { return "id" }
func (ProfileParameter) GetProfileIdColumnName() string    { return "profile_id" }
func (ProfileParameter) GetParameterIdColumnName() string  { return "parameter_id" }
func (ProfileParameter) GetDefaultValueColumnName() string { return "default_value" }
func (ProfileParameter) GetRequiredColumnName() string     { return "required" }
func (ProfileParameter) GetCreatedAtColumnName() string    { return "created_at" }
func (ProfileParameter) GetUpdatedAtColumnName() string    { return "updated_at" }
func (ProfileParameter) GetUpdatedByColumnName() string    { return "updated_by" }

// Preload helper
func (ProfileParameter) GetProfilePreload() string   { return "Profile" }
func (ProfileParameter) GetParameterPreload() string { return "Parameter" }

type ProfileParameterUpdate struct {
	DefaultValue *string    `gorm:"column:default_value;type:string;default:null" json:"default_value,omitempty"`
	Required     *bool      `gorm:"column:required;type:bool;default:false" json:"required,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy    *string    `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
}

func (ProfileParameterUpdate) TableName() string     { return USPProfileParameterTableName }
func (ProfileParameterUpdate) GetEntityName() string { return USPProfileParameterEntityName }

func NewProfileParameterUpdate(
	defaultValue *string,
	required *bool,
) *ProfileParameterUpdate {
	var pp = ProfileParameterUpdate{}
	if defaultValue != nil {
		pp.DefaultValue = defaultValue
	}
	if required != nil {
		pp.Required = required
	}

	pp.UpdatedAt = new(time.Time)
	*pp.UpdatedAt = time.Now()
	return &pp
}
