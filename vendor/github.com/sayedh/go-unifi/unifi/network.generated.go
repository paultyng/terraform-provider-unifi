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

type Network struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	AutoScaleEnabled                              bool                            `json:"auto_scale_enabled"`
	DHCPDBootEnabled                              bool                            `json:"dhcpd_boot_enabled"`
	DHCPDBootFilename                             string                          `json:"dhcpd_boot_filename,omitempty"` // .{1,256}
	DHCPDBootServer                               string                          `json:"dhcpd_boot_server"`             // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$|(?=^.{3,253}$)(^((?!-)[a-zA-Z0-9-]{1,63}(?<!-)\.)+[a-zA-Z]{2,63}$)|[a-zA-Z0-9-]{1,63}|^$
	DHCPDConflictChecking                         bool                            `json:"dhcpd_conflict_checking"`
	DHCPDDNS1                                     string                          `json:"dhcpd_dns_1"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNS2                                     string                          `json:"dhcpd_dns_2"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNS3                                     string                          `json:"dhcpd_dns_3"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNS4                                     string                          `json:"dhcpd_dns_4"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDDNSEnabled                               bool                            `json:"dhcpd_dns_enabled"`
	DHCPDEnabled                                  bool                            `json:"dhcpd_enabled"`
	DHCPDGateway                                  string                          `json:"dhcpd_gateway"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDGatewayEnabled                           bool                            `json:"dhcpd_gateway_enabled"`
	DHCPDIP1                                      string                          `json:"dhcpd_ip_1"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDIP2                                      string                          `json:"dhcpd_ip_2"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDIP3                                      string                          `json:"dhcpd_ip_3"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDLeaseTime                                int                             `json:"dhcpd_leasetime,omitempty"`
	DHCPDMAC1                                     string                          `json:"dhcpd_mac_1"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	DHCPDMAC2                                     string                          `json:"dhcpd_mac_2"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	DHCPDMAC3                                     string                          `json:"dhcpd_mac_3"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	DHCPDNtp1                                     string                          `json:"dhcpd_ntp_1"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDNtp2                                     string                          `json:"dhcpd_ntp_2"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDNtpEnabled                               bool                            `json:"dhcpd_ntp_enabled"`
	DHCPDStart                                    string                          `json:"dhcpd_start"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDStop                                     string                          `json:"dhcpd_stop"`  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDTFTPServer                               string                          `json:"dhcpd_tftp_server,omitempty"`
	DHCPDTimeOffset                               int                             `json:"dhcpd_time_offset,omitempty"` // ^0$|^-?([1-9]([0-9]{1,3})?|[1-7][0-9]{4}|[8][0-5][0-9]{3}|86[0-3][0-9]{2}|86400)$
	DHCPDTimeOffsetEnabled                        bool                            `json:"dhcpd_time_offset_enabled"`
	DHCPDUnifiController                          string                          `json:"dhcpd_unifi_controller"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDV6AllowSlaac                             bool                            `json:"dhcpdv6_allow_slaac"`
	DHCPDV6DNS1                                   string                          `json:"dhcpdv6_dns_1,omitempty"`
	DHCPDV6DNS2                                   string                          `json:"dhcpdv6_dns_2,omitempty"`
	DHCPDV6DNS3                                   string                          `json:"dhcpdv6_dns_3,omitempty"`
	DHCPDV6DNS4                                   string                          `json:"dhcpdv6_dns_4,omitempty"`
	DHCPDV6DNSAuto                                bool                            `json:"dhcpdv6_dns_auto"`
	DHCPDV6Enabled                                bool                            `json:"dhcpdv6_enabled"`
	DHCPDV6LeaseTime                              int                             `json:"dhcpdv6_leasetime,omitempty"`
	DHCPDV6Start                                  string                          `json:"dhcpdv6_start,omitempty"`
	DHCPDV6Stop                                   string                          `json:"dhcpdv6_stop,omitempty"`
	DHCPDWPAdUrl                                  string                          `json:"dhcpd_wpad_url,omitempty"`
	DHCPDWins1                                    string                          `json:"dhcpd_wins_1"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDWins2                                    string                          `json:"dhcpd_wins_2"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPDWinsEnabled                              bool                            `json:"dhcpd_wins_enabled"`
	DHCPRelayEnabled                              bool                            `json:"dhcp_relay_enabled"`
	DHCPguardEnabled                              bool                            `json:"dhcpguard_enabled"`
	DPIEnabled                                    bool                            `json:"dpi_enabled"`
	DPIgroupID                                    string                          `json:"dpigroup_id"` // [\d\w]+|^$
	DomainName                                    string                          `json:"domain_name"` // (?=^.{3,253}$)(^((?!-)[a-zA-Z0-9-]{1,63}(?<!-)\.)+[a-zA-Z]{2,63}$)|^$|[a-zA-Z0-9-]{1,63}
	Enabled                                       bool                            `json:"enabled"`
	ExposedToSiteVPN                              bool                            `json:"exposed_to_site_vpn"`
	GatewayDevice                                 string                          `json:"gateway_device"`         // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	GatewayType                                   string                          `json:"gateway_type,omitempty"` // default|switch
	IGMPFastleave                                 bool                            `json:"igmp_fastleave"`
	IGMPGroupmembership                           int                             `json:"igmp_groupmembership,omitempty"` // [2-9]|[1-9][0-9]{1,2}|[1-2][0-9]{3}|3[0-5][0-9]{2}|3600|^$
	IGMPMaxresponse                               int                             `json:"igmp_maxresponse,omitempty"`     // [1-9]|1[0-9]|2[0-5]|^$
	IGMPMcrtrexpiretime                           int                             `json:"igmp_mcrtrexpiretime,omitempty"` // [0-9]|[1-9][0-9]{1,2}|[1-2][0-9]{3}|3[0-5][0-9]{2}|3600|^$
	IGMPProxyDownstream                           bool                            `json:"igmp_proxy_downstream"`
	IGMPProxyUpstream                             bool                            `json:"igmp_proxy_upstream"`
	IGMPQuerier                                   string                          `json:"igmp_querier"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	IGMPSnooping                                  bool                            `json:"igmp_snooping"`
	IGMPSupression                                bool                            `json:"igmp_supression"`
	IPSecDhGroup                                  int                             `json:"ipsec_dh_group,omitempty"` // 2|5|14|15|16|19|20|21|25|26
	IPSecDynamicRouting                           bool                            `json:"ipsec_dynamic_routing"`
	IPSecEncryption                               string                          `json:"ipsec_encryption,omitempty"`     // aes128|aes192|aes256|3des
	IPSecEspDhGroup                               int                             `json:"ipsec_esp_dh_group,omitempty"`   // 1|2|5|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|32
	IPSecEspEncryption                            string                          `json:"ipsec_esp_encryption,omitempty"` // aes128|aes192|aes256|3des
	IPSecEspHash                                  string                          `json:"ipsec_esp_hash,omitempty"`       // sha1|md5|sha256|sha384|sha512
	IPSecEspLifetime                              string                          `json:"ipsec_esp_lifetime,omitempty"`   // ^(?:3[0-9]|[4-9][0-9]|[1-9][0-9]{2,3}|[1-7][0-9]{4}|8[0-5][0-9]{3}|86[0-3][0-9]{2}|86400)$
	IPSecHash                                     string                          `json:"ipsec_hash,omitempty"`           // sha1|md5|sha256|sha384|sha512
	IPSecIkeDhGroup                               int                             `json:"ipsec_ike_dh_group,omitempty"`   // 1|2|5|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|32
	IPSecIkeEncryption                            string                          `json:"ipsec_ike_encryption,omitempty"` // aes128|aes192|aes256|3des
	IPSecIkeHash                                  string                          `json:"ipsec_ike_hash,omitempty"`       // sha1|md5|sha256|sha384|sha512
	IPSecIkeLifetime                              string                          `json:"ipsec_ike_lifetime,omitempty"`   // ^(?:3[0-9]|[4-9][0-9]|[1-9][0-9]{2,3}|[1-7][0-9]{4}|8[0-5][0-9]{3}|86[0-3][0-9]{2}|86400)$
	IPSecInterface                                string                          `json:"ipsec_interface,omitempty"`      // wan|wan2
	IPSecKeyExchange                              string                          `json:"ipsec_key_exchange,omitempty"`   // ikev1|ikev2
	IPSecLocalIDentifier                          string                          `json:"ipsec_local_identifier,omitempty"`
	IPSecLocalIDentifierEnabled                   bool                            `json:"ipsec_local_identifier_enabled"`
	IPSecLocalIP                                  string                          `json:"ipsec_local_ip,omitempty"` // ^any$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	IPSecPeerIP                                   string                          `json:"ipsec_peer_ip,omitempty"`
	IPSecPfs                                      bool                            `json:"ipsec_pfs"`
	IPSecProfile                                  string                          `json:"ipsec_profile,omitempty"` // customized|azure_dynamic|azure_static
	IPSecRemoteIDentifier                         string                          `json:"ipsec_remote_identifier,omitempty"`
	IPSecRemoteIDentifierEnabled                  bool                            `json:"ipsec_remote_identifier_enabled"`
	IPSecSeparateIkev2Networks                    bool                            `json:"ipsec_separate_ikev2_networks"`
	IPSecTunnelIP                                 string                          `json:"ipsec_tunnel_ip,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|3[0-2])$
	IPSecTunnelIPEnabled                          bool                            `json:"ipsec_tunnel_ip_enabled"`
	IPSubnet                                      string                          `json:"ip_subnet,omitempty"`                      // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|3[0-2])$
	IPV6ClientAddressAssignment                   string                          `json:"ipv6_client_address_assignment,omitempty"` // slaac|dhcpv6
	IPV6InterfaceType                             string                          `json:"ipv6_interface_type,omitempty"`            // static|pd|single_network|none
	IPV6PDAutoPrefixidEnabled                     bool                            `json:"ipv6_pd_auto_prefixid_enabled"`
	IPV6PDInterface                               string                          `json:"ipv6_pd_interface,omitempty"` // wan|wan2
	IPV6PDPrefixid                                string                          `json:"ipv6_pd_prefixid"`            // ^$|[a-fA-F0-9]{1,4}
	IPV6PDStart                                   string                          `json:"ipv6_pd_start,omitempty"`
	IPV6PDStop                                    string                          `json:"ipv6_pd_stop,omitempty"`
	IPV6RaEnabled                                 bool                            `json:"ipv6_ra_enabled"`
	IPV6RaPreferredLifetime                       int                             `json:"ipv6_ra_preferred_lifetime,omitempty"` // ^([0-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-8][0-9]{4}|9[0-8][0-9]{3}|99[0-8][0-9]{2}|999[0-8][0-9]|9999[0-9]|[1-8][0-9]{5}|9[0-8][0-9]{4}|99[0-8][0-9]{3}|999[0-8][0-9]{2}|9999[0-8][0-9]|99999[0-9]|[1-8][0-9]{6}|9[0-8][0-9]{5}|99[0-8][0-9]{4}|999[0-8][0-9]{3}|9999[0-8][0-9]{2}|99999[0-8][0-9]|999999[0-9]|[12][0-9]{7}|30[0-9]{6}|31[0-4][0-9]{5}|315[0-2][0-9]{4}|3153[0-5][0-9]{3}|31536000)$|^$
	IPV6RaPriority                                string                          `json:"ipv6_ra_priority,omitempty"`           // high|medium|low
	IPV6RaValidLifetime                           int                             `json:"ipv6_ra_valid_lifetime,omitempty"`     // ^([0-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-8][0-9]{4}|9[0-8][0-9]{3}|99[0-8][0-9]{2}|999[0-8][0-9]|9999[0-9]|[1-8][0-9]{5}|9[0-8][0-9]{4}|99[0-8][0-9]{3}|999[0-8][0-9]{2}|9999[0-8][0-9]|99999[0-9]|[1-8][0-9]{6}|9[0-8][0-9]{5}|99[0-8][0-9]{4}|999[0-8][0-9]{3}|9999[0-8][0-9]{2}|99999[0-8][0-9]|999999[0-9]|[12][0-9]{7}|30[0-9]{6}|31[0-4][0-9]{5}|315[0-2][0-9]{4}|3153[0-5][0-9]{3}|31536000)$|^$
	IPV6SettingPreference                         string                          `json:"ipv6_setting_preference,omitempty"`    // auto|manual
	IPV6SingleNetworkInterface                    string                          `json:"ipv6_single_network_interface,omitempty"`
	IPV6Subnet                                    string                          `json:"ipv6_subnet,omitempty"`
	IPV6WANDelegationType                         string                          `json:"ipv6_wan_delegation_type,omitempty"` // pd|single_network|none
	InterfaceMtu                                  int                             `json:"interface_mtu,omitempty"`            // ^(6[89]|[7-9][0-9]|[1-9][0-9]{2,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|65500)$
	InterfaceMtuEnabled                           bool                            `json:"interface_mtu_enabled"`
	InternetAccessEnabled                         bool                            `json:"internet_access_enabled"`
	IsNAT                                         bool                            `json:"is_nat"`
	L2TpAllowWeakCiphers                          bool                            `json:"l2tp_allow_weak_ciphers"`
	L2TpInterface                                 string                          `json:"l2tp_interface,omitempty"`    // wan|wan2
	L2TpLocalWANIP                                string                          `json:"l2tp_local_wan_ip,omitempty"` // ^any$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	LocalPort                                     int                             `json:"local_port,omitempty"`        // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	LteLanEnabled                                 bool                            `json:"lte_lan_enabled"`
	MACOverride                                   string                          `json:"mac_override"` // (^$|^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$)
	MACOverrideEnabled                            bool                            `json:"mac_override_enabled"`
	MdnsEnabled                                   bool                            `json:"mdns_enabled"`
	NATOutboundIPAddresses                        []NetworkNATOutboundIPAddresses `json:"nat_outbound_ip_addresses,omitempty"`
	Name                                          string                          `json:"name,omitempty"`         // .{1,128}
	NetworkGroup                                  string                          `json:"networkgroup,omitempty"` // LAN[2-8]?
	NetworkIsolationEnabled                       bool                            `json:"network_isolation_enabled"`
	OpenVPNConfiguration                          string                          `json:"openvpn_configuration,omitempty"`
	OpenVPNConfigurationFilename                  string                          `json:"openvpn_configuration_filename,omitempty"`
	OpenVPNEncryptionCipher                       string                          `json:"openvpn_encryption_cipher,omitempty"` // AES_256_GCM|AES_256_CBC|BF_CBC
	OpenVPNInterface                              string                          `json:"openvpn_interface,omitempty"`         // wan|wan2
	OpenVPNLocalAddress                           string                          `json:"openvpn_local_address,omitempty"`     // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNLocalPort                              int                             `json:"openvpn_local_port,omitempty"`        // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	OpenVPNLocalWANIP                             string                          `json:"openvpn_local_wan_ip,omitempty"`      // ^any$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNMode                                   string                          `json:"openvpn_mode,omitempty"`              // site-to-site|client|server
	OpenVPNRemoteAddress                          string                          `json:"openvpn_remote_address,omitempty"`    // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNRemoteHost                             string                          `json:"openvpn_remote_host,omitempty"`       // [^\"\' ]+|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	OpenVPNRemotePort                             int                             `json:"openvpn_remote_port,omitempty"`       // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	OpenVPNUsername                               string                          `json:"openvpn_username,omitempty"`
	PptpcRequireMppe                              bool                            `json:"pptpc_require_mppe"`
	PptpcRouteDistance                            int                             `json:"pptpc_route_distance,omitempty"` // ^[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]$|^$
	PptpcServerIP                                 string                          `json:"pptpc_server_ip,omitempty"`      // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|(?=^.{3,253}$)(^((?!-)[a-zA-Z0-9-]{1,63}(?<!-)\.)+[a-zA-Z]{2,63}$)|^[a-zA-Z0-9-]{1,63}$
	PptpcUsername                                 string                          `json:"pptpc_username,omitempty"`       // [^\"\' ]+
	Priority                                      int                             `json:"priority,omitempty"`             // [1-4]
	Purpose                                       string                          `json:"purpose,omitempty"`              // corporate|guest|remote-user-vpn|site-vpn|vlan-only|vpn-client|wan
	RADIUSProfileID                               string                          `json:"radiusprofile_id"`
	RemoteSiteID                                  string                          `json:"remote_site_id"`
	RemoteSiteSubnets                             []string                        `json:"remote_site_subnets,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|30)$|^$
	RemoteVPNDynamicSubnetsEnabled                bool                            `json:"remote_vpn_dynamic_subnets_enabled"`
	RemoteVPNSubnets                              []string                        `json:"remote_vpn_subnets,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|3[0-2])$|^$
	ReportWANEvent                                bool                            `json:"report_wan_event"`
	RequireMschapv2                               bool                            `json:"require_mschapv2"`
	RouteDistance                                 int                             `json:"route_distance,omitempty"`     // ^[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]$|^$
	SettingPreference                             string                          `json:"setting_preference,omitempty"` // auto|manual
	SingleNetworkLan                              string                          `json:"single_network_lan,omitempty"`
	UidPolicyEnabled                              bool                            `json:"uid_policy_enabled"`
	UidPolicyName                                 string                          `json:"uid_policy_name,omitempty"`
	UidPublicGatewayPort                          int                             `json:"uid_public_gateway_port,omitempty"` // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	UidTrafficRulesAllowedIPsAndHostnames         []string                        `json:"uid_traffic_rules_allowed_ips_and_hostnames,omitempty"`
	UidTrafficRulesEnabled                        bool                            `json:"uid_traffic_rules_enabled"`
	UidVPNCustomRouting                           []string                        `json:"uid_vpn_custom_routing,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|3[0-2])$
	UidVPNDefaultDNSSuffix                        string                          `json:"uid_vpn_default_dns_suffix,omitempty"`
	UidVPNMasqueradeEnabled                       bool                            `json:"uid_vpn_masquerade_enabled"`
	UidVPNMaxConnectionTimeSeconds                int                             `json:"uid_vpn_max_connection_time_seconds,omitempty"` // ^[1-9][0-9]*$
	UidVPNSyncPublicIP                            bool                            `json:"uid_vpn_sync_public_ip"`
	UidVPNType                                    string                          `json:"uid_vpn_type,omitempty"` // openvpn|wireguard
	UidWorkspaceUrl                               string                          `json:"uid_workspace_url,omitempty"`
	UpnpLanEnabled                                bool                            `json:"upnp_lan_enabled"`
	UserGroupID                                   string                          `json:"usergroup_id"`
	VLAN                                          int                             `json:"vlan,omitempty"` // [2-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|400[0-9]|401[0-8]|^$
	VLANEnabled                                   bool                            `json:"vlan_enabled"`
	VPNClientConfigurationRemoteIPOverride        string                          `json:"vpn_client_configuration_remote_ip_override,omitempty"`
	VPNClientConfigurationRemoteIPOverrideEnabled bool                            `json:"vpn_client_configuration_remote_ip_override_enabled"`
	VPNClientDefaultRoute                         bool                            `json:"vpn_client_default_route"`
	VPNClientPullDNS                              bool                            `json:"vpn_client_pull_dns"`
	VPNProtocol                                   string                          `json:"vpn_protocol,omitempty"`       // TCP|UDP
	VPNType                                       string                          `json:"vpn_type,omitempty"`           // auto|ipsec-vpn|openvpn-client|openvpn-server|openvpn-vpn|pptp-client|l2tp-server|pptp-server|uid-server|wireguard-server|wireguard-client
	VrrpIPSubnetGw1                               string                          `json:"vrrp_ip_subnet_gw1,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|30)$
	VrrpIPSubnetGw2                               string                          `json:"vrrp_ip_subnet_gw2,omitempty"` // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|30)$
	VrrpVrid                                      int                             `json:"vrrp_vrid,omitempty"`          // [1-9]|[1-9][0-9]
	WANDHCPCos                                    int                             `json:"wan_dhcp_cos,omitempty"`       // [0-7]|^$
	WANDHCPOptions                                []NetworkWANDHCPOptions         `json:"wan_dhcp_options,omitempty"`
	WANDHCPv6PDSize                               int                             `json:"wan_dhcpv6_pd_size,omitempty"` // ^(4[89]|5[0-9]|6[0-4])$|^$
	WANDNS1                                       string                          `json:"wan_dns1"`                     // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WANDNS2                                       string                          `json:"wan_dns2"`                     // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WANDNS3                                       string                          `json:"wan_dns3"`                     // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WANDNS4                                       string                          `json:"wan_dns4"`                     // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	WANDNSPreference                              string                          `json:"wan_dns_preference,omitempty"` // auto|manual
	WANDsliteRemoteHost                           string                          `json:"wan_dslite_remote_host,omitempty"`
	WANEgressQOS                                  int                             `json:"wan_egress_qos,omitempty"`          // [1-7]|^$
	WANGateway                                    string                          `json:"wan_gateway,omitempty"`             // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	WANGatewayV6                                  string                          `json:"wan_gateway_v6"`                    // ^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	WANIP                                         string                          `json:"wan_ip,omitempty"`                  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	WANIPAliases                                  []string                        `json:"wan_ip_aliases,omitempty"`          // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([8-9]|[1-2][0-9]|3[0-2])$|^$
	WANIPV6                                       string                          `json:"wan_ipv6"`                          // ^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	WANIPV6DNS1                                   string                          `json:"wan_ipv6_dns1"`                     // ^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	WANIPV6DNS2                                   string                          `json:"wan_ipv6_dns2"`                     // ^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$|^$
	WANIPV6DNSPreference                          string                          `json:"wan_ipv6_dns_preference,omitempty"` // auto|manual
	WANLoadBalanceType                            string                          `json:"wan_load_balance_type,omitempty"`   // failover-only|weighted
	WANLoadBalanceWeight                          int                             `json:"wan_load_balance_weight,omitempty"` // ^$|[1-9]|[1-9][0-9]
	WANNetmask                                    string                          `json:"wan_netmask,omitempty"`             // ^((128|192|224|240|248|252|254)\.0\.0\.0)|(255\.(((0|128|192|224|240|248|252|254)\.0\.0)|(255\.(((0|128|192|224|240|248|252|254)\.0)|255\.(0|128|192|224|240|248|252|254)))))$
	WANNetworkGroup                               string                          `json:"wan_networkgroup,omitempty"`        // WAN[2]?|WAN_LTE_FAILOVER
	WANPppoePasswordEnabled                       bool                            `json:"wan_pppoe_password_enabled"`
	WANPppoeUsernameEnabled                       bool                            `json:"wan_pppoe_username_enabled"`
	WANPrefixlen                                  int                             `json:"wan_prefixlen,omitempty"` // ^([1-9]|[1-8][0-9]|9[0-9]|1[01][0-9]|12[0-8])$|^$
	WANProviderCapabilities                       NetworkWANProviderCapabilities  `json:"wan_provider_capabilities,omitempty"`
	WANSmartqDownRate                             int                             `json:"wan_smartq_down_rate,omitempty"` // [0-9]{1,6}|1000000
	WANSmartqEnabled                              bool                            `json:"wan_smartq_enabled"`
	WANSmartqUpRate                               int                             `json:"wan_smartq_up_rate,omitempty"` // [0-9]{1,6}|1000000
	WANType                                       string                          `json:"wan_type,omitempty"`           // disabled|dhcp|static|pppoe|dslite
	WANTypeV6                                     string                          `json:"wan_type_v6,omitempty"`        // disabled|slaac|dhcpv6|static
	WANUsername                                   string                          `json:"wan_username"`                 // [^"' ]+|^$
	WANVLAN                                       int                             `json:"wan_vlan,omitempty"`           // [0-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|40[0-8][0-9]|409[0-4]|^$
	WANVLANEnabled                                bool                            `json:"wan_vlan_enabled"`
	WireguardClientConfigurationFile              string                          `json:"wireguard_client_configuration_file,omitempty"`
	WireguardClientConfigurationFilename          string                          `json:"wireguard_client_configuration_filename,omitempty"`
	WireguardClientMode                           string                          `json:"wireguard_client_mode,omitempty"` // file|manual
	WireguardClientPeerIP                         string                          `json:"wireguard_client_peer_ip,omitempty"`
	WireguardClientPeerPort                       int                             `json:"wireguard_client_peer_port,omitempty"` // ^([1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5])$
	WireguardClientPeerPublicKey                  string                          `json:"wireguard_client_peer_public_key,omitempty"`
	WireguardClientPresharedKey                   string                          `json:"wireguard_client_preshared_key,omitempty"`
	WireguardClientPresharedKeyEnabled            bool                            `json:"wireguard_client_preshared_key_enabled"`
	WireguardInterface                            string                          `json:"wireguard_interface,omitempty"`    // wan|wan2
	WireguardLocalWANIP                           string                          `json:"wireguard_local_wan_ip,omitempty"` // ^any$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	WireguardPublicKey                            string                          `json:"wireguard_public_key,omitempty"`
	XAuthKey                                      string                          `json:"x_auth_key,omitempty"`
	XCaCrt                                        string                          `json:"x_ca_crt,omitempty"`
	XCaKey                                        string                          `json:"x_ca_key,omitempty"`
	XDhKey                                        string                          `json:"x_dh_key,omitempty"`
	XIPSecPreSharedKey                            string                          `json:"x_ipsec_pre_shared_key,omitempty"` // [^\"\' ]+
	XOpenVPNPassword                              string                          `json:"x_openvpn_password,omitempty"`
	XOpenVPNSharedSecretKey                       string                          `json:"x_openvpn_shared_secret_key,omitempty"` // [0-9A-Fa-f]{512}
	XPptpcPassword                                string                          `json:"x_pptpc_password,omitempty"`            // [^\"\' ]+
	XServerCrt                                    string                          `json:"x_server_crt,omitempty"`
	XServerKey                                    string                          `json:"x_server_key,omitempty"`
	XSharedClientCrt                              string                          `json:"x_shared_client_crt,omitempty"`
	XSharedClientKey                              string                          `json:"x_shared_client_key,omitempty"`
	XWANPassword                                  string                          `json:"x_wan_password"` // [^"' ]+|^$
	XWireguardPrivateKey                          string                          `json:"x_wireguard_private_key,omitempty"`
}

func (dst *Network) UnmarshalJSON(b []byte) error {
	type Alias Network
	aux := &struct {
		DHCPDLeaseTime                 emptyStringInt `json:"dhcpd_leasetime"`
		DHCPDTimeOffset                emptyStringInt `json:"dhcpd_time_offset"`
		DHCPDV6LeaseTime               emptyStringInt `json:"dhcpdv6_leasetime"`
		IGMPGroupmembership            emptyStringInt `json:"igmp_groupmembership"`
		IGMPMaxresponse                emptyStringInt `json:"igmp_maxresponse"`
		IGMPMcrtrexpiretime            emptyStringInt `json:"igmp_mcrtrexpiretime"`
		IPSecDhGroup                   emptyStringInt `json:"ipsec_dh_group"`
		IPSecEspDhGroup                emptyStringInt `json:"ipsec_esp_dh_group"`
		IPSecIkeDhGroup                emptyStringInt `json:"ipsec_ike_dh_group"`
		IPV6RaPreferredLifetime        emptyStringInt `json:"ipv6_ra_preferred_lifetime"`
		IPV6RaValidLifetime            emptyStringInt `json:"ipv6_ra_valid_lifetime"`
		InterfaceMtu                   emptyStringInt `json:"interface_mtu"`
		InternetAccessEnabled          *bool          `json:"internet_access_enabled"`
		LocalPort                      emptyStringInt `json:"local_port"`
		OpenVPNLocalPort               emptyStringInt `json:"openvpn_local_port"`
		OpenVPNRemotePort              emptyStringInt `json:"openvpn_remote_port"`
		PptpcRouteDistance             emptyStringInt `json:"pptpc_route_distance"`
		Priority                       emptyStringInt `json:"priority"`
		RouteDistance                  emptyStringInt `json:"route_distance"`
		UidPublicGatewayPort           emptyStringInt `json:"uid_public_gateway_port"`
		UidVPNMaxConnectionTimeSeconds emptyStringInt `json:"uid_vpn_max_connection_time_seconds"`
		VLAN                           emptyStringInt `json:"vlan"`
		VrrpVrid                       emptyStringInt `json:"vrrp_vrid"`
		WANDHCPCos                     emptyStringInt `json:"wan_dhcp_cos"`
		WANDHCPv6PDSize                emptyStringInt `json:"wan_dhcpv6_pd_size"`
		WANEgressQOS                   emptyStringInt `json:"wan_egress_qos"`
		WANLoadBalanceWeight           emptyStringInt `json:"wan_load_balance_weight"`
		WANPrefixlen                   emptyStringInt `json:"wan_prefixlen"`
		WANSmartqDownRate              emptyStringInt `json:"wan_smartq_down_rate"`
		WANSmartqUpRate                emptyStringInt `json:"wan_smartq_up_rate"`
		WANVLAN                        emptyStringInt `json:"wan_vlan"`
		WireguardClientPeerPort        emptyStringInt `json:"wireguard_client_peer_port"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.DHCPDLeaseTime = int(aux.DHCPDLeaseTime)
	dst.DHCPDTimeOffset = int(aux.DHCPDTimeOffset)
	dst.DHCPDV6LeaseTime = int(aux.DHCPDV6LeaseTime)
	dst.IGMPGroupmembership = int(aux.IGMPGroupmembership)
	dst.IGMPMaxresponse = int(aux.IGMPMaxresponse)
	dst.IGMPMcrtrexpiretime = int(aux.IGMPMcrtrexpiretime)
	dst.IPSecDhGroup = int(aux.IPSecDhGroup)
	dst.IPSecEspDhGroup = int(aux.IPSecEspDhGroup)
	dst.IPSecIkeDhGroup = int(aux.IPSecIkeDhGroup)
	dst.IPV6RaPreferredLifetime = int(aux.IPV6RaPreferredLifetime)
	dst.IPV6RaValidLifetime = int(aux.IPV6RaValidLifetime)
	dst.InterfaceMtu = int(aux.InterfaceMtu)
	dst.InternetAccessEnabled = emptyBoolToTrue(aux.InternetAccessEnabled)
	dst.LocalPort = int(aux.LocalPort)
	dst.OpenVPNLocalPort = int(aux.OpenVPNLocalPort)
	dst.OpenVPNRemotePort = int(aux.OpenVPNRemotePort)
	dst.PptpcRouteDistance = int(aux.PptpcRouteDistance)
	dst.Priority = int(aux.Priority)
	dst.RouteDistance = int(aux.RouteDistance)
	dst.UidPublicGatewayPort = int(aux.UidPublicGatewayPort)
	dst.UidVPNMaxConnectionTimeSeconds = int(aux.UidVPNMaxConnectionTimeSeconds)
	dst.VLAN = int(aux.VLAN)
	dst.VrrpVrid = int(aux.VrrpVrid)
	dst.WANDHCPCos = int(aux.WANDHCPCos)
	dst.WANDHCPv6PDSize = int(aux.WANDHCPv6PDSize)
	dst.WANEgressQOS = int(aux.WANEgressQOS)
	dst.WANLoadBalanceWeight = int(aux.WANLoadBalanceWeight)
	dst.WANPrefixlen = int(aux.WANPrefixlen)
	dst.WANSmartqDownRate = int(aux.WANSmartqDownRate)
	dst.WANSmartqUpRate = int(aux.WANSmartqUpRate)
	dst.WANVLAN = int(aux.WANVLAN)
	dst.WireguardClientPeerPort = int(aux.WireguardClientPeerPort)

	return nil
}

type NetworkNATOutboundIPAddresses struct {
	IPAddress       string   `json:"ip_address"`                  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	IPAddressPool   []string `json:"ip_address_pool,omitempty"`   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])-(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$
	Mode            string   `json:"mode,omitempty"`              // all|ip_address|ip_address_pool
	WANNetworkGroup string   `json:"wan_network_group,omitempty"` // WAN|WAN2
}

func (dst *NetworkNATOutboundIPAddresses) UnmarshalJSON(b []byte) error {
	type Alias NetworkNATOutboundIPAddresses
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

type NetworkWANDHCPOptions struct {
	OptionNumber int    `json:"optionNumber,omitempty"` // ([1-9]|[1-8][0-9]|9[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])
	Value        string `json:"value,omitempty"`
}

func (dst *NetworkWANDHCPOptions) UnmarshalJSON(b []byte) error {
	type Alias NetworkWANDHCPOptions
	aux := &struct {
		OptionNumber emptyStringInt `json:"optionNumber"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.OptionNumber = int(aux.OptionNumber)

	return nil
}

type NetworkWANProviderCapabilities struct {
	DownloadKilobitsPerSecond int `json:"download_kilobits_per_second,omitempty"` // ^[1-9][0-9]*$
	UploadKilobitsPerSecond   int `json:"upload_kilobits_per_second,omitempty"`   // ^[1-9][0-9]*$
}

func (dst *NetworkWANProviderCapabilities) UnmarshalJSON(b []byte) error {
	type Alias NetworkWANProviderCapabilities
	aux := &struct {
		DownloadKilobitsPerSecond emptyStringInt `json:"download_kilobits_per_second"`
		UploadKilobitsPerSecond   emptyStringInt `json:"upload_kilobits_per_second"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.DownloadKilobitsPerSecond = int(aux.DownloadKilobitsPerSecond)
	dst.UploadKilobitsPerSecond = int(aux.UploadKilobitsPerSecond)

	return nil
}

func (c *Client) listNetwork(ctx context.Context, site string) ([]Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/networkconf", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getNetwork(ctx context.Context, site, id string) (*Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/networkconf/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteNetwork(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/networkconf/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createNetwork(ctx context.Context, site string, d *Network) (*Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/networkconf", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateNetwork(ctx context.Context, site string, d *Network) (*Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/networkconf/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
