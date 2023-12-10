package main

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessDetail struct {
	Pid        int32
	Name       string
	Cpu        float64
	Rss        uint64
	Vms        uint64
	ReadIO     uint64
	WriteIO    uint64
	LocalIP    string
	RemoteIP   string
	LocalPort  uint32
	RemotePort uint32
}

func PollProcess() {
	// Get a list of all running processes
	pids, err := process.Pids()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pid := range pids {
		// Skip pid 0 as it is reserved and won't have corresponding process details
		if pid == 0 {
			continue
		}

		// Get process details
		p, err := process.NewProcess(pid)
		if err != nil {
			fmt.Printf("Error while fetching process %d: %s\n", pid, err)
			continue
		}

		// Get process name
		name, err := p.Name()
		if err != nil {
			fmt.Printf("Error while fetching name for process %d: %s\n", pid, err)
			continue
		}

		// Get CPU percent
		cpuPercent, err := p.CPUPercent()
		if err != nil {
			fmt.Printf("Error while fetching CPU percent for process %d: %s\n", pid, err)
			continue
		}

		// Get memory info
		memInfo, err := p.MemoryInfo()
		if err != nil {
			fmt.Printf("Error while fetching memory info for process %d: %s\n", pid, err)
			continue
		}

		// Get IO counters
		ioCounters, err := p.IOCounters()
		if err != nil {
			fmt.Printf("Error while fetching IO counters for process %d: %s\n", pid, err)
			continue
		}

		// Get network connections
		conns, err := p.Connections()
		if err != nil {
			fmt.Printf("Error while fetching connections for process %d: %s\n", pid, err)
			continue
		}

		// Print process information
		fmt.Printf("Process ID: %d\n", pid)
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("CPU Percent: %.2f%%\n", cpuPercent)
		fmt.Printf("Memory: %d bytes (RSS), %d bytes (VMS)\n", memInfo.RSS, memInfo.VMS)
		fmt.Printf("IO Counters: ReadCount: %d, WriteCount: %d\n", ioCounters.ReadCount, ioCounters.WriteCount)

		// Print network information
		for _, conn := range conns {
			fmt.Printf("Network Connection: Local Address: %s, Remote Address: %s\n", conn.Laddr.IP, conn.Raddr.IP)
			fmt.Printf("Ports used: Local: %d, Remote: %d\n", conn.Laddr.Port, conn.Raddr.Port)
		}

		fmt.Println("---------------------------------------------------------")
	}
}
