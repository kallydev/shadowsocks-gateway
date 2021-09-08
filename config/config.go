package config

import (
	"errors"
	"math"
	"net"
	"strconv"
	"strings"
)

var (
	ErrInvalidPort = errors.New("invalid port")
)

type Config struct {
	Bind   *Bind  `yaml:"bind"`
	Remote string `yaml:"remote"`
}

type Bind struct {
	IP      string `yaml:"ip"`
	Port    string `yaml:"port"`
	Network string `yaml:"network"`
}

func (config *Config) IP() net.IP {
	ip := net.ParseIP(config.Bind.IP)
	if ip == nil {
		panic("")
	}
	return ip
}

func (config *Config) Port() (start, end int) {
	splits := strings.Split(config.Bind.Port, "-")
	if len(splits) == 0 || len(splits) > 2 {
		panic(ErrInvalidPort)
	} else if len(splits) == 2 {
		start, err := strconv.Atoi(splits[0])
		if err != nil || start > math.MaxUint16 {
			panic(ErrInvalidPort)
		}
		end, err := strconv.Atoi(splits[1])
		if err != nil || end > math.MaxUint16 {
			panic(ErrInvalidPort)
		}
		return start, end
	} else {
		port, err := strconv.Atoi(config.Bind.Port)
		if err != nil || port > math.MaxUint16 {
			panic(ErrInvalidPort)
		}
		return port, port
	}
}
