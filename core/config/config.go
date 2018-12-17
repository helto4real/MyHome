package config

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/helto4real/MyHome/core/contracts"
	c "github.com/helto4real/MyHome/core/contracts"
	o "github.com/helto4real/MyHome/helpers/os"
	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	osHelper c.IOS
}

// Open configuration from disk.
func (a *Configuration) Open() (*c.Config, error) {
	file, err := os.Open(a.ConfigPath("myhome.yaml"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return a.OpenReader(file)
}

// Open configuration from a reader.
func (a *Configuration) OpenReader(r io.Reader) (*c.Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return getRawConfig(data)
}

func getRawConfig(data []byte) (*c.Config, error) {
	//log.Print("\r\n", string(data))
	config := &c.Config{}
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (a *Configuration) ConfigPath(p string) string {
	homepath := a.osHelper.HomePath()
	if homepath != "" {
		config := path.Join(homepath, ".config")
		return path.Join(config, "myhome", p)
	}
	return ""
}

func NewConfiguration() *Configuration {
	return &Configuration{
		osHelper: o.NewOsHelper()}
}

func NewConfigurationMock(osHelper contracts.IOS) *Configuration {
	return &Configuration{
		osHelper: osHelper}
}
