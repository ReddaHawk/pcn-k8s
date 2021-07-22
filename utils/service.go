package utils

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
	// all services that are active in the cluster
	Services map[string]Service
)

type Backend struct {
	IP   string
	Port int32
}

// TODO: are VIP and Vport good names?
// What about IP and Port?
type Service struct {
	UID      types.UID
	Name     string
	Type     string
	VIP      string
	Vport    int32
	Proto    string
	NodePort int32
	Backends map[Backend]bool
}