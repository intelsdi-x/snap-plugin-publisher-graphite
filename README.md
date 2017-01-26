# Snap publisher plugin - Graphite

This plugin publishes metrics to graphite.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Task Manifest Config](#task-manifest-config)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

### System Requirements
* The Snap daemon is running
* A running version of [Graphite](https://graphite.readthedocs.org/en/latest/) reachable by the Snap daemon is required for this plugin to successfully publish data

### Installation
#### Download Graphite plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [Github Releases](https://github.com/intelsdi-x/snap-plugin-publisher-graphite/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-publisher-graphite

Clone repo into `$GOPATH/src/github/intelsdi-x/`:  
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-graphite
```
Build the plugin by running `make` in the repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

## Documentation

### Task Manifest Config
A Task Manifest that includes the publishing to Graphite will require configuration data in order for the plugin to establish a connection. Config arguments include:
* "server" (required) - the IP of graphite host.
* "port" (optional) - 2003 by default.
* "prefix_tags" (optional) - coma separated list of metric tags used to add prefix on the published name, "plugin_running_on" by default.
* "prefix" (optional) - change the published name prefix, happens before prefixes from "prefix_tags".

Note: Tags in the graphite plugin are joined with the metric name without any kind of escaping and that could introduce an unintentional seperator. For example, spaces in tag could lead to a wrong request for the graphite plugin, as the metric name will be in the format %s %s %s. 
To avoid that, the graphite plugin handles special characters that are not legal in a graphite namespace. Those special characters include: , /(){}. Currently, the graphite publisher checks for all the illegal characters in the metric name and replaces each of those illegal characters with legal characters. 

### Examples

An example of running the [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil), [statistics processor plugin](https://github.com/intelsdi-x/snap-plugin-processor-statistics), and publishing data to Graphite is available in the example
Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap) in one of the following ways:
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* manually: `sudo snapteld -l 1 -t 0 &`


Download and load Snap plugins (paths to binary files for Linux/amd64):
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-graphite/latest/linux/x86_64/snap-plugin-publisher-graphite
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-processor-statistics/latest/linux/x86_64/snap-plugin-processor-statistics
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-psutil/latest/linux/x86_64/snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-publisher-graphite
$ snaptel plugin load snap-plugin-processor-statistics
$ snaptel plugin load snap-plugin-collector-psutil
```

Create a [task manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) (see [exemplary tasks](examples/tasks/)),
for example `psutil-statistics-graphite-simple.json` with following content:
```json
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "1s"
  },
  "max-failures": 10,
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/psutil/load/load1": {},
        "/intel/psutil/load/load15": {},
        "/intel/psutil/load/load5": {},
        "/intel/psutil/vm/available": {},
        "/intel/psutil/vm/free": {},
        "/intel/psutil/vm/used": {}
      },
      "process": [
        {
          "plugin_name": "statistics",
          "publish": [
            {
              "plugin_name": "graphite",
              "config": {
                "server": "127.0.0.1"
              }
            }
          ]
        }
      ]
    }
  }
}
```

Create a task:
```
$ snaptel task create -t psutil-statistics-graphite.json
```

Watch created task:
```
$ snaptel task watch <task_id>
```

To stop previously created task:
```
$ snaptel task stop <task_id>
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-graphite/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-graphite/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: [Cody Roseborough](https://github.com/ircody)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
