// Package utils
/**
 * @Author: trying
 * @Description:
 * @File:  disk.go
 * @Version: 1.0.0
 * @Date: 2024/12/27 14:51
 */

package utils

import (
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStatus struct {
	All uint64 `json:"all"`

	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// DiskUsage disk usage of path/disk
func DiskUsage(path string) (diskStats DiskStatus) {
	if stat, err := disk.Usage(path); err == nil {
		diskStats.All = stat.Total
		diskStats.Free = stat.Free
		diskStats.Used = stat.Used
	}
	return
}
