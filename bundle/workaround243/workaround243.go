package workaround243

import (
	"os"
)

var HasNonzeroLastTrit = false

func LogBundle(bundle string)  {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("nonzero243.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//log.Fatal(err)
	}
	if _, err := f.Write([]byte(bundle + "\n")); err != nil {
		//log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		//log.Fatal(err)
	}
	HasNonzeroLastTrit = false
}

