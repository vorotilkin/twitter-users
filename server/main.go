package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"twitter-users/infrastructure/repositories/user"
	"twitter-users/interfaces"
	"twitter-users/pkg/configuration"
	"twitter-users/pkg/database"
	pkgGrpc "twitter-users/pkg/grpc"
	"twitter-users/pkg/migration"
	"twitter-users/proto"
	"twitter-users/usecases"
)

type config struct {
	Grpc struct {
		Server pkgGrpc.Config
	}
	Db        database.Config
	Migration migration.Config
}

func newConfig(configuration *configuration.Configuration) (*config, error) {
	c := new(config)
	err := configuration.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	opts := []fx.Option{
		fx.Provide(zap.NewProduction),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(configuration.New),
		fx.Provide(newConfig),
		fx.Provide(func(c *config) pkgGrpc.Config {
			return c.Grpc.Server
		}),
		fx.Provide(func(c *config) database.Config {
			return c.Db
		}),
		fx.Provide(database.New),
		fx.Provide(func(c *config) migration.Config { return c.Migration }),
		fx.Provide(fx.Annotate(func(c *config) string { return c.Db.PostgresDSN() }, fx.ResultTags(`name:"dsn"`))),
		fx.Provide(fx.Annotate(pkgGrpc.NewServer,
			fx.As(new(grpc.ServiceRegistrar)),
			fx.As(new(interfaces.Hooker)))),
		fx.Provide(fx.Annotate(user.NewRepository, fx.As(new(usecases.UsersRepository)))),
		fx.Provide(fx.Annotate(usecases.NewUsersServer, fx.As(new(proto.UsersServer)))),
		fx.Invoke(func(lc fx.Lifecycle, server interfaces.Hooker) {
			lc.Append(fx.Hook{
				OnStart: server.OnStart,
				OnStop:  server.OnStop,
			})
		}),
		fx.Invoke(fx.Annotate(migration.Do, fx.ParamTags("", "", `name:"dsn"`))),
		fx.Invoke(proto.RegisterUsersServer),
	}

	app := fx.New(opts...)
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	<-app.Done()

	err = app.Stop(context.Background())
	if err != nil {
		panic(err)
	}
}
