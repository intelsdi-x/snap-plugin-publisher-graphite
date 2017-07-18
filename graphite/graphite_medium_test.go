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
	. "net/http"

	"os"
	"time"
  "testing"
  "strings"

  . "github.com/smartystreets/goconvey/convey"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
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

  Convey("snap plugin Graphite integration testing with Graphite", t, func() {
    ip := &GraphitePublisher{}
    var retention string

		if strings.HasPrefix(os.Getenv("GRAPHITE_VERSION"), "0.") {
			retention = "default"
		} else {
			retention = "autogen"
		}
		config := plugin.Config{
			"host":          os.Getenv("SNAP_GRAPHITE_HOST"),
			"skip-verify":   false,
			"retention":     retention,
			"isMultiFields": false,
			"debug":         false,
			"log-level":     "debug",
			"precision":     "s",
		}

		config["scheme"] = HTTP
		config["port"] = int64(80)
		tests(HTTP, config)

		config["scheme"] = UDP
		config["port"] = int64(4444)
		tests(UDP, config)
  })
}
