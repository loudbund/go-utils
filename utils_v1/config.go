package utils_v1

import (
	"github.com/larspensjo/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

// 结构体1：
type uConfig struct {
	cfgCache     map[string]*config.Config // 配置数据缓存配置
	cfgCacheLock sync.RWMutex
}

// 对外函数1
func Config() *uConfig {
	return &uConfig{
		cfgCache: map[string]*config.Config{},
	}
}

// 对外函数2：读取字符串内容配置
func (u *uConfig) GetCfgString(cfgFile, section, option string) (value string, err error) {
	cfg := u.readCfgFile(cfgFile)
	return cfg.String(section, option)
}

// 对外函数3：读取int内容配置
func (u *uConfig) GetCfgInt(cfgFile, section, option string) (value int, err error) {
	cfg := u.readCfgFile(cfgFile)
	return cfg.Int(section, option)
}

// 内部函数1：读取默认配置
func (u *uConfig) readCfgFile(cfgFile string) *config.Config {
	u.cfgCacheLock.Lock()
	_, ok := u.cfgCache[cfgFile]
	u.cfgCacheLock.Unlock()

	if !ok {
		c, err1 := config.ReadDefault(cfgFile)
		if err1 != nil {
			log.Panic("读取配置文件失败:" + cfgFile)
		}
		u.cfgCacheLock.Lock()
		u.cfgCache[cfgFile] = c
		u.cfgCacheLock.Unlock()
	}
	return u.cfgCache[cfgFile]
}
