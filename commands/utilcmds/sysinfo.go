package utilcmds

import (
	"fmt"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func init() {
	utils.CreateCommand("sysinfo", "ADMIN", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		hostInfo, err := host.Info()
		if err != nil {
			return nil, err
		}

		cpuInfo, err := cpu.Info()
		if err != nil {
			return nil, err
		}
		cpu := cpuInfo[0]

		memInfo, err := mem.VirtualMemory()
		if err != nil {
			return nil, err
		}

		var infoParse []string
		infoParse = append(infoParse, "HOST:")
		infoParse = append(infoParse, "Hostname: "+hostInfo.Hostname)
		infoParse = append(infoParse, "OS: "+hostInfo.OS)
		infoParse = append(infoParse, "Platform: "+hostInfo.Platform)
		infoParse = append(infoParse, "Kernel Version: "+hostInfo.KernelVersion)
		infoParse = append(infoParse, "Kernel Architecture: "+hostInfo.KernelArch)
		infoParse = append(infoParse, fmt.Sprintf("Uptime: %d", hostInfo.Uptime))
		infoParse = append(infoParse, "")

		infoParse = append(infoParse, "CPU:")
		infoParse = append(infoParse, "Model: "+cpu.ModelName)
		infoParse = append(infoParse, "")

		infoParse = append(infoParse, "MEMORY:")
		infoParse = append(infoParse, fmt.Sprintf("Total Memory: %d", memInfo.Total))
		infoParse = append(infoParse, fmt.Sprintf("Used Memory: %d", memInfo.Used))

		response := &waProto.Message{
			Conversation: proto.String(strings.Join(infoParse, "\n")),
		}
		return response, nil
	})
}
