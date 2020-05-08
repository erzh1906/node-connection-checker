package main

import (
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/go-co-op/gocron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"node-connection-checker/statsd"
	"os"
	"strings"
	"time"
	"github.com/nu7hatch/gouuid"
)

var configFileName = flag.String("config", "/opt/node-checker/conf/scheduler.yml", "Path to configuration file")

type Config struct {
	CheckInterval uint64   `yaml:"check_interval"`
	CheckTimeout  int64    `yaml:"check_timeout_ms"`
	TargetList    []string `yaml:"targets"`
}

func DefaultConfig() Config {
	return Config{
		CheckInterval: 5,
		CheckTimeout:  1500,
		TargetList:    []string{"web-1", "web-2", "web-3"},
	}
}

func ReadConfig(configFileName string, config interface{}) error {
	configYaml, err := ioutil.ReadFile(configFileName)
	if err != nil {
		raven.CaptureError(err, nil)
		return fmt.Errorf("can't read file [%s] [%s]", configFileName, err.Error())
	}
	err = yaml.Unmarshal(configYaml, config)
	if err != nil {
		raven.CaptureError(err, nil)
		return fmt.Errorf("can't parse config file [%s] [%s]", configFileName, err.Error())
	}
	return nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return defaultValue
	}
}

func CheckTarget(Target string, CheckTimeout int64) {
	requestId, uuidErr := uuid.NewV4()
	if uuidErr != nil {
		raven.CaptureError(uuidErr, nil)
		fmt.Println(uuidErr)
	}
	log.Printf("CheckTarget: [%s] [%s] check started", Target, requestId.String())

	URL := "http://" + Target + ":20000"
	t := time.Now()
	Timeout := time.Duration(CheckTimeout) * time.Millisecond
	StatsdPrefix := getEnv("METRIC_PREFIX", "production")
	StatsdSource := strings.Split(getEnv("METRIC_SOURCE", "server"), ".")[0]
	StatsdTarget := strings.Split(Target, ".")[0]

	client := &http.Client{Timeout: Timeout}

	req, reqErr := http.NewRequest("GET", URL, nil)
	if reqErr != nil {
		raven.CaptureError(reqErr, nil)
		log.Println(reqErr)
		return
	}
	req.Header.Set("Request-ID", requestId.String())
	resp, getErr := client.Do(req)

	if getErr != nil || resp.StatusCode != 200 {
		raven.CaptureError(getErr, nil)
		Counter := strings.Join([]string{StatsdPrefix, StatsdSource, StatsdTarget, "error.count"}, ".")
		statsd.StatCount(Counter, 1)
		log.Printf("CheckTarget: [%s] [%s] failed", URL, requestId)
		client.CloseIdleConnections()
		return
	}

	defer func() {
		resp.Body.Close()
		Timer := strings.Join([]string{StatsdPrefix, StatsdSource, StatsdTarget, "latency"}, ".")
		statsd.StatTime(Timer, time.Since(t))
	}()

	Counter := strings.Join([]string{StatsdPrefix, StatsdSource, StatsdTarget, "http_200.count"}, ".")

	statsd.StatCount(Counter, 1)
	log.Printf("CheckTarget: [%s] [%s] finished successfully", URL, requestId.String())
}

func CheckTargets(Targets []string, Timeout int64) {
	for _, Target := range Targets {
		go CheckTarget(Target, Timeout)
	}
}

func init()  {
	sentryDSN := getEnv("SENTRY_DSN", "https://123:abc@sentry.example.com:443/1")
	raven.SetDSN(sentryDSN)
}

func main() {
	flag.Parse()
	config := DefaultConfig()
	err := ReadConfig(*configFileName, &config)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Println(err)
		os.Exit(1)
	}
	s := gocron.NewScheduler(time.UTC)
	t := time.Date(2019, time.November, 10, 15, 0, 0, 0, time.UTC)
	s.Every(config.CheckInterval).Seconds().StartImmediately().StartAt(t).Do(CheckTargets, config.TargetList, config.CheckTimeout)
	log.Println("main: started")
	<-s.Start()
}
