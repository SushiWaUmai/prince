package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var GrayscaleCommand = utils.CreateImgCmd(imaging.Grayscale)

func init() {
	utils.CreateCommand("gray", "USER", GrayscaleCommand)
}
