/*
 * router API
 *
 * router API generated from router.yang
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Ports struct {
	// Port Name
	Name string `json:"name,omitempty"`
	// UUID of the port
	Uuid string `json:"uuid,omitempty"`
	// Status of the port (UP or DOWN)
	Status string `json:"status,omitempty"`
	// Peer name, such as a network interfaces (e.g., 'veth0') or another cube (e.g., 'br1:port2')
	Peer   string   `json:"peer,omitempty"`
	Tcubes []string `json:"tcubes,omitempty"`
	// IP address and prefix of the port
	Ip string `json:"ip,omitempty"`
	// Additional IP addresses for the port
	Secondaryip []PortsSecondaryip `json:"secondaryip,omitempty"`
	// MAC address of the port
	Mac string `json:"mac,omitempty"`
}
