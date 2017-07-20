// +build medium

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
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	//Do a ping to make sure the docker image actually came up. Otherwise this can fail Travis builds
	for i := 0; i < 3; i++ {
		resp, err := http.Get("http://" + os.Getenv("SNAP_GRAPHITE_HOST") + ":80/ping")
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			//Try again after 3 second
			time.Sleep(3 * time.Second)
		} else {
			//Give the run.sh time to create the test database
			time.Sleep(5 * time.Second)
			return
		}
	}
	//If we got here, we failed to get to the server
	panic("Unable to connect to Graphite host. Aborting test.")
}

func TestGraphitePublisher(t *testing.T) {
	ip := &GraphitePublisher{}
	config := plugin.Config{
		"server":      os.Getenv("SNAP_GRAPHITE_HOST"),
		"port":        int64(80),
		"prefix":      "medium_test_prefix",
		"prefix_tags": "medium_test_prefix_tag_1,medium_test_prefix_tag_2",
	}
	tags := map[string]string{"zone": "red"}
	mcfg := map[string]interface{}{"field": "abc123"}
	metric1 := plugin.Metric{
		Namespace: plugin.NewNamespace("test1"),
		Timestamp: time.Now(),
		Config:    mcfg,
		Tags:      tags,
		Unit:      "someunit",
		Data:      99,
	}
	metric2 := plugin.Metric{
		Namespace: plugin.NewNamespace("test2"),
		Timestamp: time.Now(),
		Config:    mcfg,
		Tags:      tags,
		Unit:      "someunit2",
		Data:      200,
	}
	Convey("Snap plugin Graphite integration testing with Graphite", t, func() {

		Convey("Publish with correct data", func() {
			Convey("no metrics", func() {
				err := ip.Publish([]plugin.Metric{}, config)
				So(err, ShouldBeNil)
			})
			Convey("single metric", func() {
				err := ip.Publish([]plugin.Metric{metric1}, config)
				So(err, ShouldBeNil)
			})
			Convey("two integer metrics", func() {
				err := ip.Publish([]plugin.Metric{metric1, metric2}, config)
				So(err, ShouldBeNil)
			})
		})

	})
}
func TestWrongConfig(t *testing.T) {
	ip := &GraphitePublisher{}

	tags := map[string]string{"zone": "red"}
	mcfg := map[string]interface{}{"field": "abc123"}
	metrics := []plugin.Metric{plugin.Metric{
		Namespace: plugin.NewNamespace("test1"),
		Timestamp: time.Now(),
		Config:    mcfg,
		Tags:      tags,
		Unit:      "someunit",
		Data:      99,
	}}
	config := plugin.Config{
		"prefix":      "medium_test_prefix",
		"prefix_tags": "medium_test_prefix_tag_1,medium_test_prefix_tag_2",
	}
	Convey("Incorrect config ", t, func() {
		Convey("nil server and port", func() {
			config["server"] = nil
			config["port"] = nil
			err := ip.Publish(metrics, config)
			So(err, ShouldNotBeNil)
		})

		Convey("nil port", func() {
			err := ip.Publish(metrics, config)
			config["server"] = os.Getenv("SNAP_GRAPHITE_HOST")
			config["port"] = nil
			So(err, ShouldNotBeNil)
		})

		Convey("wrong server ip", func() {
			config["server"] = "6.6.6.6"
			config["port"] = int64(80)
			err := ip.Publish(metrics, config)
			So(err, ShouldNotBeNil)
		})

		Convey("wrong port", func() {
			config["server"] = os.Getenv("SNAP_GRAPHITE_HOST")
			config["port"] = int64(666)
			err := ip.Publish(metrics, config)
			So(err, ShouldNotBeNil)
		})

	})
}
