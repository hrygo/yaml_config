### yaml_config usage

```go
package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hrygo/yaml_config"
)

func TestConfig(t *testing.T) {
	// CreateYamlFactory 创建一个yaml配置文件工厂
	// relativePath 相对工作目录的配置文件存储目录
	// fileName 配置文件名称
	// project 使用此模块的工程的名称（用以处理单元测试的路径问题）
	yc := yaml_config.CreateYamlFactory("test", "config.yml", "yaml_config")
	yc.ConfigFileChangeListen()

	word := yc.GetString("hello")
	assert.True(t, word == "word")
	word = yc.GetString("hello")
	assert.True(t, word == "word")

	cc := yc.Clone("test")
	cc.ConfigFileChangeListen()
	
	f := cc.GetFloat64("foo")
	assert.True(t, 1.0 == f)
	f = cc.GetFloat64("foo")
	assert.True(t, 1.0 == f)

}
```
