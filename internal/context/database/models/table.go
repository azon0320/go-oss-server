package models

const (
	TableNameBucket  = "buckets"
	TableNameKeyPair = "key_pairs"
	TableNameObject  = "objects"
)

type Bucket struct{}

func (Bucket) TableName() string { return "buckets" }

type KeyPair struct{}

func (KeyPair) TableName() string { return "key_pairs" }

type Object struct{}

func (Object) TableName() string { return "objects" }
