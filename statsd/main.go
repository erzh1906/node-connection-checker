package statsd

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var queue = make(chan string, 100)

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return defaultValue
	}
}

func StatCount(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|c", metric, value)
}

func StatTime(metric string, took time.Duration) {
	queue <- fmt.Sprintf("%s:%d|ms", metric, took/1e6)
}

func StatGauge(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|g", metric, value)
}

func StatsdSender(StatsdAddress string) {
	for s := range queue {
		if conn, err := net.Dial("udp", StatsdAddress); err == nil {
			io.WriteString(conn, s)
			conn.Close()
		}
	}
}

func init() {
	StatsdAddress := getEnv("STATSD_ADDR", "127.0.0.1:8125")
	go StatsdSender(StatsdAddress)
}
