package utilcmds

import (
	"errors"
	"net"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func IPCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	pipeString, _ := utils.GetTextContext(pipe)
	if pipeString == "" && len(args) <= 0 {
		return utils.CreateTextMessage("Please specify a url"), errors.New("No url provided")
	}

	var url string
	if pipeString != "" {
		url = pipeString
	} else {
		url = args[0]
	}

	ips, err := net.LookupIP(url)
	if err != nil {
		return nil, err
	}

	var ipParse []string

	ipParse = append(ipParse, "IPv4:")
	for _, ip := range ips {
		if IsIPv4(ip.String()) {
			ipParse = append(ipParse, ip.String())
		}
	}
	ipParse = append(ipParse, "")

	ipParse = append(ipParse, "IPv6:")
	for _, ip := range ips {
		if IsIPv6(ip.String()) {
			ipParse = append(ipParse, ip.String())
		}
	}

	return utils.CreateTextMessage(strings.Join(ipParse, "\n")), nil
}

func init() {
	utils.CreateCommand("ip", "USER", IPCommand)
}
