package mcserver

import (
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
	fields[name] = proc.Pid

	var cpuUsage float64
	var cpuNum = float64(1)
	if cu, err := proc.CPUPercent(); err==nil {
		if info, err := cpu.Info(); err==nil {
			cpuNum = float64(len(info))
		}
		cpuUsage = cu/cpuNum
	}
	fields[name + "_cpu_usage"] = cpuUsage

	var memUsage = float64(0)
	var rss uint64 = 0
	if meminfo, err := proc.MemoryInfo(); err==nil {
		rss = meminfo.RSS
		if vm,err := mem.VirtualMemory(); err==nil {
			memUsage = float64(rss)/float64(vm.Total)
		}
	}
	fields[name + "_mem_rss"] = rss
	fields[name + "_mem_usage"] = memUsage

	tags["type"] = "system"
	acc.AddFields("elibot", fields, tags)
}

func (m *McServer) Gather(acc telegraf.Accumulator) error {
	monitorProcesses := make(map[string]*process.Process)
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

	if apiServerProcess, ok := monitorProcesses[APISERVER]; ok {
		GatherMonitoringProcessStat("api_server", apiServerProcess, acc)
	}

	if mcServerProcess, ok := monitorProcesses[MCSERVER]; ok {
		GatherMonitoringProcessStat("mcserver", mcServerProcess, acc)
	}

	if cobotProcess, ok := monitorProcesses[COBOT]; ok {
		GatherMonitoringProcessStat("cobot", cobotProcess, acc)
	}

	if paramServerProcess, ok := monitorProcesses[PARAMSERVER]; ok {
		GatherMonitoringProcessStat("param_server", paramServerProcess, acc)
	}

	return nil
}

func init() {
	inputs.Add("mcserver", func() telegraf.Input {
		return &McServer{}
	})
} 
