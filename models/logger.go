package models

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[pigo🐷] ", log.LstdFlags|log.Lshortfile)
