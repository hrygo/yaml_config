package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hrygo/yaml_config"
)

func TestConfig1(t *testing.T) {
	yc := yaml_config.CreateYamlFactory("", "test.yaml", "yaml_config")
	yc.ConfigFileChangeListen()

	word := yc.GetString("hello")
	assert.True(t, word == "word")

	word = yc.GetString("hello")
	assert.True(t, word == "word")
}

func TestConfig2(t *testing.T) {
	yc := yaml_config.CreateYamlFactory("test", "", "yaml_config")
	yc.ConfigFileChangeListen()

	word := yc.GetString("hello")
	assert.True(t, word == "word")

	word = yc.GetString("hello")
	assert.True(t, word == "word")
}

func TestConfig3(t *testing.T) {
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

func TestConfig4(t *testing.T) {
	yc := yaml_config.CreateYamlFactory("", "test.yaml", "yaml_config")
	yc.ConfigFileChangeListen()

	word := yc.GetFloat64("foo")
	assert.True(t, 1.0 == word)
	word = yc.GetFloat64("foo")
	assert.True(t, 1.0 == word)
}
