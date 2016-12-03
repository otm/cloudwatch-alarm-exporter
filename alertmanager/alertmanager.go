package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// AlertPath is the path to post alarms to
	AlertPath = "/api/v1/alerts"
)

var (
	// AlertPathURL holds the base URL to the
	AlertPathURL *url.URL
)

func init() {
	var err error
	AlertPathURL, err = url.Parse(AlertPath)
	logFatalOnError(err)
}

// Alerter returns an Alerts slice
type Alerter interface {
	Alerts() Alerts
}

// Labels holds labels that uniqly identifies an alarm
type Labels map[string]string

// Annotations holds additional information for an alarm
type Annotations map[string]string

// Alert defines an alertmanager alarm
type Alert struct {
	Labels       Labels      `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorUrl"`
}

// Alerts implements the Alerter interface
func (a Alert) Alerts() Alerts {
	return Alerts{a}
}

// Alerts is a slice of alert
type Alerts []Alert

// Alerts implements the Alerter interface
func (a Alerts) Alerts() Alerts {
	return a
}

// AlertManager client
type AlertManager struct {
	URL             url.URL
	Alerter         Alerter
	RefreshInterval time.Duration
}

// New returns a new AlertManger
func New(URL url.URL, alerter Alerter, refreshInterval time.Duration) *AlertManager {
	a := AlertManager{
		URL:             URL,
		Alerter:         alerter,
		RefreshInterval: refreshInterval,
	}
	go a.start()
	return &a
}

func (am AlertManager) start() {
	if am.Alerter == nil {
		return
	}
	for {
		am.Trigger(am.Alerter)
		time.Sleep(time.Duration(am.RefreshInterval))
	}
}

// Trigger sends the Alert to the AlertManager
func (am AlertManager) Trigger(a Alerter) error {
	client := &http.Client{}

	b, err := json.Marshal(a.Alerts())
	if err != nil {
		return fmt.Errorf("error marshaling alerts: %v", err)
	}

	req, err := http.NewRequest("POST", am.URL.ResolveReference(AlertPathURL).String(), bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("error creating alert request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error posting alerts: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading alert response: %v", err)
	}

	log.Printf("Alert response: %s", body)
	return nil
}

func logFatalOnError(err error, prefix ...string) {
	if err == nil {
		return
	}

	if len(prefix) == 0 {
		log.Fatal(err)
	}

	log.Fatalf("%s: %v", strings.Join(prefix, ": "), err)

}
