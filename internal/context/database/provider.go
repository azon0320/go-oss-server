package database

import "time"

type RetrieveOption struct {
	Bucket       *string
	AccessKey    string
	AccessSecret string
}

type SetUpOption struct {
	RetrieveOption
	ExpiredAt *time.Time
}

type Provider interface {
	Init() error
	Save() error
	SetBucket(bucketName string) error
	PutObject(objectName, filename string, options SetUpOption) error
	GetObject(objectName string, option RetrieveOption) (string, error)
	ListObject(objectPrefix string, option RetrieveOption) (interface{}, error)
	RemoveObject(objectName string, option RetrieveOption) error
}
