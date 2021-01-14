package apollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type conf struct {
	env, appId, cluster, namespace, server string
}

var (
	defaultConf = &conf{
		cluster:   "default",
		namespace: "application",
	}
)

func StartWithMeta(c *MetaServerConfig) error {
	err := SetMetaServer(c)
	if err != nil{
		return err
	}
	appId, env, cluster := parseOsArgs()
	return startWithCluster(appId, env, cluster)
}

//从环境变量中读取app.id 及 env
func Start() error {
	appId, env, cluster := parseOsArgs()
	return startWithCluster(appId, env, cluster)
}

func start(appId, env string) error {
	return startWithCluster(appId, env, "default")
}

func startWithCluster(appId, env, cluster string) error {

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("recover: %v", err)
		}
	}()

	defaultConf.appId = appId
	defaultConf.env = env
	defaultConf.cluster = cluster
	if defaultConf.env != "" {

		defaultConf.server = metaServer.GetServer(defaultConf.env)

		//url, ok := metaServer[defaultConf.env]
		//if ok {
		//	defaultConf.server = url
		//}
	}

	if defaultConf.appId == "" {
		return fmt.Errorf("app.id not define")
	}

	if defaultConf.env == "" {
		return fmt.Errorf("env not define")
	}

	if defaultConf.cluster == "" {
		defaultConf.cluster = "default"
	}

	logger.Infof("start config with %+v", *defaultConf)

	server := configServer{}

	no := notify{
		notifications: make(map[string]int),
	}

	no.put(defaultConf.namespace, -1)

	config := &Config{
		conf:   defaultConf,
		server: &server,
		notify: &no,
		nCache: make(map[string]*cache),
	}

	//启动第一次获取配置
	err := server.updateServers(defaultConf)
	if err != nil {
		logger.Warnf("get meta servers fail ,try to get congfig from local, err: %v", err)
		err = loadFromLocal(config)
		if err != nil {
			logger.Errorf("get config from local fail, err: %v", err)
			return err
		}
	}

	//默认初始化 application 命名空间的配置
	err = config.updateConfig(defaultConf.namespace)
	if err != nil {
		logger.Warnf("updateConfig failed, err: %v\n", err)
		err = loadFromLocal(config)
		if err != nil {
			logger.Errorf("loadFromLocal failed, err: %v\n", err)
			return err
		}
	}
	go config.doNotify()
	go config.doUpdateMeta()
	defaultConfig = config
	return nil
}

func loadFromLocal(config *Config) error {
	f, err := os.Open(getFileName(config.conf))
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return unmarshalData(d, config, defaultConf.namespace)
}

type cf struct {
	AppId string `json:"app.id,omitempty"`
	Env   string `json:"env,omitempty"`
}

//从文件中读取app.id及env
//格式如下
//{
//"app.id":"SampleApp",
//"env":"DEV"
//}
func StartWithFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	res := &cf{}
	err = json.NewDecoder(f).Decode(&res)
	if err != nil {
		return err
	}
	return start(res.AppId, res.Env)
}
