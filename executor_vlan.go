package commonLib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var TitleVlanBrief = []string{"Name", "VLAN", "Status", "Ports"}

//var TitleVlanDetail = []string{"VLAN", "Type", "SAID", "MTU", "Parent", "RingNo", "BridgeNo", "Stp", "BrdgMode", "Trans1", "Trans2"}
var TitleVlanDetail = []string{"Name", "VLAN", "VRF", "MAC", "Status", "Address", "Ports"}

func GetRemovedTagPort(portname string) string {
	raw := strings.Fields(portname)
	if len(raw) == 0 {
		return ""
	}
	raw = strings.Split(raw[0], ".")
	return raw[0]
}

func VlanExtract(raw string) []string {
	re := regexp.MustCompile("Name (.+?) VID (.+?) VRF (.+?) MAC (.+?) \\((.+?\\))")
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func IpaddressInVlanExtract(raw string) []string {
	sp_data := strings.Split(raw, "IPv4 addr")
	if len(sp_data) > 0 {
		return sp_data[1:]
	}
	return sp_data
}

func VlanPortExtract(raw string) []string {
	sp_data := strings.Split(raw, "IPvs addr")
	if len(sp_data) > 1 {
		ssp_data := strings.Split(sp_data[0], "IPv4 addr")
		if len(ssp_data) > 1 {
			sssp_data := strings.Split(ssp_data[0], "\n\r")
			return sssp_data[1:]
		} else {
			ssp_data = strings.Split(sp_data[0], "\n\r")
			return ssp_data[1:]
		}
	} else {
		ssp_data := strings.Split(raw, "IPv4 addr")
		if len(ssp_data) > 1 {
			sssp_data := strings.Split(ssp_data[0], "\n\r")
			return sssp_data[1:]
		} else {
			ssp_data = strings.Split(sp_data[0], "\n\r")
			return ssp_data[1:]
		}
	}
	//return []string{}
}

func VlanPortDetailExtract(raw string, tagged string) []string {
	var tmp_ret []string
	var ret []string
	sp_data := strings.Split(raw, "IPvs addr")
	if len(sp_data) > 1 {
		ssp_data := strings.Split(sp_data[0], "IPv4 addr")
		if len(ssp_data) > 1 {
			sssp_data := strings.Split(ssp_data[0], "\n\r")
			tmp_ret = sssp_data[1:]
		} else {
			ssp_data = strings.Split(sp_data[0], "\n\r")
			tmp_ret = ssp_data[1:]
		}
	} else {
		ssp_data := strings.Split(raw, "IPv4 addr")
		if len(ssp_data) > 1 {
			sssp_data := strings.Split(ssp_data[0], "\n\r")
			tmp_ret = sssp_data[1:]
		} else {
			ssp_data = strings.Split(sp_data[0], "\n\r")
			tmp_ret = ssp_data[1:]
		}
	}
	for _, port := range tmp_ret {
		if tagged == "tagged" && strings.Contains(port, ".") {
			ret = append(ret, port)
		} else if tagged == "untagged" && !strings.Contains(port, ".") {
			ret = append(ret, port)
		}
	}
	return ret
}

// GetVlanBriefData??? ?????? vlan??? ????????? ???????????????.
func GetVlanBriefData() (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// Send msg and return value

	cmd := uint8(LOXILIGHT_VBR_SHOW_ALL)
	_, hdr := MakeMessage(cmd, "")
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res, err

}

func ParseVlanBridgeData(res string) [][]string {
	var row_data []string
	var port_data []string
	var data [][]string
	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	// ????????? ????????? VLAN ID ?????? ??????.
	for _, rd := range raw {
		row_data = VlanExtract(rd)
		port_data = VlanPortExtract(rd)
		if len(row_data) != 0 {
			var tmp_row []string
			if len(port_data) != 0 {
				for i, port := range port_data {
					if i == 0 {
						tmp_row = []string{row_data[0], row_data[1], row_data[4][0:2], port} // Name , Vid , Status, Ports
					} else {
						tmp_row = []string{" ", " ", " ", port}
					}
					data = append(data, tmp_row)
				}
			} else {
				tmp_row = []string{row_data[0], row_data[1], row_data[4][0:2], ""} // Name , Vid , Status, Empty port
				data = append(data, tmp_row)
			}

		}
	}
	return data
}

// ShowVlanBriefConfig ??? ?????? vlan ??? ????????? ????????? ???????????????.
func ShowVlanBriefConfig() {
	// Get_data
	res, _ := GetVlanBriefData()
	data := ParseVlanBridgeData(res)
	// Make a table to display
	makeTable(TitleVlanBrief, data)
}

// GetVlanIdConfig ??? ??? ?????? vlan??? ????????? ????????? ???????????????.
func GetVlanIdConfig(vlan_id int) (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// Send msg and return value
	var hdr []byte
	msg := fmt.Sprintf("%d", vlan_id)
	cmd := uint8(LOXILIGHT_VBR_SHOW_ID)
	_, hdr = MakeMessage(cmd, msg)
	res := SendMessage(sock, hdr)
	CloseConnection(sock)
	return res, err

}

func ParseVlanIdConfig(res string) [][]string {

	var row_data []string
	var port_data []string
	var data [][]string
	row_data = VlanExtract(res)
	port_data = VlanPortExtract(res)
	ip_data := IpaddressInVlanExtract(res)
	if len(row_data) != 0 {
		var tmp_row []string
		if len(port_data) != 0 {
			for i, port := range port_data {
				if i == 0 {
					tmp_row = []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4][0:2], "", port} // Name , Vid , Status, Ports
				} else {
					tmp_row = []string{" ", " ", " ", " ", " ", " ", port}
				}
				data = append(data, tmp_row)
			}
		} else {
			tmp_row = []string{row_data[0], row_data[1], row_data[4][0:2], ""} // Name , Vid , Status, Empty port
			data = append(data, tmp_row)
		}

		if len(data) > len(ip_data) {
			for i, dt := range data {
				if len(ip_data) == i {
					break
				}
				dt[5] = strings.TrimSpace(ip_data[i])
			}
		} else {
			for i, ip := range ip_data {
				ip = strings.TrimSpace(ip)
				if len(data) < i {
					tmp_row = []string{" ", " ", " ", " ", " ", ip, " "}
					data = append(data, tmp_row)
				} else {
					data[i][5] = ip
				}

			}
		}

	}
	return data
}

// GetVlanIdModel ??? Vlan??? ?????? ????????? Model??? ???????????? ???????????????.
func GetVlanIdBriefModel() (VlansReturnModel, error) {
	var Vlan []VlanModel
	var VlansReturn VlansReturnModel
	res, err := GetVlanBriefData()
	if err != nil {
		return VlansReturn, err
	}
	for _, v := range strings.Split(res, "\r\n") {
		raw_data := VlanExtract(v)
		port_data := VlanPortExtract(v)
		ip_data := IpaddressInVlanExtract(v)
		if len(raw_data) > 0 {
			vlan_id, _ := strconv.Atoi(raw_data[1])
			vrf_id, _ := strconv.Atoi(raw_data[2])
			output := VlanModel{
				VlanName: raw_data[0],
				VlanID:   vlan_id,
				Vrf:      vrf_id,
				Mac:      raw_data[3],
			}
			// Port Member??? ???????????? ??????
			for _, p := range port_data {
				if p != "" {
					intf_name := strings.Fields(p)
					if len(intf_name) > 1 {
						var port_out VlanMembersModel
						// .??? ???????????? ????????? ????????? ?????? ???????????? ???????????????
						// vxlan??? ????????? ?????? ????????? ?????? ????????? ????????? ???????????? ??????.
						if len(strings.Split(intf_name[0], ".")) < 2 {
							port_out.If_name = intf_name[0]
							port_out.Tagging = "untagged"
						} else {
							port_out.If_name = intf_name[0]
							port_out.Tagging = "tagged"
						}
						output.Member = append(output.Member, port_out)
					}

				}

			}

			// IP????????? ???????????? ??????
			for _, ip := range ip_data {
				if ip != "" {
					for _, real_ip := range strings.Split(ip, "\n\r") {
						if real_ip != "" {
							output.IPPrefix = append(output.IPPrefix, real_ip)
						}

					}
				}
			}
			Vlan = append(Vlan, output)
		}
	}
	VlansReturn.Attr = Vlan
	return VlansReturn, err
}

func GetVlanIdModel(vlan_id int) (VlanModel, error) {
	var attr VlanModel
	res, err := GetVlanIdConfig(vlan_id)
	if err != nil {
		return attr, err
	}
	raw_data := VlanExtract(res)
	port_data := VlanPortExtract(res)
	ip_data := IpaddressInVlanExtract(res)
	if len(raw_data) > 0 {
		vrf_id, _ := strconv.Atoi(raw_data[2])
		attr.VlanName = raw_data[0]
		attr.VlanID = vlan_id
		attr.Vrf = vrf_id
		attr.Mac = raw_data[3]

		for _, p := range port_data {
			if p != "" {
				intf_name := strings.Fields(p)
				if len(intf_name) > 1 {
					var port_out VlanMembersModel
					// .??? ???????????? ????????? ????????? ?????? ???????????? ???????????????
					// vxlan??? ????????? ?????? ????????? ?????? ????????? ????????? ???????????? ??????.
					if len(strings.Split(intf_name[0], ".")) < 2 {
						port_out.If_name = intf_name[0]
						port_out.Tagging = "untagged"
					} else {
						port_out.If_name = intf_name[0]
						port_out.Tagging = "tagged"
					}
					attr.Member = append(attr.Member, port_out)
				}

			}

		}

		// IP????????? ???????????? ??????
		for _, ip := range ip_data {
			if ip != "" {
				for _, real_ip := range strings.Split(ip, "\n\r") {
					if real_ip != "" {
						attr.IPPrefix = append(attr.IPPrefix, real_ip)
					}

				}
			}
		}

	}
	return attr, err
}

// ShowVlanIdConfig ??? ??? ?????? vlan ??? ????????? ????????? ???????????????.
func ShowVlanIdConfig(vlan_id int) {
	res, _ := GetVlanIdConfig(vlan_id)
	data := ParseVlanIdConfig(res)
	// Make a table to display
	makeTable(TitleVlanDetail, data)
}

// AddVlanBridge??? vlan bridge??? ???????????????.
func AddVlanBridge(vlan_id int) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	msg := fmt.Sprintf("%d", vlan_id)
	// make Message
	cmd := uint8(LOXILIGHT_VBR_ADD)
	_, hdr := MakeMessage(cmd, msg)
	// send msg and return value

	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// DelVlanBridge??? vlan bridge??? ???????????????.
func DelVlanBridge(vlan_id int) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	msg := fmt.Sprintf("%d", vlan_id)
	// make Messgae
	cmd := uint8(LOXILIGHT_VBR_DEL)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// AddVlanMember??? ??? ????????? vlan bridge??? port??? ???????????????.
func AddVlanMember(vlan_id int, interface_name string, tagging_status string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	msg := fmt.Sprintf("%d %s %s", vlan_id, interface_name, tagging_status)
	// make Message
	cmd := uint8(LOXILIGHT_VBR_PORT_ADD)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}

// DelVlanMember??? ??? ????????? vlan bridge??? port??? ???????????????.
func DelVlanMember(vlan_id int, interface_name string, tagging_status string) error {

	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	msg := fmt.Sprintf("%d %s %s", vlan_id, interface_name, tagging_status)
	// make Header
	cmd := uint8(LOXILIGHT_VBR_PORT_DEL)
	_, hdr := MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}
