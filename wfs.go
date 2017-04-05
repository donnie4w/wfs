/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package main

import (
	"flag"
	"fmt"

	. "github.com/donnie4w/wfs/conf"
	"github.com/donnie4w/wfs/httpserver"
)

func main() {
	ParseFlag()
	flag.Parse()
	fmt.Println("wfs start ,listen:", CF.Port)
	httpserver.Start()
}
