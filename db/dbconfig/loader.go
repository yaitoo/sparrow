package dbconfig

import (
	"errors"
	"io/ioutil"
	"regexp"
	"sync"

	"github.com/yaitoo/sparrow/db/model"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const fileRegexp = `\/*(.+\/)*(.+)\.(.+)$`

var (
	confObjGlobal model.Config //= model.Config{}

	once                 sync.Once
	fileLocationGlobal   FileLocation
	ErrInvalidConfig     = errors.New("Invalid config")
	ErrInvalidConfigPath = errors.New("Invalid config path")
)

type FileLocation string

func (fl FileLocation) String() string {
	return string(fl)
}
func Initonfig(fileLocation FileLocation) {
	fileLocationGlobal = fileLocation
}

type FileContent struct {
	FileLocation
	FileName      string
	FilePath      string
	FileExtention string
	WatchCB       func(fsnotify.Event)
}

func GetConfigObj() model.Config {

	once.Do(func() {
		fileContent, ok := parseFilelocation(fileLocationGlobal)
		if ok == false {
			panic(ErrInvalidConfigPath)
		} else {
			go MonitorConfigFile(fileContent)
			conf, err := readConfig(fileLocationGlobal)
			if err != nil {
				//return nil, err
				panic(err)
			}
			confObjGlobal = conf
		}
	})
	//out = confObjGlobal
	return confObjGlobal
}

func parseFilelocation(fileLocation FileLocation) (FileContent, bool) {
	var parseRegexp = regexp.MustCompile(fileRegexp)
	var matchGroup = parseRegexp.FindStringSubmatch(fileLocation.String())
	if len(matchGroup) < 3 {
		return FileContent{}, false
	}
	return FileContent{
		FileLocation:  fileLocation,
		FilePath:      matchGroup[1],
		FileName:      matchGroup[2],
		FileExtention: matchGroup[3],
		WatchCB: func(event fsnotify.Event) {
			conf, err := readConfig(fileLocationGlobal)
			if err != nil {
				//return nil, err
				panic(err)
			}
			confObjGlobal = conf
		},
	}, true
}

func MonitorConfigFile(fileContent FileContent) {
	// if file occur change, read file and  update confObj
	// and send log
	done := make(chan bool)

	v := viper.New()
	v.SetConfigType(fileContent.FileExtention)
	v.SetConfigName(fileContent.FileName)
	v.AddConfigPath(fileContent.FilePath) //"$GOPATH/src/geax/wallet_v2/")
	v.AddConfigPath("./")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	v.WatchConfig()

	v.OnConfigChange(fileContent.WatchCB)
	<-done
}

func readConfig(fileLocation FileLocation) (model.Config, error) {
	yamlContent, err := ioutil.ReadFile(fileLocation.String())
	if err != nil {
		return model.Config{}, err
	}
	confObj := model.Config{}
	if err := deserialize(yamlContent, &confObj); err != nil {
		return model.Config{}, err
	}
	if confObj.Validate() == false {
		return confObj, ErrInvalidConfig
	}
	return confObj, nil
}

func deserialize(content []byte, confObj *model.Config) error {
	if err := yaml.Unmarshal(content, &confObj); err != nil {
		return err
	}
	return nil
}
