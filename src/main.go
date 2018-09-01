package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cnjack/throttle"
	"github.com/gin-gonic/gin"
)

func extractIPHeader(c *gin.Context) string {
	var requestorIP string
	if behindProxy {
		requestorIP = c.Request.Header.Get("X-Forwarded-For")
	} else {
		requestor := c.Request.RemoteAddr
		requestorSplit := strings.Split(requestor, ":")
		requestorIP = requestorSplit[0]
	}
	return requestorIP
}

func handleIPRequest(c *gin.Context) {
	requestorIP := extractIPHeader(c)
	c.JSON(http.StatusOK, map[string]string{"ip": requestorIP})
	return
}

func handleIPVerification(c *gin.Context) {
	expectedIP := c.Query("ip")
	if expectedIP == "" {
		c.String(http.StatusBadRequest, "Expected IP not provided")
		return
	}

	requestorIP := extractIPHeader(c)

	expectedIPRegex, err := regexp.Compile(expectedIP)
	if err != nil {
		c.String(http.StatusBadRequest, "There is a problem with your regexp.")
		return
	}

	if expectedIPRegex.MatchString(requestorIP) == true {
		c.JSON(http.StatusOK, map[string]bool{"result": true})
		return
	} else {
		c.JSON(http.StatusOK, map[string]bool{"result": false})
		return
	}
}

func handleHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "")
	return
}

var hostFlag string
var portFlag string
var throttleFlag int64
var behindProxy bool

func init() {
	flag.StringVar(&hostFlag, "host", "127.0.0.1", "Host the server should run on")
	flag.StringVar(&portFlag, "port", "9000", "Port the server should run on")
	flag.Int64Var(&throttleFlag, "throttle", 60, "Requests per minute allowed from IP")
	flag.BoolVar(&behindProxy, "proxy", false, "Whether the server is behind a proxy")
	flag.Parse()
}

func main() {
	router := gin.Default()
	router.Use(throttle.Policy(&throttle.Quota{
		Limit:  uint64(throttleFlag),
		Within: time.Minute,
	}))
	router.GET("/", handleIPRequest)
	router.GET("/health", handleHealthCheck)
	router.GET("/verify", handleIPVerification)
	connString := fmt.Sprintf("%s:%s", hostFlag, portFlag)
	router.Run(connString)
}
