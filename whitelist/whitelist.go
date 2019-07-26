package whitelist

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/strings"
	"log"
	"net"
)

var IpNetworkWhiteList *[]net.IPNet
var IpAddrWhiteList *[]string

func InWhitelist(ipaddr string) bool {
	for _, item := range *IpNetworkWhiteList {
		ip := net.ParseIP(ipaddr)
		if item.Contains(ip) {
			return true
		}
	}

	for _, item := range *IpAddrWhiteList {
		if item == ipaddr {
			return true

		}
	}

	return false
}

func Load(records []string) {
	var ipNetworkWhiteList  []net.IPNet
	var ipAddrWhiteList []string

	empty := []string{
		"0.0.0.0/0",
	}

	if len(records) == 0 {
		records = empty
	}

	for _ , record := range records {
		if strings.Index(record, "/") != -1 {
			_, ipnet, err := net.ParseCIDR(record)
			if err != nil {
				log.Println("[error] net.ParseCIDR", record,err)
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



}
