package config_utils

import (
    "github.com/larspensjo/config"
    log "github.com/sirupsen/logrus"
)

// 全局变量1： 配置数据缓存配置
var cfgCache = map[string]*config.Config{}


// 函数1：读取默认配置
func readCfgFile(cfgFile string) *config.Config {
    if _, ok := cfgCache[cfgFile]; !ok {
        c, err1 := config.ReadDefault(cfgFile)
        if err1 != nil {
            log.Panic("读取配置文件失败:"+cfgFile)
        }
        cfgCache[cfgFile] = c
    }
    return cfgCache[cfgFile]
}

// 函数2：读取字符串内容配置
func GetCfgString(cfgFile, section, option string) (value string, err error) {
    cfg := readCfgFile(cfgFile)
    return cfg.String(section, option)
}

// 函数3：读取int内容配置
func GetCfgInt(cfgFile, section, option string) (value int, err error) {
    cfg := readCfgFile(cfgFile)
    return cfg.Int(section, option)
}
