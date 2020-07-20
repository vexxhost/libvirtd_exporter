// Copyright 2019 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collectors

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type VersionCollector struct {
	prometheus.Collector

	Connection *libvirt.Connect

	Version *prometheus.Desc
}

func NewVersionCollector(conn *libvirt.Connect) (*VersionCollector, error) {
	return &VersionCollector{
		Connection: conn,

		Version: prometheus.NewDesc(
			"libvirtd_info",
			"Version details for LibvirtD",
			[]string{"driver", "driver_version", "version"}, nil,
		),
	}, nil
}

func (c *VersionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Version
}

func (c *VersionCollector) Collect(ch chan<- prometheus.Metric) {
	hypervisorType, err := c.Connection.GetType()
	if err != nil {
		log.Errorln(err)
		return
	}

	hypervisorVersion, err := c.Connection.GetVersion()
	if err != nil {
		log.Errorln(err)
		return
	}

	libvirtVersion, err := c.Connection.GetLibVersion()
	if err != nil {
		log.Errorln(err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.Version,
		prometheus.CounterValue,
		float64(1),
		hypervisorType,
		versionToString(hypervisorVersion),
		versionToString(libvirtVersion),
	)
}

func versionToString(version uint32) string {
	major := version / 1000000
	minor := (version - (major * 1000000)) / 1000
	release := version - (major * 1000000) - (minor * 1000)

	return fmt.Sprintf("%d.%d.%d", major, minor, release)
}
