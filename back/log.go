package main

import (
	"back/interfaces"
	"log"
)

func errLog(e interfaces.CustomError) {
	log.Fatalf("%v, %v", e.Error(), e.Unwrap().Error())
}
