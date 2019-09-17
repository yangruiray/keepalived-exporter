package pkg

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func NewKeepalivedMetrics() *KeepalivedMetrics {
	return &KeepalivedMetrics{
		KeepalivedStatusRunning: prometheus.NewDesc(
			"keepalived_status_running",
			"keepalived running status on node",
			[]string{"host"},
			prometheus.Labels{"state": "running"},
		),
		KeepalivedStatusSleeping: prometheus.NewDesc(
			"keepalived_status_sleeping",
			"keepalived sleeping status on node",
			[]string{"host"},
			prometheus.Labels{"state": "sleeping"},
		),
		KeepalivedStatusWaiting: prometheus.NewDesc(
			"keepalived_status_waiting",
			"keepalived waiting status on node",
			[]string{"host"},
			prometheus.Labels{"state": "waiting"},
		),
		KeepalivedStatusZombie: prometheus.NewDesc(
			"keepalived_status_zombie",
			"keepalived zombie status on node",
			[]string{"host"},
			prometheus.Labels{"state": "zombie"},
		),
		KeepalivedStatusOther: prometheus.NewDesc(
			"keepalived_status_other",
			"keepalived other status on node",
			[]string{"host"},
			prometheus.Labels{"state": "other"},
		),
		KeepalivedVIP: prometheus.NewDesc(
			"keepalived_vip_ready",
			"keepalived vip on node",
			[]string{"host"},
			nil,
		),
	}
}

// Exec to get state of keepalived
func (m *KeepalivedMetrics) GetState() (keepalivedStatus map[string]*KeepalivedProcStatus, keepalivedVip map[string]int) {

	kpalivedStatGauge := updateKeepalivedStatus()
	kpalivedVipGauge := updateKeepalivedVIP()
	hostName := accquireHostname()

	keepalivedStatus = map[string]*KeepalivedProcStatus{
		fmt.Sprintf("%v", hostName): kpalivedStatGauge,
	}

	keepalivedVip = map[string]int{
		fmt.Sprintf("%v", hostName): kpalivedVipGauge,
	}

	return
}

// Collect func is for write gauge value to channel
func (m *KeepalivedMetrics) Collect(ch chan<- prometheus.Metric) {
	keepalivedStatus, keepalivedVIP := m.GetState()

	for host, values := range keepalivedStatus {
		ch <- prometheus.MustNewConstMetric(
			m.KeepalivedStatusRunning,
			prometheus.GaugeValue,
			float64(values.Running),
			host,
		)
	}

	for host, values := range keepalivedStatus {
		ch <- prometheus.MustNewConstMetric(
			m.KeepalivedStatusSleeping,
			prometheus.GaugeValue,
			float64(values.Sleeping),
			host,
		)
	}

	for host, values := range keepalivedStatus {
		ch <- prometheus.MustNewConstMetric(
			m.KeepalivedStatusWaiting,
			prometheus.GaugeValue,
			float64(values.Waiting),
			host,
		)
	}

	for host, values := range keepalivedStatus {
		ch <- prometheus.MustNewConstMetric(
			m.KeepalivedStatusZombie,
			prometheus.GaugeValue,
			float64(values.Zombie),
			host,
		)
	}

	for host, values := range keepalivedStatus {
		ch <- prometheus.MustNewConstMetric(
			m.KeepalivedStatusOther,
			prometheus.GaugeValue,
			float64(values.Other),
			host,
		)
	}

	for host, values := range keepalivedVIP {
		ch <- prometheus.MustNewConstMetric(
			m.KeepalivedVIP,
			prometheus.GaugeValue,
			float64(values),
			host,
		)
	}
}

// Write describe to channel
func (m *KeepalivedMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.KeepalivedStatusRunning
	ch <- m.KeepalivedStatusSleeping
	ch <- m.KeepalivedStatusWaiting
	ch <- m.KeepalivedStatusZombie
	ch <- m.KeepalivedStatusOther
	ch <- m.KeepalivedVIP
}
