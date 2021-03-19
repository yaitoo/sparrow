package dbconfig_test

import (
	"testing"

	"github.com/yaitoo/sparrow/db/dbconfig"
)

const fileLocation dbconfig.FileLocation = "./test/db.yaml"

/* func TestLoadInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			fmt.Println(err)
			if err.Error() != dbconfig.ErrInvalidConfig.Error() {
				t.Error(err)
			}
		}
	}()
	var fileLocation dbconfig.FileLocation = "./test/db_fail.yaml"
	dbconfig.Initonfig(fileLocation)
	dbconfig.GetConfigObj()
} */

func TestLoad(t *testing.T) {
	dbconfig.Initonfig(fileLocation)
	// ch := make(chan int)

	obj := dbconfig.GetConfigObj()
	if len(obj.Versions) == 0 {
		t.Error("")
	}

	//time.Sleep(30 * time.Second)

	/* 	for {
		notify := <-ch
		fmt.Println(notify)
	} */
	/* 	select {
	   	case notify := <-ch:
	   		fmt.Println(notify)
	   	default:
	   		fmt.Println("no msg")
	   	}

	   	time.Sleep(30 * time.Second) */
}
