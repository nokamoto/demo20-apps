package authorizer

import (
	"io/ioutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/nokamoto/demo20-apis/cloud/api"
)

// ConfigLoader loads AuthConfig.
type ConfigLoader struct{}

// Read loads AuthConfig from the file.
func (c *ConfigLoader) Read(filename string) (*api.AuthzConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg *api.AuthzConfig
	return cfg, jsonpb.UnmarshalString(string(bytes), cfg)
}
