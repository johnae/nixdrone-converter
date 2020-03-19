package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func New() converter.Plugin {
	return &plugin{}
}

type plugin struct{}

func (p *plugin) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	if !strings.HasSuffix(req.Repo.Config, ".nix") {
		return nil, nil
	}

	tmpfile, err := ioutil.TempFile("", "pipeline.*.nix")
	if err != nil {
		logrus.Errorf("error creating temporary file: %+v", err)
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	config := req.Config.Data
	data := []byte(config)

	// Nix won't accept stdin unfortunately
	err = ioutil.WriteFile(tmpfile.Name(), data, 0644)
	if err != nil {
		logrus.Errorf("error writing to %s: %+v", tmpfile, err)
		return nil, err
	}

	cmd := exec.Command("nix", "eval", "-f", tmpfile.Name(), "--json", "")

	env := os.Environ()
	appendEnv := func(env []string, val interface{}) []string {
		v := reflect.ValueOf(val)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			envVar := fmt.Sprintf("DRONE_%s=%v", strcase.ToScreamingSnake(t.Field(i).Name), v.Field(i).Interface())
			env = append(env, envVar)
		}
		return env
	}
	env = appendEnv(env, req.Repo)
	env = appendEnv(env, req.Build)

	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("nix eval failed: %+v", err)
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(out, &m)
	if err != nil {
		logrus.Errorf("json unmarshal error: %+v", err)
		return nil, err
	}

	yml, err := yaml.Marshal(&m)
	if err != nil {
		logrus.Errorf("yaml marshal error: %+v", err)
		return nil, err
	}

	return &drone.Config{
		Data: string(yml),
	}, nil
}
