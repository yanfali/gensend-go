package main

import (
	"log"
	"testing"
)

func TestGenerate1000IterationsNoCollisionsSamePlainText(*testing.T) {
	test := "abc123"
	url := &UrlGenerator{}
	for i := 0; i < 1000; i++ {
		resulta := url.Generate(test)
		resultb := url.Generate(test)
		//log.Printf("%q %q", resulta, resultb)
		if resulta == resultb {
			log.Fatalf("Unexpected %q === %q", resulta, resultb)
		}
	}
}
