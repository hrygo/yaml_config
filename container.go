package yaml_config

import (
	"strings"
	"sync"

	"go.uber.org/zap"
)

// 定义一个全局键值对存储容器
var sMap sync.Map

// 需外部调用方将其初始化
var log *zap.Logger

// SetLogger 必须执行此操作设置logger
func SetLogger(l *zap.Logger) {
	log = l
}

// CreateContainer 创建一个容器工厂
func CreateContainer(prefix string) *containers {
	return &containers{prefix: prefix}
}

// 定义一个容器结构体
type containers struct {
	prefix string
}

// Set  1.以键值对的形式将代码注册到容器
func (c *containers) Set(key string, value interface{}) (res bool) {
	if _, exists := c.KeyIsExists(c.prefix + key); exists == false {
		sMap.Store(c.prefix+key, value)
		res = true
	} else {
		log.Sugar().Infof("key to set is exists：%s" + key)
	}
	return
}

// Delete  2.删除
func (c *containers) Delete(key string) {
	sMap.Delete(c.prefix + key)
}

// Get 3.传递键，从容器获取值
func (c *containers) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(c.prefix + key); exists {
		return value
	}
	return nil
}

// KeyIsExists 4. 判断键是否被注册
func (c *containers) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(c.prefix + key)
}

// FuzzyDelete 按照键的前缀模糊删除容器中注册的内容
func (c *containers) FuzzyDelete() {
	sMap.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, c.prefix) {
				sMap.Delete(keyName)
			}
		}
		return true
	})
}
