package vpnunlimited

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.VPNUnlimitedServer

func (hts hostToServer) toHostsSlice() (hosts []string) {
	hosts = make([]string, 0, len(hts))
	for host := range hts {
		hosts = append(hosts, host)
	}
	return hosts
}

func (hts hostToServer) adaptWithIPs(hostToIPs map[string][]net.IP) {
	for host, IPs := range hostToIPs {
		server := hts[host]
		server.IPs = IPs
		hts[host] = server
	}
	for host, server := range hts {
		if len(server.IPs) == 0 {
			delete(hts, host)
		}
	}
}

func (hts hostToServer) toServersSlice() (servers []models.VPNUnlimitedServer) {
	servers = make([]models.VPNUnlimitedServer, 0, len(hts))
	for _, server := range hts {
		servers = append(servers, server)
	}
	return servers
}
