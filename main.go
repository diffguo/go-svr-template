package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-svr-template/common"
	"go-svr-template/common/log"
	"go-svr-template/models"
	"go-svr-template/views"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Config     common.Configure
	GWaitGroup sync.WaitGroup
)

const SERVERNAME = "TemplateServer"

func initConfig() error {
	var etcdHost string
	var configFilename string
	flag.StringVar(&configFilename, "c", "", "config file")
	flag.StringVar(&etcdHost, "etcd", "", "config file")
	ParseCommandLineParam()

	flag.Parse()
	var err error

	if configFilename != "" {
		err = common.LoadCfgFromFile(configFilename, &Config)
		if nil != err {
			fmt.Printf("get config from file(%s) err: %s \n", configFilename, err.Error())
			return err
		} else {
			fmt.Printf("get config from file ok: %v \n", Config)
		}
	}

	return err
}

func initDB(env string) error {
	dbConf := Config.MysqlSetting["MysqlInstance0"]
	_, err := models.InitGormDbPool(&dbConf, env != "online")
	if err != nil {
		return err
	}

	return nil
}

var httpSvr *http.Server
var ServerRunning = true

func main() {
	var env = os.Getenv("GOENV")
	if env == "" {
		env = "online"
	}

	env = strings.ToLower(env)
	fmt.Printf("Start %s In %s Env\n", SERVERNAME, env)

	err := initConfig()
	if nil != err {
		fmt.Println("initConfig err :", err)
		return
	}

	err = initDB(env)
	if nil != err {
		fmt.Println("initDB err :", err)
		return
	}

	_, err = log.InitLog(Config.LogSetting.LogDir, Config.LogSetting.LogFile, Config.LogSetting.LogLevel)
	if nil != err {
		fmt.Println("initLog err :", err)
		return
	}

	//启动用户订阅号
	go func() {
		views.InitUserWeChat()
	}()

	if CheckAndExecCmd() {
		fmt.Println("exec cmd finish!")
		return
	}

	log.Info("Init Finished，Going To StartServer")
	fmt.Println("Init Finished，Going To StartServer")

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(log.GLog.Log.Writer()))
	router.Use(common.GinLogger(3 * time.Second))
	router.Use(common.Cors())
	router.Use(common.CheckAuth())

	addRoute(router)
	log.Info("Run Server")

	go HandleSignal()
	go StartWork()

	GWaitGroup = sync.WaitGroup{}
	GWaitGroup.Add(1)

	httpSvr = &http.Server{Addr: Config.Listen, Handler: router}
	err = httpSvr.ListenAndServe()
	if err != nil {
		log.Errorf("%s", err.Error())
	}

	GWaitGroup.Wait()
}

func HandleSignal(signals ...os.Signal) {
	sig := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = append(signals, syscall.SIGTERM)
	}

	signal.Notify(sig, signals...)

	s := <-sig
	ServerRunning = false

	log.Infof("gin: graceful exit action from signal [%s]", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	httpSvr.Shutdown(ctx)
	cancel()

	log.Infof("gin: bye!")
}

