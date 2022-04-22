package commonLib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var TitleConntrack = []string{"Date", "Time", "Source IP", "Port", "Destinaion IP", "Port", "Protocol", "Timeout", "Total timeout", "Status", "Statistic[Count:byte]"}

// ConntrackExtract는 데이터에서 Conntrack 부분의 정보를 추출합니다.
func ConntrackExtract(raw string) []string {
	re := regexp.MustCompile(`(.+?) (.+?) # (.+?):(.+?) <-> (.+?):(.+?) \((.+?)\) (.+?)/(.+?) (.+?) \[(.+?)\]`)
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match

}

func ConntrackDetailExtract(raw string) []string {
	re := regexp.MustCompile(`(.+?):(.+?) <-> (.+?):(.+?) \\((.+?)\\) (.+?)/(.+?) (.+)`)
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

func ConntrackStatExtract(raw string) []string {
	re := regexp.MustCompile(`(.+?):(.+)`)
	match := re.FindStringSubmatch(raw)
	if len(match) > 0 {
		return match[1:]
	}
	return match
}

// GetConntrackConfig는 모든 Conntrack의 정보를 리턴합니다.
func GetConntrackConfig() (string, error) {
	sock, err := GetConnection(LoxilightMgmtIp)
	if err != nil {
		fmt.Println("Please check your Core APP and CLI network status")
		return "", err
	}
	// Send msg and return value

	cmd := uint8(LOXILIGHT_CT_SHOW)
	_, hdr := MakeMessage(cmd, "")
	res := SendMessage(sock, hdr)
	CloseConnection(sock)
	return res, err

}

func GetConntrackModel() (ConntrackReturnModel, error) {
	var Conntrack []ConntrackModel
	var ConntrackReturn ConntrackReturnModel

	connt, err := GetConntrackConfig()
	if err != nil {
		return ConntrackReturn, err
	}
	for _, v := range strings.Split(connt, "\r\n") {
		raw_data := ConntrackExtract(v)
		if len(raw_data) > 10 {

			// raw_data[3] 에서 : 기준으로 나눠서 아이피 포트
			SrcIP := raw_data[2]
			SrcPort, _ := strconv.Atoi(strings.TrimSpace(raw_data[3]))
			// raw_data[5] 에서 : 기준으로 나눠서 아이피 포트
			DstIP := raw_data[4]
			DstPort, _ := strconv.Atoi(strings.TrimSpace(raw_data[5]))

			// raw_data[10] 에서 스탯 추출하기
			PacketStats := ConntrackStatExtract(raw_data[10])
			var PacketCount int
			var PacketByte int
			if len(PacketStats) > 1 {
				PacketCount, _ = strconv.Atoi(strings.TrimSpace(PacketStats[0]))
				PacketByte, _ = strconv.Atoi(strings.TrimSpace(PacketStats[1]))
			} else {
				PacketCount = 0
				PacketByte = 0
			}
			output := ConntrackModel{
				Date:            raw_data[0],
				Time:            raw_data[1],
				SourceIp:        SrcIP,
				SourcePort:      SrcPort,
				DestinationIp:   DstIP,
				DestinationPort: DstPort,
				Protocol:        raw_data[6],
				Timeout:         raw_data[7],
				ExpireTime:      raw_data[8],
				Status:          raw_data[9],
				PacketCount:     PacketCount,
				PacketBtye:      PacketByte,
			}
			Conntrack = append(Conntrack, output)
		}
	}
	ConntrackReturn.Attr = Conntrack
	return ConntrackReturn, err

}

// ParseConntrackData은 라인별로 추출이 가능하게 도와줍니다.
func ParseConntrackData(res string) [][]string {
	var row_data []string
	var data [][]string
	// Parse the response to Data
	raw := strings.Split(res, "\r\n")
	for _, rd := range raw {
		row_data = ConntrackExtract(rd)
		if len(row_data) > 10 {
			tmp_row := []string{row_data[0], row_data[1], row_data[2], row_data[3], row_data[4], row_data[5], row_data[6], row_data[7], row_data[8], row_data[9], row_data[10]}
			data = append(data, tmp_row)
		}
	}
	return data
}

// ShowConntrackConfig 는 모든 Conntrack 의 간략한 정보를 CLI 테이블로 변환하여 보여줍니다.
func ShowConntrackConfig() {
	// Get_data
	res, _ := GetConntrackConfig()
	data := ParseConntrackData(res)
	// Make a table to display
	makeTable(TitleConntrack, data)
}
