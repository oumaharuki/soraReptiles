package ctrl

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func IndexGetDoReptiles(req *http.Request, r render.Render) {
	start := 1
	end := 1
	bol := Work(start, end)

	if bol {
		r.JSON(200, map[string]interface{}{
			"insert_id": "",
		})
		return
	}

}
