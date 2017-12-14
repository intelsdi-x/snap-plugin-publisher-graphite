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

func Test_createKey(t *testing.T) {
	type args struct {
		m          plugin.Metric
		enableTags bool
	}
	tests := []struct {
		description string
		args        args
		matchKey    string
		matchTags   []string
	}{
		{
			description: "static namespace",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar"},
						plugin.NamespaceElement{Value: "value"},
					},
				},
				enableTags: false,
			},
			matchKey:  "test.bar.value",
			matchTags: []string{},
		},
		{
			description: "dynamic namespace",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar", Name: "foo"},
						plugin.NamespaceElement{Value: "value"},
					},
				},
				enableTags: false,
			},
			matchKey:  "test.bar.value",
			matchTags: []string{},
		},
		{
			description: "TagsEnabled, static namespace",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar"},
						plugin.NamespaceElement{Value: "value"},
					},
				},
				enableTags: true,
			},
			matchKey:  "test.bar.value",
			matchTags: []string{},
		},
		{
			description: "TagsEnabled, dynamic namespace",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar", Name: "foo"},
						plugin.NamespaceElement{Value: "value"},
					},
				},
				enableTags: true,
			},
			matchKey: "test.value",
			matchTags: []string{
				"foo=bar",
			},
		},
		{
			description: "TagsEnabled, dynamic namespace with tags",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar", Name: "foo"},
						plugin.NamespaceElement{Value: "value"},
					},
					Tags: map[string]string{
						"test": "tag",
					},
				},
				enableTags: true,
			},
			matchKey: "test.value",
			matchTags: []string{
				"foo=bar",
				"test=tag",
			},
		},
		{
			description: "TagsEnabled, dynamic namespace with plugin_running_on tag",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar", Name: "foo"},
						plugin.NamespaceElement{Value: "value"},
					},
					Tags: map[string]string{
						"plugin_running_on": "example",
					},
				},
				enableTags: true,
			},
			matchKey: "test.value",
			matchTags: []string{
				"foo=bar",
				"host=example",
			},
		},
		{
			description: "TagsEnabled, dynamic namespace with dynamic namespace matching tag",
			args: args{
				m: plugin.Metric{
					Namespace: plugin.Namespace{
						plugin.NamespaceElement{Value: "test"},
						plugin.NamespaceElement{Value: "bar", Name: "foo"},
						plugin.NamespaceElement{Value: "value"},
					},
					Tags: map[string]string{
						"foo": "bar1",
					},
				},
				enableTags: true,
			},
			matchKey: "test.value",
			matchTags: []string{
				"foo=bar1",
			},
		},
	}
	Convey("Ensure proper graphite metric name conversion", t, func() {
		for _, tt := range tests {
			Convey(tt.description, func() {
				key := createKey(tt.args.m, []string{}, tt.args.enableTags)
				if tt.args.enableTags {
					for _, tag := range tt.matchTags {
						So(key, ShouldContainSubstring, tag)
					}
				}
				So(strings.SplitN(key, ";", 2)[0], ShouldEqual, tt.matchKey)
			})
		}
	})
}
