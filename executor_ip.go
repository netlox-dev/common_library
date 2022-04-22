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

var TitleRouteTable = []string{"Host-IP", "Flags", "VRF", "NextHop", "Is Install", "HW Install", "interface"}
var TitleNeighborTable = []string{"VRF", "Next hop", "Source MAC", "Destinaion MAC", "Interface", "Is Install", "HW Install"}

func RouteExtract(raw string) []string {
	re1 := regexp.MustCompile("Host-IP (.+?) Flags (.+?) VRF (.+?) \\((.+?)\\) HW\\((.+?)\\)")
	re2 := regexp.MustCompile("Host-IP (.+?) VRF (.+?) \\((.+?)\\) HW\\((.+?)\\) nexthop via (.+?) dev (.+?) ")
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
	re := regexp.MustCompile("VRF (.+?) nh-ip (.+?) smac (.+?) dmac (.+?) dev (.+?) (.+?) \\((.+?)\\)")

	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func GetIpNeighborConfig() (string, error) {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	var hdr []byte
	// make data format
	cmd := uint8(LOXILIGHT_NH_SHOW_ALL)
	_, hdr = MakeMessage(cmd, "")

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res, err

}

func GetIpNeighborOneConfig(interface_name string) (string, error) {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
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
	return res, err

}
func GetIpNeighborModel() (EthernetNeighborReturnModel, error) {
	var EthernetNeighbor []EthernetNeighborCLIModel
	var EthernetNeighborReturn EthernetNeighborReturnModel
	res, err := GetIpNeighborConfig()
	if err != nil {
		return EthernetNeighborReturn, err
	}
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
	return EthernetNeighborReturn, err
}

func GetIpNeighborOneModel(interface_name string) (EthernetNeighborReturnModel, error) {
	var EthernetNeighbor []EthernetNeighborCLIModel
	var EthernetNeighborReturn EthernetNeighborReturnModel
	res, err := GetIpNeighborOneConfig(interface_name)
	if err != nil {
		return EthernetNeighborReturn, err
	}
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
	return EthernetNeighborReturn, err

}

func GetIpRouteConfig() (string, error) {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
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

	return res, err
}

//GetIPRouteModel 은 모든 Routing Table을 struct로 제공합니다.
func GetIPRouteModel() (IProuteReturnModel, error) {
	var IProuteCLI []IProuteCLIModel
	var IProuteReturn IProuteReturnModel

	route, err := GetIpRouteConfig()
	if err != nil {
		return IProuteReturn, err
	}
	for _, v := range strings.Split(route, "\r\n") {
		raw_data := RouteExtract(v)
		if len(raw_data) > 5 {
			VrfID, _ := strconv.Atoi(raw_data[1])
			output := IProuteCLIModel{
				IPPrefix:   raw_data[0],
				Vrf:        VrfID,
				IsInstall:  raw_data[2],
				IsHareWare: raw_data[3],
				NextHop:    raw_data[4],
				Interface:  raw_data[5],
			}
			IProuteCLI = append(IProuteCLI, output)
		} else if len(raw_data) > 4 {

			VrfID, _ := strconv.Atoi(raw_data[2])
			output := IProuteCLIModel{
				IPPrefix:   raw_data[0],
				Flags:      raw_data[1],
				Vrf:        VrfID,
				IsInstall:  raw_data[3],
				IsHareWare: raw_data[4],
			}
			IProuteCLI = append(IProuteCLI, output)
		}
	}
	IProuteReturn.Attr = IProuteCLI

	return IProuteReturn, err
}

// ShowIpRouteConfig 는 모든 routing table을 호출합니다.
func ShowIpRouteConfig(interface_name string) {
	res, _ := GetIpRouteConfig()
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
		if len(row_data) > 5 {
			tmp_row := []string{row_data[0], "", row_data[1], row_data[4], row_data[2], row_data[3], row_data[5]}
			data = append(data, tmp_row)

		} else if len(row_data) > 4 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], "", row_data[3], row_data[4], ""}
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
	// 각각의 기준이 VxLAN ID 한개 기준.
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
	res, _ := GetIpNeighborOneConfig(interface_name)
	data := ParseIpNeighborConfig(res)
	makeTable(TitleNeighborTable, data)
}

// AddIpRouteConfig 는 스태틱 라우팅 테이블 엔트리를 추가합니다.
func AddIpRouteConfig(ip_address string, prefix int, next_hop string) error {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
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
		err = errors.New(res)
	}
	return err
}

// DelIpRouteConfig 는 스태틱 라우팅 테이블 엔트리를 삭제합니다.
func DelIpRouteConfig(ip_address string, prefix int) error {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	// make data format
	var cmd uint8
	v4 := IsV4(ip_address)
	if v4 {
		cmd = uint8(LOXILIGHT_RT_DEL)
	} else {
		cmd = uint8(LOXILIGHT_IPV6_RT_DEL)
	}
	msg := fmt.Sprintf("%s/%d", ip_address, prefix)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// AddIpAddressConfig 는 인터페이스에 IP 주소를 추가합니다.
func AddIpAddressConfig(interface_name string, ip_address string, prefix int) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
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
		err = errors.New(res)
	}
	return err
}

// DelIpAddressConfig 는 인터페이스에 IP 주소를 삭제합니다.
func DelIpAddressConfig(interface_name string, ip_address string, prefix int) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
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
		err = errors.New(res)
	}
	return err
}

// AddIpNeighborConfig 는 스태틱 neighbor정보를 추가합니다.
func AddIpNeighborConfig(interface_name string, ip_address string, mac_address string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
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
		err = errors.New(res)
	}
	return err
}

// DelIpNeighborConfig 는 스태틱 neighbor정보를 삭제합니다.
func DelIpNeighborConfig(interface_name string, ip_address string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
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
		err = errors.New(res)
	}
	return err
}
