package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/disintegration/imaging"
)

func init() {
	utils.CreateCommand("gray", utils.CreateImgCmd(imaging.Grayscale))
}
