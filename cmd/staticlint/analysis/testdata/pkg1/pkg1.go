package pkg1

import "os"

func osExitCheckFunc() {
	os.Exit(0) // want "declaration os.Exit shouldn't be used"
}