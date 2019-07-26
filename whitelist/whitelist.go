package whitelist

import (
	"log"
	"net"
	"strings"
)

var IpNetworkWhiteList *[]net.IPNet
var IpAddrWhiteList *[]string

func InWhitelist(ipaddr string) bool {
	log.Println("[debug] ipaddr", ipaddr)
	ip := net.ParseIP(ipaddr)

	if IpAddrWhiteList != nil {
		for _, item := range *IpAddrWhiteList {
			if item == ipaddr {
				return true

			}
		}
	}

	if IpNetworkWhiteList != nil && ip != nil {
		for _, item := range *IpNetworkWhiteList {
			if item.Contains(ip) {
				return true
			}
		}
	}

	return false
}

func Load(records []string) {
	var ipNetworkWhiteList []net.IPNet
	var ipAddrWhiteList []string

	empty := []string{
		"127.0.0.1",
		"0.0.0.0/0",
	}

	if len(records) == 0 {
		records = empty
	}

	for _, record := range records {
		if strings.Index(record, "/") != -1 {
			_, ipnet, err := net.ParseCIDR(record)
			if err != nil {
				log.Println("[error] net.ParseCIDR", record, err)
				continue
			} else {
				ipNetworkWhiteList = append(ipNetworkWhiteList, *ipnet)
			}
		} else {
			ipa := net.ParseIP(record)
			if ipa == nil {
				log.Println("[error] net.ParseIP", record)
				continue
			} else {
				ipAddrWhiteList = append(ipAddrWhiteList, record)
			}
		}
	}

	IpNetworkWhiteList = &ipNetworkWhiteList
	IpAddrWhiteList = &ipAddrWhiteList

	log.Println("[debug] network whitelist", IpNetworkWhiteList)
	log.Println("[debug] ipaddr whitelist", IpAddrWhiteList)
}
