package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

/*
CREATE TABLE public.profiles (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	name VARCHAR(64) NOT NULL,
	msg_type INT8 NOT NULL DEFAULT 1:::INT8,
	return_commands BOOL NULL,
	return_events BOOL NULL,
	return_params BOOL NULL,
	return_unique_key_sets BOOL NULL,
	allow_partial BOOL NULL,
	send_resp BOOL NULL,
	first_level_only BOOL NULL,
	max_depth INT2 NULL,
	tags STRING[] NULL,
	created_at TIMESTAMP NOT NULL DEFAULT 'now()':::STRING::TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT 'now()':::STRING::TIMESTAMP,
	status public.status NOT NULL DEFAULT 'ENABLE':::public.status,
	updated_by VARCHAR(255) NULL,
	CONSTRAINT profiles_pkey PRIMARY KEY (id ASC),
    UNIQUE INDEX profiles_name_idx (name ASC) STORING (
        msg_type,
        return_commands,
        return_events,
        return_params,
        return_unique_key_sets,
        allow_partial,
        send_resp,
        first_level_only,
        max_depth,
        tags,
        created_at,updated_at,
        ),
    INDEX profiles_status_idx (status ASC)
);
COMMENT ON COLUMN public.profiles.name IS 'profile name can be the route path with replace `/` to `_`, example: system_resources';
COMMENT ON COLUMN public.profiles.tags IS 'some tags for categories api, can be: [<consumer_name>, <public or private>';
*/

const USPProfileTableName = "profiles"
const USPProfileEntityName = "Profile"

type Profile struct {
	Id        *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`

	Name                string         `gorm:"column:name" json:"name"`
	MsgType             int            `gorm:"column:msg_type" json:"msg_type"`
	ReturnCommands      bool           `gorm:"column:return_commands" json:"return_commands"`
	ReturnEvents        bool           `gorm:"column:return_events" json:"return_events"`
	ReturnParams        bool           `gorm:"column:return_params" json:"return_params"`
	ReturnUniqueKeySets bool           `gorm:"column:return_unique_key_sets" json:"return_unique_key_sets"`
	AllowPartial        bool           `gorm:"column:allow_partial;type:bool;default:null" json:"allow_partial"`
	SendResp            bool           `gorm:"column:send_resp;type:bool;default:null" json:"send_resp"`
	FirstLevelOnly      bool           `gorm:"column:first_level_only;type:bool;default:null" json:"first_level_only"`
	MaxDepth            int            `gorm:"column:max_depth;type:int2;default:null" json:"max_depth"`
	Tags                pq.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	Status              string         `gorm:"column:status" json:"status"`
	UpdatedBy           string         `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
	Description         string         `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`

	// Relationships
	ProfileParameters []*ProfileParameter `gorm:"foreignKey:ProfileId;references:Id" json:"profile_parameters,omitempty"`
}

func (Profile) TableName() string     { return USPProfileTableName }
func (Profile) GetEntityName() string { return USPProfileEntityName }

func (Profile) GetProfileParameterPreload() string          { return "ProfileParameters" }
func (Profile) GetProfileParameterParameterPreload() string { return "ProfileParameters.Parameter" }

func (Profile) GetIdColumnName() string                  { return "id" }
func (Profile) GetProfileNameColumnName() string         { return "name" }
func (Profile) GetMsgTypeColumnName() string             { return "msg_type" }
func (Profile) GetReturnCommandsColumnName() string      { return "return_commands" }
func (Profile) GetReturnEventsColumnName() string        { return "return_events" }
func (Profile) GetReturnParamsColumnName() string        { return "return_params" }
func (Profile) GetReturnUniqueKeySetsColumnName() string { return "return_unique_key_sets" }
func (Profile) GetAllowPartialColumnName() string        { return "allow_partial" }
func (Profile) GetSendRespColumnName() string            { return "send_resp" }
func (Profile) GetFirstLevelOnlyColumnName() string      { return "first_level_only" }
func (Profile) GetMaxDepthColumnName() string            { return "max_depth" }
func (Profile) GetTagsColumnName() string                { return "tags" }
func (Profile) GetStatusColumnName() string              { return "status" }
func (Profile) GetCreatedAtColumnName() string           { return "created_at" }
func (Profile) GetUpdatedAtColumnName() string           { return "updated_at" }
func (Profile) GetUpdatedByColumnName() string           { return "updated_by" }
func (Profile) GetDescriptionColumnName() string         { return "description" }

type ProfileUpdate struct {
	Name                *string        `gorm:"column:name" json:"name"`
	MsgType             *int           `gorm:"column:msg_type" json:"msg_type"`
	ReturnCommands      *bool          `gorm:"column:return_commands" json:"return_commands"`
	ReturnEvents        *bool          `gorm:"column:return_events" json:"return_events"`
	ReturnParams        *bool          `gorm:"column:return_params" json:"return_params"`
	ReturnUniqueKeySets *bool          `gorm:"column:return_unique_key_sets" json:"return_unique_key_sets"`
	AllowPartial        *bool          `gorm:"column:allow_partial;type:bool;default:null" json:"allow_partial"`
	SendResp            *bool          `gorm:"column:send_resp;type:bool;default:null" json:"send_resp"`
	FirstLevelOnly      *bool          `gorm:"column:first_level_only;type:bool;default:null" json:"first_level_only"`
	MaxDepth            *int           `gorm:"column:max_depth;type:int2;default:null" json:"max_depth"`
	Tags                pq.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	Status              *string        `gorm:"column:status" json:"status"`
	UpdatedAt           *time.Time     `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy           *string        `gorm:"column:updated_by;type:varchar(255);default:null" json:"updated_by,omitempty"`
	Description         *string        `gorm:"column:description;type:varchar(255);default:null" json:"description,omitempty"`
	Parameters          []ParameterRef `json:"parameters,omitempty" gorm:"-"`
}

type ParameterRef struct {
	Id string `json:"id"`
}

func (ProfileUpdate) TableName() string     { return USPProfileTableName }
func (ProfileUpdate) GetEntityName() string { return USPProfileEntityName }

func NewProfileUpdate() *ProfileUpdate {
	var profileUpdate = ProfileUpdate{
		Name:                nil,
		MsgType:             nil,
		ReturnCommands:      nil,
		ReturnEvents:        nil,
		ReturnParams:        nil,
		ReturnUniqueKeySets: nil,
		AllowPartial:        nil,
		SendResp:            nil,
		FirstLevelOnly:      nil,
		MaxDepth:            nil,
		Tags:                nil,
		Status:              nil,
		UpdatedBy:           nil,
		Description:         nil,
	}

	now := time.Now()
	profileUpdate.UpdatedAt = &now

	return &profileUpdate
}
