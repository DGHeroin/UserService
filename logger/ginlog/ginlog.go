package ginlog

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "time"
)

func GetLogger(logger *logrus.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 开始时间
        startTime := time.Now()
        // 处理请求
        c.Next()
        // 结束时间
        endTime := time.Now()
        // 执行时间
        latencyTime := endTime.Sub(startTime)
        // 请求方式
        reqMethod := c.Request.Method
        // 请求路由
        reqUri := c.Request.RequestURI
        // 状态码
        statusCode := c.Writer.Status()
        // 请求IP
        clientIP := c.ClientIP()
        // 日志格式
        // logger.WithFields(logrus.Fields{
        //     "status_code":  statusCode,
        //     "latency_time": latencyTime,
        //     "client_ip":    clientIP,
        //     "req_method":   reqMethod,
        //     "req_uri":      reqUri,
        // }).Info(time.Now().Format("2006/01/02 15:04:05.000"))
        logger.Infof("%v| %v | %v| %v| %v| %v", time.Now().Format("2006/01/02 15:04:05.000"),
            statusCode,
            latencyTime,
            clientIP,
            reqMethod,
            reqUri,
        )
    }
}
