package logging

import (
	"encoding/json"
)

var (
	sourceLogName string
)

func SetSourceLogName(srcLog string) {
	sourceLogName = srcLog
}

// Structured Logging
/*
{
    "service_src": "http-configure-backend",
    "mac_address": "000000000000",
    "event_uuid": "00000000-0000-0000-0000-000000000000",
    "original": "hifpt_group_acl",
    "topic": "public.hifpt.prod.cpe.configure.v1",
    "method": "GET",
    "action": "/api/v1/health",
    "http_code": 200,
    "status_key": "success",
    "message": "Request is successful",
    "response" :{},
    "request": {}
}
*/

type serviceLogging struct {
	ServiceSrc string `json:"service_src"`
	MacAddress string `json:"mac_address"`
	EventUUID  string `json:"event_uuid"`
	Original   string `json:"original"`
	Topic      string `json:"topic"`
	Method     string `json:"method"`
	Action     string `json:"action"`
	HttpCode   int    `json:"http_code"`
	StatusKey  string `json:"status_key"`
	Message    string `json:"message"`
	Response   string `json:"response"`
	Request    string `json:"request"`
}

func (l *serviceLogging) ToJsonString() string {
	jsonData, _ := json.Marshal(l)
	return string(jsonData)
}

func (l *serviceLogging) Set(elements ...LogElement) {
	for _, e := range elements {
		e.f(l)
	}
}

type LogElement struct {
	f func(*serviceLogging)
}

func MacAddress(mac string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.MacAddress = mac
	}}
}

func EventUUID(uuid string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.EventUUID = uuid
	}}
}

func Original(original string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Original = original
	}}
}

func Topic(topic string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Topic = topic
	}}
}

func Method(method string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Method = method
	}}
}

func Action(action string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Action = action
	}}
}

func HttpCode(code int) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.HttpCode = code
	}}
}

func StatusKey(key string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.StatusKey = key
	}}
}

func Message(msg string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Message = msg
	}}
}

func Response(resp string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Response = resp
	}}
}

func Request(req string) LogElement {
	return LogElement{func(l *serviceLogging) {
		l.Request = req
	}}
}

func NewServiceLogging(
	mac,
	uuid,
	original string,
	elements ...LogElement,
) *serviceLogging {
	l := &serviceLogging{
		ServiceSrc: sourceLogName,
		MacAddress: mac,
		EventUUID:  uuid,
		Original:   original,
	}
	l.Set(elements...)
	return l
}

func NewHealthcheckLogging() *serviceLogging {
	return &serviceLogging{
		ServiceSrc: sourceLogName,
		StatusKey:  "Success",
		Message:    "Healthcheck successfully",
	}
}
