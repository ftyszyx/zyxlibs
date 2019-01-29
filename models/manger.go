package models

import "github.com/ftyszyx/libs/beego/logs"

var allModels map[string]ModelInterface //存储所有的数据

func GetModel(modelname string) ModelInterface {
	return allModels[modelname]
}

func GetAllModel() map[string]ModelInterface {
	return allModels
}

func ResetAllModel() {
	allModels = make(map[string]ModelInterface)
}

//刷新
func RefrshCache(modelname string) {
	model := GetModel(modelname)
	if model != nil {
		logs.Info("clear  cache:%s", modelname)
		model.ClearCache()
	}
}

func RefrshAllCache() {
	logs.Info("clear all cache")
	for _, value := range allModels {
		value.ClearCache()
	}
}
