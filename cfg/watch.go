// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"log"
	"sync"
	"time"
)

//ReloadInterval  the frequency at configs are checked
var ReloadInterval = 30 * time.Second

var once sync.Once

func startWatch() {
	go watch()
}

//watch configs
func watch() {

	t := time.NewTicker(ReloadInterval)
	defer t.Stop()
	for {
		<-t.C

		mutex.RLock()
		for _, c := range configs {
			modTime, err := c.reader.ModTime(context.Background())

			if err != nil {
				log.Printf("[cfg]%s ModTime %v\n", c.Name, err)
				continue
			}

			if modTime != c.modTime {
				bytes, err := c.reader.Read(context.Background())
				if err != nil {
					log.Printf("[cfg]%s Read %v\n", c.Name, err)
					continue
				}

				c.Bytes = bytes
				c.modTime = modTime

				c.fireHandlers()
			}
		}
		mutex.RUnlock()
	}

}
