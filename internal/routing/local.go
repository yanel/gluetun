package routing

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	ErrLinkLocalNotFound     = errors.New("local link not found")
	ErrSubnetDefaultNotFound = errors.New("default subnet not found")
	ErrSubnetLocalNotFound   = errors.New("local subnet not found")
)

type LocalNetwork struct {
	IPNet         *net.IPNet
	InterfaceName string
	IP            net.IP
}

type LocalSubnetGetter interface {
	LocalSubnet() (defaultSubnet net.IPNet, err error)
}

func (r *Routing) LocalSubnet() (defaultSubnet net.IPNet, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return defaultSubnet, fmt.Errorf("cannot list routes: %w", err)
	}

	defaultLinkIndex := -1
	for _, route := range routes {
		if route.Dst == nil {
			defaultLinkIndex = route.LinkIndex
			break
		}
	}
	if defaultLinkIndex == -1 {
		return defaultSubnet, fmt.Errorf("%w: in %d route(s)", ErrLinkDefaultNotFound, len(routes))
	}

	for _, route := range routes {
		if route.Gw != nil || route.LinkIndex != defaultLinkIndex {
			continue
		}
		defaultSubnet = *route.Dst
		r.logger.Info("local subnet found: " + defaultSubnet.String())
		return defaultSubnet, nil
	}

	return defaultSubnet, fmt.Errorf("%w: in %d routes", ErrSubnetDefaultNotFound, len(routes))
}

type LocalNetworksGetter interface {
	LocalNetworks() (localNetworks []LocalNetwork, err error)
}

func (r *Routing) LocalNetworks() (localNetworks []LocalNetwork, err error) {
	links, err := r.netLinker.LinkList()
	if err != nil {
		return localNetworks, fmt.Errorf("cannot list links: %w", err)
	}

	localLinks := make(map[int]struct{})

	for _, link := range links {
		if link.Attrs().EncapType != "ether" {
			continue
		}

		localLinks[link.Attrs().Index] = struct{}{}
		r.logger.Info("local ethernet link found: " + link.Attrs().Name)
	}

	if len(localLinks) == 0 {
		return localNetworks, fmt.Errorf("%w: in %d links", ErrLinkLocalNotFound, len(links))
	}

	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return localNetworks, fmt.Errorf("cannot list routes: %w", err)
	}

	for _, route := range routes {
		if route.Gw != nil || route.Dst == nil {
			continue
		} else if _, ok := localLinks[route.LinkIndex]; !ok {
			continue
		}

		var localNet LocalNetwork

		localNet.IPNet = route.Dst
		r.logger.Info("local ipnet found: " + localNet.IPNet.String())

		link, err := r.netLinker.LinkByIndex(route.LinkIndex)
		if err != nil {
			return localNetworks, fmt.Errorf("cannot find link at index %d: %w", route.LinkIndex, err)
		}

		localNet.InterfaceName = link.Attrs().Name

		ip, err := r.assignedIP(localNet.InterfaceName)
		if err != nil {
			return localNetworks, err
		}

		localNet.IP = ip

		localNetworks = append(localNetworks, localNet)
	}

	if len(localNetworks) == 0 {
		return localNetworks, fmt.Errorf("%w: in %d routes", ErrSubnetLocalNotFound, len(routes))
	}

	return localNetworks, nil
}
