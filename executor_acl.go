package commonLib

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var TitleAcl = []string{"UUID", "Filter", "Action", "Statistic"}

// AclExtract 데이터에서 ACL 부분의 정보를 추출합니다.
func AclExtract(raw string) []string {
	re := regexp.MustCompile("acl:(.+?) filter-(.+?),act-(.+?),stats-(.+)")
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

// GetAclModel 는 필터 스트링을 ACL모델로 변경합니다.
func GetAclModel() (ACLReturnModels, error) {

	var Acl []ACLReturnModel
	var AclReturn ACLReturnModels
	acls, err := GetAclConfig()
	if err != nil {
		return AclReturn, err
	}
	// filter info
	// T0,ip,ipprot:0x13,dip:0.2.3.1/23,sip:1.31.4.1/32,sport:1234,dport:123,iif:1,act-drop
	// Action type :  act-snat: 0x1020304:1234, act-dnat: 0x1020304:1234, act-allow, act-drop
	// Split with comma(,)
	// table, --- , Action
	for _, acl := range strings.Split(acls, "\r\n") {
		filter_val := AclExtract(acl)
		//0 : uuid
		//1 : filter
		//2 : action
		//3 : statistic
		if len(filter_val) > 0 {
			output := ACLReturnModel{}
			output.UUID = filter_val[0]
			output.getFilters(filter_val[1])
			output.getActions(filter_val[2])
			output.getStats(filter_val[3])

			Acl = append(Acl, output)
		}

	}
	AclReturn.Attr = Acl
	return AclReturn, err
}

// GetAclConfig 모든 Nat의 정보를 리턴합니다.
func GetAclConfig() (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// Send msg and return value

	cmd := uint8(LOXILIGHT_ACL_SHOW)
	_, hdr := MakeMessage(cmd, "")
	res := SendMessage(sock, hdr)
	CloseConnection(sock)
	return res, err

}

// ParseAclConfig 라인별로 추출이 가능하게 도와줍니다.
func ParseAclConfig(res string) [][]string {
	var row_data []string
	var data [][]string
	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data = AclExtract(rd)
		if len(row_data) != 0 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3]}
			data = append(data, tmp_row)
		}
	}
	return data
}

// ShowAclConfig 는 모든 Nat 의 간략한 정보를 CLI 테이블로 변환하여 보여줍니다.
func ShowAclConfig() {
	// Get_data
	res, _ := GetAclConfig()
	data := ParseAclConfig(res)
	// Make a table to display
	makeTable(TitleAcl, data)
}

// AddAclConfig 는 스태틱 라우팅 테이블 엔트리를 추가합니다.
func AddAclConfig(acl ACLModel) (string, error) {
	var returnuuid string
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return returnuuid, err
	}
	cmd := uint8(LOXILIGHT_ACL_ADD)

	// make data format
	msg, err := MakeAclMessage(acl)
	if err != nil {
		fmt.Println("Please Check Match values")
		return returnuuid, err
	}
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if IsValidUUID(res) {
		returnuuid = strings.TrimSpace(res)
	} else {
		err = errors.New(res)
	}
	return returnuuid, err
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")
	return r.MatchString(uuid)
}

func MakeAclMessage(acl ACLModel) (matches string, err error) {
	if acl.Filters.Port != 0 {
		matches += fmt.Sprintf(" in-port hs%d ", acl.Filters.Port)
	}
	if acl.Filters.EthernetType != "" {
		matches += fmt.Sprintf(" eth-type %s ", acl.Filters.EthernetType)
	}
	if acl.Filters.SourceIp != "" {
		if acl.Filters.EthernetType == "" {
			err = errors.New("Set ethernet type(eth-type) first ")
			return "", err
		}
		matches += fmt.Sprintf(" src-ip %s ", acl.Filters.SourceIp)
	}
	if acl.Filters.DestinationIp != "" {
		if acl.Filters.EthernetType == "" {
			err = errors.New("Set ethernet type(eth-type) first ")
			return "", err
		}
		matches += fmt.Sprintf(" dst-ip %s ", acl.Filters.DestinationIp)
	}
	if acl.Filters.Protocol != 0 {
		if acl.Filters.EthernetType == "" {
			err = errors.New("Set ethernet type(eth-type) first ")
			return "", err
		}
		matches += fmt.Sprintf(" ip-proto %d ", acl.Filters.Protocol)
	}

	if acl.Filters.SourcePort != 0 {
		if acl.Filters.EthernetType == "" {
			err = errors.New("Set ethernet type(eth-type) first ")
			return "", err
		}
		matches += fmt.Sprintf(" sport %d ", acl.Filters.SourcePort)
	}

	if acl.Filters.DestinationPort != 0 {
		if acl.Filters.EthernetType == "" {
			err = errors.New("Set ethernet type(eth-type) first ")
			return "", err
		}
		matches += fmt.Sprintf(" dport %d ", acl.Filters.DestinationPort)
	}

	if matches == "" {
		return "", err
	} else {
		matches = "match " + matches
	}

	if len(acl.Actions) == 0 {
		// Default is Drop
		matches += " action drop "
	} else if len(acl.Actions) == 1 {
		// Drop , Allow
		if acl.Actions[0].ActionName == "drop" || acl.Actions[0].ActionName == "allow" {
			matches += fmt.Sprintf(" action %s ", acl.Actions[0].ActionName)
		} else {
			// Dnat, Snat
			matches += fmt.Sprintf(" action %s %s port %d ", acl.Actions[0].ActionName, acl.Actions[0].IpAddress, acl.Actions[0].Port)
		}
	} else if len(acl.Actions) > 1 {
		// LB part
		for _, action := range acl.Actions {
			matches += fmt.Sprintf(" action %s %s port %d weight %d ", action.ActionName, action.IpAddress, action.Port, action.Weight)
		}
	}

	return matches, err
}

// DelAclConfig 는 스태틱 라우팅 테이블 엔트리를 삭제합니다.
func DelAclConfig(UUID string) error {
	var err error = nil
	if strings.TrimSpace(UUID) == "" {
		err = errors.New("Empty UUID come in ")
		return err
	}
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return err
	}
	cmd := uint8(LOXILIGHT_ACL_DEL)

	msg := UUID
	_, hdr := MakeMessage(cmd, msg)
	//fmt.Println(hdr, msg, cmd)
	// send msg and return value
	res := SendMessage(sock, hdr)
	if len(res) != 0 {
		err = errors.New(res)
	}
	return err
}
