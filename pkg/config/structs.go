package config

import (
	"net/http"
	"time"
)

type Config struct {
	Mysql         *MysqlConfig `json:"mysql" yaml:"mysql"`
	LogConfig     *LogConfig   `json:"logConfig" yaml:"logConfig"`
	ServerAddress string       `json:"serverAddress" description:"服务器地址"`
	UploadPath    string       `json:"uploadPath" yaml:"uploadPath" description:"文件上传路径"`
	TokenPeriod   int          `json:"tokenPeriod" yaml:"tokenPeriod" description:"系统Token有效期"`
	LoginConfig   LoginConfig  `json:"loginConfig" yaml:"loginConfig" description:"登录配置"`
	Email         *EmailConfig `json:"email"  yaml:"email" description:"邮箱配置"`
}
type EmailConfig struct {
	Username   string `json:"username" yaml:"username" description:"用户名"`
	Password   string `json:"password" yaml:"password" description:"密码"`
	SmtpServer string `json:"smtpServer" yaml:"smtpServer" description:"SMTP服务"`
	SmtpPort   int    `json:"smtpPort" yaml:"smtpPort" description:"SMTP服务"`
}
type LoginConfig struct {
	MFA             bool `json:"mfa" yaml:"mfa" description:"多因素认证开启"`
	FaceRecognition bool `json:"faceRecognition" yaml:"faceRecognition" description:"人脸识别"`
}
type LogConfig struct {
	//Filename is the file to write logs to.  Backup log files will be retained
	//in the same directory.  It uses <process name>-lumberjack.log in
	//os.TempDir() if empty.
	Filename string `json:"filename" yaml:"filename"`

	//MaxSize is the maximum size in megabytes of the log file before it gets
	//rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	//MaxAge is the maximum number of days to retain old log files based on the
	//timestamp encoded in their filename.  Note that a day is defined as 24
	//hours and may not exactly correspond to calendar days due to daylight
	//savings, leap seconds, etc. The default is not to remove old log files
	//based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`

	//MaxBackups is the maximum number of old log files to retain.  The default
	//is to retain all old log files (though MaxAge may still cause them to get
	//deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	//LocalTime determines if the time used for formatting the timestamps in
	//backup files is the computer's local time.  The default is to use UTC
	//time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	//Compress determines if the rotated log files should be compressed
	//using gzip. The default is not to perform compression.
	Compress   bool `json:"compress" yaml:"compress"`
	Production bool `json:"production" yaml:"production"`
}
type MysqlConfig struct {
	Host     string `json:"host" yaml:"host"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Dbname   string `json:"dbname" yaml:"dbname"`
	//utf8
	Charset string `json:"charset" yaml:"charset"`
	//Local
	Loc string `json:"loc" yaml:"loc"`
	//string 类型字段的默认长度
	DefaultStringSize uint `json:"defaultStringSize" yaml:"defaultStringSize"`
	//禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	DisableDatetimePrecision bool `json:"disableDatetimePrecision" yaml:"disableDatetimePrecision"`
	//重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	DontSupportRenameIndex bool `json:"dontSupportRenameIndex" yaml:"dontSupportRenameIndex"`
	//用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	DontSupportRenameColumn bool `json:"dontSupportRenameColumn" yaml:"dontSupportRenameColumn"`
	//根据当前 MySQL 版本自动配置
	SkipInitializeWithVersion bool `json:"skipInitializeWithVersion" yaml:"skipInitializeWithVersion"`
}

func (m *MysqlConfig) Default() {
	if len(m.Charset) == 0 {
		m.Charset = "utf8"
	}
	if len(m.Loc) == 0 {
		m.Loc = "Local"
	}
	if m.DefaultStringSize == 0 {
		m.DefaultStringSize = 255
	}
}

type ElasticsearchConfig struct {
	Addresses []string `json:"addresses" yaml:"addresses"` //A list of Elasticsearch nodes to use.
	Username  string   `json:"username" yaml:"username"`   //Username for HTTP Basic Authentication.
	Password  string   `json:"password" yaml:"password"`   //Password for HTTP Basic Authentication.

	APIKey                 string `json:"apiKey" yaml:"apiKey"`                                 //Base64-encoded token for authorization; if set, overrides username/password and service token.
	ServiceToken           string `json:"serviceToken" yaml:"serviceToken"`                     //Service token for authorization; if set, overrides username/password.
	CertificateFingerprint string `json:"certificateFingerprint" yaml:"certificateFingerprint"` //SHA256 hex fingerprint given by Elasticsearch on first launch.

	Header http.Header `json:"header" yaml:"header"` // HTTP request header.

	//PEM-encoded certificate authorities.
	//When set, an empty certificate pool will be created, and the certificates will be appended to it.
	//The option is only valid when the transport is not specified, or when it's http.Transport.
	CACert string `json:"caCert" yaml:"caCert"`

	RetryOnStatus        []int `json:"retryOnStatus" yaml:"retryOnStatus"`               //List of status codes for retry. Default: 502, 503, 504.
	DisableRetry         bool  `json:"disableRetry" yaml:"disableRetry"`                 //Default: false.
	EnableRetryOnTimeout bool  `json:"enableRetryOnTimeout" yaml:"enableRetryOnTimeout"` //Default: false.
	MaxRetries           int   `json:"maxRetries" yaml:"maxRetries"`                     //Default: 3.

	CompressRequestBody  bool `json:"compressRequestBody" yaml:"compressRequestBody"`   //Default: false.
	DiscoverNodesOnStart bool `json:"discoverNodesOnStart" yaml:"discoverNodesOnStart"` //Discover nodes when initializing the client. Default: false.

	DiscoverNodesInterval time.Duration `json:"discoverNodesInterval" yaml:"discoverNodesInterval"` //Discover nodes periodically. Default: disabled.

	EnableMetrics           bool `json:"enableMetrics" yaml:"enableMetrics"`                     //Enable the metrics collection.
	EnableDebugLogger       bool `json:"enableDebugLogger" yaml:"enableDebugLogger"`             //Enable the debug logging.
	EnableCompatibilityMode bool `json:"enableCompatibilityMode" yaml:"enableCompatibilityMode"` //Enable sends compatibility header

	DisableMetaHeader    bool `json:"disableMetaHeader" yaml:"disableMetaHeader"` //Disable the additional "X-Elastic-client-Meta" HTTP header.
	UseResponseCheckOnly bool `json:"useResponseCheckOnly" yaml:"useResponseCheckOnly"`
}
