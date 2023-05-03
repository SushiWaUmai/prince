package commands

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/ipinfo/go/v2/ipinfo"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var ipClient = ipinfo.NewClient(nil, nil, "")

func init() {
	createCommand("ipinfo", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		if len(args) <= 0 {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please specify a ip address"),
			})
			return
		}

		ipAddress := args[0]

		if !IsIPv4(ipAddress) || !IsIPv6(ipAddress) {
			ips, err := net.LookupIP(ipAddress)
			if err != nil || len(ips) == 0 {
				log.Println("Failed to get IPs")
				return
			}

			ipAddress = ips[0].String()
		}

		info, err := ipClient.GetIPInfo(net.ParseIP(ipAddress))
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed fetch info"),
			})
			return
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

		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: proto.String(strings.Join(infoParse, "\n")),
		})

		if err != nil {
			log.Println(err)
		}
	})

}
