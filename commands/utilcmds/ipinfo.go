package utilcmds

import (
	"errors"
	"net"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"github.com/ipinfo/go/v2/ipinfo"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

var ipClient = ipinfo.NewClient(nil, nil, "")

func IPInfoCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	pipeString, _ := utils.GetTextContext(pipe)
	if pipeString == "" && len(args) <= 0 {
		return utils.CreateTextMessage("Please specify a ip address"), errors.New("No ip address specified")
	}

	var ipAddress string
	if pipeString != "" {
		ipAddress = pipeString
	} else {
		ipAddress = args[0]
	}

	if !IsIPv4(ipAddress) || !IsIPv6(ipAddress) {
		ips, err := net.LookupIP(ipAddress)
		if err != nil || len(ips) == 0 {
			return nil, err
		}

		ipAddress = ips[0].String()
	}

	info, err := ipClient.GetIPInfo(net.ParseIP(ipAddress))
	if err != nil {
		return nil, err
	}

	var infoParse []string
	infoParse = append(infoParse, "IP: "+info.IP.String())

	infoParse = append(infoParse, "")

	infoParse = append(infoParse, "Timezone: "+info.Timezone)
	infoParse = append(infoParse, "Country: "+info.CountryName)
	infoParse = append(infoParse, "City: "+info.City)
	infoParse = append(infoParse, "Postal: "+info.Postal)
	infoParse = append(infoParse, "Location: "+info.Location)
	infoParse = append(infoParse, "Organization: "+info.Org)

	return utils.CreateTextMessage(strings.Join(infoParse, "\n")), nil
}

func init() {
	utils.CreateCommand("ipinfo", "USER", "Find the information of an IP address", IPInfoCommand)
}
