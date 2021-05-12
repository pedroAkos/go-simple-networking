package neti

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration struct {
	iface string
	ip    string
	port  int

	buffSize int
}

func (c Configuration) String() string {
	return fmt.Sprintf("{iface: %v, ip: %v, port: %v, buffsize: %v}", c.iface, c.ip, c.port, c.buffSize)
}

func SetPFlags() {
	pflag.String("net.ip", "", "IP address to bind")
	pflag.String("net.iface", "lo0", "Network interface to bind")
	pflag.Int("net.port", 10000, "Port to bind")
}

func LoadConfiguration(viper *viper.Viper) Configuration {
	c := Configuration{}
	c.ip = viper.GetString("net.ip")
	c.port = viper.GetInt("net.port")
	c.iface = viper.GetString("net.iface")

	viper.SetDefault("net.buffSize", 1024)
	c.buffSize = viper.GetInt("net.buffSize")

	if c.ip == "" {
		var err error
		if c.ip, err = GetInterfaceIpv4Addr(c.iface); err != nil {
			panic(err)
		}
	}

	return c
}

func (c Configuration) Address() string {
	return fmt.Sprintf("%s:%d", c.ip, c.port)
}

func (c Configuration) BuffSize() int {
	return c.buffSize
}
