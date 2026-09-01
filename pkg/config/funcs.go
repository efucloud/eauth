package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/efucloud/common"
	"github.com/glebarez/sqlite"
	mysqldriver "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/big"
	"os"
	"path"
	"strings"
	"time"
)

func maskMySQLDSN(dsn string) string {
	if dsn == "" {
		return ""
	}
	atIdx := strings.Index(dsn, "@")
	if atIdx <= 0 {
		return dsn
	}
	cred := dsn[:atIdx]
	colonIdx := strings.Index(cred, ":")
	if colonIdx < 0 {
		return dsn
	}
	maskedCred := cred[:colonIdx+1] + "******"
	return maskedCred + dsn[atIdx:]
}

// createDBConnection  create database connection
func createDBConnection() (err error) {
	if ApplicationConfig.Mysql != nil {
		Logger.Info("database is Mysql")
		ApplicationConfig.Mysql.Default()
		c := ApplicationConfig.Mysql
		dsnConfig := mysqldriver.Config{
			User:      c.User,
			Passwd:    c.Password,
			Net:       "tcp",
			Addr:      c.Host,
			DBName:    c.Dbname,
			ParseTime: true,
			Loc:       time.Local,
			Params: map[string]string{
				"charset": c.Charset,
			},
			AllowNativePasswords: true,
		}
		dsn := dsnConfig.FormatDSN()
		Logger.Infof("database connection: %s", maskMySQLDSN(dsn))
		DBConnect, err = gorm.Open(gormmysql.Open(dsn), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
		if err == nil {
			sqlDB, _ := DBConnect.DB()
			//SetMaxIdleConns 设置空闲连接池中连接的最大数量
			sqlDB.SetMaxIdleConns(1)
			//SetMaxOpenConns 设置打开数据库连接的最大数量。
			sqlDB.SetMaxOpenConns(100)
			//SetConnMaxLifetime 设置了连接可复用的最大时间。
			sqlDB.SetConnMaxLifetime(5 * time.Minute)
		} else {
			Logger.Errorf("database connect failed, err: %s", err.Error())
		}
	} else {
		Logger.Info("database is sqlite")
		DBConnect, err = gorm.Open(sqlite.Open("eauth.db"), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
	}
	if err != nil {
		Logger.Errorf("create database connect failed, err: %s", err.Error())
	}
	return err
}
func logConfig(conf *LogConfig) {

	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
	})
	var encoderConfig zapcore.EncoderConfig
	if conf.Production {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()

}

func (c *Config) Init() {
	if c.LogConfig == nil {
		c.LogConfig = new(LogConfig)
		c.LogConfig.Filename = "./log/eauth.log"
		c.LogConfig.MaxAge = 30
		c.LogConfig.MaxSize = 1
		c.LogConfig.MaxBackups = 10
		c.LogConfig.Compress = false
	}
	if c.TokenPeriod < 8 {
		c.TokenPeriod = 8
	}
	logConfig(c.LogConfig)
	Logger.Infof("build info GoVersion: %s", GoVersion)
	Logger.Infof("build info Commit: %s", Commit)
	Logger.Infof("build info BuildDate: %s", BuildDate)
	Logger.Infof("build info ApplicationName: %s", ApplicationName)
	c.ServerAddress = strings.TrimSuffix(c.ServerAddress, "/")
	if len(c.UploadPath) == 0 {
		c.UploadPath, _ = os.Getwd()
		c.UploadPath = path.Join(c.UploadPath, "uploads")
	} else {
		c.UploadPath = strings.TrimSuffix(c.UploadPath, "/")
	}
	if err := createDBConnection(); err != nil {
		Logger.Fatalf("create database connect failed, err: %s", err.Error())
	}

	pa := fmt.Sprintf("%s/%s", ApplicationConfig.UploadPath, UserAvatars)
	if !common.PathExists(pa) {
		_ = os.MkdirAll(pa, os.ModePerm)
	}

}

func GenerateRsaKeys(bitSize int, expireInYears int) (certificate string, privateKey string) {
	//https://stackoverflow.com/questions/64104586/use-golang-to-get-rsa-key-the-same-way-openssl-genrsa
	//https://stackoverflow.com/questions/43822945/golang-can-i-create-x509keypair-using-rsa-key
	commonName := "eauth"
	//Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	//Encode private key to PKCS#1 ASN.1 PEM.
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	tml := x509.Certificate{
		//you can add any attr that you need
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(expireInYears, 0, 0),
		//you have to generate a different serial number each execution
		SerialNumber: big.NewInt(123456),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"efucloud", "efu-cloud"},
		},
		BasicConstraintsValid: true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	//Generate a pem block with the certificate
	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return string(certPem), string(privateKeyPem)
}
