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

package main

import (
	"net/http"

	"github.com/libvirt/libvirt-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"opendev.org/vexxhost/libvirtd_exporter/collectors"
)

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9474").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		libvirtURI = kingpin.Flag(
			"libvirt.uri",
			"Libvirt Connection URI",
		).Default("qemu:///system").String()
		libvirtNova = kingpin.Flag(
			"libvirt.nova",
			"Parse Libvirt Nova metadata",
		).Bool()
	)

	kingpin.Version(version.Print("libvirtd_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	conn, err := libvirt.NewConnect(*libvirtURI)
	if err != nil {
		log.Fatalln(err)
		return
	}

	versionCollector, err := collectors.NewVersionCollector(conn)
	if err != nil {
		log.Fatalln(err)
	}

	domainStats, err := collectors.NewDomainStatsCollector(*libvirtNova, conn)
	if err != nil {
		log.Fatalln(err)
	}

	prometheus.MustRegister(domainStats)
	prometheus.MustRegister(versionCollector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Libvirtd Exporter</title></head>
			<body>
			<h1>Libvirtd Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>Build</h2>
			<pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
			</body>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
