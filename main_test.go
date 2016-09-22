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

package main

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/intelsdi-x/snap/control"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/plugin/helper"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	PluginName = "snap-plugin-publisher-graphite"
	PluginType = "publisher"
	PluginPath = path.Join("build", runtime.GOOS, "x86_64", PluginName)
)

func TestGraphitePluginLoad(t *testing.T) {
	// Helper plugin trigger build if possible for this plugin
	helper.BuildPlugin(PluginType, PluginName)
	Convey("ensure plugin loads and responds", t, func() {
		c := control.New(control.GetDefaultConfig())
		c.Start()
		gp, err := core.NewRequestedPlugin(PluginPath)
		So(err, ShouldBeNil)
		_, err = c.Load(gp)
		So(err, ShouldBeNil)
	})
}

func TestMain(t *testing.T) {
	Convey("ensure plugin loads and responds", t, func() {
		os.Args = []string{"", "{\"NoDaemon\": true}"}
		So(func() { main() }, ShouldNotPanic)
	})
}
