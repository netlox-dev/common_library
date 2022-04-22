package commonLib

import (
	"regexp"
	"strconv"
	"strings"
)

type HeartbeatReturnModel struct {
	ServerVersion   string `json:"server_version,omitempty"`
	ResetGUID       string `json:"reset_GUID,omitempty"`
	ResetTime       string `json:"reset_time,omitempty"`
	RoutesAvailable int    `json:"routes_available,omitempty"`
}

type ConfigResetStatusModel struct {
	ResetStatus string `json:"reset_status,omitempty"`
}

type RouteModel struct {
	Cmd         string `json:"cmd,omitempty"`
	IPPrefix    string `json:"ip_prefix"`
	IfName      string `json:"ifname,omitempty"`
	NextHopType string `json:"nexthop_type,omitempty"`
	NextHop     string `json:"nexthop"`
	MACAddress  string `json:"mac_address,omitempty"`
	Vnid        int    `json:"vnid,omitempty"`
	Error_code  int    `json:"error_code,omitempty"`
	Error_msg   string `json:"error_msg,omitempty"`
}

type RouteReturnModel struct {
	Failed []RouteModel `json:"failed,omitempty"`
}

type InterfaceModel struct {
	AdminState string `json:"admin-state"`
}

type InterfaceReturnModel struct {
	Port string         `json:"port"`
	Attr InterfaceModel `json:"attr"`
}

type VlanModel struct {
	VlanName string             `json:"vlan_name"`
	VlanID   int                `json:"vlan_id"`
	Vrf      int                `json:"vrf"`
	Mac      string             `json:"mac_addr"`
	IPPrefix []string           `json:"ip_addr,omitempty"`
	Member   []VlanMembersModel `json:"member,omitempty"`
}

type VlanReturnModel struct {
	VlanID int       `json:"vlan_id"`
	Attr   VlanModel `json:"attr"`
}

type VlansReturnModel struct {
	Attr []VlanModel `json:"attr"`
}

type VlanMemberModel struct {
	Tagging string `json:"tagging_mode"`
}

type VlanMemberReturnModel struct {
	VlanID  int             `json:"vlan_id"`
	If_name string          `json:"if_name"`
	Attr    VlanMemberModel `json:"attr"`
}

type VlanMembersModel struct {
	If_name string `json:"if_name"`
	Tagging string `json:"tagging_mode"`
}

type VlanMembersReturnModel struct {
	VlanID int                `json:"vlan_id"`
	Attr   []VlanMembersModel `json:"attr"`
}

type VlanNeighborReturnModel struct {
	VlanID  int    `json:"vlan_id"`
	Ip_addr string `json:"ip_addr"`
}

type VlanNeighborsModel struct {
	Ip_addr string `json:"ip_addr"`
}

type VlanNeighborsReturnModel struct {
	VlanID int                  `json:"vlan_id"`
	Attr   []VlanNeighborsModel `json:"attr"`
}

type VlansPerVnetReturnModel struct {
	Vnet_id string              `json:"vnet_id,omitempty"`
	Attr    []VlansPerVnetModel `json:"attr"`
}

type VlansPerVnetModel struct {
	VlanID   int    `json:"vlan_id"`
	IPPrefix string `json:"ip_prefix,omitempty"`
}

type TunnelDecapModel struct {
	IPAddr string `json:"ip_addr"`
}

type TunnelDecapReturnModel struct {
	TunnelType string           `json:"tunnel_type"`
	Attr       TunnelDecapModel `json:"attr"`
}

type VnetModel struct {
	Vnid int `json:"vnid"`
}

type VnetReturnModel struct {
	VnetName string    `json:"vnet_id"`
	Attr     VnetModel `json:"attr"`
}

type EthernetIpaddressModel struct {
	Ip_addr string `json:"ip_addr"`
	Prefix  int    `json:"prefix"`
}

type EthernetNeighborModel struct {
	Ip_addr string `json:"ip_addr"`
	Mac     string `json:"mac_addr"`
}

type EthernetNeighborCLIModel struct {
	Vrf       int    `json:"vrf"`
	IpAddress string `json:"ip_addr"`
	SrcMac    string `json:"src_mac_addr"`
	DstMac    string `json:"dst_mac_addr"`
	Interface string `json:"interface"`
	Status    string `json:"status"`
	IsInstall string `json:"is_install"`
}

type EthernetNeighborReturnModel struct {
	Attr []EthernetNeighborCLIModel `json:"attr"`
}

type EthernetNeighborDeleteModel struct {
	Ip_addr string `json:"ip_addr"`
}

type EthernetMacModel struct {
	VlanID int    `json:"vlan_id"`
	Mac    string `json:"mac_addr"`
}

type IProuteReturnModel struct {
	Attr []IProuteCLIModel `json:"attr"`
}

type IProuteCLIModel struct {
	IPPrefix   string `json:"ip_prefix"`
	Flags      string `json:"flags,omitempty"`
	Vrf        int    `json:"vrf"`
	Interface  string `json:"interface,omitempty"`
	NextHop    string `json:"next_hop,omitempty"`
	IsInstall  string `json:"is_install"`
	IsHareWare string `json:"hw_install"`
}

type IProuteModel struct {
	Ip_addr  string `json:"ip_addr"`
	Prefix   int    `json:"prefix"`
	Next_hop string `json:"next_hop"`
}

type IProuteDeleteModel struct {
	Ip_addr string `json:"ip_addr"`
	Prefix  int    `json:"prefix"`
}

type VxlanBridgesModel struct {
	Uif_name string `json:"uif"`
	VlanID   *int   `json:"vlan_id,omitempty"`
}

type VxlanPeerModel struct {
	PeerIP string `json:"peer_ip"`
}

type VxlanMemberModel struct {
	If_name string `json:"if_name"`
	Tagging string `json:"tagging_mode"`
}

type VxlanMACModel struct {
	Mac     string `json:"mac_addr"`
	Ip_addr string `json:"ip_addr"`
}

type VxlanModel struct {
	VxlanName string `json:"vxlan_name"`
	VxlanID   int    `json:"vxlan_id"`
	Vrf       int    `json:"vrf"`
	LocalIP   string `json:"local_ip_addr"`
	Udp       string `json:"udp"`
	Uif       string `json:"uif"`
	Mac       string `json:"mac_addr"`
}

type VxlanReturnModel struct {
	Attr []VxlanModel `json:"attr"`
}

type VlanFDBModel struct {
	Mac    string `json:"mac_addr"`
	VlanID int    `json:"vxlan_id"`
	Vrf    int    `json:"vrf"`
	Ifname string `json:"if_name"`
	Status string `json:"status"`
}

type VlanFDBReturnModel struct {
	Attr []VlanFDBModel `json:"attr"`
}

type VxlanFDBModel struct {
	Mac     string `json:"mac_addr"`
	VxlanID int    `json:"vxlan_id"`
	Vrf     int    `json:"vrf"`
	DstIP   string `json:"destination_ip"`
	VtepID  string `json:"vtep_id"`
	Status  string `json:"status"`
}

type VxlanFDBReturnModel struct {
	Attr []VxlanFDBModel `json:"attr"`
}

type IpsecModel struct {
	Source_ip      string `json:"source_ip"`
	Destination_ip string `json:"destination_ip"`
	Protocol       string `json:"protocol"`
	SPI            string `json:"spi"`
	Direction      string `json:"direction"`
	Reqid          string `json:"reqid"`
	Security_key   string `json:"security_key"`
	Key_length     int    `json:"key_length"`
	Interface      string `json:"interface"`
}

type IpsecDeleteModel struct {
	Source_ip      string `json:"source_ip"`
	Destination_ip string `json:"destination_ip"`
	Protocol       string `json:"protocol"`
	SPI            string `json:"spi"`
	Direction      string `json:"direction"`
}

type LoxilightMgmtReturnModel struct {
	LoxilightMgmtIp string `json:"loxilightmgmtIp"`
}
type PingRequestModel struct {
	IpAddress string `json:"ip_addr"`
	VnetId    string `json:"vnet_id"`
	Count     string `json:"count"`
}

type PingReturnModel struct {
	PacketsTransmitted string `json:"packets_transmitted"`
	PacketsReceived    string `json:"packets_received"`
	MinRTT             string `json:"min_rtt"`
	MaxRTT             string `json:"max_rtt"`
	AvgRTT             string `json:"avg_rtt"`
}

type ErrorInner struct {
	Code    int      `json:"code"`
	SubCode *int     `json:"sub-code,omitempty"`
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"`
	Details string   `json:"details,omitempty"`
}

type ErrorModel struct {
	Error ErrorInner `json:"error"`
}

type MissingValueError struct {
	Field string
}

type InvalidFormatError struct {
	Field   string
	Message string
}

type ACLReturnModels struct {
	Attr []ACLReturnModel `json:"attr"`
}

type ACLReturnModel struct {
	UUID      string            `json:"uuid"`
	Table     string            `json:"table_id"`
	Filters   ACLFilterModel    `json:"filters"`
	Actions   []ACLValueModel   `json:"actions"`
	Statistic ACLStatisticModel `json:"statistic"`
}

type ACLModel struct {
	Filters ACLFilterModel  `json:"filters,omitempty"`
	Actions []ACLValueModel `json:"actions"`
}

type ACLFilterModel struct {
	EthernetType    string `json:"ethernet_type,omitempty"`
	DestinationIp   string `json:"destination_ip,omitempty"`
	SourceIp        string `json:"source_ip,omitempty"`
	DestinationPort int    `json:"destination_port,omitempty"`
	SourcePort      int    `json:"source_port,omitempty"`
	Protocol        int    `json:"protocol,omitempty"`
	Port            int    `json:"port,omitempty"`
}

type ACLStatisticModel struct {
	Packet int `json:"packets"`
	Bytes  int `json:"bytes"`
}

type ACLValueModel struct {
	ActionName string `json:"action_name"`
	IpAddress  string `json:"ip_address,omitempty"`
	Port       int    `json:"port,omitempty"`
	Weight     int    `json:"weight,omitempty"`
}

type ConntrackModel struct {
	Date            string `json:"date"`
	Time            string `json:"time"`
	DestinationIp   string `json:"destination_ip"`
	SourceIp        string `json:"source_ip"`
	DestinationPort int    `json:"destination_port"`
	SourcePort      int    `json:"source_port"`
	Protocol        string `json:"protocol"`
	Timeout         string `json:"timeout"`
	ExpireTime      string `json:"expire_time"`
	Status          string `json:"status"`
	PacketBtye      int    `json:"byte"`
	PacketCount     int    `json:"count"`
}

type ConntrackReturnModel struct {
	Attr []ConntrackModel `json:"attr"`
}

type VifModel struct {
	IfName   string   `json:"if_name"`
	Mac      string   `json:"mac_addr"`
	Port     int      `json:"port"`
	Flags    int      `json:"flags"`
	Status   string   `json:"status"`
	Rvid     int      `json:"rvid"`
	IPPrefix []string `json:"ip_addr,omitempty"`
}

type VifReturnModel struct {
	Attr []VifModel `json:"attr"`
}

func (e *MissingValueError) Error() string {
	return "JSON missing field: " + (*e).Field
}

func (e *InvalidFormatError) Error() string {
	return (*e).Message
}

func (a *ACLReturnModel) getFilters(filter_val string) {
	//T0,ip,sip:0.2.3.4/24,ipprot:0x6,iif:2
	filters := strings.Split(filter_val, ",")
	for _, filter := range filters {
		id_val := strings.Split(filter, ":")
		if len(id_val) == 1 && id_val[0] == "T0" {
			a.Table = id_val[0]
		} else if len(id_val) == 1 {
			a.Filters.EthernetType = id_val[0]
		} else if len(id_val) == 2 {
			switch id_val[0] {
			case "ipprot":
				id_val[1] = strings.Replace(id_val[1], "0x", "", -1)
				proto, _ := strconv.ParseInt(id_val[1], 16, 64)
				a.Filters.Protocol = int(proto)
			case "dip":
				a.Filters.DestinationIp = id_val[1]
			case "sip":
				a.Filters.SourceIp = id_val[1]
			case "dport":
				a.Filters.DestinationPort, _ = strconv.Atoi(id_val[1])
			case "sport":
				a.Filters.SourcePort, _ = strconv.Atoi(id_val[1])
			case "iif":
				a.Filters.Port, _ = strconv.Atoi(id_val[1])
			}
		} else {
			return
		}
	}

}

func (a *ACLReturnModel) getActions(actions string) {
	action_list := strings.Split(actions, ",")

	for _, action := range action_list {
		values := strings.Split(action, ":")
		output := ACLValueModel{}
		if len(values) == 1 {
			output.ActionName = values[0]
		} else if len(values) == 3 {
			output.ActionName = values[0]
			output.IpAddress = values[1]
			output.Port, _ = strconv.Atoi(values[2])
		} else if len(values) == 4 {
			output.ActionName = values[0]
			output.IpAddress = values[1]
			output.Port, _ = strconv.Atoi(values[2])
			output.Weight, _ = strconv.Atoi(values[3])
		}
		a.Actions = append(a.Actions, output)
	}

}

func (a *ACLReturnModel) getStats(stats string) {
	re := regexp.MustCompile("(.+)pkts,(.+)bytes")
	match := re.FindStringSubmatch(stats)
	if len(match) > 1 {
		a.Statistic.Packet, _ = strconv.Atoi(match[0])
		a.Statistic.Bytes, _ = strconv.Atoi(match[1])
	}
}
