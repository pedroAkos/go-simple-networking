package neti

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration is the configuration for the network
type Configuration struct {
	iface string
	ip    string
	port  int

	buffSize int

	logger *logrus.Logger
}

// String returns a string representation of the configuration
func (c Configuration) String() string {
	return fmt.Sprintf("{iface: %v, ip: %v, port: %v, buffsize: %v}", c.iface, c.ip, c.port, c.buffSize)
}

// SetPFlags sets the pflags for the configuration
func SetPFlags() {
	pflag.String("net.ip", "", "IP address to bind")
	pflag.String("net.iface", "lo0", "Network interface to bind")
	pflag.Int("net.port", 10000, "Port to bind")
	pflag.String("net.loglvl", "info", "Log Level for network")
}

// LoadConfiguration loads the configuration from viper
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

	c.logger = logrus.New()
	level := viper.GetString("net.loglvl")
	switch level {
	case "debug":
		c.logger.Level = logrus.DebugLevel
	}

	return c
}

// Address returns the addresses for the configuration
func (c Configuration) Address() string {
	return fmt.Sprintf("%s:%d", c.ip, c.port)
}

// BuffSize returns the buffer size for the configuration (used for UDP connections)
func (c Configuration) BuffSize() int {
	return c.buffSize
}

// WithPort returns a new configuration with the port changed
func (c Configuration) WithPort(port int) string {
	return fmt.Sprintf("%v:%v", c.ip, port)
}
