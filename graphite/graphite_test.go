// +build small

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
	"strings"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphitePlugin(t *testing.T) {

	Convey("Create Graphite publisher", t, func() {
		gp := &GraphitePublisher{}
		Convey("So publisher should not be nil", func() {
			So(gp, ShouldNotBeNil)
		})

		Convey("Publisher should be of type graphitePublisher", func() {
			So(gp, ShouldHaveSameTypeAs, &GraphitePublisher{})
		})

		configPolicy, err := gp.GetConfigPolicy()
		Convey("Should return a config policy", func() {

			Convey("configPolicy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)

				Convey("and retrieving config policy should not error", func() {
					So(err, ShouldBeNil)

					Convey("config policy should be a cpolicy.ConfigPolicy", func() {
						So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
					})

					testConfig := make(plugin.Config)
					testConfig["server"] = "localhost"
					testConfig["port"] = int64(8080)

					server, err := testConfig.GetString("server")
					Convey("So testConfig should return the right server config", func() {
						So(err, ShouldBeNil)
						So(server, ShouldEqual, "localhost")
					})

					port, err := testConfig.GetInt("port")
					Convey("So testConfig should return the right port config", func() {
						So(err, ShouldBeNil)
						So(port, ShouldEqual, int64(8080))

					})
				})
			})
		})

		Convey("Check for Illegal chars and replace", func() {
			key := "testing if this (string) has / any illegal, {characters}"
			illegal := "(), /{}"
			expected := "testing_if_this_[string]_has_|_any_illegal;_[characters]"
			r := strings.NewReplacer(" ", "_",
				",", ";",
				"(", "[",
				")", "]",
				"/", "|",
				"{", "[",
				"}", "]")

			if strings.ContainsAny(key, illegal) {
				key = r.Replace(key)
			}
			So(key, ShouldEqual, expected)
		})
	})
}
