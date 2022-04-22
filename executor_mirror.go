package commonLib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.nanomsg.org/mangos/v3"
)

func MirrorExtract(raw string) []string {
	re1 := regexp.MustCompile("Host-IP (.+?) Flags (.+?) VRF (.+?) \\((.+?)\\) HW\\((.+?)\\)")
	match1 := re1.FindStringSubmatch(raw)
	if len(match1) > 0 {
		return match1[1:]
	}
	return match1
}

func GetMirrorConfig() string {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return ""
	}
	// make data format
	var hdr []byte
	cmd := uint8(LOXILIGHT_MIRROR_SHOW)
	_, hdr = MakeMessage(cmd, "")

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res
}

func GetMirrorModel() EthernetNeighborReturnModel {
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

func ParseMirrorConfig(res string) [][]string {

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

// AddMirrorSpanConfig 는 SPAN 미러링 아웃풋 포트 및 설정을 추가합니다.
func AddMirrorSpanConfig(MirrorId int, OutInterfaceName string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	cmd := uint8(LOXILIGHT_MIRROR_ADD)

	// make data format
	msg := fmt.Sprintf("%d span %s *", MirrorId, OutInterfaceName)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddMirrorRspanConfig 는 RSPAN미러링 아웃풋 포트 및 설정을 추가합니다.
func AddMirrorRspanConfig(MirrorId int, OutInterfaceName string, VlanId int) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	cmd := uint8(LOXILIGHT_MIRROR_ADD)

	// make data format
	msg := fmt.Sprintf("%d rspan %s %d", MirrorId, OutInterfaceName, VlanId)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelMirrorConfig 는 미러링 아웃풋 포트 및 설정을 삭제합니다.
func DelMirrorConfig(MirrorId int) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	cmd := uint8(LOXILIGHT_MIRROR_DEL)

	// make data format
	msg := fmt.Sprintf("%d", MirrorId)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// AddMirrorPortConfig 는 미러링 포트를 추가합니다.
func AddMirrorPortConfig(MirrorId int, InterfaceName string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	cmd := uint8(LOXILIGHT_MIRROR_PORT_ADD)

	// make data format
	msg := fmt.Sprintf("%d %s", MirrorId, InterfaceName)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

// DelMirrorPortConfig 는 미러링 포트를 삭제합니다.
func DelMirrorPortConfig(MirrorId int, InterfaceName string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	cmd := uint8(LOXILIGHT_MIRROR_PORT_DEL)

	// make data format
	msg := fmt.Sprintf("%d %s", MirrorId, InterfaceName)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}
