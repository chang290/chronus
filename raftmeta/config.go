package raftmeta

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
)

const (
	DefaultNumPendingProposals = 1000
	DefaultAddr                = "127.0.0.1:2347"
	DefaultElectionTick        = 100
	DefaultHeartbeatTick       = 1
	DefaultMaxSizePerMsg       = 4096
	DefaultMaxInflightMsgs     = 256
)

type IPRange struct {
	Lower, Upper net.IP
}

type Config struct {
	NumPendingProposals int    `toml:"num-pending-proposals"`
	Tracing             bool   `toml:"tracing"`
	MyAddr              string `toml:"my-addr"`
	Peers               []Peer `toml:"peers"`
	RaftId              uint64 `toml:"raft-id"`
	TickTimeMs          int    `toml:"tick-time-ms"`
	ElectionTick        int    `toml:"election-tick"`
	HeartbeatTick       int    `toml:"heartbeat-tick"`
	MaxSizePerMsg       uint64 `toml:"max-size-per-msg"`
	MaxInflightMsgs     int    `toml:"max-inflight-msgs"`

	WalDir              string `toml:"wal-dir"`
	SnapshotIntervalSec int    `toml:"snapshot-interval"`
	ChecksumIntervalSec int    `toml:"checksum-interval"`
	RetentionAutoCreate bool   `toml:"retention-auto-create"`
}

type Peer struct {
	Addr   string `toml:"addr"`
	RaftId uint64 `toml:"raft-id"`
}

// NewConfig returns an instance of Config with defaults.
func NewConfig() Config {
	return Config{
		NumPendingProposals: DefaultNumPendingProposals,
		Tracing:             false,
		MyAddr:              DefaultAddr,
		RaftId:              1,
		Peers: []Peer{
			{Addr: "127.0.0.1:2348", RaftId: 2},
			{Addr: "127.0.0.1:2349", RaftId: 3},
		},
		TickTimeMs:          20,
		ElectionTick:        DefaultElectionTick,
		HeartbeatTick:       DefaultHeartbeatTick,
		MaxSizePerMsg:       DefaultMaxSizePerMsg,
		MaxInflightMsgs:     DefaultMaxInflightMsgs,
		WalDir:              "./wal",
		SnapshotIntervalSec: 60,
		ChecksumIntervalSec: 10,
		RetentionAutoCreate: true,
	}
}

// FromTomlFile loads the config from a TOML file.
func (c *Config) FromTomlFile(fpath string) error {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	// Handle any potential Byte-Order-Marks that may be in the config file.
	// This is for Windows compatibility only.
	// See https://github.com/influxdata/telegraf/issues/1378 and
	// https://github.com/influxdata/influxdb/issues/8965.
	bom := unicode.BOMOverride(transform.Nop)
	bs, _, err = transform.Bytes(bom, bs)
	if err != nil {
		return err
	}
	return c.FromToml(string(bs))
}

// FromToml loads the config from TOML.
func (c *Config) FromToml(input string) error {
	_, err := toml.Decode(input, c)
	return err
}
