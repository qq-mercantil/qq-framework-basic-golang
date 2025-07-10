package cache

import (
	"go.uber.org/fx"
)

func CacheModule(serviceName string) fx.Option {
	return fx.Module(
        "cache",
        fx.Provide(fx.Annotate(
            func(config ICacheProvider) (*CacheClient, error) {
                return NewCacheClient(CacheOptions{
                    Config:      config,
                    ServiceName: serviceName,
                })
            }, 
            fx.As(new(ICacheClient)),
        )),
    )
}
