// Package firewall defines a configurator used to change the state
// of the firewall as well as do some light routing changes.
package firewall

import (
	"context"
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/command"
)

var _ Configurator = (*Config)(nil)

// Configurator allows to change firewall rules and modify network routes.
type Configurator interface {
	Enabler
	VPNConnectionSetter
	PortAllower
	OutboundSubnetsSetter
}

type Config struct { //nolint:maligned
	runner           command.Runner
	logger           Logger
	iptablesMutex    sync.Mutex
	ip6tablesMutex   sync.Mutex
	defaultInterface string
	defaultGateway   net.IP
	localNetworks    []routing.LocalNetwork
	localIP          net.IP

	// Fixed state
	ipTables        string
	ip6Tables       string
	customRulesPath string

	// State
	enabled           bool
	vpnConnection     models.Connection
	vpnIntf           string
	outboundSubnets   []net.IPNet
	allowedInputPorts map[uint16]string // port to interface mapping
	stateMutex        sync.Mutex
}

// NewConfig creates a new Config instance and returns an error
// if no iptables implementation is available.
func NewConfig(ctx context.Context, logger Logger,
	runner command.Runner, defaultInterface string,
	defaultGateway net.IP, localNetworks []routing.LocalNetwork,
	localIP net.IP) (config *Config, err error) {
	iptables, err := findIptablesSupported(ctx, runner)
	if err != nil {
		return nil, err
	}

	return &Config{
		runner:            runner,
		logger:            logger,
		allowedInputPorts: make(map[uint16]string),
		ipTables:          iptables,
		ip6Tables:         findIP6tablesSupported(ctx, runner),
		customRulesPath:   "/iptables/post-rules.txt",
		// Obtained from routing
		defaultInterface: defaultInterface,
		defaultGateway:   defaultGateway,
		localNetworks:    localNetworks,
		localIP:          localIP,
	}, nil
}
