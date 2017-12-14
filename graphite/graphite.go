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
	"fmt"
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	graphite "github.com/marpaia/graphite-golang"
	log "github.com/sirupsen/logrus"
)

const (
	Name    = "graphite"
	Version = 7
)

var (
	illegal     = "(), /{}"
	replacement = strings.NewReplacer(" ", "_",
		",", ";",
		"(", "[",
		")", "]",
		"/", "|",
		"{", "[",
		"}", "]")
)

type GraphitePublisher struct {
}

func (f *GraphitePublisher) Publish(metrics []plugin.Metric, cfg plugin.Config) error {

	logger := getLogger(cfg)
	logger.Debug("Publishing started")
	var tagsForPrefix []string

	logger.Debug("publishing %v metrics to %v", len(metrics), cfg)
	server, err := cfg.GetString("server")
	if err != nil {
		return err
	}
	port, err := cfg.GetInt("port")
	if err != nil {
		return err
	}
	tagConfigs, err := cfg.GetString("prefix_tags")
	if err == nil {
		tagsForPrefix = strings.Split(tagConfigs, ",")
	}
	pre, err := cfg.GetString("prefix")
	if err != nil {
		pre = ""
	}
	enableTags, err := cfg.GetBool("enable_tags")
	if err != nil {
		return err
	}

	logger.Debug("Attempting to connect to %s:%d", server, port)
	gite, err := graphite.NewGraphiteWithMetricPrefix(server, int(port), pre)
	if err != nil {
		logger.Errorf("Error Connecting to graphite at %s:%d. Error: %v", server, port, err)
		return fmt.Errorf("Error Connecting to graphite at %s:%d. Error: %v", server, port, err)
	}
	defer gite.Disconnect()
	logger.Debug("Connected to %s:%s successfully", server, port)

	giteMetrics := make([]graphite.Metric, len(metrics))
	for i, m := range metrics {
		key := createKey(m, tagsForPrefix, enableTags)

		data := fmt.Sprintf("%v", m.Data)
		logger.Debug("Metric ready to send %s:%s", key, data)

		giteMetrics[i] = graphite.NewMetric(key, data, m.Timestamp.Unix())
	}

	err = gite.SendMetrics(giteMetrics)
	if err != nil {
		logger.Errorf("Unable to send metrics. Error: %s", err)
		return fmt.Errorf("Unable to send metrics. Error: %s", err)
	}
	logger.Debug("Metrics sent to Graphite.")

	return nil
}

func (f *GraphitePublisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	policy.AddNewStringRule([]string{""}, "server", true)
	policy.AddNewIntRule([]string{""}, "port", false, plugin.SetDefaultInt(2003))
	policy.AddNewStringRule([]string{""}, "prefix_tags", false, plugin.SetDefaultString("plugin_running_on"))
	policy.AddNewStringRule([]string{""}, "prefix", false)
	policy.AddNewStringRule([]string{""}, "log-level", false)
	policy.AddNewBoolRule([]string{""}, "enable_tags", false, plugin.SetDefaultBool(false))

	return *policy, nil
}

func getLogger(cfg plugin.Config) *log.Entry {
	logger := log.WithFields(log.Fields{
		"plugin-name":    Name,
		"plugin-version": Version,
		"plugin-type":    "publisher",
	})

	log.SetLevel(log.WarnLevel)

	levelValue, err := cfg.GetString("log-level")
	if err == nil {
		if level, err := log.ParseLevel(strings.ToLower(levelValue)); err == nil {
			log.SetLevel(level)
		} else {
			log.WithFields(log.Fields{
				"value":             strings.ToLower(levelValue),
				"acceptable values": "warn, error, debug, info",
			}).Warn("Invalid config value")
		}
	}
	return logger
}

func replaceDynamicElements(m plugin.Metric) ([]string, map[string]string) {
	tags := map[string]string{}
	ns := m.Namespace.Strings()

	isDynamic, indexes := m.Namespace.IsDynamic()
	if isDynamic {
		for i, j := range indexes {
			// The second return value from IsDynamic(), in this case `indexes`, is the index of
			// the dynamic element in the unmodified namespace. However, here we're deleting
			// elements, which is problematic when the number of dynamic elements in a namespace is
			// greater than 1. Therefore, we subtract i (the loop iteration) from j
			// (the original index) to compensate.
			//
			// Remove "data" from the namespace and create a tag for it
			ns = append(ns[:j-i], ns[j-i+1:]...)
			tags[m.Namespace[j].Name] = m.Namespace[j].Value
		}
	}
	return ns, tags
}

func createKey(m plugin.Metric, tagsForPrefix []string, enableTags bool) string {
	var ns []string
	var tags map[string]string
	if enableTags {
		ns, tags = replaceDynamicElements(m)
		// Process the tags for this metric
		for k, v := range m.Tags {
			// Convert the standard tag describing where the plugin is running to "host"
			if k == "plugin_running_on" {
				// Unless the "host" tag is already being used
				if _, ok := m.Tags["host"]; !ok {
					k = "host"
				}
			}
			tags[k] = v
		}
	} else {
		ns = m.Namespace.Strings()
	}
	key := strings.Join(ns, ".")

	if strings.ContainsAny(key, illegal) {
		key := replacement.Replace(key)
		log.Info("Metric after replacement is %s", key)
	}

	if enableTags {
		for k, v := range tags {
			key = key + ";" + k + "=" + v
		}
	}

	for _, tag := range tagsForPrefix {
		nextTag, ok := m.Tags[tag]
		if ok {
			key = nextTag + "." + key
		}
	}

	return key
}
