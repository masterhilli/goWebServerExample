package Logger

import (
	"log"
	"os"
)

var (
	LOGGER *log.Logger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
)
