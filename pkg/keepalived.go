package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/procfs"
)

const (
	mountPath             = "/host/proc"
	defaultKeepalivedConf = "/etc/keepalived/keepalived.conf"
)

// accquire all procs on node
func accquireProcsOnNode() procfs.Procs {
	newFS, err := procfs.NewFS(mountPath)
	if err != nil {
		glog.Fatalf("mount path can not be found")
		os.Exit(1)
	}

	fs, err := newFS.AllProcs()
	if err != nil {
		glog.Errorf("failed to open procfs: %v", err)
		return nil
	}

	return fs
}

func accquireHostname() string {
	std, err := exec.Command("bash", "-c", "hostname", "-f").CombinedOutput()
	if err != nil {
		glog.Errorf("failed to get current node hostname")
		return ""
	}

	hostname := strings.TrimSuffix(string(std), "\n")
	return hostname
}

func readKeepalivedFile() string {
	file, err := ioutil.ReadFile(defaultKeepalivedConf)
	if err != nil {
		glog.Errorf("fail to open keepalived config on path %v, err: %v", defaultKeepalivedConf, err)
		return ""
	}

	return string(file)
}

// parse keepalived procs stat on node
func parseKeepalivedStatus() *KeepalivedProcStatus {
	kpprocstate := &KeepalivedProcStatus{
		Comm:     "keepalived",
		Running:  0,
		Sleeping: 0,
		Waiting:  0,
		Zombie:   0,
		Other:    0,
	}

	procSet := accquireProcsOnNode()
	for _, pid := range procSet {
		process, err := pid.Stat()
		if err != nil {
			glog.Infof("error occurs on open proc %v, err: %v", pid, err)
		}

		if process.Comm == "keepalived" {
			switch process.State {
			case "R":
				kpprocstate.Running++
			case "S":
				kpprocstate.Sleeping++
			case "W":
				kpprocstate.Waiting++
			case "Z":
				kpprocstate.Zombie++
			case "O":
				kpprocstate.Other++
			default:
				glog.Infof("can not get state of process: %v", process.Comm)
			}
		}
	}

	return kpprocstate
}

// parse keepalived config and return ip collection of keepalived VIP
func parseKeepalivedVIP(inputContent string) []string {
	var ipCollection []string
	vipMap := make(map[string]bool)

	vipPattern := `(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	vipObject := regexp.MustCompile(vipPattern)
	matchVIP := vipObject.FindAllStringSubmatch(inputContent, -1)

	for _, value := range matchVIP {
		// TODO should we ignore localhost?
		if value[0] != "127.0.0.1" {
			if _, ok := vipMap[value[0]]; !ok {
				vipMap[value[0]] = true
				ipCollection = append(ipCollection, value[0])
			}
		}
	}

	if len(ipCollection) == 0 {
		glog.Infof("keepalived ip address is empty")
		return nil
	}

	return ipCollection
}

// check VIP on node or not
func checkKeepalivedVipReady() bool {
	var vipCheckArray []bool
	keepalivedContent := readKeepalivedFile()
	currentHost := accquireHostname()

	if keepalivedContent == "" {
		glog.Fatalf("keepalived config parse failed, return with empty string")
		return false
	}

	IPCollection := parseKeepalivedVIP(keepalivedContent)

	for _, v := range IPCollection {
		cmd := fmt.Sprintf("ip addr | grep %v 2>/dev/null", v)
		std, err := exec.Command("bash", "-c", cmd).Output()

		if err != nil {
			// glog.Errorf("exec command: [%v], error occurs: %v", cmd, err)
			return false
		}

		currentVIP := parseKeepalivedVIP(string(std))
		if len(currentVIP) == 0 {
			glog.Errorf("current host: \"%v\" lost vip", currentHost)
			return false
		}

		if parseKeepalivedVIP(string(std))[0] == v {
			vipCheckArray = append(vipCheckArray, true)
		} else {
			vipCheckArray = append(vipCheckArray, false)
		}
	}

	for _, v := range vipCheckArray {
		if v == false {
			glog.Infof("current host lost vip")
			return false
		}
	}

	if len(vipCheckArray) == len(IPCollection) {
		return true
	}

	return false
}

// upgrade status of keepalived VIP
func updateKeepalivedVIP() int {
	if checkKeepalivedVipReady() == true {
		return 1
	}
	return 0
}

// upgrade status of keepalived status
func updateKeepalivedStatus() *KeepalivedProcStatus {
	procState := parseKeepalivedStatus()
	if procState == nil {
		glog.Errorf("can not catch proc keepalived state, return nil")
		return nil
	}
	return procState
}
