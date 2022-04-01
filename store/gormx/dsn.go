package gormx

import (
	"fmt"
	"strings"
)

type Dsn struct {
	User     string            `json:"user" toml:"user"`
	Password string            `json:"password" toml:"password"`
	Proto    string            `json:"proto" toml:"proto"`
	Host     string            `json:"host" toml:"host"`
	Schema   string            `json:"schema" toml:"schema"`
	Options  map[string]string `json:"options" toml:"options"`
}

func (d *Dsn) String() string {
	var url = fmt.Sprintf("%s:%s@%s(%s)/%s", d.User, d.Password, d.Proto, d.Host, d.Schema)
	if len(d.Options) == 0 {
		return url
	}
	url += "?"
	for k, v := range d.Options {
		url += fmt.Sprintf("%s=%s&", k, v)
	}
	url = strings.TrimSuffix(url, "&")
	return url
}

func (d *Dsn) UseDefault() string {
	if d.Options == nil {
		d.Options = make(map[string]string)
	}
	var ok bool
	_, ok = d.Options["charset"]
	if !ok {
		d.Options["charset"] = "utf8mb4"
	}
	_, ok = d.Options["parseTime"]
	if !ok {
		d.Options["parseTime"] = "True"
	}
	_, ok = d.Options["loc"]
	if !ok {
		d.Options["loc"] = "Local"
	}
	return d.String()
}
