package config

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
}

func (m *MysqlConfig) Default() {
	if len(m.Charset) == 0 {
		m.Charset = "utf8"
	}
}
