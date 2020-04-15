package database

import "time"

type Options struct {
	ExpiredAt *time.Duration
}

type Provider interface {
	Init() error
	Save() error
	SetBucket(bucketName string) error
	PutObject(objectName, filename string, options Options) error
	GetObject(objectName string) (string, error)
	ListObject(objectPrefix string) map[string]string
	RemoveObject(objectName string) error
}
