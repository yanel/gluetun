package resolver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"time"
)

type Repeat interface {
	Resolve(ctx context.Context, host string, settings RepeatSettings) (IPs []net.IP, err error)
}

type repeat struct {
	resolver *net.Resolver
}

func NewRepeat(address string) Repeat {
	return &repeat{
		resolver: newResolver(address),
	}
}

type RepeatSettings struct {
	MaxDuration     time.Duration
	BetweenDuration time.Duration
	MaxNoNew        int
	// Maximum consecutive DNS resolution failures
	MaxFails int
	SortIPs  bool
}

func (r *repeat) Resolve(ctx context.Context, host string, settings RepeatSettings) (ips []net.IP, err error) {
	timedCtx, cancel := context.WithTimeout(ctx, settings.MaxDuration)
	defer cancel()

	noNewCounter := 0
	failCounter := 0
	uniqueIPs := make(map[string]struct{})

	for err == nil {
		// TODO
		// - one resolving every 100ms for round robin DNS responses
		// - one every second for time based DNS cycling responses
		noNewCounter, failCounter, err = r.resolveOnce(ctx, timedCtx, host, settings, uniqueIPs, noNewCounter, failCounter)
	}

	if len(uniqueIPs) == 0 {
		return nil, err
	}

	ips = uniqueIPsToSlice(uniqueIPs)

	if settings.SortIPs {
		sort.Slice(ips, func(i, j int) bool {
			return bytes.Compare(ips[i], ips[j]) < 1
		})
	}

	return ips, nil
}

var (
	ErrMaxNoNew = errors.New("reached the maximum number of no new update")
	ErrMaxFails = errors.New("reached the maximum number of consecutive failures")
)

func (r *repeat) resolveOnce(ctx, timedCtx context.Context, host string,
	settings RepeatSettings, uniqueIPs map[string]struct{}, noNewCounter, failCounter int) (
	newNoNewCounter, newFailCounter int, err error) {
	IPs, err := r.lookupIPs(timedCtx, host)
	if err != nil {
		failCounter++
		if settings.MaxFails > 0 && failCounter == settings.MaxFails {
			return noNewCounter, failCounter, fmt.Errorf("%w: %d failed attempts resolving %s: %s",
				ErrMaxFails, settings.MaxFails, host, err)
		}
		// it's fine to fail some of the resolutions
		return noNewCounter, failCounter, nil
	}
	failCounter = 0 // reset the counter if we had no error

	anyNew := false
	for _, IP := range IPs {
		key := IP.String()
		if _, ok := uniqueIPs[key]; !ok {
			anyNew = true
			uniqueIPs[key] = struct{}{}
		}
	}

	if !anyNew {
		noNewCounter++
	}

	if settings.MaxNoNew > 0 && noNewCounter == settings.MaxNoNew {
		// we reached the maximum number of resolutions without
		// finding any new IP address to our unique IP addresses set.
		return noNewCounter, failCounter,
			fmt.Errorf("%w: %d times no updated for %d IP addresses found",
				ErrMaxNoNew, noNewCounter, len(uniqueIPs))
	}

	timer := time.NewTimer(settings.BetweenDuration)
	select {
	case <-timer.C:
		return noNewCounter, failCounter, nil
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return noNewCounter, failCounter, ctx.Err()
	case <-timedCtx.Done():
		if err := ctx.Err(); err != nil {
			// timedCtx was canceled from its parent context
			return noNewCounter, failCounter, err
		}
		return noNewCounter, failCounter,
			fmt.Errorf("reached the timeout: %w", timedCtx.Err())
	}
}

func (r *repeat) lookupIPs(ctx context.Context, host string) (ips []net.IP, err error) {
	addresses, err := r.resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	ips = make([]net.IP, 0, len(addresses))
	for i := range addresses {
		ip := addresses[i].IP
		if ip == nil {
			continue
		}
		ips = append(ips, ip)
	}
	return ips, nil
}
