package pkg

import "github.com/prometheus/client_golang/prometheus"

type KeepalivedMetrics struct {
	KeepalivedStatusRunning  *prometheus.Desc
	KeepalivedStatusSleeping *prometheus.Desc
	KeepalivedStatusWaiting  *prometheus.Desc
	KeepalivedStatusZombie   *prometheus.Desc
	KeepalivedStatusOther    *prometheus.Desc
	KeepalivedVIP            *prometheus.Desc
}

type KeepalivedProcStatus struct {
	Comm     string
	Running  int
	Sleeping int
	Waiting  int
	Zombie   int
	Other    int
}
