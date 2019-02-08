package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
	"github.com/jamiealquiza/envy"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var instance = flag.String("instance", "https://sipis.example.com", "Instance")
var every = flag.String("every", "15m", "Update time")
var myClient = &http.Client{Timeout: 10 * time.Second}

var (
	instanceCountMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_instance_count",
			Help: "Instance Count",
		},
		[]string{"instance"},
	)
	instanceIdleCountMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_instance_idle_count",
			Help: "Instance Idle Count",
		},
		[]string{"instance"},
	)
	instanceRegisteredCountMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_instance_registered_count",
			Help: "Instance Registered Count",
		},
		[]string{"instance"},
	)
	instanceRegisteringCountMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_instance_registering_count",
			Help: "Instance Registering Count",
		},
		[]string{"instance"},
	)
	instanceUnauthorizedCountMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_instance_unauthorized_count",
			Help: "Instance Unauthorized Count",
		},
		[]string{"instance"},
	)
	instanceErrorCountMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_instance_error_count",
			Help: "Instance Error Count",
		},
		[]string{"instance"},
	)
	serverUptimeInSecondsMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_server_uptime_in_seconds",
			Help: "Server Uptime in Seconds",
		},
		[]string{"instance"},
	)
	serverMessageLoopQueueSizeMeasurement = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sipis_server_message_loop_queue_size",
			Help: "Server Message Loop Queue Size",
		},
		[]string{"instance"},
	)
)

type Server struct {
	Name   	    					string
	VersionString       	string
	VersionNumber       	string
	Build 								string
	UptimeInSeconds 		 	float64
	MessageLoopQueueSize	float64
}

type CountInState struct {
	Idle 					float64
	Registered 		float64
	Registering 	float64
	Unauthorized	float64
	Error 				float64
}

type Instance struct {
	Count      		float64
	CountInState	CountInState
}

type Measurement struct {
	Server 	 Server
	Instance Instance
}

func init() {
	prometheus.MustRegister(instanceCountMeasurement)
	prometheus.MustRegister(instanceIdleCountMeasurement)
	prometheus.MustRegister(instanceRegisteredCountMeasurement)
	prometheus.MustRegister(instanceRegisteringCountMeasurement)
	prometheus.MustRegister(instanceUnauthorizedCountMeasurement)
	prometheus.MustRegister(instanceErrorCountMeasurement)
	prometheus.MustRegister(serverUptimeInSecondsMeasurement)
	prometheus.MustRegister(serverMessageLoopQueueSizeMeasurement)
}

func getMeasurement(instance string, target interface{}) error {
	var url = fmt.Sprintf("%s/stats/summary/json", instance)

	r, err := myClient.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func collectSample() {
	log.Println("Collecting sample...")
	sipisMeasurement := new(Measurement)
	getMeasurement(*instance, sipisMeasurement)

	instanceCountMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Instance.Count)
	instanceIdleCountMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Instance.CountInState.Idle)
	instanceRegisteredCountMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Instance.CountInState.Registered)
	instanceRegisteringCountMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Instance.CountInState.Registering)
	instanceUnauthorizedCountMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Instance.CountInState.Unauthorized)
	instanceErrorCountMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Instance.CountInState.Error)
	serverUptimeInSecondsMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Server.UptimeInSeconds)
	serverMessageLoopQueueSizeMeasurement.With(prometheus.Labels{"instance": *instance}).Set(sipisMeasurement.Server.MessageLoopQueueSize)
}

func main() {
	envy.Parse("SIPIS")
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())

	collectSample()
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %s", *every), collectSample)
	c.Start()

	log.Printf("Listening on %s!", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
