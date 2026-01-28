package client

import "fmt"

const version = "1.1.0"

func getUserAgent() string {
	return fmt.Sprintf("morisawafonts-webfont-go/%s", version)
}
