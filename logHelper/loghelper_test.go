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

package loghelper

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func testMeta() *plugin.PluginMeta {
	return plugin.NewPluginMeta("test", 1, plugin.PublisherPluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

type logTestHelper struct {
	key      string
	value    string
	cvalue   ctypes.ConfigValue
	loglevel log.Level
}

func getHelperStr(key, value string, loglevel log.Level) logTestHelper {
	return logTestHelper{
		key:      key,
		value:    value,
		cvalue:   ctypes.ConfigValueStr{Value: value},
		loglevel: loglevel,
	}
}

func getHelperBool(key string, value bool, loglevel log.Level) logTestHelper {
	return logTestHelper{
		key:      key,
		value:    "",
		cvalue:   ctypes.ConfigValueBool{Value: value},
		loglevel: loglevel,
	}
}

func TestGetLogger(t *testing.T) {
	Convey("Pass in debug config", t, func() {
		dth := getHelperBool("debug", true, log.DebugLevel)
		helpTestGetLogger(dth, t)
		dth = getHelperBool("debug", false, log.WarnLevel)
		helpTestGetLogger(dth, t)
		dth = getHelperStr("debug", "invalid", log.WarnLevel)
		helpTestGetLogger(dth, t)
	})

	Convey("Pass in loglevel config", t, func() {
		lth := getHelperStr("log-level", "warn", log.WarnLevel)
		helpTestGetLogger(lth, t)
		lth = getHelperStr("log-level", "error", log.ErrorLevel)
		helpTestGetLogger(lth, t)
		lth = getHelperStr("log-level", "debug", log.DebugLevel)
		helpTestGetLogger(lth, t)
		lth = getHelperStr("log-level", "info", log.InfoLevel)
		helpTestGetLogger(lth, t)
		lth = getHelperStr("log-level", "invalid", log.WarnLevel)
		helpTestGetLogger(lth, t)
	})
}

func helpTestGetLogger(lth logTestHelper, t *testing.T) {
	Convey("with "+lth.key+"='"+lth.value+"' loglevel should be: "+fmt.Sprintf("%v", lth.loglevel), func() {
		testConfig := make(map[string]ctypes.ConfigValue)
		testConfig[lth.key] = lth.cvalue
		GetLogger(testConfig, testMeta())
		So(log.GetLevel(), ShouldEqual, lth.loglevel)
	})
}
