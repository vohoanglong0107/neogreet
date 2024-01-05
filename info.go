package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/jaypipes/ghw"
	"golang.org/x/sys/unix"
)

type Info interface {
	CPU() string
	GPU() string
	OS() string
	Memory() string
	Disk() string
}

type SystemInfo struct {
}

func (info *SystemInfo) CPU() string {
	cpuInfo, err := ghw.CPU()
	if err != nil {
		return err.Error()
	}
	if len(cpuInfo.Processors) == 0 {
		return "No CPU found"
	}
	return cpuInfo.Processors[0].Model
}

func (info *SystemInfo) GPU() string {
	gpuInfo, err := ghw.GPU()
	if err != nil {
		return err.Error()
	}
	if len(gpuInfo.GraphicsCards) == 0 || gpuInfo.GraphicsCards[0].DeviceInfo == nil {
		return "No GPU found"
	}
	gpu := gpuInfo.GraphicsCards[0].DeviceInfo
	return gpu.Vendor.Name + " " + gpu.Product.Name
}

func (info *SystemInfo) OS() string {
	osRelease, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return err.Error()
	}
	r, err := regexp.Compile(`PRETTY_NAME="(.*)"`)
	if err != nil {
		return err.Error()
	}
	matches := r.FindStringSubmatch(string(osRelease))
	return strings.Join(matches[1:], "")
}

func (info *SystemInfo) Memory() string {
	memoryInfo, err := ghw.Memory()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%vMB", int64(math.Ceil(float64(memoryInfo.TotalUsableBytes)/float64(1024*1024))))
}

func (info *SystemInfo) Disk() string {
	var stat unix.Statfs_t
	unix.Statfs("/var/home", &stat)
	toHuman := func(block uint64) uint64 {
		return block * uint64(stat.Bsize) / 1024 / 1024 / 1024
	}
	available := toHuman(stat.Bavail)
	total := toHuman(stat.Blocks)
	used := total - available
	return fmt.Sprintf("%vGB / %vGB (%.2f%%)", used, total, float64(used)/float64(total)*100)
}

func (info *SystemInfo) Product() string {
	productName, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_name")
	if err != nil {
		return err.Error()
	}
	productVersion, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_version")
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(productName)) + " " + strings.TrimSpace(string(productVersion))
}
