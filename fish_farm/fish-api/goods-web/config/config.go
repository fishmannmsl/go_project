package config

//用户微服务后端接口信息配置
type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

//JWT加密密钥配置
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

//阿里云	Sms信息模板配置
type AliSmsConfig struct {
	AliKey          string `mapstructure:"key" json:"key"`
	AliSecret       string `mapstructure:"secret" json:"secret"`
	AliSignName     string `mapstructure:"sign_name" json:"sign_name"`
	AliTemplateCode string `mapstructure:"template_code" json:"template_code"`
	AliPhoneNumbers string `mapstructure:"phone_numbers" json:"phone_numbers"`
}

//consul 配置
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

//Redis 信息配置
type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ServerConfig struct {
	Name       string        `mapstructure:"name" json:"name"`
	Port       int           `mapstructure:"port" json:"port"`
	UserInfo   UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo    JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo  RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
