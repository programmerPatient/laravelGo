/*
 * @Description:配置文件加载核心代码
 * @Author: mali
 * @Date: 2022-09-07 09:25:12
 * @LastEditTime: 2023-03-10 14:42:14
 * @LastEditors: VSCode
 * @Reference:
 */
package config

import (
	"os"

	"github.com/laravelGo/core/helper"
	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper"
)

// viper 库实例
var viper *viperlib.Viper

// ConfigFunc 动态加载配置信息的函数
type ConfigFunc func() map[string]interface{}

// ConfigFuncs 先加载到此数组，防止config比env先加载 ，loadConfig 再动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	//初始化viper
	viper = viperlib.New()
	// 配置类型，支持 "json", "toml", "yaml", "yml", "properties","props", "prop", "env", "dotenv"
	// 默认使用env格式
	viper.SetConfigType("env")
	//环境变量配置文件查找的路径，相对于 main.go
	viper.AddConfigPath(".")
	//设置环境变量前缀，用以区分 Go 的系统环境变量
	viper.SetEnvPrefix("appenv")
	//读取环境变量（支持 flags）
	viper.AutomaticEnv()
	//初始化ConfigFuncs
	ConfigFuncs = make(map[string]ConfigFunc)
}

func InitConfig() {
	//加载env配置文件
	loadEnv()
	//加载config
	loadConfig()
}

/**
 * @Author: mali
 * @Func:
 * @Description: 加载env配置文件
 * @Param:
 * @Return:
 * @Example:
 * @param {string} envPrefix 配置文件名前缀
 */
func loadEnv(envPrefix ...string) {
	// 默认加载 .env 文件，如果有传参 --env=name 的话，加载 .env.name 文件
	envPath := ".env"
	if len(envPrefix) > 0 {
		filepath := envPath + "." + envPrefix[0]
		//判断是否为合法的文件属性
		if _, err := os.Stat(filepath); err != nil {
			envPath = filepath
		}
	}
	// 加载 env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 监控 .env 文件，变更时重新加载
	viper.WatchConfig()
}

// Env 读取环境变量，支持默认值
func Env(name string, defaultValue ...interface{}) interface{} {
	return internalGet(name, defaultValue...)
}

//加载config
func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

//添加配置信息
func AddConfig(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

//获取配置
func internalGet(path string, defaultValue ...interface{}) interface{} {
	// config 或者环境变量不存在的情况
	if !viper.IsSet(path) || helper.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetFloat64 获取 float64 类型的配置信息
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetUint64 获取 uint6 类型的配置信息
func GetUint64(path string, defaultValue ...interface{}) uint64 {
	return cast.ToUint64(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}

// GetStringMapInt 获取结构数据
func GetStringMapInt(path string) map[string]int {
	map_string := viper.GetStringMapString(path)
	return_data := map[string]int{}
	for k, v := range map_string {
		return_data[k] = cast.ToInt(v)
	}
	return return_data
}
