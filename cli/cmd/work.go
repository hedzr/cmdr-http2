/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmd

import (
	"fmt"
	"os"
	"strings"
)

// ESC CSI escape code
const ESC = 27

var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)

func clearLines(lineCount int) {
	_, _ = fmt.Fprint(os.Stdout, strings.Repeat(clear, lineCount))
}
