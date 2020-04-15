package database

import (
	"errors"
	"fmt"
	"github.com/dormao/go-oss-server/internal/context/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type YamlStore map[string]YamlFile

type YamlFile struct {
	Object string `yaml:"object"`
	Path   string `yaml:"path"`
	Bucket string `yaml:"bucket"`
}

type YamlProvider struct {
	Store  YamlStore
	Bucket string
}

func (prov *YamlProvider) Init() error {
	file, err := ioutil.ReadFile(config.Config.Provider.FilePath)
	if err != nil {
		prov.Store = make(map[string]YamlFile, 0)
		return err
	}
	var dat YamlStore
	err = yaml.Unmarshal(file, &dat)
	if err != nil {
		prov.Store = make(map[string]YamlFile, 0)
		return err
	}

	prov.Store = dat
	return nil
}

func (prov *YamlProvider) Save() error {
	dat, err := yaml.Marshal(prov.Store)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(config.Config.Provider.FilePath, dat, 0755)
}

func (prov *YamlProvider) SetBucket(bucketName string) error {
	return nil
}

func (prov *YamlProvider) GetObject(objectName string) (path string, err error) {
	for _, f := range prov.Store {
		if f.Object == objectName {
			path = f.Path
			break
		}
	}
	if len(path) == 0 {
		err = errors.New(fmt.Sprintf("object %s not found", objectName))
	}
	return path, err
}

func (prov *YamlProvider) RemoveObject(objectName string) error {
	var key string
	for k, f := range prov.Store {
		if f.Object == objectName {
			key = k
			break
		}
	}
	if key != "" {
		delete(prov.Store, key)
		return nil
	} else {
		return errors.New(fmt.Sprintf("object %s not found", key))
	}
}

func (prov YamlProvider) ListObject(prefix string) map[string]string {
	var result = make(map[string]string, 0)
	for k, f := range prov.Store {
		if strings.HasPrefix(f.Object, prefix) {
			result[k] = f.Path
		}
	}
	return result
}

func (prov *YamlProvider) PutObject(objectName, filename string, options Options) error {
	prov.RemoveObject(objectName)
	prov.Store[objectName] = YamlFile{
		Object: objectName,
		Path:   filename,
		Bucket: prov.Bucket,
	}
	return nil
}

func NewYamlDataProvider() *YamlProvider {
	return &YamlProvider{
		Store:  nil,
		Bucket: config.Config.Bucket,
	}
}
