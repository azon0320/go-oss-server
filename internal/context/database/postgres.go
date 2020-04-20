package database

import (
	"errors"
	"fmt"
	"github.com/dormao/go-oss-server/internal/context/config"
	"github.com/dormao/go-oss-server/internal/context/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
	"time"
)

type KeyPairBucket struct {
	models.FieldBucket
	KeyPair *models.FieldKeyPair `gorm:"FOREIGNKEY:KeyPairId"`
}

type PostgresProvider struct {
	Bucket      string
	ResourceURI string
	DB          *gorm.DB
}

func (p *PostgresProvider) Init() error {
	db, err := gorm.Open("postgres", p.ResourceURI)
	if err != nil {
		return err
	}
	p.DB = db
	// auto migrate
	if config.Config.Provider.DbAutoMigrate {
		db.AutoMigrate(&models.FieldKeyPair{})
		db.AutoMigrate(&models.FieldBucket{})
		db.AutoMigrate(&models.FieldObject{})
	}
	// Init operation
	var keypair models.FieldKeyPair
	notFound := db.First(&keypair, "name = ?", "root").RecordNotFound()
	if notFound {
		keypair.Name = "root"
		keypair.AccessKey = config.Config.AccessKey
		keypair.AccessSecret = config.Config.AccessSecret
		db.Save(&keypair)
	}
	var bucket models.FieldBucket
	notFound = db.First(&bucket, "name = ?", config.Config.Bucket).RecordNotFound()
	if notFound {
		bucket.Name = config.Config.Bucket
		bucket.KeyPairId = keypair.ID
		db.Save(&bucket)
	}
	p.Bucket = bucket.Name
	return nil
}

func (p *PostgresProvider) SetBucket(bucket string) error {
	notFound := p.DB.First(&models.FieldBucket{}, "name = ?", bucket).RecordNotFound()
	if notFound {
		return errors.New(fmt.Sprintf("bucket (%s) not found in database", bucket))
	}
	p.Bucket = bucket
	return nil
}

func (p *PostgresProvider) PutObject(object, filename string, option SetUpOption) error {
	bucket, err := p.validateAuth(option.RetrieveOption)
	if err != nil { return err }
	var obj models.FieldObject
	obj.Name = object
	obj.Path = filename
	obj.BucketId = bucket.ID
	if option.ExpiredAt != nil {
		obj.ExpiredAt = option.ExpiredAt
	}
	return p.DB.FirstOrCreate(&obj, "objects.name = ?", object).Error
}

func (p *PostgresProvider) GetObject(object string, option RetrieveOption) (string, error) {
	type AssociationObject struct {
		models.FieldObject
		Bucket models.FieldBucket `gorm:"FOREIGNKEY:BucketId" json:"bucket"`
	}
	var bkName = p.Bucket
	if option.Bucket != nil {
		bkName = *option.Bucket
	}
	var db = p.DB
	var associationRecord AssociationObject
	notFound := db.
		Set("gorm:auto_preload", true).
		Joins("INNER JOIN buckets ON buckets.id = objects.bucket_id").
		Where("buckets.name = ? AND objects.name = ?", bkName, object).
		First(&associationRecord).RecordNotFound()
	if notFound {
		return "", errors.New(fmt.Sprintf("object (%s) not found", object))
	} else if associationRecord.ExpiredAt != nil && associationRecord.ExpiredAt.After(time.Now()) {
		return "", errors.New(fmt.Sprintf("object (%s) expired", object))
	}
	return associationRecord.Path, nil
}

func (p *PostgresProvider) RemoveObject(objectName string, option RetrieveOption) error {
	buck, err := p.validateAuth(option)
	if err != nil { return err }
	return p.DB.Delete(&models.FieldObject{}, "objects.name = ? and objects.bucket_id = ?", objectName, buck.ID).Error
}

func (p *PostgresProvider) Save() error {
	return nil
}

func (p *PostgresProvider) ListObject(objectPrefix string, option RetrieveOption) (interface{}, error) {
	buck, err := p.validateAuth(option)
	if err != nil { return map[string]string{} , err}
	var db = p.DB.Joins("INNER JOIN buckets ON buckets.id = objects.bucket_id")
	var objectes []models.FieldObject
	db.Where(`buckets.name = ? AND objects.name ~ ?`, buck.Name, fmt.Sprintf("^%s/", strings.Trim(objectPrefix, "/"))).Find(&objectes)
	var result = make([]string, 0)
	for _, v := range objectes {
		result = append(result, v.Name)
	}
	return result, nil
}

func (p *PostgresProvider) validateAuth(option RetrieveOption) (*KeyPairBucket, error){
	var bkName = p.Bucket
	if option.Bucket != nil {
		bkName = *option.Bucket
	}
	var bucket KeyPairBucket
	var db = p.DB.Set("gorm:auto_preload", true)
	notFound := db.First(&bucket, "buckets.name = ?", bkName).RecordNotFound()
	if notFound {
		return nil, errors.New(fmt.Sprintf("bucket (%s) not found", bkName))
	} else if bucket.KeyPair == nil {
		return nil, errors.New(fmt.Sprintf("access denied by bucket (%s): unbound access keypair", bkName))
	} else if bucket.KeyPair.AccessKey != option.AccessKey || bucket.KeyPair.AccessSecret != option.AccessSecret {
		return nil, errors.New(fmt.Sprintf("access denied by bucket (%s)", bkName))
	}
	if option.AccessKey != bucket.KeyPair.AccessKey || option.AccessSecret != bucket.KeyPair.AccessSecret {
		return nil, errors.New(fmt.Sprintf("access denied by bucket (%s)", p.Bucket))
	}else{
		return &bucket, nil
	}
}
