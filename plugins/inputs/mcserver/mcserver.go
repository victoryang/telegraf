package mcserver

import (
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

const (
	APISERVER = "api-server"
	MCSERVER = "mcserver"
	COBOT = "cobot"
	PARAMSERVER = "param-server"
)

var monitorProcesses map[string]*process.Process = nil

func getMonitorProcesses() {
	monitorProcesses = make(map[string]*process.Process)

	if procs, err := process.Processes(); err == nil {
		for _,proc := range procs {
			n,err := proc.Name();
			if err!=nil {
				continue
			}

			switch n {
			case APISERVER, MCSERVER, COBOT, PARAMSERVER:
				monitorProcesses[n] = proc
			default:
			}
		}
	}
}

type McServer struct {

}

var sampleConfig = `
  ## No Configuration

`

func (_ *McServer) SampleConfig() string {
	return sampleConfig
}

func (_ *McServer) Description() string {
	return "Monitor mcserver system status"
}

func GatherMonitoringProcessStat(name string, proc *process.Process, acc telegraf.Accumulator) {
	fields := make(map[string]interface{})
	tags := make(map[string]string)

	// get mapped field name
	mapped := process_map[name]

	// get process status
	if s, err := proc.Status(); err==nil {
		fields[mapped] = proc.Pid
	} else {
		fields[mapped] = 0
	}

	var cpuUsage float64
	var cpuNum = float64(1)
	if cu, err := proc.Percent(time.Duration(0)); err==nil {
		if info, err := cpu.Info(); err==nil {
			cpuNum = float64(len(info))
		}
		cpuUsage = cu/cpuNum
	}
	fields[mapped + "_cpu_usage"] = cpuUsage

	var memUsage = float64(0)
	var rss uint64 = 0
	if meminfo, err := proc.MemoryInfo(); err==nil {
		rss = meminfo.RSS
		if vm,err := mem.VirtualMemory(); err==nil {
			memUsage = float64(rss)/float64(vm.Total)
		}
	}
	fields[mapped + "_mem_rss"] = rss
	fields[mapped + "_mem_usage"] = memUsage

	tags["type"] = "system"
	acc.AddFields("elibot", fields, tags)
}

var process_map = map[string]string{
	APISERVER: "api_server",
	MCSERVER:  "mcserver",
	COBOT: "cobot",
	PARAMSERVER: "param_server",
}

func (m *McServer) Gather(acc telegraf.Accumulator) error {
	if monitorProcesses == nil {
		getMonitorProcesses()
	}

	for name,_ := range process_map {
		if proc, ok := monitorProcesses[name]; ok {
			GatherMonitoringProcessStat(name, proc, acc)
		}
	}

	return nil
}

func init() {
	inputs.Add("mcserver", func() telegraf.Input {
		return &McServer{}
	})
}