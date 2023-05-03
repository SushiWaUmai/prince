package commands

import (
	"context"
	"log"
	"net"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("ip", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		if len(args) <= 0 {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please specify a url"),
			})
			return
		}

		url := args[0]

		ips, err := net.LookupIP(url)
		if err != nil {
			log.Println("Failed to get IPs")
			return
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

		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: proto.String(strings.Join(ipParse, "\n")),
		})

		if err != nil {
			log.Println(err)
		}
	})

}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
