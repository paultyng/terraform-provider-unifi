// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type Network struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	DHCPRelayEnabled        bool     `json:"dhcp_relay_enabled"`
	DHCPDBootEnabled        bool     `json:"dhcpd_boot_enabled"`
	DHCPDBootFilename       string   `json:"dhcpd_boot_filename,omitempty"` // .{1,256}
	DHCPDBootServer         string   `json:"dhcpd_boot_server"`             // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$|(?=^.{3,253}$)(^((?!-)[a-zA-Z0-9-]{1,63}(?<!-)\.)+[a-zA-Z]{2,63}$)|[a-zA-Z0-9-]{1,63}|^$
	DHCPDDNS1               string   `json:"dhcpd_dns_1"`                   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNS2               string   `json:"dhcpd_dns_2"`                   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNS3               string   `json:"dhcpd_dns_3"`                   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNS4               string   `json:"dhcpd_dns_4"`                   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNSEnabled         bool     `json:"dhcpd_dns_enabled"`
	DHCPDEnabled            bool     `json:"dhcpd_enabled"`
	DHCPDGateway            string   `json:"dhcpd_gateway"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDGatewayEnabled     bool     `json:"dhcpd_gateway_enabled"`
	DHCPDIP1                string   `json:"dhcpd_ip_1"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDIP2                string   `json:"dhcpd_ip_2"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDIP3                string   `json:"dhcpd_ip_3"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDLeaseTime          int      `json:"dhcpd_leasetime,omitempty"`
	DHCPDMAC1               string   `json:"dhcpd_mac_1,omitempty"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	DHCPDMAC2               string   `json:"dhcpd_mac_2,omitempty"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	DHCPDMAC3               string   `json:"dhcpd_mac_3,omitempty"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	DHCPDNtp1               string   `json:"dhcpd_ntp_1"`           // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDNtp2               string   `json:"dhcpd_ntp_2"`           // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDNtpEnabled         bool     `json:"dhcpd_ntp_enabled"`
	DHCPDStart              string   `json:"dhcpd_start"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDStop               string   `json:"dhcpd_stop"`  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDTftpServer         string   `json:"dhcpd_tftp_server,omitempty"`
	DHCPDTimeOffset         int      `json:"dhcpd_time_offset,omitempty"` // ^0$|^-?([1-9]([0-9]{1,3})?|[1-7][0-9]{4}|[8][0-5][0-9]{3}|86[0-3][0-9]{2}|86400)$
	DHCPDTimeOffsetEnabled  bool     `json:"dhcpd_time_offset_enabled"`
	DHCPDUnifiController    string   `json:"dhcpd_unifi_controller"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDWins1              string   `json:"dhcpd_wins_1"`           // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDWins2              string   `json:"dhcpd_wins_2"`           // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDWinsEnabled        bool     `json:"dhcpd_wins_enabled"`
	DHCPDWPAdUrl            string   `json:"dhcpd_wpad_url,omitempty"`
	DHCPDv6DNS1             string   `json:"dhcpdv6_dns_1,omitempty"`
	DHCPDv6DNS2             string   `json:"dhcpdv6_dns_2,omitempty"`
	DHCPDv6DNS3             string   `json:"dhcpdv6_dns_3,omitempty"`
	DHCPDv6DNS4             string   `json:"dhcpdv6_dns_4,omitempty"`
	DHCPDv6DNSAuto          bool     `json:"dhcpdv6_dns_auto"`
	DHCPDv6Enabled          bool     `json:"dhcpdv6_enabled"`
	DHCPDv6LeaseTime        int      `json:"dhcpdv6_leasetime,omitempty"`
	DHCPDv6Start            string   `json:"dhcpdv6_start,omitempty"`
	DHCPDv6Stop             string   `json:"dhcpdv6_stop,omitempty"`
	DHCPguardEnabled        bool     `json:"dhcpguard_enabled"`
	DomainName              string   `json:"domain_name,omitempty"` // (?=^.{3,253}$)(^((?!-)[a-zA-Z0-9-]{1,63}(?<!-)\.)+[a-zA-Z]{2,63}$)|^$|[a-zA-Z0-9-]{1,63}
	DPIEnabled              bool     `json:"dpi_enabled"`
	DPIgroupID              string   `json:"dpigroup_id"` // [\d\w]+|^$
	Enabled                 bool     `json:"enabled"`
	ExposedToSiteVpn        bool     `json:"exposed_to_site_vpn"`
	IgmpFastleave           bool     `json:"igmp_fastleave"`
	IgmpGroupmembership     int      `json:"igmp_groupmembership,omitempty"` // [2-9]|[1-9][0-9]{1,2}|[1-2][0-9]{3}|3[0-5][0-9]{2}|3600|^$
	IgmpMaxresponse         int      `json:"igmp_maxresponse,omitempty"`     // [1-9]|1[0-9]|2[0-5]|^$
	IgmpMcrtrexpiretime     int      `json:"igmp_mcrtrexpiretime,omitempty"` // [0-9]|[1-9][0-9]{1,2}|[1-2][0-9]{3}|3[0-5][0-9]{2}|3600|^$
	IgmpQuerier             string   `json:"igmp_querier"`                   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	IgmpSnooping            bool     `json:"igmp_snooping"`
	IgmpSupression          bool     `json:"igmp_supression"`
	IPSubnet                string   `json:"ip_subnet,omitempty"`      // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|30)$
	IPsecDhGroup            int      `json:"ipsec_dh_group,omitempty"` // 2|5|14|15|16|19|20|21|25|26
	IPsecDynamicRouting     bool     `json:"ipsec_dynamic_routing"`
	IPsecEncryption         string   `json:"ipsec_encryption,omitempty"`   // aes128|aes192|aes256|3des
	IPsecEspDhGroup         int      `json:"ipsec_esp_dh_group,omitempty"` // 1|2|5|14|15|16|17|18
	IPsecHash               string   `json:"ipsec_hash,omitempty"`         // sha1|md5|sha256|sha384|sha512
	IPsecIkeDhGroup         int      `json:"ipsec_ike_dh_group,omitempty"` // 1|2|5|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|32
	IPsecInterface          string   `json:"ipsec_interface,omitempty"`    // wan|wan2
	IPsecKeyExchange        string   `json:"ipsec_key_exchange,omitempty"` // ikev1|ikev2
	IPsecLocalIP            string   `json:"ipsec_local_ip,omitempty"`     // ^any$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	IPsecPeerIP             string   `json:"ipsec_peer_ip,omitempty"`      // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	IPsecPfs                bool     `json:"ipsec_pfs"`
	IPsecProfile            string   `json:"ipsec_profile,omitempty"`       // customized|azure_dynamic|azure_static
	IPV6InterfaceType       string   `json:"ipv6_interface_type,omitempty"` // static|pd|none
	IPV6PDInterface         string   `json:"ipv6_pd_interface,omitempty"`   // wan|wan2
	IPV6PDPrefixid          string   `json:"ipv6_pd_prefixid"`              // ^$|[a-fA-F0-9]{1,4}
	IPV6PDStart             string   `json:"ipv6_pd_start,omitempty"`
	IPV6PDStop              string   `json:"ipv6_pd_stop,omitempty"`
	IPV6RaEnabled           bool     `json:"ipv6_ra_enabled"`
	IPV6RaPreferredLifetime int      `json:"ipv6_ra_preferred_lifetime,omitempty"` // ^([0-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-8][0-9]{4}|9[0-8][0-9]{3}|99[0-8][0-9]{2}|999[0-8][0-9]|9999[0-9]|[1-8][0-9]{5}|9[0-8][0-9]{4}|99[0-8][0-9]{3}|999[0-8][0-9]{2}|9999[0-8][0-9]|99999[0-9]|[1-8][0-9]{6}|9[0-8][0-9]{5}|99[0-8][0-9]{4}|999[0-8][0-9]{3}|9999[0-8][0-9]{2}|99999[0-8][0-9]|999999[0-9]|[12][0-9]{7}|30[0-9]{6}|31[0-4][0-9]{5}|315[0-2][0-9]{4}|3153[0-5][0-9]{3}|31536000)$|^$
	IPV6RaPriority          string   `json:"ipv6_ra_priority,omitempty"`           // high|medium|low
	IPV6RaValidLifetime     int      `json:"ipv6_ra_valid_lifetime,omitempty"`     // ^([0-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-8][0-9]{4}|9[0-8][0-9]{3}|99[0-8][0-9]{2}|999[0-8][0-9]|9999[0-9]|[1-8][0-9]{5}|9[0-8][0-9]{4}|99[0-8][0-9]{3}|999[0-8][0-9]{2}|9999[0-8][0-9]|99999[0-9]|[1-8][0-9]{6}|9[0-8][0-9]{5}|99[0-8][0-9]{4}|999[0-8][0-9]{3}|9999[0-8][0-9]{2}|99999[0-8][0-9]|999999[0-9]|[12][0-9]{7}|30[0-9]{6}|31[0-4][0-9]{5}|315[0-2][0-9]{4}|3153[0-5][0-9]{3}|31536000)$|^$
	IPV6Subnet              string   `json:"ipv6_subnet,omitempty"`
	IsNAT                   bool     `json:"is_nat"`
	L2TpInterface           string   `json:"l2tp_interface,omitempty"` // wan|wan2
	LteLanEnabled           bool     `json:"lte_lan_enabled"`
	Name                    string   `json:"name,omitempty"` // .{1,128}
	NATOutboundIP           string   `json:"nat_outbound_ip,omitempty"`
	NetworkGroup            string   `json:"networkgroup,omitempty"`           // LAN[2-8]?
	OpenVPNLocalAddress     string   `json:"openvpn_local_address,omitempty"`  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNLocalPort        int      `json:"openvpn_local_port,omitempty"`     // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	OpenVPNMode             string   `json:"openvpn_mode,omitempty"`           // site-to-site|client|server
	OpenVPNRemoteAddress    string   `json:"openvpn_remote_address,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNRemoteHost       string   `json:"openvpn_remote_host,omitempty"`    // [^\"\' ]+|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNRemotePort       int      `json:"openvpn_remote_port,omitempty"`    // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	PptpcRequireMppe        bool     `json:"pptpc_require_mppe"`
	PptpcRouteDistance      int      `json:"pptpc_route_distance,omitempty"` // ^[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]$|^$
	PptpcServerIP           string   `json:"pptpc_server_ip,omitempty"`      // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|(?=^.{3,253}$)(^((?!-)[a-zA-Z0-9-]{1,63}(?<!-)\.)+[a-zA-Z]{2,63}$)|^[a-zA-Z0-9-]{1,63}$
	PptpcUsername           string   `json:"pptpc_username,omitempty"`       // [^\"\' ]+
	Priority                int      `json:"priority,omitempty"`             // [1-4]
	Purpose                 string   `json:"purpose,omitempty"`              // corporate|guest|remote-user-vpn|site-vpn|vlan-only|vpn-client|wan
	RADIUSprofileID         string   `json:"radiusprofile_id"`
	RemoteSiteID            string   `json:"remote_site_id"`
	RemoteSiteSubnets       []string `json:"remote_site_subnets,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|30)$|^$
	RemoteVpnSubnets        []string `json:"remote_vpn_subnets,omitempty"`  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|30)$|^$
	ReportWanEvent          bool     `json:"report_wan_event"`
	RequireMschapv2         bool     `json:"require_mschapv2"`
	RouteDistance           int      `json:"route_distance,omitempty"` // ^[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]$|^$
	UpnpLanEnabled          bool     `json:"upnp_lan_enabled"`
	UserGroupID             string   `json:"usergroup_id"`
	VLAN                    int      `json:"vlan,omitempty"` // [2-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|400[0-9]|^$
	VLANEnabled             bool     `json:"vlan_enabled"`
	VpnClientDefaultRoute   bool     `json:"vpn_client_default_route"`
	VpnClientPullDNS        bool     `json:"vpn_client_pull_dns"`
	VpnType                 string   `json:"vpn_type,omitempty"`                // auto|ipsec-vpn|openvpn-vpn|pptp-client|l2tp-server|pptp-server
	WanDHCPv6PDSize         int      `json:"wan_dhcpv6_pd_size,omitempty"`      // ^(4[89]|5[0-9]|6[0-4])$|^$
	WanDNS1                 string   `json:"wan_dns1"`                          // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WanDNS2                 string   `json:"wan_dns2"`                          // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WanDNS3                 string   `json:"wan_dns3"`                          // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WanDNS4                 string   `json:"wan_dns4"`                          // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WanEgressQOS            int      `json:"wan_egress_qos,omitempty"`          // [1-7]|^$
	WanGateway              string   `json:"wan_gateway"`                       // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WanGatewayV6            string   `json:"wan_gateway_v6"`                    // ^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	WanIP                   string   `json:"wan_ip,omitempty"`                  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	WanIPV6                 string   `json:"wan_ipv6"`                          // ^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	WanLoadBalanceType      string   `json:"wan_load_balance_type,omitempty"`   // failover-only|weighted
	WanLoadBalanceWeight    int      `json:"wan_load_balance_weight,omitempty"` // [1-9]|[1-9][0-9]
	WanNetmask              string   `json:"wan_netmask,omitempty"`             // ^((128|192|224|240|248|252|254)\.0\.0\.0)|(255\.(((0|128|192|224|240|248|252|254)\.0\.0)|(255\.(((0|128|192|224|240|248|252|254)\.0)|255\.(0|128|192|224|240|248|252|254)))))$
	WanNetworkGroup         string   `json:"wan_networkgroup,omitempty"`        // WAN[2]?|WAN_LTE_FAILOVER
	WanPrefixlen            int      `json:"wan_prefixlen,omitempty"`           // ^([1-9]|[1-8][0-9]|9[0-9]|1[01][0-9]|12[0-8])$|^$
	WanSmartqDownRate       int      `json:"wan_smartq_down_rate,omitempty"`    // [0-9]{1,6}|1000000
	WanSmartqEnabled        bool     `json:"wan_smartq_enabled"`
	WanSmartqUpRate         int      `json:"wan_smartq_up_rate,omitempty"` // [0-9]{1,6}|1000000
	WanType                 string   `json:"wan_type,omitempty"`           // disabled|dhcp|static|pppoe
	WanTypeV6               string   `json:"wan_type_v6,omitempty"`        // disabled|dhcpv6|static
	WanUsername             string   `json:"wan_username,omitempty"`       // [^"' ]+
	WanVLAN                 int      `json:"wan_vlan,omitempty"`           // [0-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|40[0-8][0-9]|409[0-4]|^$
	WanVLANEnabled          bool     `json:"wan_vlan_enabled"`
	XIPsecPreSharedKey      string   `json:"x_ipsec_pre_shared_key,omitempty"`      // [^\"\' ]+
	XOpenVPNSharedSecretKey string   `json:"x_openvpn_shared_secret_key,omitempty"` // [0-9A-Fa-f]{512}
	XPptpcPassword          string   `json:"x_pptpc_password,omitempty"`            // [^\"\' ]+
	XWanPassword            string   `json:"x_wan_password,omitempty"`              // [^"' ]+
}
