package configurations

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configurator struct {
	l logrus.FieldLogger
}

func NewConfigurator(l logrus.FieldLogger) *Configurator {
	return &Configurator{l}
}

type Configuration struct {
	Worlds []WorldConfiguration `yaml:"worlds"`
}

type WorldConfiguration struct {
	Name              string `yaml:"name"`
	Flag              string `yaml:"flag"`
	ServerMessage     string `yaml:"serverMessage"`
	EventMessage      string `yaml:"eventMessage"`
	WhyAmIRecommended string `yaml:"whyAmIRecommended"`
}

func (c *Configurator) GetConfiguration() (*Configuration, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		c.l.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		c.l.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return con, nil
}

func (c Configuration) GetWorldConfiguration(index byte) (*WorldConfiguration, error) {
	if len(c.Worlds) > 0 && index < byte(len(c.Worlds)) {
		w := &WorldConfiguration{}
		w = &c.Worlds[index]
		return w, nil
	}
	return nil, errors.New(fmt.Sprintf("Index out of bounds: %d", index))
}
