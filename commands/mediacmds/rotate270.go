package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var Rotate270Command = utils.CreateImgCmd(imaging.Rotate270)

func init() {
	utils.CreateCommand("rotate270", "USER", "Rotates an image by 270 degrees", Rotate270Command)
}
