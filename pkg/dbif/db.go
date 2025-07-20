package dbif

import (
	"fmt"

	"github.com/abc1763613206/nabili/pkg/cdn"
	"github.com/abc1763613206/nabili/pkg/geoip"
	"github.com/abc1763613206/nabili/pkg/ip2location"
	"github.com/abc1763613206/nabili/pkg/ip2region"
	"github.com/abc1763613206/nabili/pkg/ipip"
	"github.com/abc1763613206/nabili/pkg/qqwry"
	"github.com/abc1763613206/nabili/pkg/zxipv6wry"
)

type QueryType uint

const (
	TypeIPv4 = iota
	TypeIPv6
	TypeDomain
)

type DB interface {
	Find(query string, params ...string) (result fmt.Stringer, err error)
	Name() string
}

var (
	_ DB = &qqwry.QQwry{}
	_ DB = &zxipv6wry.ZXwry{}
	_ DB = &ipip.IPIPFree{}
	_ DB = &geoip.GeoIP{}
	_ DB = &ip2region.Ip2Region{}
	_ DB = &ip2location.IP2Location{}
	_ DB = &cdn.CDN{}
)
