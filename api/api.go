package api

import (
	"net/http"
	"io/fs"
	"github.com/SushiWaUmai/prince/frontend"
)

func CreateAPI() {
	assetsFs, err := fs.Sub(frontend.Assets(), "build")
	if err != nil {
		panic(err)
	}


	http.Handle("/", http.FileServer(http.FS(assetsFs)))

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
