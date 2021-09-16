package tc

type Config struct {
	SecretID  string `json:"secretID" toml:"secretID"`
	SecretKey string `json:"secretKey" toml:"secretLey"`
	AppID     uint64 `json:"appid" toml:"appid"`
	AppSecret string `json:"appSecret" toml:"appSecret"`
}
