package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AlertManagerRequest struct {
	Version         string      `json:"version"`
	GroupKey        string      `json:"groupKey"`
	TruncatedAlerts int         `json:"truncatedAlerts"`
	AlertStatus     AlertStatus `json:"alertStatus"`
	Receiver        string      `json:"receiver"`
	GroupLabels     struct {
		AlertName string `json:"alertname"`
	}
	CommonLabels struct {
		AlertName string `json:"alertname"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL string  `json:"externalURL"`
	Alerts      []Alert `json:"alerts"`
}

type Alert struct {
	Status      AlertStatus       `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	StartsAt     string `json:"startsAt"`
	EndsAt       string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
	Fingerprint  string `json:"fingerprint"`
}

type AlertStatus string

const (
	Resolved AlertStatus = "resolved"
	Firing   AlertStatus = "firing"
)

type AppriseRequest struct {
	Title  string    `json:"title,omitempty"`
	Body   string    `json:"body"`
	Type   MsgType   `json:"type,omitempty"`
	Tag    string    `json:"tag,omitempty"`
	Format MsgFormat `json:"format,omitempty"`
}

type MsgType string
type MsgFormat string

const (
	TypeInfo    MsgType = "info"
	TypeSuccess MsgType = "success"
	TypeWarning MsgType = "warning"
	TypeError   MsgType = "error"
)

const (
	FormatText     MsgFormat = "text"
	FormatMarkdown MsgFormat = "markdown"
	FormatHTML     MsgFormat = "html"
)

// Configuration flags
var (
	tag           = flag.String("tag", getEnv("TAG", "all"), "apprise.io notification tag")
	url           = flag.String("apprise.url", os.Getenv("APPRISE_URL"), "apprise.io API URL")
	listenAddress = flag.String("listen.address", getEnv("LISTEN_ADDRESS", ":8080"), "Address:Port to listen on")
)

// Application Entry
func main() {
	log.Printf("Listening on %s", *listenAddress)
	log.Printf("Apprise.url: %s", *url)
	log.Printf("Apprise.tag: %s", *tag)

	if *url == "" {
		log.Fatalf("apprise.url flag is required")
	}

	http.Handle("/metrics", promhttp.Handler())

	log.Fatalf("Failed to liisten on HTTP: %v",
		http.ListenAndServe(*listenAddress, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Printf("%s - [%s] %s", req.Host, req.Method, req.URL.RawPath)
			body, err := io.ReadAll(req.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
			}

			alert := AlertManagerRequest{}
			err = json.Unmarshal(body, &alert)
			if err != nil {
				log.Printf("Error parsing body: %v", err)
				return
			}

			notify(&alert)

		})))
}

// Notifies the Apprise server
func notify(alertManagerRequest *AlertManagerRequest) {
	groupedAlerts := make(map[AlertStatus][]Alert)
	for _, alert := range alertManagerRequest.Alerts {
		groupedAlerts[alert.Status] = append(groupedAlerts[alert.Status], alert)
	}

	message := AppriseRequest{
		Type:   TypeInfo,
		Format: FormatText,
		Tag:    *tag,
	}

	var body strings.Builder
	for status, alerts := range groupedAlerts {
		if status == Firing {
			body.WriteString("🔴 ")
		} else {
			body.WriteString("✅ ")
		}
		body.WriteString(fmt.Sprintf("[%s:%d] %s \n", strings.ToUpper(string(status)), len(alerts), alertManagerRequest.CommonLabels.AlertName))
		body.WriteString(fmt.Sprintf("\t\tSummary: %s \n", alertManagerRequest.CommonAnnotations.Summary))
		body.WriteString(fmt.Sprintf("\t\tURL:  %s \n", alertManagerRequest.ExternalURL))
	}
	message.Body = body.String()

	payload, _ := json.Marshal(message)
	_, err := http.Post(*url, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("Error notifing to Apprise: %v", err)
	}

}

// Helper function to retrieve environment variables or return the fallback value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
