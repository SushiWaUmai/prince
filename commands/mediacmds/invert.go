package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

func init() {
	utils.CreateCommand("invert", utils.CreateImgCmd(imaging.Invert))
}