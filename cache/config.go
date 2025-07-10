package cache

type ICacheProvider interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetDB() int
}