package commands

import (
	"errors"
	"net"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("ip", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		pipeString, _ := GetTextContext(pipe)
		if pipeString == "" && len(args) <= 0 {
			response := &waProto.Message{
				Conversation: proto.String("Please specify a url"),
			}
			return response, errors.New("No url provided")
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

		response := &waProto.Message{
			Conversation: proto.String(strings.Join(ipParse, "\n")),
		}
		return response, nil
	})

}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
