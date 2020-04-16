package database

import (
	"errors"
	"fmt"
	"github.com/dormao/go-oss-server/internal/context/config"
	"github.com/dormao/go-oss-server/internal/context/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

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
	type KeyPairBucket struct {
		models.FieldBucket
		KeyPair *models.FieldKeyPair `gorm:"FOREIGNKEY:KeyPairId" json:"key_pair"`
	}
	var bkName = p.Bucket
	if option.Bucket != nil {
		bkName = *option.Bucket
	}
	var bucket KeyPairBucket
	p.DB.Set("gorm:autoload", true).Joins("INNER JOIN key_pairs ON key_pairs.id = buckets.key_pair_id")
	var notFound = p.DB.First(&bucket, "buckets.name = ?", bkName).RecordNotFound()
	if notFound {
		return errors.New(fmt.Sprintf("bucket (%s) not found", bkName))
	} else if bucket.KeyPair != nil && (bucket.KeyPair.AccessKey != option.AccessKey || bucket.KeyPair.AccessSecret != option.AccessSecret) {
		return errors.New(fmt.Sprintf("access denied by bucket (%s)", bkName))
	}
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
	db.Set("gorm:autoload", true)
	db.Joins(arraysImplode(" ", []string{
		"INNER JOIN buckets ON buckets.id = objects.bucket_id",
	})).Where("buckets.name = ? AND objects.name = ?", bkName, object)
	notFound := db.First(&associationRecord).RecordNotFound()
	if notFound {
		return "", errors.New(fmt.Sprintf("object (%s) not found", object))
	} else if associationRecord.ExpiredAt != nil && associationRecord.ExpiredAt.After(time.Now()) {
		return "", errors.New(fmt.Sprintf("object (%s) expired", object))
	}
	return associationRecord.Path, nil
}

func (p *PostgresProvider) RemoveObject(objectName string) error {
	return p.DB.Delete(&models.FieldObject{}, "objects.name = ?", objectName).Error
}

func (p *PostgresProvider) Save() error {
	return nil
}

func (p *PostgresProvider) ListObject(objectPrefix string) map[string]string {
	// TODO Regex Query
	return map[string]string{}
}

func arraysImplode(glue string, stack []string) (out string) {
	for index, v := range stack {
		if index == len(stack)-1 {
			out += v
		} else {
			out += v + glue
		}
	}
	return out
}
