package yaml_config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hrygo/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 由于 vipver 包本身对于文件的变化事件有一个bug，相关事件会被回调两次
// 常年未彻底解决，相关的 issue 清单：https://github.com/spf13/viper/issues?q=OnConfigChange
// 设置一个内部全局变量，记录配置文件变化时的时间点，如果两次回调事件事件差小于1秒，我们认为是第二次回调事件，而不是人工修改配置文件
// 这样就避免了 viper 包的这个bug

var lastChangeTime time.Time

type YmlConfig interface {
	ConfigFileChangeListen()
	Clone(fileName string) YmlConfig
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetBool(keyName string) bool
	GetInt(keyName string) int
	GetInt32(keyName string) int32
	GetInt64(keyName string) int64
	GetFloat64(keyName string) float64
	GetDuration(keyName string) time.Duration
	GetStringSlice(keyName string) []string
}

func init() {
	lastChangeTime = time.Now()
}

// CreateYamlFactory 创建一个yaml配置文件工厂
// relativePath 相对工作目录的配置文件存储目录
// fileName 配置文件名称
// project 使用此模块的工程的名称（用以处理单元测试的路径问题）
func CreateYamlFactory(relativePath string, fileName string, project string) YmlConfig {

	yamlConfig := viper.New()
	// 配置文件所在目录
	var basePath = BasePath(project)
	if len(relativePath) > 0 && "." != relativePath && "./" != relativePath {
		if strings.HasPrefix(relativePath, "./") {
			relativePath = fileName[1:]
		}
		basePath += relativePath
	}
	yamlConfig.AddConfigPath(basePath)
	var prefix string
	// 需要读取的文件名,默认为：config
	if len(fileName) == 0 {
		yamlConfig.SetConfigName("config")
		prefix = "config_"
	} else {
		yamlConfig.SetConfigName(fileName)
		prefix = fileName + "_"
	}
	// 设置配置文件类型(后缀)为 yaml(兼容yml)
	yamlConfig.SetConfigType("yaml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatalf("Config file init error：%v", err.Error())
	}
	v := CreateContainer(prefix)

	return &ymlLoader{
		viper: yamlConfig,
		mu:    new(sync.Mutex),
		c:     v,
	}
}

type ymlLoader struct {
	viper *viper.Viper
	mu    *sync.Mutex
	c     *containers
}

// ConfigFileChangeListen 监听文件变化
func (y *ymlLoader) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Milliseconds() >= 10 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
				log.Warnf("[YAML] Config file changed, reload!")
			}
		}
	})
	y.viper.WatchConfig()
}

// keyIsCache 判断相关键是否已经缓存
func (y *ymlLoader) keyIsCache(keyName string) bool {
	if _, exists := y.c.KeyIsExists(keyName); exists {
		return true
	} else {
		return false
	}
}

// 对键值进行缓存
func (y *ymlLoader) cache(keyName string, value interface{}) bool {
	// 避免瞬间缓存键、值时，程序提示键名已经被注册的日志输出
	y.mu.Lock()
	defer y.mu.Unlock()
	if _, exists := y.c.KeyIsExists(keyName); exists {
		return true
	}
	return y.c.Set(keyName, value)
}

// 通过键获取缓存的值
func (y *ymlLoader) getValueFromCache(keyName string) interface{} {
	return y.c.Get(keyName)
}

// 清空已经缓存的配置项信息
func (y *ymlLoader) clearCache() {
	y.c.FuzzyDelete()
}

// Clone 允许 clone 一个相同功能的结构体
func (y *ymlLoader) Clone(fileName string) YmlConfig {
	// 这里存在一个深拷贝，需要注意，避免拷贝的结构体操作对原始结构体造成影响
	var ymlC = *y
	var ymlConfViper = *(y.viper)
	(&ymlC).viper = &ymlConfViper
	(&ymlC).viper.SetConfigName(fileName)
	(&ymlC).c = CreateContainer(fileName)
	if err := (&ymlC).viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件Clone失败：%v", zap.Error(err))
	}
	return &ymlC
}

// Get 一个原始值
func (y *ymlLoader) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetString 字符串格式返回值
func (y *ymlLoader) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

// GetBool 布尔格式返回值
func (y *ymlLoader) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt 整数格式返回值
func (y *ymlLoader) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt32 整数格式返回值
func (y *ymlLoader) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt64 整数格式返回值
func (y *ymlLoader) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetFloat64 小数格式返回值
func (y *ymlLoader) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetDuration 时间单位格式返回值
func (y *ymlLoader) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetStringSlice 字符串切片数格式返回值
func (y *ymlLoader) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}

var basePath string

func BasePath(project string) string {
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		pl, cl := len(project), len(curPath)
		if pl != 0 && cl > pl && len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			i := strings.Index(curPath, project)
			if i > 0 {
				basePath = curPath[:i] + project
			}
		} else {
			basePath = curPath
		}
		return basePath + "/"
	} else {
		log.Fatalf("Running directory has no permission!")
	}
	return ""
}
