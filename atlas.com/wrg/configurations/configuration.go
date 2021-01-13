package configurations

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

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

func GetConfiguration() (*Configuration, error) {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	c := &Configuration{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return c, nil
}

func (c Configuration) GetWorldConfiguration(index byte) (*WorldConfiguration, error) {
	w := &WorldConfiguration{}
	if len(c.Worlds) > 0 && index < byte(len(c.Worlds)) {
		w = &c.Worlds[index]
	}
	return w, nil
}
