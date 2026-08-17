package utils

import (
	"os"
)

var (
	CurrentDir, _ = os.Getwd()
)

const (
	Yes = "yes"

	No = "no"

	PathSeparator = string(os.PathSeparator)

	SpaceSeparator = " "

	NewLineSeparator = "\n"

	NotAvailable = -1

	Empty = ""
)
