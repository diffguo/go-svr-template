// @title Swagger go-svr-template API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/diffguo/gocom"
	"github.com/diffguo/gocom/cache"
	"github.com/diffguo/gocom/log"
	"github.com/gin-gonic/gin"
	"go-svr-template/apis"
	"go-svr-template/docs"
	"go-svr-template/models"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Config     gocom.Configure
	GWaitGroup sync.WaitGroup
)

const SERVERNAME = "go-svr-templateServer"
const GinAuthKey = "GinAuthKeyContent123456789012345" // 32字节长

func initConfig(confFilePath string) error {
	var etcdHost string
	var configFilename string
	flag.StringVar(&configFilename, "c", "", "config file")
	flag.StringVar(&etcdHost, "etcd", "", "config file")
	ParseCommandLineParam()

	flag.Parse()
	var err error

	if configFilename == "" {
		configFilename = confFilePath
	}

	if configFilename != "" {
		err = gocom.LoadCfgFromFile(configFilename, &Config)
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
	dbConf := Config.MysqlSetting["MysqlInstance"]
	_, err := models.InitGormDbPool(&dbConf, env != GoEnvOnline)
	if err != nil {
		return err
	}

	return nil
}

func initCache() {
	redisConf := Config.RedisSetting["RedisInstance"]
	cache.InitRedis([]string{redisConf.RedisConn}, redisConf.RedisPassword, redisConf.RedisDb, redisConf.PoolSize,
		redisConf.MinIdleConnNum, 3, false)
}

func swaggerInfo() {
	docs.SwaggerInfo.Title = SERVERNAME + " API"
	docs.SwaggerInfo.Description = "请使用前全局替换go-svr-template，以及更换SERVERNAME"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

var httpSvr *http.Server
var ServerRunning = true

const GoEnvOnline = "online"
const GoEnvTest = "test"

func main() {
	swaggerInfo()

	err := initConfig("")
	if nil != err {
		fmt.Println("initConfig err :", err)
		return
	}

	env := Config.Environment
	if env == "" {
		env = GoEnvTest
	}

	env = strings.ToLower(env)
	fmt.Printf("Start %s In %s Env\n", SERVERNAME, env)

	err = initDB(env)
	if nil != err {
		fmt.Println("initDB err :", err)
		return
	}

	initCache()

	err = gocom.InitGinAuth(GinAuthKey)
	if nil != err {
		fmt.Println("InitGinAuth err :", err)
		return
	}

	_, err = log.InitLog(Config.LogSetting.LogDir, Config.LogSetting.LogFile, Config.LogSetting.LogLevel, Config.LogSetting.LogSize)
	if nil != err {
		fmt.Println("initLog err :", err)
		return
	}

	//启动用户订阅号
	go func() {
		apis.InitUserWeChat()
	}()

	if CheckAndExecCmd() {
		fmt.Println("exec cmd finish!")
		return
	}

	log.Info("Init Finished，Going To StartServer")
	fmt.Println("Init Finished，Going To StartServer")

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(log.GLog.Log.Writer()))
	router.Use(gocom.GinLogger(3 * time.Second))
	router.Use(gocom.Cors())
	router.Use(gocom.CheckAuth())

	addRoute(router)
	log.Info("Run Server")

	go HandleSignal()
	go StartWork()

	GWaitGroup = sync.WaitGroup{}
	GWaitGroup.Add(1)

	httpSvr = &http.Server{Addr: Config.Listen,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = httpSvr.ListenAndServe()
	if err != nil {
		log.Errorf("%s", err.Error())
	}

	GWaitGroup.Wait()
}

// 优雅退出，以后可以考虑优雅重启
func HandleSignal(signals ...os.Signal) {
	sig := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	}

	signal.Notify(sig, signals...)

	s := <-sig
	ServerRunning = false

	log.Infof("io: graceful exit action from signal [%s]", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	httpSvr.Shutdown(ctx)
	cancel()

	log.Infof("io: bye!")
}
