package app

import (
	"fmt"
	"time"
)

func log(msg string) {
	t := time.Now()

	s := "[" + t.Format("2006-01-02 15:04:05") + "] " + msg
	fmt.Println(s)
}
