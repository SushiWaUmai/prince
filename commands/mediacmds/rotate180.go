package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var Rotate180Command = utils.CreateImgCmd(imaging.Rotate180)

func init() {
	utils.CreateCommand("rotate180", "USER", "Rotates an image by 180 degrees", Rotate180Command)
}
