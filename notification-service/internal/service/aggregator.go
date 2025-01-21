package service

import (
	"bytes"
	"net/smtp"
	"strconv"
	"sync"
	"time"
)

type aggregatorData struct {
	sender    string
	count     int
	timer     *time.Timer
	recipient string
}

type Aggregator struct {
	smtpHost string
	duration time.Duration
	mu       sync.Mutex
	storage  map[string]*aggregatorData
}

func NewAggregator(h string, d time.Duration) *Aggregator {
	return &Aggregator{
		smtpHost: h,
		duration: d,
		storage:  make(map[string]*aggregatorData),
	}
}

func (a *Aggregator) AddMessage(sender, recipient string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	key := sender + "_" + recipient
	v, ok := a.storage[key]
	if !ok {
		v = &aggregatorData{sender: sender, count: 0, recipient: recipient}
		v.timer = time.AfterFunc(a.duration, func() {
			a.mu.Lock()
			defer a.mu.Unlock()
			if v.count > 0 {
				a.sendEmail(v.sender, v.recipient, v.count)
			}
			delete(a.storage, key)
		})
		a.storage[key] = v
	}
	v.count++
	v.timer.Reset(a.duration)
}

func (a *Aggregator) sendEmail(sender, recipient string, count int) {
	var buf bytes.Buffer
	buf.WriteString("From: notification-service\n")
	buf.WriteString("To: " + recipient + "\n")
	buf.WriteString("Subject: New Messages\n")
	buf.WriteString("\n")
	buf.WriteString("You have " + strconv.Itoa(count) + " new messages from @" + sender + "\n")
	smtp.SendMail(a.smtpHost, nil, "no-reply@service", []string{recipient}, buf.Bytes())
}
