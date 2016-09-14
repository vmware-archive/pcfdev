package ui

import (
	"fmt"
	"io"
)

type UI struct {
	Stdout io.Writer
}

func (u *UI) PrintHelpText(domain string) error {
	_, err := fmt.Fprintf(
		u.Stdout,
		` _______  _______  _______    ______   _______  __   __
|       ||       ||       |  |      | |       ||  | |  |
|    _  ||       ||    ___|  |  _    ||    ___||  |_|  |
|   |_| ||       ||   |___   | | |   ||   |___ |       |
|    ___||      _||    ___|  | |_|   ||    ___||       |
|   |    |     |_ |   |      |       ||   |___  |     |
|___|    |_______||___|      |______| |_______|  |___|
is now running.
To begin using PCF Dev, please run:
    cf login -a https://api.%s --skip-ssl-validation
Apps Manager URL: https://%s
Admin user => Email: admin / Password: admin
Regular user => Email: user / Password: pass
`,
		domain,
		domain,
	)
	return err
}
