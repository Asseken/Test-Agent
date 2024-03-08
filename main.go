package main

import (
	"Agent/ipinfo"
	"Agent/user"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"io"
	"net/http"
	"runtime"
	"time"
)

// 定义 SystemInfo 结构体来存储系统信息
type SystemInfo struct {
	// 内存信息
	Memory struct {
		Total       uint64  `json:"total"`        // 总内存
		Available   uint64  `json:"available"`    // 可用内存
		Used        uint64  `json:"used"`         // 已用内存
		UsedPercent float64 `json:"used_percent"` // 内存使用率
	} `json:"memory"`
	// CPU 信息
	CPU struct {
		TotalUsage   float64 `json:"total_usage"`  // CPU 总使用率
		Cores        int     `json:"cores"`        // CPU 核心数
		VendorID     string  `json:"vendor_id"`    // CPU 厂商
		Frequency    uint64  `json:"frequency"`    // CPU 频率
		NowFrequency uint64  `json:"Nowfrequency"` // CPU 当前频率
		ModelName    string  `json:"model_name"`   // CPU 型号
	} `json:"cpu"`
	// 磁盘信息
	Disk []struct {
		MountPoint  string  `json:"mount_point"`  // 挂载点
		Total       uint64  `json:"total"`        // 总容量
		Used        uint64  `json:"used"`         // 已用容量
		Free        uint64  `json:"free"`         // 可用容量
		UsedPercent float64 `json:"used_percent"` // 磁盘使用率
	} `json:"disk"`
	// 交换分区信息
	Swap struct {
		Total       uint64  `json:"total"`        // 交换分区总容量
		Available   uint64  `json:"available"`    // 交换分区可用容量
		Used        uint64  `json:"used"`         // 交换分区已用容量
		UsedPercent float64 `json:"used_percent"` // 交换分区使用率
	} `json:"swap"`
	// 开机时间及运行时间
	BootTime struct {
		StartTime string `json:"start_time"` // 开机时间
		RunTime   string `json:"runtime"`    // 运行时间
	} `json:"BootTime"`
	// 内核信息
	Kernel struct {
		SystemPlatform  string `json:"sysplatform"`      // 系统平台
		SystemStructure string `json:"sysstructure"`     // 系统架构
		KernelVersion   string `json:"kernel_version"`   // 内核版本
		Platform        string `json:"platform"`         // 平台信息
		PlatformFamily  string `json:"platform_family"`  // 平台族
		PlatformVer     string `json:"platform_version"` // 平台版本
	} `json:"kernel"`
	// IP 地址信息
	Ipinfo struct {
		Ipv4 string `json:"ipv4"` // IPv4 地址
		Ipv6 string `json:"ipv6"` // IPv6 地址
	} `json:"IpInfo"`
	GetIpInfo     []ipinfo.NetworkInfo `json:"Networkinfo"` // 获取网络信息
	TcpCount      int                  `json:"tcpCount"`    // TCP 连接数
	UdpCount      int                  `json:"udpCount"`    // UDP 连接数
	user.UserInfo `json:"UserInfo"`    // 用户信息
}

//	func main() {
//		var sysInfo SystemInfo
//		GetMem(&sysInfo)
//		cputime(&sysInfo)
//		GetCpu(&sysInfo)
//		GetSwap(&sysInfo)
//		GetDisk(&sysInfo)
//		boottime(&sysInfo)
//		kernel(&sysInfo)
//		//user.Userinfo()
//		sysInfo.UserInfo = user.Userinfo()
//		//ipinfo.Getinfonetwork()
//		sysInfo.GetIpInfo = ipinfo.Getinfonetwork()
//		//fmt.Println("公网IPv4地址:", getPublicIP("https://api.ipify.org"))
//		//fmt.Println("公网IPv6地址:", getPublicIP("https://api6.ipify.org"))
//		sysInfo.Ipinfo.Ipv4 = getPublicIP("https://api.ipify.org")
//		sysInfo.Ipinfo.Ipv6 = getPublicIP("https://api6.ipify.org")
//
//		// Marshal struct to JSON
//		jsonData, err := json.Marshal(sysInfo)
//		if err != nil {
//			fmt.Println("Error marshaling to JSON:", err)
//			return
//		}
//
//		// Print JSON data
//		fmt.Println(string(jsonData))
//	}
func main() {
	//http.HandleFunc("0.0.0.0/info", systemInfoHandler)
	//http.ListenAndServe(":1234", nil)
	http.HandleFunc("/api/sys", systemInfoHandler)

	// 监听在系统的 IP 地址上，端口为 8080
	// 如果您的系统有多个网络接口，请确保选择正确的 IP 地址
	ipAddress := "0.0.0.0" // 0.0.0.0 表示监听所有网络接口
	port := "2086"
	serverAddress := ipAddress + ":" + port
	fmt.Printf("Server listening on %s\n", serverAddress)
	http.ListenAndServe(serverAddress, nil)
}

// 处理系统信息请求的处理器函数
func systemInfoHandler(w http.ResponseWriter, r *http.Request) {
	var sysInfo SystemInfo // 假设 SystemInfo 结构体已定义

	// 填充 sysInfo 结构体的信息
	GetMem(&sysInfo)
	cputime(&sysInfo)
	GetCpu(&sysInfo)
	GetSwap(&sysInfo)
	GetDisk(&sysInfo)
	boottime(&sysInfo)
	kernel(&sysInfo)
	sysInfo.UserInfo = user.Userinfo()
	sysInfo.GetIpInfo = ipinfo.Getinfonetwork()
	sysInfo.Ipinfo.Ipv4 = getPublicIP("https://api.ipify.org")
	sysInfo.Ipinfo.Ipv6 = getPublicIP("https://api6.ipify.org")

	// 将结构体转换为 JSON 格式
	jsonData, err := json.Marshal(sysInfo)
	if err != nil {
		http.Error(w, "转换为 JSON 格式时出错", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 将 JSON 数据写入响应
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "写入 JSON 响应时出错", http.StatusInternalServerError)
		return
	}
}

// 获取内存信息
func GetMem(sysInfo *SystemInfo) {
	v, _ := mem.VirtualMemory()
	//totalGB := float64(v.Total) / 1024 / 1024 / 1024
	//availableGB := float64(v.Available) / 1024 / 1024 / 1024
	//usedGB := float64(v.Used) / 1024 / 1024 / 1024 // 已使用的内存（以GB为单位）
	sysInfo.Memory.Total = v.Total
	sysInfo.Memory.Available = v.Available
	sysInfo.Memory.Used = v.Used
	sysInfo.Memory.UsedPercent = v.UsedPercent
	//fmt.Printf("内存-总计: %.3f GB, 可用: %.3f GB, 已用: %.3f GB, 使用百分比: %.3f%%\n",
	//	totalGB, availableGB, usedGB, v.UsedPercent)
}

// 获取实时CPU使用率
func cputime(sysInfo *SystemInfo) {
	// 获取CPU实时总使用率
	totalUsage, err := cpu.Percent(time.Second, false)
	if err != nil || len(totalUsage) == 0 {
		fmt.Println("获取CPU实时总使用率时出错:", err)
		return
	}
	sysInfo.CPU.TotalUsage = totalUsage[0]
	//fmt.Printf("CPU 实时总使用率: %.2f%%\n", totalUsage[0])
}

// 获取cpu的信息
func GetCpu(sysInfo *SystemInfo) {
	info, err := cpu.Info()
	if err != nil {
		fmt.Println("获取CPU信息时出错:", err)
		return
	}
	for _, cpuInfo := range info {
		if runtime.GOOS == "windows" {
			//fmt.Println("CPU 核心数:", cpuInfo.Cores)
			sysInfo.CPU.Cores = int(cpuInfo.Cores)
		} else if runtime.GOOS == "linux" {
			// 计算唯一的CPU核心数
			coreSet := make(map[string]bool)
			for _, cpuInfo := range info {
				coreSet[cpuInfo.CoreID] = true
			}
			uniqueCores := len(coreSet)
			//fmt.Println("CPU 核心数:", uniqueCores)
			sysInfo.CPU.Cores = uniqueCores
		} else {
			fmt.Println("系统暂不支持")
		}

		//fmt.Println("CPU 厂商:", cpuInfo.VendorID)
		sysInfo.CPU.VendorID = cpuInfo.VendorID
		//fmt.Println("cpu 物理核心数:", cpuInfo.CPU)
		//fmt.Println("CPU 频率:", cpuInfo.Mhz, "MHz")
		sysInfo.CPU.Frequency = uint64(cpuInfo.Mhz)
		//fmt.Println("CPU 型号:", cpuInfo.ModelName)
		sysInfo.CPU.ModelName = cpuInfo.ModelName
		//fmt.Println("CPU 当前频率:", getCurrentCPUFrequency(), "MHz")
		sysInfo.CPU.NowFrequency = getCurrentCPUFrequency()
	}
}
func getCurrentCPUFrequency() uint64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return uint64(percent[0] * 1000) // Convert to MHz
}

// 获取交换分区信息
func GetSwap(sysInfo *SystemInfo) {
	swap, err := mem.SwapMemory()
	if err != nil {
		fmt.Println("获取交换分区信息时出错:", err)
		return
	}
	//totalGB := float64(swap.Total) / 1024 / 1024 / 1024
	//usedGB := float64(swap.Used) / 1024 / 1024 / 1024
	//freeGB := float64(swap.Free) / 1024 / 1024 / 1024
	//fmt.Printf("交换分区总计: %.3f GB, 已使用: %.3f GB, 剩余: %.3f GB, 使用百分比: %.3f%%\n",
	//	totalGB, usedGB, freeGB, swap.UsedPercent)
	sysInfo.Swap.Total = swap.Total
	sysInfo.Swap.Available = swap.Free
	sysInfo.Swap.Used = swap.Used
	sysInfo.Swap.UsedPercent = swap.UsedPercent

}

// 获取系统的硬盘信息，Windows和Linux下都可以使用 gopsutil 库来获取
// 获取系统的硬盘信息，仅收集容量大于0的分区
func GetDisk(sysInfo *SystemInfo) {
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Println("获取磁盘分区时出错:", err)
		return
	}
	for _, part := range parts {
		usageStat, err := disk.Usage(part.Mountpoint)
		if err != nil {
			fmt.Printf("获取磁盘 %s 使用情况时出错: %v\n", part.Mountpoint, err)
			continue
		}
		// 只收集容量大于0的分区
		if usageStat.Total > 0 {
			var diskInfo struct {
				MountPoint  string  `json:"mount_point"`
				Total       uint64  `json:"total"`
				Used        uint64  `json:"used"`
				Free        uint64  `json:"free"`
				UsedPercent float64 `json:"used_percent"`
			}
			diskInfo.MountPoint = part.Mountpoint
			diskInfo.Total = usageStat.Total
			diskInfo.Used = usageStat.Used
			diskInfo.Free = usageStat.Free
			diskInfo.UsedPercent = usageStat.UsedPercent

			sysInfo.Disk = append(sysInfo.Disk, diskInfo)
		}
	}
}

// 获取开机时间	运行时间
func boottime(sysInfo *SystemInfo) {
	bootTimestamp, _ := host.BootTime()
	bootTime := time.Unix(int64(bootTimestamp), 0)
	currentTime := time.Now()

	// 计算运行时间
	upTime := currentTime.Sub(bootTime)
	days := int(upTime.Hours()) / 24
	hours := int(upTime.Hours()) % 24
	minutes := int(upTime.Minutes()) % 60
	seconds := int(upTime.Seconds()) % 60
	sysInfo.BootTime.StartTime = bootTime.Local().Format("2006-01-02 15:04:05")
	//fmt.Println("开机时间:", bootTime.Local().Format("2006-01-02 15:04:05"))
	//fmt.Printf("运行时间: %d天 %02d小时 %02d分钟 %02d秒\n", days, hours, minutes, seconds)
	sysInfo.BootTime.RunTime = fmt.Sprintf("%d天 %02d小时 %02d分钟 %02d秒", days, hours, minutes, seconds)
}

// 内核版本 平台信息 系统架构amr 还是amd64
func kernel(sysInfo *SystemInfo) {
	infoStat, _ := host.Info()
	sysInfo.Kernel.SystemPlatform = runtime.GOOS
	sysInfo.Kernel.SystemStructure = runtime.GOARCH
	//fmt.Println("系统架构:", runtime.GOOS, runtime.GOARCH) //通过runtime
	sysInfo.Kernel.KernelVersion = infoStat.KernelVersion
	//fmt.Println("内核版本:", infoStat.KernelVersion)
	//fmt.Println("平台信息:", infoStat.Platform, infoStat.PlatformFamily, infoStat.PlatformVersion)
	sysInfo.Kernel.Platform = infoStat.Platform
	sysInfo.Kernel.PlatformFamily = infoStat.PlatformFamily
	sysInfo.Kernel.PlatformVer = infoStat.PlatformVersion
}

// 获取公网IPv4 ipv6地址
func getPublicIP(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取公网IP地址时出错:", err)
		return ""
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("获取公网IP地址时出错:", err)
		return ""
	}
	return string(ip)
}

// 获取系统upd tcp连接数
//func GetConnectionCount(sysInfo *SystemInfo) {
//	fmt.Println(sys.GetUDPCount())
//
//	fmt.Println(sys.GetTCPCount())
//}
