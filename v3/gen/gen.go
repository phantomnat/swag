package gen

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/swaggo/swag/v3/openapi"
)

const (
	Name = "openapi-v3"
)

type Gen struct {
	debug  Debugger
	config *Config
}

// Debugger is the interface that wraps the basic Printf method.
type Debugger interface {
	Printf(format string, v ...interface{})
}

type Config struct {
	// InstanceName is used to get distinct names for different openapi documents in the
	// same project. The default value is "openapi-v3".
	InstanceName string

	// OutputDir represents the output directory for all the generated files
	OutputDir string

	// OutputTypes define types of files which should be generated
	OutputTypes []string
}

func (c *Config) applyDefault() *Config {
	if c == nil {
		c = &Config{}
	}
	if c.InstanceName == "" {
		c.InstanceName = Name
	}
	return c
}

func New(cfg *Config) *Gen {
	g := &Gen{
		debug:  log.New(os.Stdout, "gen-v3", log.LstdFlags),
		config: cfg.applyDefault(),
	}
	return g
}

func (g *Gen) Build() error {
	p := openapi.New()

	err := g.writeYAMLSwagger(g.config, p.GetOpenAPI())
	if err != nil {
		return err
	}
	return nil
}

func (g *Gen) writeYAMLSwagger(config *Config, doc *openapi3.T) error {
	var filename = "openapi-v3.yaml"

	if config.InstanceName != Name {
		filename = config.InstanceName + "_" + filename
	}

	yamlFileName := path.Join(config.OutputDir, filename)

	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	y, err := yaml.JSONToYAML(b)
	if err != nil {
		return fmt.Errorf("cannot covert json to yaml error: %s", err)
	}

	err = g.writeFile(y, yamlFileName)
	if err != nil {
		return err
	}

	g.debug.Printf("create openapi-v3.yaml at %+v", yamlFileName)

	return nil
}

func (g *Gen) writeFile(b []byte, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(b)

	return err
}
