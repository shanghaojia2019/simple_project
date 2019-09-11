package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

//Log 日志log
var Log *logrus.Logger

//JSON 格式化json，替换golang自带的json包
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	Log = logrus.New()
}

//Info info日志写入文件
func Info(fields logrus.Fields, message string) {
	ctx := Log.WithFields(fields)
	ctx.Info(message)
}

//Error 错误日志写入文件
func Error(fields logrus.Fields, message string) error {
	ctx := Log.WithFields(fields)
	ctx.Error(message)
	return errors.New(message)
}

//NewWriter 写入文件
type NewWriter struct {
	*logrus.Entry
	Info bool
}

func (nw *NewWriter) Write(p []byte) (n int, err error) {
	if nw.Info {
		nw.Entry.Info(string(p))
		return os.Stdout.Write(p)
	}
	nw.Entry.Errorf(string(p))
	return os.Stderr.Write(p)
}

type bufferedWriter struct {
	gin.ResponseWriter
	Buffer *bytes.Buffer
}

func (g *bufferedWriter) WriteString(s string) (n int, err error) {
	g.Buffer.WriteString(s)
	return g.ResponseWriter.WriteString(s)
}

func (g *bufferedWriter) Write(data []byte) (int, error) {
	g.Buffer.Write(data)
	return g.ResponseWriter.Write(data)
}

//GinLog 日志格式定义
type GinLog struct {
	Print          bool
	Method         string
	ClientIP       string
	Path           string
	Start          time.Time
	End            time.Time
	Cost           int64
	RequestQuery   string
	RequestHeader  http.Header
	RequestBody    string
	ResponseCode   int
	ResponseHeader http.Header
	ResponseBody   string
	RequestFlag    string
}

//CreateDefaultServer 创建service API
func CreateDefaultServer(prof bool, rclient *redis.Client) *gin.Engine {
	Log.Formatter = &logrus.JSONFormatter{}
	out := Log.WithFields(logrus.Fields{"type": "gin-access"})
	gin.DefaultWriter = io.MultiWriter(out.Writer(), os.Stdout) //&NewWriter{Entry: out, Info: true} // io.MultiWriter(out.Writer(), os.Stdout) 同时输出到控制台和文件

	err := Log.WithFields(logrus.Fields{"type": "gin-error"})
	gin.DefaultErrorWriter = io.MultiWriter(err.Writer(), os.Stdout) //&NewWriter{Entry: err, Info: false}
	server := gin.New()
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	server.Use(JSONAccess(Log, rclient))
	server.Use(Cors())
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	if prof {
		ginpprof.Wrap(server)
	}
	return server
}

//Cors 跨域配置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//origin := c.Request.Header.Get("Origin")
		//var filterHost = [...]string{"http://localhost.*", "http://*.hfjy.com"}
		// filterHost 做过滤器，防止不合法的域名访问
		var isAccess = true
		// for _, v := range filterHost {
		// 	match, _ := regexp.MatchString(v, origin)
		// 	if match {
		// 		isAccess = true
		// 	}
		// }
		if isAccess {
			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept,Authorization") //允许自定义header
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")                              //允许请求方法
			// c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			// 	//c.Next()
			c.JSONP(http.StatusOK, "Options Request!")
			return
		}

		c.Next()
	}
}

//JSONAccess 日志格式定义
func JSONAccess(log *logrus.Logger, rclient *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			return
		}
		info := AccessLogger(c, rclient)
		if info.Print == false {
			return
		}
		fields := make(map[string]interface{})
		data, err := json.Marshal(info)
		if err == nil {
			err = json.Unmarshal(data, &fields)
		}
		fields["type"] = "gin-detail"
		tx := log.WithFields(logrus.Fields(fields))
		if err != nil {
			tx.Error("Log信息记录失败:", err)
		} else {
			tx.Info("")
		}
	}
}

//AccessLogger 中间件，记录所有的访问记录
func AccessLogger(c *gin.Context, rclient *redis.Client) *GinLog {
	buff := &bytes.Buffer{}
	newWriter := &bufferedWriter{c.Writer, buff}
	c.Writer = newWriter

	info := &GinLog{Print: true}
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys["gin-log"] = info
	// Start timer
	info.Start = time.Now()
	info.Path = c.Request.URL.Path
	info.Method = c.Request.Method
	info.ClientIP = c.ClientIP()
	info.RequestQuery = c.Request.URL.RawQuery
	info.RequestHeader = c.Request.Header

	body := ""
	if data, err := c.GetRawData(); err == nil {
		body = string(data)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	info.RequestBody = body

	if rclient != nil {
		info.RequestFlag = fmt.Sprintf("%s:%s", info.Path, time.Now().Format(time.RFC3339Nano))
		rclient.HSet("Request:Current", info.RequestFlag, "")
	}
	// Process request
	defer func() {
		if rclient != nil {
			rclient.HDel("Request:Current", info.RequestFlag)
		}
	}()
	c.Next()
	// Log only when path is not being skipped
	// Stop time
	if info.Print == false {
		return info
	}
	info.End = time.Now()
	info.Cost = info.End.Sub(info.Start).Nanoseconds()
	info.ResponseCode = c.Writer.Status()
	info.ResponseHeader = c.Writer.Header()
	info.ResponseBody = newWriter.Buffer.String()
	return info
}
