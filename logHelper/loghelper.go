/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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

package loghelper

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

func AddLogging(cp *cpolicy.ConfigPolicyNode) (*cpolicy.ConfigPolicyNode, error) {
	logLevel, err := cpolicy.NewStringRule("log-level", false)
	if err != nil {
		return nil, err
	}
	logLevel.Description = "Log level for plugin. Accepted values: info, warn, error, debug"
	cp.Add(logLevel)
	return cp, nil
}

func GetLogger(config map[string]ctypes.ConfigValue, pluginInfo *plugin.PluginMeta) *log.Entry {
	logger := log.WithFields(log.Fields{
		"plugin-name":    pluginInfo.Name,
		"plugin-version": pluginInfo.Version,
		"plugin-type":    pluginInfo.Type.String(),
	})

	log.SetLevel(log.WarnLevel)

	if debug, ok := config["debug"]; ok {
		if v, ok := debug.(ctypes.ConfigValueBool); ok {
			if v.Value {
				log.SetLevel(log.DebugLevel)
				return logger
			}
		} else {
			logger.WithFields(log.Fields{
				"field":         "debug",
				"type":          v,
				"expected type": "ctyps.ConfigValueBool",
			}).Error("Invalid config type")
		}
	}

	if loglevel, ok := config["log-level"]; ok {
		if v, ok := loglevel.(ctypes.ConfigValueStr); ok {
			if level, err := log.ParseLevel(strings.ToLower(v.Value)); err == nil {
				log.SetLevel(level)
			} else {
				log.WithFields(log.Fields{
					"value":             strings.ToLower(v.Value),
					"acceptable values": "warn, error, debug, info",
				}).Warn("Invalid config value")
			}
		} else {
			logger.WithFields(log.Fields{
				"field":         "log-level",
				"type":          v,
				"expected type": "ctypes.ConfigValueStr",
			}).Error("Invalid config type")
		}
	}

	return logger
}
