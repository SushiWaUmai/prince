package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

var FlipVCommand = utils.CreateImgCmd(imaging.FlipV)

func init() {
	utils.CreateCommand("flipv", "USER", FlipVCommand)
}
