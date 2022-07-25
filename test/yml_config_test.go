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

func Test_ymlLoader_Viper(t *testing.T) {
  yc := yaml_config.CreateYamlFactory("test", "", "yaml_config")
  vip := yc.Viper()

  m := vip.GetStringMap("test.map1")
  key1 := m["key1"]
  key2 := m["key2"]
  key3 := m["key3"]

  assert.True(t, key1.(string) == "value1")
  assert.True(t, key2.(string) == "value2")
  assert.True(t, key3.(int) == 3)
}
