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
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"libvirt.org/go/libvirt"
)

type VersionCollector struct {
	prometheus.Collector

	logger     *slog.Logger
	connection *libvirt.Connect

	Version *prometheus.Desc
}

func NewVersionCollector(logger *slog.Logger, connection *libvirt.Connect) *VersionCollector {
	return &VersionCollector{
		logger:     logger,
		connection: connection,

		Version: prometheus.NewDesc(
			"libvirtd_info",
			"Version details for LibvirtD",
			[]string{"driver", "driver_version", "version"}, nil,
		),
	}
}

func (c *VersionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Version
}

func (c *VersionCollector) Collect(ch chan<- prometheus.Metric) {
	alive, err := c.connection.IsAlive()
	if err != nil {
		c.logger.Error("Failed to check if connection is alive", "err", err)
		return
	}

	if !alive {
		uri, err := c.connection.GetURI()
		if err != nil {
			// NOTE(mnaser): If we get to this point, we don't have
			//               a URI and we can't reconnect, die
			c.logger.Error("Failed to get URI", "err", err)
			return
		}

		c.connection.Close()

		conn, err := libvirt.NewConnect(uri)
		if err != nil {
			c.logger.Error("Failed to reconnect", "err", err)
			return
		}
		c.connection = conn
	}

	hypervisorType, err := c.connection.GetType()
	if err != nil {
		c.logger.Error("Failed to get hypervisor type", "err", err)
		return
	}

	hypervisorVersion, err := c.connection.GetVersion()
	if err != nil {
		c.logger.Error("Failed to get hypervisor version", "err", err)
		return
	}

	libvirtVersion, err := c.connection.GetLibVersion()
	if err != nil {
		c.logger.Error("Failed to get libvirt version", "err", err)
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
