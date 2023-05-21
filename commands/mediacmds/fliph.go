package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

func init() {
	utils.CreateCommand("fliph", utils.CreateImgCmd(imaging.FlipH))
}
