package rate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"io"

	"github.com/paulbellamy/ratecounter"
	"k8s.io/client-go/pkg/api/v1"
)

// WarningData stores info about last seen warning
type WarningData struct {
	LastEvent   *v1.Event
	RateCounter *ratecounter.RateCounter
}

// WarningRates stores info per namespace
type WarningRates struct {
	counters map[string]WarningData
	debug    bool
}

// NewRateWarningCounter creates a new rate service
func NewRateWarningCounter() *WarningRates {
	debug := false
	if debugVar := os.Getenv("DEBUG"); debugVar != "" {
		debug = true
	}
	return &WarningRates{
		counters: make(map[string]WarningData),
		debug:    debug,
	}
}

// Register registers an event
func (w *WarningRates) Register(evt *v1.Event) {
	if strings.Compare(evt.Type, "Warning") != 0 {
		return
	}
	now := time.Now()
	lastMinute := now.Add(-1 * time.Minute)
	if evt.LastTimestamp.Time.After(lastMinute) {
		warningMessage := strings.Replace(strings.TrimSpace(evt.Message), "\n", "", -1)
		if w.debug {
			fmt.Printf("{\"now\" : \"%s\", \"namespace\" : \"%s\", \"message\" : \"%s\", \"ts\" : \"%s\"}\n",
				now.Format("20060102150405"), evt.Namespace, warningMessage, evt.LastTimestamp.Format("20060102150405"))
		}
		i, ok := w.counters[evt.Namespace]
		if !ok {
			i = WarningData{
				RateCounter: ratecounter.NewRateCounter(1 * time.Minute),
				LastEvent:   nil,
			}
			//w.counters[evt.Namespace] = i
		}
		i.RateCounter.Incr(1)
		i.LastEvent = evt
		w.counters[evt.Namespace] = i
	}
}

// Dump the content of the current rates
func (w *WarningRates) Dump(io io.Writer) {
	for k, v := range w.counters {

		if io == nil {
			if w.debug {
				b, _ := json.Marshal(v.LastEvent)
				fmt.Printf("{\"namespace\" : \"%s\", \"rate\" : %d , \"lastEvent\" : %s }\n", k, v.RateCounter.Rate(), string(b))
			}
		} else {
			b, _ := json.MarshalIndent(v.LastEvent, "", "    ")
			fmt.Fprintf(io, "{\"namespace\" : \"%s\", \"rate\" : %d , \"lastEvent\" : %s}", k, v.RateCounter.Rate(), string(b))
		}
	}
}

// GetWarningRate get the current rate for a project
func (w *WarningRates) GetWarningRate(namespace string) (*WarningData, bool) {
	if namespace, ok := w.counters[namespace]; ok {
		return &namespace, true
	}
	return nil, false
}
