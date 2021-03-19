# Quikstart 

## Read local config file as `Inifile`

- Create config file `./app.ini`

```
[db]
host=127.0.0.1:3306
user=yaitoo
passwd=amazing
database=sparrow
```

- Open and convert config to `Inifile`

```go

import (
    "fmt"
    "context"
    "github.com/yaitoo/sparrow/cfg"
)

c := cfg.Open(context.TODO(), "./app.ini").ToInifile()

host := c.Section("db").Value("host","")
user := c.Section("db").Value("user","")

```

- Enable HotReload feature

```go

c := cfg.Open(context.TODO(), "./app.ini")
inifile := c.ToInifile()
c.OnChanged(func(c *cfg.Config){
    inifile.TryParse(string(c.Bytes))
})

```

- Works with other config formats. eg toml,yaml...
  
```go

import "github.com/BurntSushi/toml"

c := cfg.Open(context.TODO(), "./app.ini")

var conf Config
if _, err := toml.Decode(string(c.Bytes), &conf); err != nil {
  // handle error
}

c.OnChanged(func(c *cfg.Config){
    if _, err := toml.Decode(string(c.Bytes), &conf); err != nil {
  // handle error
    }
})

```

## Read remote config content as `Inifile`

- Implement a custom `Reader` to read remote content via http/tcp...

```go

type RemoteReader struct {
    Name string
}

func NewReader(name string) *RemoteReader {
    return &RemoteReader{name:name}
}

//Read implement `Reader.Read`
func (r *RemoteReader) Read(ctx context.Context) ([]byte, error) {

   return ReadBytesFromRemote(r.name)
}

//ModTime implement `Reader.ModTime`
func (r *RemoteReader) ModTime(ctx context.Context) (int64, error) {

	return ReadLatestModTimeFromRemote(r.name)
}

```

- Open config with `RemoteReader` from remote store(eg. redis,mysql)

```go

 c := cfg.Open("cmdb.apis", cfg.WithReader(func(ctx context.Context) cfg.Reader {
		return NewReader("cmdb.apis")
  })).ToInifile()


host := c.Section("db").Value("host","")
user := c.Section("db").Value("user","")

```