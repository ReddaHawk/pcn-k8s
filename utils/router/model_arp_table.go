/*
 * router API
 *
 * router API generated from router.yang
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ArpTable struct {
	// Destination IP address
	Address string `json:"address,omitempty"`
	// Destination MAC address
	Mac string `json:"mac,omitempty"`
	// Outgoing interface
	Interface_ string `json:"interface,omitempty"`
}