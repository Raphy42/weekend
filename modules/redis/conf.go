package redis

import (
	"context"
	"strings"

	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/core/config"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/pkg/std/slice"
)

const (
	ModeLocal    = "local"
	ModeSentinel = "sentinel"
	ModeCluster  = "cluster"
)

var (
	allModes = strings.Join(slice.New(ModeLocal, ModeCluster, ModeSentinel), ",")
)

var (
	ConfMode            = config.Key("redis", "mode")
	confModeMissing     = config.MissingKeyMessage(ConfMode)
	ConfAddr            = config.Key("redis", "addr")
	confAddrMissing     = config.MissingKeyMessage(ConfAddr)
	ConfDatabase        = config.Key("redis", "database")
	confDatabaseMissing = config.MissingKeyMessage(ConfDatabase)
	ConfUsername        = config.Key("redis", "username")
	confUsernameMissing = config.MissingKeyMessage(ConfUsername)
	ConfPassword        = config.Key("redis", "password")
	confPasswordMissing = config.MissingKeyMessage(ConfPassword)
)

type ServerConfiguration struct {
	Addr     string
	Database int
	Password string
	Username string
}

type Configuration struct {
	Mode    string
	Servers []ServerConfiguration
}

func localConfigFrom(ctx context.Context, conf config.Config) (*Configuration, error) {
	addr, err := conf.String(ctx, ConfAddr)
	if err != nil {
		return nil, stacktrace.Propagate(err, confAddrMissing("host address not found in config (host:port)"))
	}

	database, err := conf.Number(ctx, ConfDatabase, 0)
	if err != nil {
		return nil, stacktrace.Propagate(err, confDatabaseMissing("database not found in config (number)"))
	}

	username, err := conf.String(ctx, ConfUsername, "")
	if err != nil {
		return nil, stacktrace.Propagate(err, confUsernameMissing("username not found in config"))
	}

	password, err := conf.String(ctx, ConfPassword, "")
	if err != nil {
		return nil, stacktrace.Propagate(err, confPasswordMissing("password not found in config"))
	}

	return &Configuration{
		Mode: ModeLocal,
		Servers: slice.New(ServerConfiguration{
			Addr:     addr,
			Database: int(database),
			Password: password,
			Username: username,
		}),
	}, nil
}

func sentinelConfigFrom(ctx context.Context, conf config.Config) (*Configuration, error) {
	return nil, errors.NotImplemented("todo")
}

func clusterConfigFrom(ctx context.Context, conf config.Config) (*Configuration, error) {
	return nil, errors.NotImplemented("todo")
}

func ConfigFrom(ctx context.Context, conf config.Config) (*Configuration, error) {
	mode, err := conf.String(ctx, ConfMode, ModeLocal)
	if err != nil {
		return nil, stacktrace.Propagate(err, confModeMissing("redis mode not found in config (wanted one of [%s])", allModes))
	}

	switch mode {
	case ModeLocal:
		return localConfigFrom(ctx, conf)
	case ModeSentinel:
		return sentinelConfigFrom(ctx, conf)
	case ModeCluster:
		return clusterConfigFrom(ctx, conf)
	default:
		return nil, stacktrace.NewErrorWithCode(errors.EUnreachable, "no such redis mode: '%s'", mode)
	}
}
