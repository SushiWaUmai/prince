package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var Rotate90Command = utils.CreateImgCmd(imaging.Rotate90)

func init() {
	utils.CreateCommand("rotate90", "USER", "Rotates an image by 90 degrees", Rotate90Command)
}
