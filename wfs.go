/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package main

import (
	"flag"

	"wfs/httpserver"

	. "wfs/conf"
)

func main() {
	ParseFlag()
	flag.Parse()
	httpserver.Start()
}
