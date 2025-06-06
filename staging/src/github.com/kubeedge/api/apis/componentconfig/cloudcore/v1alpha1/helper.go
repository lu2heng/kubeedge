/*
Copyright 2019 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"encoding/json"
	"errors"
	"os"

	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

type IptablesMgrMode string

const (
	InternalMode IptablesMgrMode = "internal"
	ExternalMode IptablesMgrMode = "external"
)

// Parse reads config file and converts YAML to CloudCoreConfig
func (c *CloudCoreConfig) Parse(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		klog.Errorf("Failed to read configfile %s: %v", filename, err)
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		klog.Errorf("Failed to unmarshal configfile %s: %v", filename, err)
		return err
	}
	return nil
}

// WriteTo converts CloudCoreConfig to yaml and write it to the file
func (c *CloudCoreConfig) WriteTo(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	fileInfo, err := os.Stat(filename)
	if err != nil {
		// If the file does not exist, the default permissions are 0644.
		return os.WriteFile(filename, data, 0644)
	}
	// Write the file using the original file permissions.
	return os.WriteFile(filename, data, fileInfo.Mode())
}

func (in *IptablesManager) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Define a secondary type so that we don't end up with a recursive call to json.Unmarshal
	type IM IptablesManager
	var out = (*IM)(in)
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	// Validate the valid enum values
	switch in.Mode {
	case InternalMode, ExternalMode:
		return nil
	default:
		in.Mode = ""
		return errors.New("invalid value for iptablesmgr mode")
	}
}
