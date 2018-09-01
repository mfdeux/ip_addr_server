package main

import (
	"net/http"
	"regexp"
	"strings"
	"time"

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
	}
	c.JSON(http.StatusOK, map[string]bool{"result": false})
	return
}

func handleStatsRequest(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]int{"last_hour": getIntervalSums(5, time.Duration(1)*time.Hour), "last_24h": getIntervalSums(5, time.Duration(24)*time.Hour), "last_7d": getIntervalSums(5, time.Duration(7*24)*time.Hour), "last_30d": getIntervalSums(5, time.Duration(30*24)*time.Hour)})
	return
}

func handleHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "")
	return
}
