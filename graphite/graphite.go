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

package graphite

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/marpaia/graphite-golang"
)

const (
	name       = "graphite"
	version    = 1
	pluginType = plugin.PublisherPluginType
)

type graphitePublisher struct {
}

func NewGraphitePublisher() *graphitePublisher {
	return &graphitePublisher{}
}

func (f *graphitePublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	logger.Println("Publishing started")
	var metrics []plugin.PluginMetricType

	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return fmt.Errorf("Unknown content type '%s'", contentType)
	}

	logger.Printf("publishing %v metrics to %v", len(metrics), config)

	server := config["graphite-server"].(ctypes.ConfigValueStr).Value
	port := config["graphite-port"].(ctypes.ConfigValueInt).Value

	logger.Printf("Attempting to connect to %s:%d", server, port)
	gite, err := graphite.NewGraphite(server, port)
	handleErr(err)
	for _, m := range metrics {
		key := strings.Join(m.Namespace(), ".")
		data := fmt.Sprintf("%v", m.Data())
		gite.SimpleSend(key, data)
		logger.Printf("Send %s, %d", key, data)
	}
	return nil
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

func (f *graphitePublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	r1, err := cpolicy.NewStringRule("graphite-server", true)
	handleErr(err)
	r1.Description = "Address of graphite server"
	config.Add(r1)

	r2, err := cpolicy.NewIntegerRule("graphite-port", true)
	handleErr(err)
	r2.Description = "Port to connect on"
	config.Add(r2)

	cp.Add([]string{""}, config)
	return cp, nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
