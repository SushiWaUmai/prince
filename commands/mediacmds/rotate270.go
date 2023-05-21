package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

func init() {
	utils.CreateCommand("rotate270", "USER", utils.CreateImgCmd(imaging.Rotate270))
}
