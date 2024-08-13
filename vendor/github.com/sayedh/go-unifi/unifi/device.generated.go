// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type Device struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	MAC string `json:"mac,omitempty"`

	Adopted                     bool                              `json:"adopted"`
	AtfEnabled                  bool                              `json:"atf_enabled,omitempty"`
	BandsteeringMode            string                            `json:"bandsteering_mode,omitempty"` // off|equal|prefer_5g
	BaresipAuthUser             string                            `json:"baresip_auth_user,omitempty"` // ^\+?[a-zA-Z0-9_.\-!~*'()]*
	BaresipEnabled              bool                              `json:"baresip_enabled,omitempty"`
	BaresipExtension            string                            `json:"baresip_extension,omitempty"` // ^\+?[a-zA-Z0-9_.\-!~*'()]*
	ConfigNetwork               DeviceConfigNetwork               `json:"config_network,omitempty"`
	ConnectedBatteryOverrides   []DeviceConnectedBatteryOverrides `json:"connected_battery_overrides,omitempty"`
	DPIEnabled                  bool                              `json:"dpi_enabled,omitempty"`
	Disabled                    bool                              `json:"disabled,omitempty"`
	Dot1XFallbackNetworkID      string                            `json:"dot1x_fallback_networkconf_id,omitempty"` // [\d\w]+|
	Dot1XPortctrlEnabled        bool                              `json:"dot1x_portctrl_enabled,omitempty"`
	EtherLighting               DeviceEtherLighting               `json:"ether_lighting,omitempty"`
	EthernetOverrides           []DeviceEthernetOverrides         `json:"ethernet_overrides,omitempty"`
	FlowctrlEnabled             bool                              `json:"flowctrl_enabled,omitempty"`
	GatewayVrrpMode             string                            `json:"gateway_vrrp_mode,omitempty"`     // primary|secondary
	GatewayVrrpPriority         int                               `json:"gateway_vrrp_priority,omitempty"` // [1-9][0-9]|[1-9][0-9][0-9]
	HeightInMeters              float64                           `json:"heightInMeters,omitempty"`
	Hostname                    string                            `json:"hostname,omitempty"` // .{1,128}
	JumboframeEnabled           bool                              `json:"jumboframe_enabled,omitempty"`
	LcmBrightness               int                               `json:"lcm_brightness,omitempty"` // [1-9]|[1-9][0-9]|100
	LcmBrightnessOverride       bool                              `json:"lcm_brightness_override,omitempty"`
	LcmIDleTimeout              int                               `json:"lcm_idle_timeout,omitempty"` // [1-9][0-9]|[1-9][0-9][0-9]|[1-2][0-9][0-9][0-9]|3[0-5][0-9][0-9]|3600
	LcmIDleTimeoutOverride      bool                              `json:"lcm_idle_timeout_override,omitempty"`
	LcmNightModeBegins          string                            `json:"lcm_night_mode_begins,omitempty"` // (^$)|(^(0[1-9])|(1[0-9])|(2[0-3])):([0-5][0-9]$)
	LcmNightModeEnds            string                            `json:"lcm_night_mode_ends,omitempty"`   // (^$)|(^(0[1-9])|(1[0-9])|(2[0-3])):([0-5][0-9]$)
	LcmSettingsRestrictedAccess bool                              `json:"lcm_settings_restricted_access,omitempty"`
	LcmTrackerEnabled           bool                              `json:"lcm_tracker_enabled,omitempty"`
	LcmTrackerSeed              string                            `json:"lcm_tracker_seed,omitempty"`              // .{0,50}
	LedOverride                 string                            `json:"led_override,omitempty"`                  // default|on|off
	LedOverrideColor            string                            `json:"led_override_color,omitempty"`            // ^#(?:[0-9a-fA-F]{3}){1,2}$
	LedOverrideColorBrightness  int                               `json:"led_override_color_brightness,omitempty"` // ^[0-9][0-9]?$|^100$
	Locked                      bool                              `json:"locked,omitempty"`
	LowpfmodeOverride           bool                              `json:"lowpfmode_override,omitempty"`
	LteApn                      string                            `json:"lte_apn,omitempty"`       // .{1,128}
	LteAuthType                 string                            `json:"lte_auth_type,omitempty"` // PAP|CHAP|PAP-CHAP|NONE
	LteDataLimitEnabled         bool                              `json:"lte_data_limit_enabled,omitempty"`
	LteDataWarningEnabled       bool                              `json:"lte_data_warning_enabled,omitempty"`
	LteExtAnt                   bool                              `json:"lte_ext_ant,omitempty"`
	LteHardLimit                int                               `json:"lte_hard_limit,omitempty"`
	LtePassword                 string                            `json:"lte_password,omitempty"`
	LtePoe                      bool                              `json:"lte_poe,omitempty"`
	LteRoamingAllowed           bool                              `json:"lte_roaming_allowed,omitempty"`
	LteSimPin                   int                               `json:"lte_sim_pin,omitempty"`
	LteSoftLimit                int                               `json:"lte_soft_limit,omitempty"`
	LteUsername                 string                            `json:"lte_username,omitempty"`
	MapID                       string                            `json:"map_id,omitempty"`
	MeshStaVapEnabled           bool                              `json:"mesh_sta_vap_enabled,omitempty"`
	MgmtNetworkID               string                            `json:"mgmt_network_id,omitempty"` // [\d\w]+
	Model                       string                            `json:"model,omitempty"`
	Name                        string                            `json:"name,omitempty"`                  // .{0,128}
	OutdoorModeOverride         string                            `json:"outdoor_mode_override,omitempty"` // default|on|off
	OutletEnabled               bool                              `json:"outlet_enabled,omitempty"`
	OutletOverrides             []DeviceOutletOverrides           `json:"outlet_overrides,omitempty"`
	OutletPowerCycleEnabled     bool                              `json:"outlet_power_cycle_enabled,omitempty"`
	PeerToPeerMode              string                            `json:"peer_to_peer_mode,omitempty"` // ap|sta
	PoeMode                     string                            `json:"poe_mode,omitempty"`          // auto|pasv24|passthrough|off
	PortOverrides               []DevicePortOverrides             `json:"port_overrides"`
	PowerSourceCtrl             string                            `json:"power_source_ctrl,omitempty"`        // auto|8023af|8023at|8023bt-type3|8023bt-type4|pasv24|poe-injector|ac|adapter|dc|rps
	PowerSourceCtrlBudget       int                               `json:"power_source_ctrl_budget,omitempty"` // [0-9]|[1-9][0-9]|[1-9][0-9][0-9]
	PowerSourceCtrlEnabled      bool                              `json:"power_source_ctrl_enabled,omitempty"`
	RADIUSProfileID             string                            `json:"radiusprofile_id,omitempty"`
	RadioTable                  []DeviceRadioTable                `json:"radio_table,omitempty"`
	ResetbtnEnabled             string                            `json:"resetbtn_enabled,omitempty"` // on|off
	RpsOverride                 DeviceRpsOverride                 `json:"rps_override,omitempty"`
	SnmpContact                 string                            `json:"snmp_contact,omitempty"`  // .{0,255}
	SnmpLocation                string                            `json:"snmp_location,omitempty"` // .{0,255}
	State                       DeviceState                       `json:"state"`
	StpPriority                 string                            `json:"stp_priority,omitempty"` // 0|4096|8192|12288|16384|20480|24576|28672|32768|36864|40960|45056|49152|53248|57344|61440
	StpVersion                  string                            `json:"stp_version,omitempty"`  // stp|rstp|disabled
	SwitchVLANEnabled           bool                              `json:"switch_vlan_enabled,omitempty"`
	Type                        string                            `json:"type,omitempty"`
	UbbPairName                 string                            `json:"ubb_pair_name,omitempty"` // .{1,128}
	Volume                      int                               `json:"volume,omitempty"`        // [0-9]|[1-9][0-9]|100
	X                           float64                           `json:"x,omitempty"`
	XBaresipPassword            string                            `json:"x_baresip_password,omitempty"` // ^[a-zA-Z0-9_.\-!~*'()]*
	Y                           float64                           `json:"y,omitempty"`
}

func (dst *Device) UnmarshalJSON(b []byte) error {
	type Alias Device
	aux := &struct {
		GatewayVrrpPriority        emptyStringInt   `json:"gateway_vrrp_priority"`
		LcmBrightness              emptyStringInt   `json:"lcm_brightness"`
		LcmIDleTimeout             emptyStringInt   `json:"lcm_idle_timeout"`
		LedOverrideColorBrightness emptyStringInt   `json:"led_override_color_brightness"`
		LteExtAnt                  booleanishString `json:"lte_ext_ant"`
		LteHardLimit               emptyStringInt   `json:"lte_hard_limit"`
		LtePoe                     booleanishString `json:"lte_poe"`
		LteSimPin                  emptyStringInt   `json:"lte_sim_pin"`
		LteSoftLimit               emptyStringInt   `json:"lte_soft_limit"`
		PowerSourceCtrlBudget      emptyStringInt   `json:"power_source_ctrl_budget"`
		StpPriority                numberOrString   `json:"stp_priority"`
		Volume                     emptyStringInt   `json:"volume"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.GatewayVrrpPriority = int(aux.GatewayVrrpPriority)
	dst.LcmBrightness = int(aux.LcmBrightness)
	dst.LcmIDleTimeout = int(aux.LcmIDleTimeout)
	dst.LedOverrideColorBrightness = int(aux.LedOverrideColorBrightness)
	dst.LteExtAnt = bool(aux.LteExtAnt)
	dst.LteHardLimit = int(aux.LteHardLimit)
	dst.LtePoe = bool(aux.LtePoe)
	dst.LteSimPin = int(aux.LteSimPin)
	dst.LteSoftLimit = int(aux.LteSoftLimit)
	dst.PowerSourceCtrlBudget = int(aux.PowerSourceCtrlBudget)
	dst.StpPriority = string(aux.StpPriority)
	dst.Volume = int(aux.Volume)

	return nil
}

type DeviceConfigNetwork struct {
	BondingEnabled bool   `json:"bonding_enabled,omitempty"`
	DNS1           string `json:"dns1,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	DNS2           string `json:"dns2,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	DNSsuffix      string `json:"dnssuffix,omitempty"`
	Gateway        string `json:"gateway,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	IP             string `json:"ip,omitempty"`      // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	Netmask        string `json:"netmask,omitempty"` // ^((128|192|224|240|248|252|254)\.0\.0\.0)|(255\.(((0|128|192|224|240|248|252|254)\.0\.0)|(255\.(((0|128|192|224|240|248|252|254)\.0)|255\.(0|128|192|224|240|248|252|254)))))$
	Type           string `json:"type,omitempty"`    // dhcp|static
}

func (dst *DeviceConfigNetwork) UnmarshalJSON(b []byte) error {
	type Alias DeviceConfigNetwork
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceConnectedBatteryOverrides struct {
	MAC string `json:"mac,omitempty"` // ^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$
}

func (dst *DeviceConnectedBatteryOverrides) UnmarshalJSON(b []byte) error {
	type Alias DeviceConnectedBatteryOverrides
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceEtherLighting struct {
	Behavior   string `json:"behavior,omitempty"`   // breath|steady
	Brightness int    `json:"brightness,omitempty"` // [1-9]|[1-9][0-9]|100
	Mode       string `json:"mode,omitempty"`       // speed|network
}

func (dst *DeviceEtherLighting) UnmarshalJSON(b []byte) error {
	type Alias DeviceEtherLighting
	aux := &struct {
		Brightness emptyStringInt `json:"brightness"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Brightness = int(aux.Brightness)

	return nil
}

type DeviceEthernetOverrides struct {
	Ifname       string `json:"ifname,omitempty"`       // eth[0-9]{1,2}
	NetworkGroup string `json:"networkgroup,omitempty"` // LAN[2-8]?|WAN[2]?
}

func (dst *DeviceEthernetOverrides) UnmarshalJSON(b []byte) error {
	type Alias DeviceEthernetOverrides
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceOutletOverrides struct {
	CycleEnabled bool   `json:"cycle_enabled,omitempty"`
	Index        int    `json:"index,omitempty"`
	Name         string `json:"name,omitempty"` // .{0,128}
	RelayState   bool   `json:"relay_state,omitempty"`
}

func (dst *DeviceOutletOverrides) UnmarshalJSON(b []byte) error {
	type Alias DeviceOutletOverrides
	aux := &struct {
		Index emptyStringInt `json:"index"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Index = int(aux.Index)

	return nil
}

type DevicePortOverrides struct {
	AggregateNumPorts            int              `json:"aggregate_num_ports,omitempty"` // [1-8]
	Autoneg                      bool             `json:"autoneg,omitempty"`
	Dot1XCtrl                    string           `json:"dot1x_ctrl,omitempty"`             // auto|force_authorized|force_unauthorized|mac_based|multi_host
	Dot1XIDleTimeout             int              `json:"dot1x_idle_timeout,omitempty"`     // [0-9]|[1-9][0-9]{1,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]
	EgressRateLimitKbps          int              `json:"egress_rate_limit_kbps,omitempty"` // 6[4-9]|[7-9][0-9]|[1-9][0-9]{2,6}
	EgressRateLimitKbpsEnabled   bool             `json:"egress_rate_limit_kbps_enabled,omitempty"`
	ExcludedNetworkIDs           []string         `json:"excluded_networkconf_ids,omitempty"`
	FecMode                      string           `json:"fec_mode,omitempty"` // rs-fec|fc-fec|default|disabled
	Forward                      string           `json:"forward,omitempty"`  // all|native|customize|disabled
	FullDuplex                   bool             `json:"full_duplex,omitempty"`
	Isolation                    bool             `json:"isolation,omitempty"`
	LldpmedEnabled               bool             `json:"lldpmed_enabled,omitempty"`
	LldpmedNotifyEnabled         bool             `json:"lldpmed_notify_enabled,omitempty"`
	MirrorPortIDX                int              `json:"mirror_port_idx,omitempty"` // [1-9]|[1-4][0-9]|5[0-2]
	NATiveNetworkID              string           `json:"native_networkconf_id,omitempty"`
	Name                         string           `json:"name,omitempty"`     // .{0,128}
	OpMode                       string           `json:"op_mode,omitempty"`  // switch|mirror|aggregate
	PoeMode                      string           `json:"poe_mode,omitempty"` // auto|pasv24|passthrough|off
	PortIDX                      int              `json:"port_idx,omitempty"` // [1-9]|[1-4][0-9]|5[0-2]
	PortKeepaliveEnabled         bool             `json:"port_keepalive_enabled,omitempty"`
	PortProfileID                string           `json:"portconf_id,omitempty"` // [\d\w]+
	PortSecurityEnabled          bool             `json:"port_security_enabled,omitempty"`
	PortSecurityMACAddress       []string         `json:"port_security_mac_address,omitempty"` // ^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$
	PriorityQueue1Level          int              `json:"priority_queue1_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue2Level          int              `json:"priority_queue2_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue3Level          int              `json:"priority_queue3_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue4Level          int              `json:"priority_queue4_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	QOSProfile                   DeviceQOSProfile `json:"qos_profile,omitempty"`
	SettingPreference            string           `json:"setting_preference,omitempty"` // auto|manual
	Speed                        int              `json:"speed,omitempty"`              // 10|100|1000|2500|5000|10000|20000|25000|40000|50000|100000
	StormctrlBroadcastastEnabled bool             `json:"stormctrl_bcast_enabled,omitempty"`
	StormctrlBroadcastastLevel   int              `json:"stormctrl_bcast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlBroadcastastRate    int              `json:"stormctrl_bcast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StormctrlMcastEnabled        bool             `json:"stormctrl_mcast_enabled,omitempty"`
	StormctrlMcastLevel          int              `json:"stormctrl_mcast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlMcastRate           int              `json:"stormctrl_mcast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StormctrlType                string           `json:"stormctrl_type,omitempty"`        // level|rate
	StormctrlUcastEnabled        bool             `json:"stormctrl_ucast_enabled,omitempty"`
	StormctrlUcastLevel          int              `json:"stormctrl_ucast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlUcastRate           int              `json:"stormctrl_ucast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StpPortMode                  bool             `json:"stp_port_mode,omitempty"`
	TaggedVLANMgmt               string           `json:"tagged_vlan_mgmt,omitempty"` // auto|block_all|custom
	VoiceNetworkID               string           `json:"voice_networkconf_id,omitempty"`
}

func (dst *DevicePortOverrides) UnmarshalJSON(b []byte) error {
	type Alias DevicePortOverrides
	aux := &struct {
		AggregateNumPorts          emptyStringInt `json:"aggregate_num_ports"`
		Dot1XIDleTimeout           emptyStringInt `json:"dot1x_idle_timeout"`
		EgressRateLimitKbps        emptyStringInt `json:"egress_rate_limit_kbps"`
		MirrorPortIDX              emptyStringInt `json:"mirror_port_idx"`
		PortIDX                    emptyStringInt `json:"port_idx"`
		PriorityQueue1Level        emptyStringInt `json:"priority_queue1_level"`
		PriorityQueue2Level        emptyStringInt `json:"priority_queue2_level"`
		PriorityQueue3Level        emptyStringInt `json:"priority_queue3_level"`
		PriorityQueue4Level        emptyStringInt `json:"priority_queue4_level"`
		Speed                      emptyStringInt `json:"speed"`
		StormctrlBroadcastastLevel emptyStringInt `json:"stormctrl_bcast_level"`
		StormctrlBroadcastastRate  emptyStringInt `json:"stormctrl_bcast_rate"`
		StormctrlMcastLevel        emptyStringInt `json:"stormctrl_mcast_level"`
		StormctrlMcastRate         emptyStringInt `json:"stormctrl_mcast_rate"`
		StormctrlUcastLevel        emptyStringInt `json:"stormctrl_ucast_level"`
		StormctrlUcastRate         emptyStringInt `json:"stormctrl_ucast_rate"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.AggregateNumPorts = int(aux.AggregateNumPorts)
	dst.Dot1XIDleTimeout = int(aux.Dot1XIDleTimeout)
	dst.EgressRateLimitKbps = int(aux.EgressRateLimitKbps)
	dst.MirrorPortIDX = int(aux.MirrorPortIDX)
	dst.PortIDX = int(aux.PortIDX)
	dst.PriorityQueue1Level = int(aux.PriorityQueue1Level)
	dst.PriorityQueue2Level = int(aux.PriorityQueue2Level)
	dst.PriorityQueue3Level = int(aux.PriorityQueue3Level)
	dst.PriorityQueue4Level = int(aux.PriorityQueue4Level)
	dst.Speed = int(aux.Speed)
	dst.StormctrlBroadcastastLevel = int(aux.StormctrlBroadcastastLevel)
	dst.StormctrlBroadcastastRate = int(aux.StormctrlBroadcastastRate)
	dst.StormctrlMcastLevel = int(aux.StormctrlMcastLevel)
	dst.StormctrlMcastRate = int(aux.StormctrlMcastRate)
	dst.StormctrlUcastLevel = int(aux.StormctrlUcastLevel)
	dst.StormctrlUcastRate = int(aux.StormctrlUcastRate)

	return nil
}

type DeviceQOSMarking struct {
	CosCode          int `json:"cos_code,omitempty"`           // [0-7]
	DscpCode         int `json:"dscp_code,omitempty"`          // 0|8|16|24|32|40|48|56|10|12|14|18|20|22|26|28|30|34|36|38|44|46
	IPPrecedenceCode int `json:"ip_precedence_code,omitempty"` // [0-7]
	Queue            int `json:"queue,omitempty"`              // [0-7]
}

func (dst *DeviceQOSMarking) UnmarshalJSON(b []byte) error {
	type Alias DeviceQOSMarking
	aux := &struct {
		CosCode          emptyStringInt `json:"cos_code"`
		DscpCode         emptyStringInt `json:"dscp_code"`
		IPPrecedenceCode emptyStringInt `json:"ip_precedence_code"`
		Queue            emptyStringInt `json:"queue"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.CosCode = int(aux.CosCode)
	dst.DscpCode = int(aux.DscpCode)
	dst.IPPrecedenceCode = int(aux.IPPrecedenceCode)
	dst.Queue = int(aux.Queue)

	return nil
}

type DeviceQOSMatching struct {
	CosCode          int    `json:"cos_code,omitempty"`           // [0-7]
	DscpCode         int    `json:"dscp_code,omitempty"`          // [0-9]|[1-5][0-9]|6[0-3]
	DstPort          int    `json:"dst_port,omitempty"`           // [0-9]|[1-9][0-9]|[1-9][0-9][0-9]|[1-9][0-9][0-9][0-9]|[1-5][0-9][0-9][0-9][0-9]|6[0-4][0-9][0-9][0-9]|65[0-4][0-9][0-9]|655[0-2][0-9]|6553[0-4]|65535
	IPPrecedenceCode int    `json:"ip_precedence_code,omitempty"` // [0-7]
	Protocol         string `json:"protocol,omitempty"`           // ([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|ah|ax.25|dccp|ddp|egp|eigrp|encap|esp|etherip|fc|ggp|gre|hip|hmp|icmp|idpr-cmtp|idrp|igmp|igp|ip|ipcomp|ipencap|ipip|ipv6|ipv6-frag|ipv6-icmp|ipv6-nonxt|ipv6-opts|ipv6-route|isis|iso-tp4|l2tp|manet|mobility-header|mpls-in-ip|ospf|pim|pup|rdp|rohc|rspf|rsvp|sctp|shim6|skip|st|tcp|udp|udplite|vmtp|vrrp|wesp|xns-idp|xtp
	SrcPort          int    `json:"src_port,omitempty"`           // [0-9]|[1-9][0-9]|[1-9][0-9][0-9]|[1-9][0-9][0-9][0-9]|[1-5][0-9][0-9][0-9][0-9]|6[0-4][0-9][0-9][0-9]|65[0-4][0-9][0-9]|655[0-2][0-9]|6553[0-4]|65535
}

func (dst *DeviceQOSMatching) UnmarshalJSON(b []byte) error {
	type Alias DeviceQOSMatching
	aux := &struct {
		CosCode          emptyStringInt `json:"cos_code"`
		DscpCode         emptyStringInt `json:"dscp_code"`
		DstPort          emptyStringInt `json:"dst_port"`
		IPPrecedenceCode emptyStringInt `json:"ip_precedence_code"`
		SrcPort          emptyStringInt `json:"src_port"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.CosCode = int(aux.CosCode)
	dst.DscpCode = int(aux.DscpCode)
	dst.DstPort = int(aux.DstPort)
	dst.IPPrecedenceCode = int(aux.IPPrecedenceCode)
	dst.SrcPort = int(aux.SrcPort)

	return nil
}

type DeviceQOSPolicies struct {
	QOSMarking  DeviceQOSMarking  `json:"qos_marking,omitempty"`
	QOSMatching DeviceQOSMatching `json:"qos_matching,omitempty"`
}

func (dst *DeviceQOSPolicies) UnmarshalJSON(b []byte) error {
	type Alias DeviceQOSPolicies
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceQOSProfile struct {
	QOSPolicies    []DeviceQOSPolicies `json:"qos_policies,omitempty"`
	QOSProfileMode string              `json:"qos_profile_mode,omitempty"` // custom|unifi_play|aes67_audio|crestron_audio_video|dante_audio|ndi_aes67_audio|ndi_dante_audio|qsys_audio_video|qsys_video_dante_audio|sdvoe_aes67_audio|sdvoe_dante_audio|shure_audio
}

func (dst *DeviceQOSProfile) UnmarshalJSON(b []byte) error {
	type Alias DeviceQOSProfile
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceRadioIDentifiers struct {
	DeviceID  string `json:"device_id,omitempty"`
	RadioName string `json:"radio_name,omitempty"`
}

func (dst *DeviceRadioIDentifiers) UnmarshalJSON(b []byte) error {
	type Alias DeviceRadioIDentifiers
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceRadioTable struct {
	AntennaGain                int                      `json:"antenna_gain,omitempty"`   // ^-?([0-9]|[1-9][0-9])
	AntennaID                  int                      `json:"antenna_id,omitempty"`     // -1|[0-9]
	BackupChannel              string                   `json:"backup_channel,omitempty"` // [0-9]|[1][0-4]|4.5|5|16|17|21|25|29|33|34|36|37|38|40|41|42|44|45|46|48|49|52|53|56|57|60|61|64|65|69|73|77|81|85|89|93|97|100|101|104|105|108|109|112|113|117|116|120|121|124|125|128|129|132|133|136|137|140|141|144|145|149|153|157|161|165|169|173|177|181|183|184|185|187|188|189|192|193|196|197|201|205|209|213|217|221|225|229|233|auto
	Channel                    string                   `json:"channel,omitempty"`        // [0-9]|[1][0-4]|4.5|5|16|17|21|25|29|33|34|36|37|38|40|41|42|44|45|46|48|49|52|53|56|57|60|61|64|65|69|73|77|81|85|89|93|97|100|101|104|105|108|109|112|113|117|116|120|121|124|125|128|129|132|133|136|137|140|141|144|145|149|153|157|161|165|169|173|177|181|183|184|185|187|188|189|192|193|196|197|201|205|209|213|217|221|225|229|233|auto
	ChannelOptimizationEnabled bool                     `json:"channel_optimization_enabled,omitempty"`
	HardNoiseFloorEnabled      bool                     `json:"hard_noise_floor_enabled,omitempty"`
	Ht                         int                      `json:"ht,omitempty"` // 20|40|80|160|240|320|1080|2160|4320
	LoadbalanceEnabled         bool                     `json:"loadbalance_enabled,omitempty"`
	Maxsta                     int                      `json:"maxsta,omitempty"`   // [1-9]|[1-9][0-9]|1[0-9]{2}|200|^$
	MinRssi                    int                      `json:"min_rssi,omitempty"` // ^-(6[7-9]|[7-8][0-9]|90)$
	MinRssiEnabled             bool                     `json:"min_rssi_enabled,omitempty"`
	Name                       string                   `json:"name,omitempty"`
	Radio                      string                   `json:"radio,omitempty"` // ng|na|ad|6e
	RadioIDentifiers           []DeviceRadioIDentifiers `json:"radio_identifiers,omitempty"`
	SensLevel                  int                      `json:"sens_level,omitempty"` // ^-([5-8][0-9]|90)$
	SensLevelEnabled           bool                     `json:"sens_level_enabled,omitempty"`
	TxPower                    string                   `json:"tx_power,omitempty"`      // [\d]+|auto
	TxPowerMode                string                   `json:"tx_power_mode,omitempty"` // auto|medium|high|low|custom
	VwireEnabled               bool                     `json:"vwire_enabled,omitempty"`
}

func (dst *DeviceRadioTable) UnmarshalJSON(b []byte) error {
	type Alias DeviceRadioTable
	aux := &struct {
		AntennaGain   emptyStringInt `json:"antenna_gain"`
		AntennaID     emptyStringInt `json:"antenna_id"`
		BackupChannel numberOrString `json:"backup_channel"`
		Channel       numberOrString `json:"channel"`
		Ht            emptyStringInt `json:"ht"`
		Maxsta        emptyStringInt `json:"maxsta"`
		MinRssi       emptyStringInt `json:"min_rssi"`
		SensLevel     emptyStringInt `json:"sens_level"`
		TxPower       numberOrString `json:"tx_power"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.AntennaGain = int(aux.AntennaGain)
	dst.AntennaID = int(aux.AntennaID)
	dst.BackupChannel = string(aux.BackupChannel)
	dst.Channel = string(aux.Channel)
	dst.Ht = int(aux.Ht)
	dst.Maxsta = int(aux.Maxsta)
	dst.MinRssi = int(aux.MinRssi)
	dst.SensLevel = int(aux.SensLevel)
	dst.TxPower = string(aux.TxPower)

	return nil
}

type DeviceRpsOverride struct {
	PowerManagementMode string               `json:"power_management_mode,omitempty"` // dynamic|static
	RpsPortTable        []DeviceRpsPortTable `json:"rps_port_table,omitempty"`
}

func (dst *DeviceRpsOverride) UnmarshalJSON(b []byte) error {
	type Alias DeviceRpsOverride
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

type DeviceRpsPortTable struct {
	Name     string `json:"name,omitempty"`      // .{0,32}
	PortIDX  int    `json:"port_idx,omitempty"`  // [1-8]
	PortMode string `json:"port_mode,omitempty"` // auto|force_active|manual|disabled
}

func (dst *DeviceRpsPortTable) UnmarshalJSON(b []byte) error {
	type Alias DeviceRpsPortTable
	aux := &struct {
		PortIDX emptyStringInt `json:"port_idx"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.PortIDX = int(aux.PortIDX)

	return nil
}

func (c *Client) listDevice(ctx context.Context, site string) ([]Device, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/stat/device", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDevice(ctx context.Context, site, id string) (*Device, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/stat/device/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDevice(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/device/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/device", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/device/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
