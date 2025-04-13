package tools

import "time"

func TimeLoadLocationWIB() (*time.Location, error) {
	return time.LoadLocation("Asia/Jakarta")
}
