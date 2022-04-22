package commonLib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var TitleInterfaceBrief = []string{"Ethernet_Interface", "VLAN", "Type", "Mode", "Status", "Reason", "Speed", "Port_Ch#"}
var TitleFdb = []string{"MAC", "VLAN", "VRF", "Interface", "Status"}
var TitleVxfdb = []string{"MAC", "VxLAN", "VRF", "Destinaion IP", "VETP id", "Status"}

/*
 It will be added
					breakout
					counters
					description
					naming_mode
					neighbor
					portchannel
					status
					transceiver
					fcip
 these information
*/

// GetFdbConfig 는 모든 L2FDB의 데이터를 리턴합니다.
func GetFdbConfig() (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// make data format
	var hdr []byte
	cmd := uint8(LOXILIGHT_L2FDB_SHOW_ALL)
	_, hdr = MakeMessage(cmd, "")
	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res, err
}

// GetFdbVlanConfig 는 vlan에 설정된 L2FDB값을 리턴합니다.
func GetFdbVlanConfig(vlan_id int) (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// make data format
	var hdr []byte

	cmd := uint8(LOXILIGHT_L2FDB_SHOW_VLAN)
	vlanId := fmt.Sprintf("%d", vlan_id)
	_, hdr = MakeMessage(cmd, vlanId)
	// send msg and return value
	res := SendMessage(sock, hdr)
	//fmt.Printf("Receive %s\n", string(res))
	return res, err
}

// GetFdbModel 는 모든 L2FDB의 데이터를 모델 형식으로 리턴합니다.
func GetFdbModel() (VlanFDBReturn VlanFDBReturnModel, err error) {
	var VlanFDB []VlanFDBModel
	vlanFDB, err := GetFdbConfig()
	if err != nil {
		return VlanFDBReturn, err
	}

	for _, v := range strings.Split(vlanFDB, "\r\n") {
		raw_data := FdbExtract(v)
		if len(raw_data) > 0 {
			VlanID, _ := strconv.Atoi(raw_data[1])
			VrfID, _ := strconv.Atoi(raw_data[2])
			output := VlanFDBModel{
				Mac:    raw_data[0],
				VlanID: VlanID,
				Vrf:    VrfID,
				Ifname: raw_data[3],
				Status: raw_data[4],
			}
			VlanFDB = append(VlanFDB, output)
		}
	}
	VlanFDBReturn.Attr = VlanFDB
	return VlanFDBReturn, err
}

// GetFdbVlanModel 는 특정 VLAN의 L2FDB의 데이터를 모델 형식으로 리턴합니다.
func GetFdbVlanModel(VlanId int) (VlanFDBReturn VlanFDBReturnModel, err error) {
	var VlanFDB []VlanFDBModel
	vlanFDB, err := GetFdbVlanConfig(VlanId)
	if err != nil {
		return VlanFDBReturn, err
	}
	for _, v := range strings.Split(vlanFDB, "\r\n") {
		raw_data := FdbExtract(v)
		if len(raw_data) > 0 {
			VlanID, _ := strconv.Atoi(raw_data[1])
			VrfID, _ := strconv.Atoi(raw_data[2])
			output := VlanFDBModel{
				Mac:    raw_data[0],
				VlanID: VlanID,
				Vrf:    VrfID,
				Ifname: raw_data[3],
				Status: raw_data[4],
			}
			VlanFDB = append(VlanFDB, output)
		}
	}
	VlanFDBReturn.Attr = VlanFDB
	return VlanFDBReturn, err

}

// GetVxFdbConfig 는  VxLAN관련 모든 FDB데이터를 리턴합니다.
func GetVxFdbConfig() (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	cmd := uint8(LOXILIGHT_VXFDB_SHOW_ALL)
	_, hdr := MakeMessage(cmd, "")
	// send msg and return value
	res := SendMessage(sock, hdr)
	//fmt.Printf("Receive %s\n", string(res))
	CloseConnection(sock)

	return res, err
}

// GetVxFdbModel 는 VxLAN관련 모든 FDB데이터를 모델 형식으로 리턴합니다.
func GetVxFdbModel() (VxlanFDBReturnModel, error) {
	var VxlanFDB []VxlanFDBModel
	var VxlanFDBReturn VxlanFDBReturnModel

	vxlanFDB, err := GetVxFdbConfig()
	if err != nil {
		return VxlanFDBReturn, err
	}
	for _, v := range strings.Split(vxlanFDB, "\r\n") {
		raw_data := VxFdbExtract(v)
		if len(raw_data) > 0 {
			VxlanID, _ := strconv.Atoi(raw_data[1])
			VrfID, _ := strconv.Atoi(raw_data[2])
			output := VxlanFDBModel{
				Mac:     raw_data[0],
				VxlanID: VxlanID,
				Vrf:     VrfID,
				DstIP:   raw_data[3],
				VtepID:  raw_data[4],
				Status:  raw_data[5],
			}
			VxlanFDB = append(VxlanFDB, output)
		}
	}
	VxlanFDBReturn.Attr = VxlanFDB
	return VxlanFDBReturn, err
}

// GetFdbVxlanConfig 는 Vxlan에 설정된 L2FDB의 값을 리턴합니다.
func GetFdbVxlanConfig(vxlan_id int) (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// make data format
	var hdr []byte
	vxlanId := fmt.Sprintf("%d", vxlan_id)
	cmd := uint8(LOXILIGHT_VXFDB_SHOW_VXLAN)
	_, hdr = MakeMessage(cmd, vxlanId)
	// send msg and return value
	res := SendMessage(sock, hdr)
	//fmt.Printf("Receive %s\n", string(res))
	return res, err
}

// GetFdbVxlanModel 는 Vxlan에 설정된 L2FDB의 모델을 리턴합니다.
func GetFdbVxlanModel(vxlan_id int) (VxlanFDBReturnModel, error) {
	var VxlanFDB []VxlanFDBModel
	var VxlanFDBReturn VxlanFDBReturnModel

	vxlanFDB, err := GetFdbVxlanConfig(vxlan_id)
	if err != nil {
		return VxlanFDBReturn, err
	}
	for _, v := range strings.Split(vxlanFDB, "\r\n") {
		raw_data := VxFdbExtract(v)
		if len(raw_data) > 0 {
			VxlanID, _ := strconv.Atoi(raw_data[1])
			VrfID, _ := strconv.Atoi(raw_data[2])
			output := VxlanFDBModel{
				Mac:     raw_data[0],
				VxlanID: VxlanID,
				Vrf:     VrfID,
				DstIP:   raw_data[3],
				VtepID:  raw_data[4],
				Status:  raw_data[5],
			}
			VxlanFDB = append(VxlanFDB, output)
		}
	}
	VxlanFDBReturn.Attr = VxlanFDB
	return VxlanFDBReturn, err
}

func ParseFdbConfig(res string) [][]string {
	var row_data []string
	var data [][]string
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data = FdbExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4]}
			data = append(data, tmp_row)

		}
	}
	return data
}

func ParseVxFdbConfig(res string) [][]string {
	var row_data2 []string
	var data2 [][]string
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data2 = VxFdbExtract(rd)
		//fmt.Println(row_data)
		if len(row_data2) != 0 {
			tmp_row := []string{row_data2[0], row_data2[1], row_data2[4], row_data2[2], row_data2[3], row_data2[5]}
			data2 = append(data2, tmp_row)

		}
	}
	return data2
}

// ShowFdbConfig 는 모든 L2FDB의 값을 보여줍니다.
func ShowFdbConfig() {
	res, _ := GetFdbConfig()
	data := ParseFdbConfig(res)
	fmt.Println("FDB Table")
	makeTable(TitleFdb, data)

	res2, _ := GetVxFdbConfig()
	data2 := ParseVxFdbConfig(res2)
	fmt.Println("VxLAN FDB Table")
	makeTable(TitleVxfdb, data2)
}
func FdbExtract(raw string) []string {
	re := regexp.MustCompile("mac (.+?) vlan (.+?) vrf (.+?) dev (.+?) \\((.+?)\\) (.+?)")

	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func VxFdbExtract(raw string) []string {
	re := regexp.MustCompile("vx-mac (.+?) vxlan (.+?) dst.r.ip (.+?) vtepid (.+?) vrf (.+?) \\((.+?)\\)")

	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}
func ParseFdbVlanConfig(res string) [][]string {

	var row_data []string
	var data [][]string
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data = FdbExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4]}
			data = append(data, tmp_row)

		}
	}
	return data
}

// ShowFdbVlanConfig는 vlan에 설정된 L2FDB의 값을 보여줍니다.
func ShowFdbVlanConfig(vlan_id int) {
	res, _ := GetFdbVlanConfig(vlan_id)
	data := ParseFdbVlanConfig(res)
	makeTable(TitleFdb, data)
}

func ParseFdbVxlanConfig(res string) [][]string {
	var row_data []string
	var data [][]string
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data = VxFdbExtract(rd)
		//fmt.Println(row_data)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[4], row_data[2], row_data[3], row_data[5]}
			data = append(data, tmp_row)

		}
	}
	return data
}

// ShowFdbVxlanConfig는 vxlan에 설정된 L2FDB의 값을 보여줍니다.
func ShowFdbVxlanConfig(vxlan_id int) {
	res, _ := GetFdbVxlanConfig(vxlan_id)
	data := ParseFdbVxlanConfig(res)
	makeTable(TitleVxfdb, data)
}

// ShowInterfaceAllCounter는 현재 모든 인터페이스의 통계 카운트를 보여줍니다. 현재는 사용하지 못합니다.
func ShowInterfaceAllCounter() {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return
	}
	// Send msg and return value

	cmd := uint8(LOXILIGHT_INTERFACE_STAT_SHOW_BRIEF)
	_, hdr := MakeMessage(cmd, "")

	res := SendMessage(sock, hdr)

	// Parse the response to Data
	// Currnet format of response is not decieded. So Make it later TODO
	data := [][]string{strings.Fields(res)}

	// Make a table to display
	makeTable(TitleFdb, data)
}

// ShowInterfaceOneCounter는 한 개의 인터페이스의 통계 카운트를 보여줍니다. 현재는 사용하지 못합니다.
func ShowInterfaceOneCounter(interface_name string) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return
	}
	// Send msg and return value
	var hdr []byte
	msg := interface_name
	cmd := uint8(LOXILIGHT_INTERFACE_STAT_SHOW_INTERFACE)
	_, hdr = MakeMessage(cmd, msg)
	res := SendMessage(sock, hdr)

	// Parse the response to Data
	// Currnet format of response is not decieded. So Make it later TODO
	data := [][]string{strings.Fields(res)}

	// Make a table to display
	makeTable(TitleFdb, data)
}

// ShowInterfaceAllStatus 는 현재 모든 인터페이스의 간략한 정보를 보여줍니다. 현재는 사용하지 못합니다.
func ShowInterfaceAllStatus() {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return
	}
	// Send msg and return value
	var hdr []byte
	cmd := uint8(LOXILIGHT_INTERFACE_SHOW_BRIEF)
	_, hdr = MakeMessage(cmd, "")
	res := SendMessage(sock, hdr)

	// Parse the response to Data
	// Currnet format of response is not decieded. So Make it later TODO
	data := [][]string{strings.Fields(res)}

	// Make a table to display
	makeTable(TitleFdb, data)
}

// ShowInterfaceOneStatus 는 한 개의 인터페이스의 간략한 정보를 보여줍니다. 현재는 사용하지 못합니다.
func ShowInterfaceOneStatus(interface_name string) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return
	}
	// Send msg and return value
	var hdr []byte
	msg := interface_name
	cmd := uint8(LOXILIGHT_INTERFACE_SHOW_INTERFACE)
	_, hdr = MakeMessage(cmd, msg)
	res := SendMessage(sock, hdr)

	// Parse the response to Data
	// Currnet format of response is not decieded. So Make it later TODO
	data := [][]string{strings.Fields(res)}

	// Make a table to display

	makeTable(TitleFdb, data)
}

// AddFdbConfig 는 vlan id를 기반으로 L2FDB 주소를 추가합니다.
func AddFdbConfig(vlan_id string, mac_address string, interface_name string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	// make data format

	cmd := uint8(LOXILIGHT_L2FDB_ADD)
	msg := fmt.Sprintf("%s %s %s", vlan_id, mac_address, interface_name)
	_, hdr := MakeMessage(cmd, msg)
	fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// DelFdbConfig 는 vlan id를 기반으로 L2FDB 주소를 삭제합니다.
func DelFdbConfig(vlan_id string, mac_address string, interface_name string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	// make data format
	cmd := uint8(LOXILIGHT_L2FDB_DEL)
	msg := fmt.Sprintf("%s %s %s", vlan_id, mac_address, interface_name)
	_, hdr := MakeMessage(cmd, msg)
	fmt.Println(hdr, msg, cmd)
	// send msg and return value

	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// AddVxfdbConfig 는 vxlan id를 기반으로 L2FDB 주소를 추가합니다.
func AddVxfdbConfig(vxlan_id int, mac_address string, ip_address string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	// make data format

	cmd := uint8(LOXILIGHT_VXFDB_ADD)
	msg := fmt.Sprintf("%d %s %s", vxlan_id, mac_address, ip_address)
	_, hdr := MakeMessage(cmd, msg)
	fmt.Println(hdr, msg, cmd)
	// send msg and return value

	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// DelVxfdbConfig 는 vxlan id를 기반으로 L2FDB 주소를 삭제합니다.
func DelVxfdbConfig(vxlan_id int, mac_address string, ip_address string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	// make data format
	cmd := uint8(LOXILIGHT_VXFDB_DEL)
	msg := fmt.Sprintf("%d %s %s", vxlan_id, mac_address, ip_address)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value

	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

func GetVifConfig() (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// make data format
	var hdr []byte
	cmd := uint8(LOXILIGHT_VIF_SHOW_ALL)
	_, hdr = MakeMessage(cmd, "")

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res, err
}
func GetVifOneConfig(IfName string) (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// make data format
	var hdr []byte
	cmd := uint8(LOXILIGHT_VIF_SHOW_NAME)
	_, hdr = MakeMessage(cmd, IfName)

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res, err
}

func VifExtract(raw string) []string {
	re := regexp.MustCompile("(.+?): mac 0x(.+?) port (.+?) flags (.+?) (.+?) rvid (.+)")
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func IpaddressInVifExtract(raw string) []string {
	sp_data := strings.Split(raw, "v4a")
	if len(sp_data) > 0 {
		return sp_data[1:]
	}
	return sp_data
}

func GetVifModel() (VifReturnModel, error) {
	var Vif []VifModel
	var VifReturn VifReturnModel
	res, err := GetVifConfig()
	if err != nil {
		return VifReturn, err
	}
	raw := strings.Split(res, "\r\n\r\n")
	// 각각의 기준이 VIF 한개 기준.
	for _, rd := range raw {
		raw_data := VifExtract(rd)
		ipaddress := IpaddressInVifExtract(rd)
		if len(raw_data) != 0 {
			port, _ := strconv.Atoi(strings.TrimSpace(raw_data[2]))
			flags, _ := strconv.Atoi(strings.TrimSpace(raw_data[3]))
			rvid, _ := strconv.Atoi(strings.TrimSpace(raw_data[5]))
			output := VifModel{
				IfName:   strings.TrimSpace(raw_data[0]),
				Mac:      raw_data[1],
				Port:     port,
				Flags:    flags,
				Status:   raw_data[4],
				Rvid:     rvid,
				IPPrefix: ipaddress,
			}
			Vif = append(Vif, output)
		}
	}
	VifReturn.Attr = Vif
	return VifReturn, err
}

func GetVifOneModel(IfName string) (VifModel, error) {
	var Vif VifModel
	res, err := GetVifOneConfig(IfName)
	if err != nil {
		return Vif, err
	}
	raw := strings.Split(res, "\r\n\r\n")
	for _, rd := range raw {
		raw_data := VifExtract(rd)
		ipaddress := IpaddressInVifExtract(rd)
		if len(raw_data) != 0 {
			port, _ := strconv.Atoi(strings.TrimSpace(raw_data[2]))
			flags, _ := strconv.Atoi(strings.TrimSpace(raw_data[3]))
			rvid, _ := strconv.Atoi(strings.TrimSpace(raw_data[5]))
			Vif = VifModel{
				IfName:   strings.TrimSpace(raw_data[0]),
				Mac:      raw_data[1],
				Port:     port,
				Flags:    flags,
				Status:   raw_data[4],
				Rvid:     rvid,
				IPPrefix: ipaddress,
			}
		}
	}
	return Vif, err
}
