// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type PortConf struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Autoneg                      bool     `json:"autoneg"`
	Dot1XCtrl                    string   `json:"dot1x_ctrl,omitempty"`             // auto|force_authorized|force_unauthorized|mac_based|multi_host
	EgressRateLimitKbps          int      `json:"egress_rate_limit_kbps,omitempty"` // 6[4-9]|[7-9][0-9]|[1-9][0-9]{2,6}
	EgressRateLimitKbpsEnabled   bool     `json:"egress_rate_limit_kbps_enabled"`
	Forward                      string   `json:"forward,omitempty"` // all|native|customize|disabled
	FullDuplex                   bool     `json:"full_duplex"`
	Isolation                    bool     `json:"isolation"`
	LldpmedEnabled               bool     `json:"lldpmed_enabled"`
	LldpmedNotifyEnabled         bool     `json:"lldpmed_notify_enabled"`
	Name                         string   `json:"name,omitempty"`
	NATiveNetworkconfID          string   `json:"native_networkconf_id"`
	OpMode                       string   `json:"op_mode,omitempty"`  // switch
	PoeMode                      string   `json:"poe_mode,omitempty"` // auto|pasv24|passthrough|off
	PortSecurityEnabled          bool     `json:"port_security_enabled"`
	PortSecurityMACAddress       []string `json:"port_security_mac_address,omitempty"` // ^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$
	PriorityQueue1Level          int      `json:"priority_queue1_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue2Level          int      `json:"priority_queue2_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue3Level          int      `json:"priority_queue3_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue4Level          int      `json:"priority_queue4_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	Speed                        int      `json:"speed,omitempty"`                     // 10|100|1000|2500|5000|10000|20000|25000|40000|50000|100000
	StormctrlBroadcastastEnabled bool     `json:"stormctrl_bcast_enabled"`
	StormctrlBroadcastastLevel   int      `json:"stormctrl_bcast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlBroadcastastRate    int      `json:"stormctrl_bcast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StormctrlMcastEnabled        bool     `json:"stormctrl_mcast_enabled"`
	StormctrlMcastLevel          int      `json:"stormctrl_mcast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlMcastRate           int      `json:"stormctrl_mcast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StormctrlType                string   `json:"stormctrl_type,omitempty"`        // level|rate
	StormctrlUcastEnabled        bool     `json:"stormctrl_ucast_enabled"`
	StormctrlUcastLevel          int      `json:"stormctrl_ucast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlUcastRate           int      `json:"stormctrl_ucast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StpPortMode                  bool     `json:"stp_port_mode"`
	TaggedNetworkconfIDs         []string `json:"tagged_networkconf_ids,omitempty"`
	VoiceNetworkconfID           string   `json:"voice_networkconf_id"`
}
