package commonLib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var TitleVxlanBrief = []string{"VxLAN-Name", "VxLAN-ID", "UIF", "MAC"}

//var title_vlan_detail = []string{"VLAN", "Type", "SAID", "MTU", "Parent", "RingNo", "BridgeNo", "Stp", "BrdgMode", "Trans1", "Trans2"}
var TitleVxlanDetail = []string{"VxLAN-Name", "VxLAN-ID", "VRF", "Local IP", "UDP port", "UIF", "MAC"}

func VxlanExtract(raw string) []string {
	re := regexp.MustCompile("VxLAN-Name (.+?) VxLAN-ID (.+?) VRF-id (.+?) Local-IP (.+?) UDP.P (.+?) UIF (.+?) VXBR-mac (.+)")
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

// GetVxlanBriefConfig 는 설정된 간략한 VxLAN 정보를 리턴합니다.
func GetVxlanBriefConfig() string {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return ""
	}
	// Send msg and return value
	cmd := uint8(LOXILIGHT_VXBR_SHOW_ALL)
	_, hdr := MakeMessage(cmd, "")
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res

}

// GetVxlanBriefModel 는 설정된 간략한 VxLAN 정보를 모델로 리턴합니다.
func GetVxlanBriefModel() VxlanReturnModel {
	var Vxlans []VxlanModel
	var VxlansReturn VxlanReturnModel
	raw_vxlan := GetVxlanBriefConfig()
	for _, r_vxlan := range strings.Split(raw_vxlan, "\r\n") {
		raw_data := VxlanExtract(r_vxlan)
		if len(raw_data) > 0 {
			VxlanID, _ := strconv.Atoi(raw_data[1])
			VrfID, _ := strconv.Atoi(raw_data[2])
			output := VxlanModel{
				VxlanName: raw_data[0],
				VxlanID:   VxlanID,
				Vrf:       VrfID,
				LocalIP:   raw_data[3],
				Udp:       raw_data[4],
				Uif:       raw_data[5],
				Mac:       raw_data[6],
			}
			Vxlans = append(Vxlans, output)
		}
	}
	VxlansReturn.Attr = Vxlans
	return VxlansReturn
}

func ParseVxlanBriefConfig(res string) [][]string {

	var row_data []string
	var data [][]string
	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	// 각각의 기준이 VxLAN ID 한개 기준.
	for _, rd := range raw {
		row_data = VxlanExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[5], row_data[6]}
			data = append(data, tmp_row)

		}
	}
	return data
}

// ShowVxlanBriefConfig 는 설정된 간략한 vxlan 정보를 보여줍니다.
func ShowVxlanBriefConfig() {
	res := GetVxlanBriefConfig()
	data := ParseVxlanBriefConfig(res)
	makeTable(TitleVxlanBrief, data)
}

// GetVxlanVtepConfig 는 설정된 간략한 vtep 정보를 리턴합니다. 현재는 작동하지 않습니다.
func GetVxlanVtepConfig() [][]string {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return [][]string{}
	}
	// Send msg and return value
	var row_data []string
	var data [][]string
	cmd := uint8(LOXILIGHT_VTEP_SHOW_ALL)
	_, hdr := MakeMessage(cmd, "")
	res := SendMessage(sock, hdr)

	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	// 각각의 기준이 VxLAN ID 한개 기준.
	for _, rd := range raw {
		row_data = VxlanExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4], row_data[5], row_data[6]}
			data = append(data, tmp_row)

		}
	}
	return data
}

// ShowVxlanVtepConfig 는 설정된 간략한 vtep 정보를 보여줍니다. 현재는 작동하지 않습니다.
func ShowVxlanVtepConfig() {
	data := GetVxlanVtepConfig()
	makeTable(TitleVxlanDetail, data)
}

// GetVxlanIdConfig 는 자세한 vxlan 정보를 리턴합니다.
func GetVxlanIdConfig(vxlan_id int) string {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return ""
	}
	// Send msg and return value
	var hdr []byte
	msg := fmt.Sprintf("%d", vxlan_id)
	cmd := uint8(LOXILIGHT_VXBR_SHOW_ID)
	_, hdr = MakeMessage(cmd, msg)
	res := SendMessage(sock, hdr)

	return res
}

// GetVxlanIdConfig 는 1개의 vxlan 정보를 모델로 리턴합니다.
func GetVxlanIdModel(vxlan_id int) VxlanModel {
	var output VxlanModel
	raw_vxlan := GetVxlanIdConfig(vxlan_id)
	raw_data := VxlanExtract(raw_vxlan)
	if len(raw_data) > 5 {
		VxlanID, _ := strconv.Atoi(raw_data[1])
		VrfID, _ := strconv.Atoi(raw_data[2])

		output = VxlanModel{
			VxlanName: raw_data[0],
			VxlanID:   VxlanID,
			Vrf:       VrfID,
			LocalIP:   raw_data[3],
			Udp:       raw_data[4],
			Uif:       raw_data[5],
			Mac:       raw_data[6],
		}
	}
	return output
}

func ParseVxlanIdConfig(res string) [][]string {

	var row_data []string
	var data [][]string
	raw := strings.Split(res, "\r\n")
	// 각각의 기준이 VxLAN ID 한개 기준.
	for _, rd := range raw {
		row_data = VxlanExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4], row_data[5], row_data[6]}
			data = append(data, tmp_row)

		}
	}
	//ip_addr := ipaddress_in_vlan_extract(res)
	return data
}

// ShowVxlanIdConfig 는 자세한 vxlan 정보를 보여줍니다.
func ShowVxlanIdConfig(vxlan_id int) {
	raw := GetVxlanIdConfig(vxlan_id)
	data := ParseVxlanIdConfig(raw)
	// Make a table to display
	makeTable(TitleVxlanDetail, data)
}

// AddVxlanBridgeWithoutVlanID는 vxlan을 추가하는 함수입니다.
func AddVxlanBridgeWithoutVlanID(vxlan_id int, interface_name string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d %s", vxlan_id, interface_name)
	// make Message
	cmd := uint8(LOXILIGHT_VXBR_ADD)
	_, hdr := MakeMessage(cmd, msg)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddVxlanBridgeWithVlanID는 vxlan을 추가하는 함수입니다.
func AddVxlanBridgeWithVlanID(vxlan_id int, interface_name string, vlan_id int) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d %s %d", vxlan_id, interface_name, vlan_id)
	// make Message
	cmd := uint8(LOXILIGHT_VXBR_ADD)
	_, hdr := MakeMessage(cmd, msg)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelVxlanBridge는 vxlan을 삭제하는 함수입니다.
func DelVxlanBridge(vxlan_id int) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d", vxlan_id)
	// make Messgae
	cmd := uint8(LOXILIGHT_VXBR_DEL)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddVxlanPeer 는 vxlan의 remote ip 주소를 추가하는 함수입니다.
func AddVxlanPeer(vxlan_id int, remote_ip string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d %s", vxlan_id, remote_ip)
	// make Message
	cmd := uint8(LOXILIGHT_VXBR_EP_RIP_ADD)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelVxlanPeer 는 vxlan의 remote ip 주소를 삭제하는 함수입니다.
func DelVxlanPeer(vxlan_id int, remote_ip string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d %s", vxlan_id, remote_ip)
	// make Header
	cmd := uint8(LOXILIGHT_VXBR_EP_RIP_DEL)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddVxlanMember는 호스트부분의 vxlan 설정을 추가하는 함수입니다.
func AddVxlanMember(vxlan_id int, interface_name string, tagging_status string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d %s %s", vxlan_id, interface_name, tagging_status)
	// make Message
	cmd := uint8(LOXILIGHT_VXBR_AP_ADD)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelVxlanMember는 호스트부분의 vxlan 설정을 삭제하는 함수입니다.
func DelVxlanMember(vxlan_id int, interface_name string, tagging_status string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	msg := fmt.Sprintf("%d %s %s", vxlan_id, interface_name, tagging_status)
	// make Header
	cmd := uint8(LOXILIGHT_VXBR_AP_DEL)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}
