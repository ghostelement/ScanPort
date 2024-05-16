package scan

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Hosts       []string `yaml:"hosts"`
	Timeout     int      `yaml:"timeout,omitempty"`
	ParallelNum int      `yaml:"parallelNum,omitempty"`
}

// 解析yaml文件
func Config(p string) (*Job, error) {
	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	c := Job{}
	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// 验证配置文件
func (c *Job) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("hosts can't be empty")
	}
	return nil
}
