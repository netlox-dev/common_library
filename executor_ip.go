package commonLib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.nanomsg.org/mangos/v3"
)

// Title Declare
var TitleIpBrief = []string{"Interface", "IP-Address", "OK?", "Method", "Status", "Protocol"}
var TitleAutoneg = []string{"Interface", "Auto-Neg Mode", "Speed", "Adv Speeds", "Type", "Adv Types", "Oper", "Admin"}
var TitleBreakout = []string{"Interface", "Current Breakout Mode"}
var TitleCounters = []string{"IFACE", "STATE", "RX_OK", "RX_BPS", "RX_UTIL", "RX_ERR", "RX_DRP", "RX_OVR", "TX_OK", "TX_BPS", "TX_UTIL", "TX_ERR", "TX_DRP", "TX_OVR"}

var TitleRouteTable = []string{"VRF", "Route", "Is Install", "NextHop", "interface"}
var TitleNeighborTable = []string{"VRF", "Next hop", "Source MAC", "Destinaion MAC", "Interface", "Is Install", "HW Install"}

func RouteExtract(raw string) []string {
	re1 := regexp.MustCompile(`VRF (.+?) (.+?) \((.+?)\) nexthop via (.+?)\tdev (.+)`)
	re2 := regexp.MustCompile(`VRF (.+?) (.+?) \((.+?)\) \((.+?)\)`)
	match1 := re1.FindStringSubmatch(raw)
	match2 := re2.FindStringSubmatch(raw)
	if len(match1) > 0 {
		return match1[1:]
	} else if len(match2) > 0 {
		return match2[1:]
	}
	return match1
}

func RouteOnlyExtract(raw string) []string {
	re := regexp.MustCompile("Host-IP (.+?) ")
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func NeighborExtract(raw string) []string {
	re := regexp.MustCompile(`VRF (.+?) nh-ip (.+?) smac (.+?) dmac (.+?) dev (.+?) ([A-Z]+)\t\((.+?)\)`)

	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func GetIpRouteConfig() string {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return ""
	}
	// make data format
	var hdr []byte
	//var row_data []string
	//var data [][]string
	cmd := uint8(LOXILIGHT_RT_SHOW_ALL)
	_, hdr = MakeMessage(cmd, "")

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res
}

func GetIpNeighborModel() EthernetNeighborReturnModel {
	var sock mangos.Socket
	var EthernetNeighbor []EthernetNeighborCLIModel
	var EthernetNeighborReturn EthernetNeighborReturnModel

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return EthernetNeighborReturn
	}
	var hdr []byte
	// make data format
	cmd := uint8(LOXILIGHT_NH_SHOW_ALL)
	_, hdr = MakeMessage(cmd, "")

	// send msg and return value
	res := SendMessage(sock, hdr)
	for _, v := range strings.Split(res, "\r\n") {
		raw_data := NeighborExtract(v)
		if len(raw_data) > 0 {
			VrfID, _ := strconv.Atoi(raw_data[0])
			output := EthernetNeighborCLIModel{
				Vrf:       VrfID,
				IpAddress: raw_data[1],
				SrcMac:    raw_data[2],
				DstMac:    raw_data[3],
				Interface: raw_data[4],
				Status:    raw_data[5],
				IsInstall: raw_data[6],
			}
			EthernetNeighbor = append(EthernetNeighbor, output)
		}
	}
	EthernetNeighborReturn.Attr = EthernetNeighbor
	return EthernetNeighborReturn
}

func GetIpNeighborOneModel(interface_name string) EthernetNeighborReturnModel {
	var sock mangos.Socket

	var EthernetNeighbor []EthernetNeighborCLIModel
	var EthernetNeighborReturn EthernetNeighborReturnModel
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return EthernetNeighborReturn
	}
	// make data format
	var hdr []byte
	//var row_data []string
	//var data [][]string
	msg := interface_name

	if interface_name == "all" {
		cmd := uint8(LOXILIGHT_NH_SHOW_ALL)
		_, hdr = MakeMessage(cmd, "")
	} else {
		cmd := uint8(LOXILIGHT_NH_SHOW_INTERFACE)
		_, hdr = MakeMessage(cmd, msg)
	}
	// send msg and return value
	res := SendMessage(sock, hdr)
	for _, v := range strings.Split(res, "\r\n") {
		raw_data := NeighborExtract(v)
		if len(raw_data) > 0 {

			VrfID, _ := strconv.Atoi(raw_data[0])
			output := EthernetNeighborCLIModel{
				Vrf:       VrfID,
				IpAddress: raw_data[1],
				SrcMac:    raw_data[2],
				DstMac:    raw_data[3],
				Interface: raw_data[4],
				Status:    raw_data[5],
				IsInstall: raw_data[6],
			}
			EthernetNeighbor = append(EthernetNeighbor, output)
		}
	}
	EthernetNeighborReturn.Attr = EthernetNeighbor
	return EthernetNeighborReturn

}

func GetIpNeighborConfig() string {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return ""
	}
	var hdr []byte
	// make data format
	cmd := uint8(LOXILIGHT_NH_SHOW_ALL)
	_, hdr = MakeMessage(cmd, "")

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res

}

func GetIpNeighborOneConfig(interface_name string) string {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return ""
	}
	// make data format
	var hdr []byte
	//var row_data []string
	//var data [][]string
	msg := interface_name

	if interface_name == "all" {
		cmd := uint8(LOXILIGHT_NH_SHOW_ALL)
		_, hdr = MakeMessage(cmd, "")
	} else {
		cmd := uint8(LOXILIGHT_NH_SHOW_INTERFACE)
		_, hdr = MakeMessage(cmd, msg)
	}
	// send msg and return value
	res := SendMessage(sock, hdr)
	return res

}

//GetIPRouteModel 은 모든 Routing Table을 struct로 제공합니다.
func GetIPRouteModel() IProuteReturnModel {
	var IProuteCLI []IProuteCLIModel
	var IProuteReturn IProuteReturnModel

	route := GetIpRouteConfig()
	for _, v := range strings.Split(route, "\r\n") {
		raw_data := RouteExtract(v)
		if len(raw_data) > 4 {
			VrfID, _ := strconv.Atoi(raw_data[0])
			output := IProuteCLIModel{
				IPPrefix:  raw_data[1],
				Vrf:       VrfID,
				IsInstall: raw_data[2],
				NextHop:   raw_data[3],
				Interface: raw_data[4],
			}
			IProuteCLI = append(IProuteCLI, output)
		} else if len(raw_data) > 3 {

			VrfID, _ := strconv.Atoi(raw_data[0])
			output := IProuteCLIModel{
				IPPrefix:  raw_data[1],
				Vrf:       VrfID,
				IsInstall: raw_data[2],
				Interface: raw_data[3],
			}
			IProuteCLI = append(IProuteCLI, output)
		}
	}
	IProuteReturn.Attr = IProuteCLI

	return IProuteReturn
}

// ShowIpRouteConfig 는 모든 routing table을 호출합니다.
func ShowIpRouteConfig(interface_name string) {
	res := GetIpRouteConfig()
	data := ParseIpRouteConfig(res)
	makeTable(TitleRouteTable, data)
}

func ParseIpRouteConfig(res string) [][]string {

	var row_data []string
	var data [][]string
	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	// 각각의 기준이 VxLAN ID 한개 기준.
	for _, rd := range raw {
		row_data = RouteExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) > 4 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4]}
			data = append(data, tmp_row)

		} else if len(row_data) > 3 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], "", row_data[3]}
			data = append(data, tmp_row)

		} else {
			continue
		}
	}
	return data

}

func ParseIpNeighborConfig(res string) [][]string {

	var row_data []string
	var data [][]string
	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data = NeighborExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) > 6 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4], row_data[5], row_data[6]}
			data = append(data, tmp_row)

		}
	}
	return data

}

// ShowIpNeighborConfig 는 모든 neighbor정보 혹은 한개의 인터페이스에 대한 neighbor정보를 호출합니다.
func ShowIpNeighborConfig(interface_name string) {

	// send msg and return value
	res := GetIpNeighborOneConfig(interface_name)
	data := ParseIpNeighborConfig(res)
	makeTable(TitleNeighborTable, data)
}

// AddIpRouteConfig 는 스태틱 라우팅 테이블 엔트리를 추가합니다.
func AddIpRouteConfig(ip_address string, prefix int, next_hop string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	var cmd uint8
	v4 := IsV4(ip_address)
	if v4 {
		cmd = uint8(LOXILIGHT_RT_ADD)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_RT_ADD)
	}
	// make data format
	msg := fmt.Sprintf("%s/%d %s", ip_address, prefix, next_hop)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelIpRouteConfig 는 스태틱 라우팅 테이블 엔트리중 넥스트홉을 명시하여 삭제합니다.
func DelIpRouteConfig(ip_address string, prefix int, nextHop string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// make data format
	var cmd uint8
	v4 := IsV4(ip_address)

	if v4 {
		cmd = uint8(LOXILIGHT_RT_DEL)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_RT_DEL)
	}
	msg := fmt.Sprintf("%s/%d %s", ip_address, prefix, nextHop)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddIpAddressConfig 는 인터페이스에 IP 주소를 추가합니다.
func AddIpAddressConfig(interface_name string, ip_address string, prefix int) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// make data format
	var cmd uint8
	v4 := IsV4(ip_address)
	if v4 {
		cmd = uint8(LOXILIGHT_IP_ADDR_ADD)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_ADDR_ADD)
	}
	msg := fmt.Sprintf("%s/%d %s", ip_address, prefix, interface_name)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelIpAddressConfig 는 인터페이스에 IP 주소를 삭제합니다.
func DelIpAddressConfig(interface_name string, ip_address string, prefix int) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// make data format
	var cmd uint8
	v4 := IsV4(ip_address)
	if v4 {
		cmd = uint8(LOXILIGHT_IP_ADDR_DEL)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_ADDR_DEL)
	}
	msg := fmt.Sprintf("%s/%d %s", ip_address, prefix, interface_name)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddIpNeighborConfig 는 스태틱 neighbor정보를 추가합니다.
func AddIpNeighborConfig(interface_name string, ip_address string, mac_address string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// make data format
	var cmd uint8
	v4 := IsV4(ip_address)
	if v4 {
		cmd = uint8(LOXILIGHT_NH_ADD)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_NH_ADD)
	}
	msg := fmt.Sprintf("%s %s %s", ip_address, interface_name, mac_address)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelIpNeighborConfig 는 스태틱 neighbor정보를 삭제합니다.
func DelIpNeighborConfig(interface_name string, ip_address string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// make data format
	var cmd uint8
	v4 := IsV4(ip_address)
	if v4 {
		cmd = uint8(LOXILIGHT_NH_DEL)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_NH_DEL)
	}
	msg := fmt.Sprintf("%s %s", ip_address, interface_name)
	_, hdr := MakeMessage(cmd, msg)
	// send msg and return value

	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}
