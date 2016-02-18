package destination

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/healthcareblocks/ec2_metrics_publisher/metadata"
	"github.com/healthcareblocks/ec2_metrics_publisher/metrics"
)

// Slack implements the Service interface, wrapping Slack Incoming WebHook functionality
type Slack struct {
	*metadata.Machine
	*ServiceEndpoint
}

// SlackMessage wraps a payload for Slack's Incoming WebHook API
// https://api.slack.com/incoming-webhooks
type SlackMessage struct {
	Text        string            `json:"text,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment is an attachment object for wrapping data
type SlackAttachment struct {
	Title  string                 `json:"title"`
	Fields []SlackAttachmentField `json:"fields,omitempty"`
}

// SlackAttachmentField contains field info for attachments
type SlackAttachmentField struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Short bool        `json:"short"`
}

// AttachmentData is a slice of SlackAttachment
type AttachmentData []SlackAttachment

// String returns the name of the service
func (s Slack) String() string {
	return "Slack"
}

// SetEndpoint sets the service endpoint
func (s *Slack) SetEndpoint(endpoint *ServiceEndpoint) {
	s.ServiceEndpoint = endpoint
}

// SetMachine sets the metadata.Machine object
func (s *Slack) SetMachine(machine *metadata.Machine) {
	s.Machine = machine
}

// SendMessage posts message to the Slack WebHook API
func (s *Slack) SendMessage(sys *metrics.System) error {
	if s.ServiceEndpoint == nil {
		return errors.New("must set service endpoint")
	}

	if s.Machine == nil || s.Machine.IsEmpty() {
		return errors.New("must set machine instance and region")
	}

	var attachments AttachmentData

	attachments.add(sys.CPU.Metrics, sys.CPU.Timestamp, "CPU")
	attachments.add(sys.Memory.Metrics, sys.Memory.Timestamp, "Memory Usage")

	for _, v := range sys.Volumes {
		attachments.add(v.Metrics, v.Timestamp, fmt.Sprintf("Volume Usage (Path: %s)", v.LookupPath))
	}

	if len(attachments) == 0 {
		return errors.New("must set at least one metric")
	}

	msg := SlackMessage{
		Text:        fmt.Sprintf("Instance %s", s.Instance),
		Attachments: attachments,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(t))
	}

	return nil
}

// adds SlackAttachmentField to AttachmentData
func (attachments *AttachmentData) add(metrics map[string]float64,
	timestamp time.Time, title string) {

	if len(metrics) > 0 {
		fields := []SlackAttachmentField{}

		for key, value := range metrics {
			fields = append(fields, createField(key, value))
		}
		sort.Sort(sortByTitle(fields))
		fields = append(fields, createField("Timestamp", timestamp))
		*attachments = append(*attachments, SlackAttachment{Title: title, Fields: fields})
	}
}

// creates SlackAttachmentField
func createField(key string, value interface{}) SlackAttachmentField {
	return SlackAttachmentField{
		Title: key,
		Value: value,
		Short: true,
	}
}

// sortByTitle is an internal slice for sorting SlackAttachmentField objects
type sortByTitle []SlackAttachmentField

// Len implements the sort.Interface for sortByTitle
func (a sortByTitle) Len() int {
	return len(a)
}

// Swap implements the sort.Interface for sortByTitle
func (a sortByTitle) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less implements the sort.Interface for sortByTitle
func (a sortByTitle) Less(i, j int) bool {
	return a[i].Title < a[j].Title
}
