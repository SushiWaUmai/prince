package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var InvertCommand = utils.CreateImgCmd(imaging.Invert)

func init() {
	utils.CreateCommand("invert", "USER", "Inverts the colors of an image", InvertCommand)
}
