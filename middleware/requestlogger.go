package middleware

import (
	"fmt"
	"time"

	"github.com/gavruk/crater"
)

func RequestLogger(req *crater.Request, res *crater.Response) {
	var now = time.Now()

	fmt.Printf("(%s) %s", now, req.URL.String())
	fmt.Println("")
}
