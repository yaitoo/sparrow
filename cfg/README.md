# Quikstart 

## Load a local config file as inifile

- create config file `./app.ini`

```
[db]
host=127.0.0.1:3306
user=yaitoo
passwd=amazing
database=sparrow
```

- open and convert config to `Inifile`

```go

import (
    "fmt"
    "context"
    "github.com/yaitoo/sparrow/cfg"
)

c := cfg.Open(context.TODO(), "./app.ini").ToInifile()

host := c.Section("db").Value("host","")
user := c.Section("db").Value("user","")

fmt.Println(host, user)

```

- enable HotReload feature

```go

c := cfg.Open(context.TODO(), "./app.ini")
inifile := c.ToInifile()
c.OnChanged(func(c *Config){
    inifile.TryParse(string(c.Bytes))
})

```

- works with other config formats. eg toml,yaml...
  
```go

import "github.com/BurntSushi/toml"

c := cfg.Open(context.TODO(), "./app.ini")

var conf Config
if _, err := toml.Decode(string(c.Bytes), &conf); err != nil {
  // handle error
}

c.OnChanged(func(c *Config){
    if _, err := toml.Decode(string(c.Bytes), &conf); err != nil {
  // handle error
})

```