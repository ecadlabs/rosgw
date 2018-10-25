package config

import (
	"context"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/ecadlabs/rosgw/conn"
	"github.com/ecadlabs/rosgw/errors"
	"gopkg.in/yaml.v2"
)

const defaultPort = "22"

type Device struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	IdentityFile string `yaml:"identity_file"`
	Timeout      string `yaml:"timeout"`

	IdentityData []byte `yaml:"-"`
}

func (d *Device) Address() string {
	port := d.Port
	if port == "" {
		port = defaultPort
	}

	return net.JoinHostPort(d.Host, port)
}

func (d *Device) Config() *conn.Config {
	return &conn.Config{
		Username: d.Username,
		Password: d.Password,
		KeyFunc:  func() ([]byte, error) { return d.IdentityData, nil },
	}
}

func (d *Device) GetTimeout() time.Duration {
	t, _ := time.ParseDuration(d.Timeout)
	return t
}

type Config struct {
	Devices map[string]*Device `yaml:"devices"`
	Common  *Device            `yaml:"common"`
	Address string             `yaml:"address"`
	MaxConn int                `yaml:"max_connections"`

	mtx sync.Mutex `yaml:"-"`
}

var idCache = make(map[string][]byte)

func readIdentityFile(name string) ([]byte, error) {
	if data, ok := idCache[name]; ok {
		return data, nil
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	idCache[name] = data
	return data, nil
}

func (d *Device) loadIdentity() error {
	if d.IdentityFile == "" || d.IdentityData != nil {
		return nil
	}

	data, err := readIdentityFile(d.IdentityFile)
	if err != nil {
		return err
	}

	d.IdentityData = data
	return nil
}

func Load(name string) (*Config, error) {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	if err := c.Common.loadIdentity(); err != nil {
		return nil, err
	}

	for addr, dev := range c.Devices {
		if err := dev.loadIdentity(); err != nil {
			return nil, err
		}

		if dev.Host == "" {
			dev.Host = addr
		}

		if dev.Port == "" {
			dev.Port = c.Common.Port
		}

		if dev.Username == "" {
			dev.Username = c.Common.Username
		}

		if dev.Password == "" {
			dev.Password = c.Common.Password
		}

		if dev.Timeout == "" {
			dev.Timeout = c.Common.Timeout
		}

		if dev.IdentityData == nil {
			dev.IdentityData = c.Common.IdentityData
		}
	}

	return &c, nil
}

func (c *Config) GetDevice(ctx context.Context, id string) (*Device, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if dev, ok := c.Devices[id]; ok {
		return dev, nil
	}

	return nil, errors.ErrDeviceNotFound
}
