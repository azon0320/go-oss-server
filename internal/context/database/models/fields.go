package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FieldBucket struct {
	Bucket
	gorm.Model
	Name      string `gorm:"COLUMN:name" json:"name"`
	KeyPairId uint   `gorm:"COLUMN:key_pair_id" json:"key_pair_id"`
}

type FieldKeyPair struct {
	KeyPair
	gorm.Model
	Name         string `gorm:"COLUMN:name" json:"name"`
	AccessKey    string `gorm:"COLUMN:accesskey" json:"access_key"`
	AccessSecret string `gorm:"COLUMN:secret" json:"secret"`
}

type FieldObject struct {
	Object
	gorm.Model
	BucketId  uint       `gorm:"COLUMN:bucket_id" json:"bucket_id"`
	Name      string     `gorm:"UNIQUE;COLUMN:name" json:"name"`
	Path      string     `gorm:"COLUMN:path" json:"path"`
	ExpiredAt *time.Time `gorm:"COLUMN:expired_at" json:"expired_at"`
}
