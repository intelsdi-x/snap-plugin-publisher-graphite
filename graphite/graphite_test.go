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
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphitePlugin(t *testing.T) {
	Convey("Meta should return Metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.PublisherPluginType)
	})

	Convey("Create Graphite publisher", t, func() {
		gp := NewGraphitePublisher()
		Convey("So publisher should not be nil", func() {
			So(gp, ShouldNotBeNil)
		})

		Convey("Publisher should be of type graphitePublisher", func() {
			So(gp, ShouldHaveSameTypeAs, &graphitePublisher{})
		})

		configPolicy, err := gp.GetConfigPolicy()
		Convey("Should return a config policy", func() {

			Convey("configPolicy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)

				Convey("and retrieving config policy should not error", func() {
					So(err, ShouldBeNil)

					Convey("config policy should be a cpolicy.ConfigPolicy", func() {
						So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
					})

					testConfig := make(map[string]ctypes.ConfigValue)
					testConfig["graphite-server"] = ctypes.ConfigValueStr{Value: "localhost"}
					testConfig["graphite-port"] = ctypes.ConfigValueInt{Value: 8080}
					cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)

					Convey("So configpolicy should process testConfig and return a config", func() {
						So(cfg, ShouldNotBeNil)
					})
					Convey("so testConfig processing should not return errors", func() {
						So(errs.HasErrors(), ShouldBeFalse)
					})
				})
			})
		})
	})
}
