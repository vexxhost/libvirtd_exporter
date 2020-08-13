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
	"encoding/xml"
	"strconv"
	"time"

	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type DomainStatsCollector struct {
	prometheus.Collector

	Connection *libvirt.Connect

	Nova bool

	DomainSeconds *prometheus.Desc

	DomainDomainState       *prometheus.Desc
	DomainDomainStateReason *prometheus.Desc

	DomainCPUTime   *prometheus.Desc
	DomainCPUUser   *prometheus.Desc
	DomainCPUSystem *prometheus.Desc

	DomainBalloonCurrent *prometheus.Desc
	DomainBalloonMaximum *prometheus.Desc

	DomainVcpuState *prometheus.Desc
	DomainVcpuTime  *prometheus.Desc

	DomainNetRxBytes *prometheus.Desc
	DomainNetRxPkts  *prometheus.Desc
	DomainNetRxErrs  *prometheus.Desc
	DomainNetRxDrop  *prometheus.Desc
	DomainNetTxBytes *prometheus.Desc
	DomainNetTxPkts  *prometheus.Desc
	DomainNetTxErrs  *prometheus.Desc
	DomainNetTxDrop  *prometheus.Desc

	DomainBlockRdReqs     *prometheus.Desc
	DomainBlockRdBytes    *prometheus.Desc
	DomainBlockRdTimes    *prometheus.Desc
	DomainBlockWrReqs     *prometheus.Desc
	DomainBlockWrBytes    *prometheus.Desc
	DomainBlockWrTimes    *prometheus.Desc
	DomainBlockFlReqs     *prometheus.Desc
	DomainBlockFlTimes    *prometheus.Desc
	DomainBlockErrors     *prometheus.Desc
	DomainBlockAllocation *prometheus.Desc
	DomainBlockCapacity   *prometheus.Desc
	DomainBlockPhysical   *prometheus.Desc
}

type NovaFlavorMetadata struct {
	Name string `xml:"name,attr"`
}

type NovaOwnerMetadata struct {
	UUID string `xml:"uuid,attr"`
}

type NovaMetadata struct {
	Seconds      float64            `xml:"omitempty"`
	CreationTime string             `xml:"creationTime"`
	Flavor       NovaFlavorMetadata `xml:"flavor"`
	User         NovaOwnerMetadata  `xml:"owner>user"`
	Project      NovaOwnerMetadata  `xml:"owner>project"`
}

// nolint:funlen
func NewDomainStatsCollector(nova bool, connection *libvirt.Connect) (*DomainStatsCollector, error) {
	return &DomainStatsCollector{
		Connection: connection,
		Nova:       nova,

		DomainSeconds: prometheus.NewDesc(
			"libvirtd_domain_seconds",
			"seconds since creation time",
			[]string{"uuid", "instance_type", "user_id", "project_id"}, nil,
		),

		DomainDomainState: prometheus.NewDesc(
			"libvirtd_domain_domain_state",
			"state of the VM (virDomainState enum)",
			[]string{"uuid"}, nil,
		),
		DomainDomainStateReason: prometheus.NewDesc(
			"libvirtd_domain_domain_state_reason",
			"reason for entering given state (virDomain*Reason enum)",
			[]string{"uuid"}, nil,
		),

		DomainCPUTime: prometheus.NewDesc(
			"libvirtd_domain_cpu_time",
			"total cpu time spent for this domain in nanoseconds",
			[]string{"uuid"}, nil,
		),
		DomainCPUUser: prometheus.NewDesc(
			"libvirtd_domain_cpu_user",
			"user cpu time spent in nanoseconds",
			[]string{"uuid"}, nil,
		),
		DomainCPUSystem: prometheus.NewDesc(
			"libvirtd_domain_cpu_system",
			"system cpu time spent in nanoseconds",
			[]string{"uuid"}, nil,
		),

		DomainBalloonCurrent: prometheus.NewDesc(
			"libvirtd_domain_balloon_current",
			"the memory in kiB currently used",
			[]string{"uuid"}, nil,
		),
		DomainBalloonMaximum: prometheus.NewDesc(
			"libvirtd_domain_balloon_maximum",
			"the maximum memory in kiB allowed",
			[]string{"uuid"}, nil,
		),

		DomainVcpuState: prometheus.NewDesc(
			"libvirtd_domain_vcpu_state",
			"state of the virtual CPU (virVcpuState enum)",
			[]string{"uuid", "vcpu"}, nil,
		),
		DomainVcpuTime: prometheus.NewDesc(
			"libvirtd_domain_vcpu_time",
			"virtual cpu time spent",
			[]string{"uuid", "vcpu"}, nil,
		),

		DomainNetRxBytes: prometheus.NewDesc(
			"libvirtd_domain_net_rx_bytes",
			"bytes received",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetRxPkts: prometheus.NewDesc(
			"libvirtd_domain_net_rx_packets",
			"packets received",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetRxErrs: prometheus.NewDesc(
			"libvirtd_domain_net_rx_errors",
			"receive errors",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetRxDrop: prometheus.NewDesc(
			"libvirtd_domain_net_rx_drop",
			"receive packets dropped",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetTxBytes: prometheus.NewDesc(
			"libvirtd_domain_net_tx_bytes",
			"bytes transmitted",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetTxPkts: prometheus.NewDesc(
			"libvirtd_domain_net_tx_packets",
			"packets transmitted",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetTxErrs: prometheus.NewDesc(
			"libvirtd_domain_net_tx_errors",
			"transmission errors",
			[]string{"uuid", "interface"}, nil,
		),
		DomainNetTxDrop: prometheus.NewDesc(
			"libvirtd_domain_net_tx_drop",
			"transmit packets dropped",
			[]string{"uuid", "interface"}, nil,
		),

		DomainBlockRdReqs: prometheus.NewDesc(
			"libvirtd_domain_block_read_requests",
			"number of read requests",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockRdBytes: prometheus.NewDesc(
			"libvirtd_domain_block_read_bytes",
			"number of read bytes",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockRdTimes: prometheus.NewDesc(
			"libvirtd_domain_block_read_times",
			"total time (ns) spent on reads",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockWrReqs: prometheus.NewDesc(
			"libvirtd_domain_block_write_requests",
			"number of written requests",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockWrBytes: prometheus.NewDesc(
			"libvirtd_domain_block_write_bytes",
			"number of written bytes",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockWrTimes: prometheus.NewDesc(
			"libvirtd_domain_block_write_times",
			"total time (ns) spent on writes",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockFlReqs: prometheus.NewDesc(
			"libvirtd_domain_block_flush_requests",
			"total flush requests",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockFlTimes: prometheus.NewDesc(
			"libvirtd_domain_block_flush_times",
			"total time (ns) spent on cache flushing",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockAllocation: prometheus.NewDesc(
			"libvirtd_domain_block_allocation",
			"offset of the highest written sector",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockCapacity: prometheus.NewDesc(
			"libvirtd_domain_block_capacity",
			"logical size in bytes of the block device backing image",
			[]string{"uuid", "device", "path"}, nil,
		),
		DomainBlockPhysical: prometheus.NewDesc(
			"libvirtd_domain_block_physical",
			"physical size in bytes of the container of the backing image",
			[]string{"uuid", "device", "path"}, nil,
		),
	}, nil
}

func (c *DomainStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.describeNova(ch)
	c.describeState(ch)
	c.describeCPU(ch)
	c.describeBalloon(ch)
	c.describeVcpu(ch)
	c.describeNet(ch)
	c.describeBlock(ch)
}

func (c *DomainStatsCollector) describeNova(ch chan<- *prometheus.Desc) {
	if c.Nova {
		ch <- c.DomainSeconds
	}
}

func (c *DomainStatsCollector) describeState(ch chan<- *prometheus.Desc) {
	ch <- c.DomainDomainState
	ch <- c.DomainDomainStateReason
}

func (c *DomainStatsCollector) describeCPU(ch chan<- *prometheus.Desc) {
	ch <- c.DomainCPUTime
	ch <- c.DomainCPUUser
	ch <- c.DomainCPUSystem
}

func (c *DomainStatsCollector) describeBalloon(ch chan<- *prometheus.Desc) {
	ch <- c.DomainBalloonCurrent
	ch <- c.DomainBalloonMaximum
}

func (c *DomainStatsCollector) describeVcpu(ch chan<- *prometheus.Desc) {
	ch <- c.DomainVcpuState
	ch <- c.DomainVcpuTime
}

func (c *DomainStatsCollector) describeNet(ch chan<- *prometheus.Desc) {
	ch <- c.DomainNetRxBytes
	ch <- c.DomainNetRxPkts
	ch <- c.DomainNetRxErrs
	ch <- c.DomainNetRxDrop
	ch <- c.DomainNetTxBytes
	ch <- c.DomainNetTxPkts
	ch <- c.DomainNetTxErrs
	ch <- c.DomainNetTxDrop
}

func (c *DomainStatsCollector) describeBlock(ch chan<- *prometheus.Desc) {
	ch <- c.DomainBlockRdReqs
	ch <- c.DomainBlockRdBytes
	ch <- c.DomainBlockRdTimes
	ch <- c.DomainBlockWrReqs
	ch <- c.DomainBlockWrBytes
	ch <- c.DomainBlockWrTimes
	ch <- c.DomainBlockFlReqs
	ch <- c.DomainBlockFlTimes
	ch <- c.DomainBlockAllocation
	ch <- c.DomainBlockCapacity
	ch <- c.DomainBlockPhysical
}

func (c *DomainStatsCollector) Collect(ch chan<- prometheus.Metric) {
	alive, err := c.Connection.IsAlive()
	if err != nil {
		log.Errorln(err)
		return
	}

	if !alive {
		uri, err := c.Connection.GetURI()
		if err != nil {
			// NOTE(mnaser): If we get to this point, we don't have
			//               a URI and we can't reconnect, die
			log.Fatalln(err)
			return
		}

		c.Connection.Close()

		conn, err := libvirt.NewConnect(uri)
		if err != nil {
			log.Errorln(err)
			return
		}
		c.Connection = conn
	}

	stats, err := c.Connection.GetAllDomainStats(
		[]*libvirt.Domain{},
		libvirt.DOMAIN_STATS_STATE|libvirt.DOMAIN_STATS_CPU_TOTAL|libvirt.DOMAIN_STATS_BALLOON|
			libvirt.DOMAIN_STATS_VCPU|libvirt.DOMAIN_STATS_INTERFACE|libvirt.DOMAIN_STATS_BLOCK,
		0,
	)

	defer func(stats []libvirt.DomainStats) {
		for _, stat := range stats {
			err := stat.Domain.Free()
			if err != nil {
				log.Errorln(err)
			}
		}
	}(stats)

	if err != nil {
		log.Errorln(err)
		return
	}

	for _, stat := range stats {
		uuid, err := stat.Domain.GetUUIDString()
		if err != nil {
			log.Errorln(err)
			continue
		}

		c.collectNova(uuid, stat, ch)
		c.collectState(uuid, stat, ch)
		c.collectCPU(uuid, stat, ch)
		c.collectBalloon(uuid, stat, ch)
		c.collectVcpu(uuid, stat, ch)
		c.collectNet(uuid, stat, ch)
		c.collectBlock(uuid, stat, ch)
	}
}

func (c *DomainStatsCollector) collectNova(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	if c.Nova {
		metadata, err := c.getNovaMetadata(stat.Domain)

		if err != nil {
			log.Errorln(err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.DomainSeconds,
				prometheus.CounterValue,
				metadata.Seconds, uuid, metadata.Flavor.Name, metadata.User.UUID, metadata.Project.UUID,
			)
		}
	}
}

func (c *DomainStatsCollector) collectState(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.DomainDomainState,
		prometheus.GaugeValue,
		float64(stat.State.State), uuid,
	)
	ch <- prometheus.MustNewConstMetric(
		c.DomainDomainStateReason,
		prometheus.GaugeValue,
		float64(stat.State.Reason), uuid,
	)
}

func (c *DomainStatsCollector) collectCPU(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	if stat.Cpu != nil {
		ch <- prometheus.MustNewConstMetric(
			c.DomainCPUTime,
			prometheus.CounterValue,
			float64(stat.Cpu.Time), uuid,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainCPUUser,
			prometheus.CounterValue,
			float64(stat.Cpu.User), uuid,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainCPUSystem,
			prometheus.CounterValue,
			float64(stat.Cpu.System), uuid,
		)
	}
}

func (c *DomainStatsCollector) collectBalloon(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.DomainBalloonCurrent,
		prometheus.GaugeValue,
		float64(stat.Balloon.Current), uuid,
	)
	ch <- prometheus.MustNewConstMetric(
		c.DomainBalloonMaximum,
		prometheus.GaugeValue,
		float64(stat.Balloon.Maximum), uuid,
	)
}

func (c *DomainStatsCollector) collectVcpu(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	for vcpu, vcpuStats := range stat.Vcpu {
		ch <- prometheus.MustNewConstMetric(
			c.DomainVcpuState,
			prometheus.GaugeValue,
			float64(vcpuStats.State), uuid, strconv.Itoa(vcpu),
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainVcpuTime,
			prometheus.CounterValue,
			float64(vcpuStats.Time), uuid, strconv.Itoa(vcpu),
		)
	}
}

func (c *DomainStatsCollector) collectNet(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	for _, netStats := range stat.Net {
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetRxBytes,
			prometheus.CounterValue,
			float64(netStats.RxBytes), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetRxPkts,
			prometheus.CounterValue,
			float64(netStats.RxPkts), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetRxErrs,
			prometheus.CounterValue,
			float64(netStats.RxErrs), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetRxDrop,
			prometheus.CounterValue,
			float64(netStats.RxDrop), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetTxBytes,
			prometheus.CounterValue,
			float64(netStats.TxBytes), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetTxPkts,
			prometheus.CounterValue,
			float64(netStats.TxPkts), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetTxErrs,
			prometheus.CounterValue,
			float64(netStats.TxErrs), uuid, netStats.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainNetTxDrop,
			prometheus.GaugeValue,
			float64(netStats.TxDrop), uuid, netStats.Name,
		)
	}
}

func (c *DomainStatsCollector) collectBlock(uuid string, stat libvirt.DomainStats, ch chan<- prometheus.Metric) {
	for device, blockStats := range stat.Block {
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockRdReqs,
			prometheus.CounterValue,
			float64(blockStats.RdReqs), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockRdBytes,
			prometheus.CounterValue,
			float64(blockStats.RdBytes), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockRdTimes,
			prometheus.CounterValue,
			float64(blockStats.RdTimes), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockWrReqs,
			prometheus.CounterValue,
			float64(blockStats.WrReqs), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockWrBytes,
			prometheus.CounterValue,
			float64(blockStats.WrBytes), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockWrTimes,
			prometheus.CounterValue,
			float64(blockStats.WrTimes), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockFlReqs,
			prometheus.CounterValue,
			float64(blockStats.FlReqs), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockFlTimes,
			prometheus.CounterValue,
			float64(blockStats.FlTimes), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockAllocation,
			prometheus.GaugeValue,
			float64(blockStats.Allocation), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockCapacity,
			prometheus.GaugeValue,
			float64(blockStats.Capacity), uuid, strconv.Itoa(device), blockStats.Path,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DomainBlockPhysical,
			prometheus.GaugeValue,
			float64(blockStats.Physical), uuid, strconv.Itoa(device), blockStats.Path,
		)
	}
}

func (c *DomainStatsCollector) getNovaMetadata(domain *libvirt.Domain) (*NovaMetadata, error) {
	data, err := domain.GetMetadata(
		libvirt.DOMAIN_METADATA_ELEMENT,
		"http://openstack.org/xmlns/libvirt/nova/1.0",
		libvirt.DOMAIN_AFFECT_LIVE,
	)
	if err != nil {
		return nil, err
	}

	m := &NovaMetadata{}
	err = xml.Unmarshal([]byte(data), &m)

	if err != nil {
		return nil, err
	}

	// Parse creationTime from Nova format: "%Y-%m-%d %H:%M:%S"
	layout := "2006-01-02 15:04:05"
	creationTime, err := time.Parse(layout, m.CreationTime)

	if err != nil {
		return nil, err
	}

	m.Seconds = time.Since(creationTime).Seconds()

	return m, nil
}
