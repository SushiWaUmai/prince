package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var FlipHCommand = utils.CreateImgCmd(imaging.FlipH)

func init() {
	utils.CreateCommand("fliph", "USER", "Flips an image horizontally", FlipHCommand)
}
