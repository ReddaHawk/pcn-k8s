package controllers

import (
	"k8s.io/apimachinery/pkg/types"
)

type Pod struct {
	UID    types.UID
	Name   string
	Ip 	   string
}

var (
	Pods map[string]Pod
	Services map[string]Service
)

type Backend struct {
	IP   string
	Port int32
}

type ServicePortKey struct {
	Port  int32
	Proto string
}

type ServicePort struct {
	Port     int32
	Proto    string
	Name     string
	Nodeport int32
	Backends map[Backend]bool // backends implementing this port
}

type Service struct {
	UID   types.UID
	Name  string
	Type  string
	VIP   string
	ExternalTrafficPolicy string
	Ports map[ServicePortKey]ServicePort // different ports exposed by the service
}

