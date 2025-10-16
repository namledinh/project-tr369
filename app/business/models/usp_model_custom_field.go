package models

import (
	"time"

	"github.com/google/uuid"
)

/*
CREATE TABLE public.model_custom_fields (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	model_id UUID NOT NULL,
	key VARCHAR(255) NOT NULL,
	value VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	CONSTRAINT model_custom_fields_pkey PRIMARY KEY (id ASC),
	CONSTRAINT model_custom_fields_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id),
	INDEX model_custom_fields_model_id_idx (model_id ASC) STORING (key, value, created_at, updated_at),
    UNIQUE INDEX model_custom_fields_model_id_key_idx (model_id ASC, key ASC) STORING (value, created_at, updated_at)
);
COMMENT ON COLUMN public.model_custom_fields.key IS 'Data Model';
COMMENT ON COLUMN public.model_custom_fields.value IS 'TR181';
*/

const USPModelCustomFieldTableName = "model_custom_fields"
const USPModelCustomFieldEntityName = "ModelCustomField"

type ModelCustomField struct {
	Id        *uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"-"`
	ModelId   *uuid.UUID `gorm:"column:model_id;type:uuid;not null" json:"-"`
	Key       string     `gorm:"column:key;type:varchar(255);not null" json:"key"`
	Value     string     `gorm:"column:value;type:varchar(255);not null" json:"value"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (ModelCustomField) TableName() string     { return USPModelCustomFieldTableName }
func (ModelCustomField) GetEntityName() string { return USPModelCustomFieldEntityName }

func (ModelCustomField) GetIdColumnName() string        { return "id" }
func (ModelCustomField) GetModelIdColumnName() string   { return "model_id" }
func (ModelCustomField) GetKeyColumnName() string       { return "key" }
func (ModelCustomField) GetValueColumnName() string     { return "value" }
func (ModelCustomField) GetCreatedAtColumnName() string { return "created_at" }
func (ModelCustomField) GetUpdatedAtColumnName() string { return "updated_at" }
