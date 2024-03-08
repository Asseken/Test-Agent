package ipinfo

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

type NetworkInfo struct {
	MACAddress string `json:"Mac"`
	IPAddress  string `json:"ip"`
	Name       string `json:"Networkname"`
	Address    string `json:"addr"`
	Gateway    string `json:"gateway"`
}

func Getinfonetwork() []NetworkInfo {
	var networkInfos []NetworkInfo

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("错误:", err)
		os.Exit(1)
	}

	for _, iface := range interfaces {
		var networkInfo NetworkInfo
		networkInfo.Name = iface.Name
		networkInfo.MACAddress = iface.HardwareAddr.String()

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("错误:", err)
			continue
		}

		for _, addr := range addrs {
			networkInfo.Address = addr.String()
			break
		}

		if runtime.GOOS == "linux" {
			gateways, err := net.InterfaceAddrs()
			if err != nil {
				fmt.Println("错误:", err)
				continue
			}

			for _, gateway := range gateways {
				ipNet, ok := gateway.(*net.IPNet)
				if !ok {
					continue
				}

				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					networkInfo.Gateway = ipNet.IP.String()
					break // Get only the first gateway address
				}
			}
		}

		networkInfos = append(networkInfos, networkInfo)
	}

	return networkInfos
}
