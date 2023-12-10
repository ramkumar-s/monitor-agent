package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	// other necessary imports
)

// DataStruct represents the structure of the data to be collected.
type DataStruct struct {
	// define your fields here
	Pid       int32
	Name      string
	Cpu       float64
	Rss       uint64
	Vms       uint64
	ReadIO    uint64
	WriteIO   uint64
	Timestamp time.Time
}

func collectData(dataChan chan<- DataStruct) {

	for {
		pids, err := process.Pids()
		if err != nil {
			panic(err)
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
			name, err := p.Cmdline()
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

			data := DataStruct{
				Pid:       pid,
				Name:      name,
				Cpu:       cpuPercent,
				Rss:       memInfo.RSS,
				Vms:       memInfo.VMS,
				ReadIO:    ioCounters.ReadCount,
				WriteIO:   ioCounters.WriteCount,
				Timestamp: time.Now(), // Current time
			}

			// Send the data to the channel
			dataChan <- data

		}
		// Collect data and create a DataStruct instance

		// Wait for 10 seconds before the next collection
		time.Sleep(10 * time.Second)
	}
}

func writeToCSV(dataChan <-chan DataStruct) {
	// Open or create the CSV file
	file, err := os.OpenFile("data.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header to the CSV file
	header := []string{"Timestamp", "Pid", "Name", "Cpu", "Rss", "Vms", "ReadIO", "WriteIO"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}
	writer.Flush() // Ensure the header is written immediately

	for data := range dataChan {
		// Convert data to a slice of strings (or any format suitable for CSV)
		record := []string{
			data.Timestamp.Format(time.RFC3339),       // Convert time.Time to string
			strconv.FormatInt(int64(data.Pid), 10),    // Convert int32 to string
			data.Name,                                 // String can be used directly
			strconv.FormatFloat(data.Cpu, 'f', 2, 64), // Convert float64 to string with 2 decimal precision
			strconv.FormatUint(data.Rss, 10),          // Convert uint64 to string
			strconv.FormatUint(data.Vms, 10),          // Convert uint64 to string
			strconv.FormatUint(data.ReadIO, 10),       // Convert uint64 to string
			strconv.FormatUint(data.WriteIO, 10),      // Convert uint64 to string
		}
		// Write to CSV
		if err := writer.Write(record); err != nil {
			panic(err)
		}
		writer.Flush()
	}
}

func main() {
	dataChan := make(chan DataStruct)

	// Start the data collection goroutine
	go collectData(dataChan)

	// Start the CSV writing goroutine
	go writeToCSV(dataChan)

	// Keep the main goroutine alive
	select {}
}
