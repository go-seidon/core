package local

import (
	"fmt"
	"strings"
)

type LocalConfig struct {
	StorageDir string
}

type LocalStorageOption interface {
	Apply(c *LocalConfig) error
}

type withNormalStorageDir struct {
	storageDir string
}

func (o *withNormalStorageDir) Apply(c *LocalConfig) error {
	if o.storageDir == "" {
		return fmt.Errorf("invalid storage directory")
	}
	sDir := strings.ToLower(o.storageDir)
	sDir = strings.TrimSuffix(sDir, "/")
	c.StorageDir = sDir
	return nil
}

func WithNormalStorageDir(storageDir string) LocalStorageOption {
	return &withNormalStorageDir{
		storageDir: storageDir,
	}
}
