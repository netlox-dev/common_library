package commonLib

import (
	"errors"
	"fmt"

	"go.nanomsg.org/mangos/v3"
)

func GetIpsecConfig() string {
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
	msg := ""

	cmd := uint8(LOXILIGHT_IPSEC_POLICY_SHOW)
	_, hdr = MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	CloseConnection(sock)

	return res
}

func ShowIpsecConfig() {
	var sock mangos.Socket
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return
	}
	// make data format
	var hdr []byte
	//var row_data []string
	//var data [][]string
	msg := ""

	cmd := uint8(LOXILIGHT_IPSEC_POLICY_SHOW)
	_, hdr = MakeMessage(cmd, msg)

	// send msg and return value
	res := SendMessage(sock, hdr)
	fmt.Printf("-----------------------------------------------------------------------------------\n")
	fmt.Printf("%s", string(res))
	fmt.Printf("-----------------------------------------------------------------------------------\n")

}

func AddIpsecConfig(source_ip string, destination_ip string, protocol string, spi string, direction string, reqid string, key string, key_len string, interface_name string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// choose cmd
	cmd := uint8(LOXILIGHT_IPSEC_POLICY_ADD)
	// make data format
	msg := fmt.Sprintf("%s %s %s %s %s %s %s %s %s", source_ip, destination_ip, protocol, spi, direction, reqid, key, key_len, interface_name)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}

func DelIpsecConfig(source_ip string, destination_ip string, protocol string, spi string, direction string) error {
	var returnError error = nil
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnError
	}
	// make data format
	cmd := uint8(LOXILIGHT_IPSEC_POLICY_DEL)
	msg := fmt.Sprintf("%s %s %s %s %s", source_ip, destination_ip, protocol, spi, direction)
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		returnError = errors.New(res)
	}
	return returnError
}
