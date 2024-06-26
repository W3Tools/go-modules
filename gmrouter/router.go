package gmrouter

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ApiContext *gin.Context
}

var routerEngine *gin.Engine
var debugMode bool

const (
	readTimeout  = time.Second * 60
	writeTimeout = time.Second * 60
)

func InitRouter(basePath string, debug bool) *gin.RouterGroup {
	debugMode = debug
	mode := gin.ReleaseMode
	loggerHandler := func(*gin.Context) {}

	if debugMode {
		mode = gin.DebugMode
		loggerHandler = gin.Logger()
	}

	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(crosSet)
	r.Use(loggerHandler)
	r.NoRoute(noRouteSet)
	r.NoMethod(methodNotAllowed)

	r.GET("/health", healthCheck)
	r.GET("/favicon.ico", faviconIcon)

	group := r.Group(basePath)

	if routerEngine == nil {
		routerEngine = r
	}

	return group
}

func RunningApi(host string, port int64) error {
	if reflect.ValueOf(host).IsZero() {
		host = "127.0.0.1"
	}

	if reflect.ValueOf(port).IsZero() {
		port = 8080
	}
	address := fmt.Sprintf("%s:%d", host, port)
	httpServer := &http.Server{
		Addr:           address,
		Handler:        routerEngine,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	if debugMode {
		ginpprof.Wrap(routerEngine)
	}

	fmt.Printf("service run at http://%s\n", address)

	err := httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("httpServer listen err %v", err)
	}
	return nil
}

func crosSet(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Apitoken, Authorization, Token")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

func healthCheck(ctx *gin.Context) {
	r := Router{ApiContext: ctx}
	r.ApiResponseOk("ok")
}

func noRouteSet(ctx *gin.Context) {
	r := Router{ApiContext: ctx}
	r.ApiResponse(http.StatusNotFound, "404 page not found.")
}

func methodNotAllowed(ctx *gin.Context) {
	r := Router{ApiContext: ctx}
	r.ApiResponse(http.StatusMethodNotAllowed, "Method not allowed.")
}

func faviconIcon(ctx *gin.Context) {
	data, err := os.ReadFile("./static/favicon.ico")
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.Header("Content-Type", "image/x-icon")
	ctx.Data(http.StatusOK, "image/x-icon", data)
}
