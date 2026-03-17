package dtos

import (
	"context"
	"crypto"
	"crypto/x509"
	"database/sql/driver"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"math"
	"strings"
)

type ExistResponse struct {
	Name  string `json:"name" description:"邮箱或者手机号"`
	Exist bool   `json:"exist" description:"存在"`
}

// TableListPagination 分页信息
type TableListPagination struct {
	//总数
	Total uint `json:"total" description:"总数"`
	//每页数量
	PageSize uint `json:"pageSize" description:"每页数量"`
	//当前页
	Current uint `json:"current" description:"当前页"`
}

// AuthedUserInfo 用户信息
type AuthedUserInfo struct {
	//主键
	ID uint `json:"id" validate:"required" description:"记录ID"`
	//用户名
	Username string `json:"username" validate:"required,max=255" description:"用户名"`
	//昵称，如中文名
	Nickname string `json:"nickname" validate:"max=255" description:"昵称"`
	//工号
	JobNumber string `json:"jobNumber" validate:"max=255" description:"工号"`
	//系统角色
	Role string `json:"role" validate:"oneof=admin view edit none" enum:"admin|view|edit|none" description:"系统角色"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
	//邮箱
	Email string `json:"email" validate:"max=255" description:"邮箱"`
	//手机号码
	Phone string `json:"phone" validate:"max=50" description:"电话"`
	//默认语言
	Language string `json:"language" validate:"oneof=zh en" enum:"zh|en" description:"默认语言"`
	//头像
	Avatar string `json:"avatar" validate:"required" description:"头像"`
	//有设置密码
	HasPassword bool `json:"hasPassword" description:"有设置密码"`
	//密码强度
	PasswordStrength string `json:"passwordStrength" description:"密码强度"`
}

// LoginByUsername 用户名密码登录
type LoginByUsername struct {
	//自动登录,12小时，1周，15天,1个月，半年
	//-|12h|1w|15d|1m|0.5y
	RememberMe string `json:"rememberMe" enum:"-|12h|1w|15d|1m|0.5y" description:"自动登录,12小时，1周，15天,1个月，半年"`
	//用户名,邮箱、手机号码、工号
	Username string `json:"username" validate:"required" description:"用户名"`
	//密码
	Password string `json:"password" validate:"required" description:"密码"`
	//账户来源
	Source string `json:"source" description:"账户来源"`
}

// LoginByFaceIdData 人脸识别数据登录
type LoginByFaceIdData struct {
	//自动登录,12小时，1周，15天,1个月，半年
	//-|12h|1w|15d|1m|0.5y
	RememberMe string `json:"rememberMe" enum:"-|12h|1w|15d|1m|0.5y" description:"自动登录,12小时，1周，15天,1个月，半年"`
	//用户名,邮箱、手机号码、工号
	Username string `json:"username" validate:"required" description:"用户名"`
	//人脸识别数据
	FaceIdData ArrayFloat64 `json:"faceIdData" validate:"required"  description:"人脸识别数据"`
	//账户来源
	Source string `json:"source" description:"账户来源"`
}
type ThirdAuthProfile struct {
	//所属用户
	UserId uint `json:"userId" description:"所属用户"`
	//认证提供商
	Provider string `json:"provider" description:"所属用户"`
	//第三方认证信息中的ID
	LoginID string `json:"loginId" description:"第三方认证信息中的ID"`
	//第三方认证信息中的用户名
	Username string `json:"username" description:"第三方认证信息中的用户名"`
	//第三方认证信息中的邮箱
	Email string `json:"email" description:"第三方认证信息中的邮箱"`
	//手机号码
	Phone string `json:"phone" description:"手机号码"`
	//第三方认证信息中的别名
	Nickname string `json:"nickname" description:"第三方认证信息中的别名"`
	//第三方认证信息中的头像
	Avatar string `json:"avatar" description:"第三方认证信息中的头像"`
	//第三方认证返回的所有用户信息
	Properties JsonMap `json:"properties" description:"第三方认证返回的所有用户信息"`
	//用户主页
	Home string `json:"home" description:"用户主页"`
}
type LoginByLDAP struct {
	//自动登录,12小时，1周，15天,1个月，半年
	RememberMe string `json:"rememberMe" enum:"12h|1w|15d|1m|0.5y" description:"自动登录,12小时，1周，15天,1个月，半年"`
	//绑定的系统用户
	BindId uint `json:"bindId" description:"绑定的系统用户或者组织用户"`
	//用户名
	Username string `json:"username" validate:"required" description:"用户名"`
	//密码
	Password string `json:"password" validate:"required" description:"密码"`
}

// RegisterByOIDC OIDC认证时注册新用户时补全的信息
type RegisterByOIDC struct {
	//请求码
	Code string `json:"code" validate:"required" description:"请求码"`
	//默认用户名
	Username string `json:"username" validate:"required" description:"默认用户名"`
	// 密码
	Password string `json:"password" validate:"required" description:"密码"`
	//昵称
	Nickname string `json:"nickname" validate:"required" description:"昵称"`
	//手机号码
	Phone string `json:"phone" validate:"required" description:"手机号码"`
	//邮箱
	Email string `json:"email" validate:"required" description:"邮箱"`
}

func (t *RegisterByOIDC) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(t)
	if err != nil {
		var lines []string
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, v := range errs.Translate(trans) {
				lines = append(lines, v)
			}
			if len(lines) > 0 {
				err = errors.New(strings.Join(lines, "\n"))
			}
		}
	}
	return
}

// LoginByOIDC OIDC登录
type LoginByOIDC struct {
	RememberMe  string `json:"rememberMe" enum:"-|12h|1w|15d|1m|0.5y" description:"自动登录,12小时，1周，15天,1个月，半年"`
	Code        string `json:"code" description:"OIDC提供商返回的Code"`
	Provider    string `json:"provider" validate:"required" description:"OIDC认证提供商"`
	RedirectUri string `json:"redirectUri" description:"重定向地址"`
}

// OidcRequestToken oidc认证获取token请求结构
type OidcRequestToken struct {
	//客户端ID
	ClientId string `schema:"client_id" json:"client_id" validate:"required"`
	//客户端密钥
	ClientSecret string `schema:"client_secret" json:"client_secret" validate:"required"`
	//类型
	GrantType string `schema:"grant_type" json:"grant_type" validate:"required"`
	//请求码
	Code string `schema:"code" json:"code" validate:"required" `
	//重定向地址
	RedirectUri string `schema:"redirect_uri" json:"redirect_uri"`
	//PKCE验证码
	CodeVerifier string `schema:"code_verifier" json:"code_verifier"`
	//PKCE验证码(camel风格兼容)
	CodeVerifierCamel string `schema:"codeVerifier" json:"codeVerifier"`
}

func (md *OidcRequestToken) GetCodeVerifier() string {
	if len(md.CodeVerifier) > 0 {
		return md.CodeVerifier
	}
	return md.CodeVerifierCamel
}

// AccessTokenResponse 登录获取token
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token,omitempty" validate:"required"`
	TokenType    string `json:"token_type ,omitempty" validate:"required"`
	ExpiresIn    int64  `json:"expires_in,omitempty" validate:"required"`
	RefreshToken string `json:"refresh_token,omitempty" validate:"required"`
	IDToken      string `json:"id_token,omitempty" validate:"required"`
	Timestamp    int64  `json:"timestamp,omitempty" validate:"required" description:"token有效时间"`
	Need         bool   `json:"need,omitempty" description:"需要补充信息"`
	ID           uint   `json:"id" description:"用户ID"`
	Code         string `json:"code,omitempty" description:""`
	Username     string `json:"username,omitempty" description:"默认用户名"`
	Nickname     string `json:"nickname,omitempty" description:"昵称"`
	Phone        string `json:"phone,omitempty" description:"手机号码"`
	Email        string `json:"email,omitempty" description:"邮箱"`
	Mfa          bool   `json:"mfa" description:"需要MFA验证码"`
	Secret       string `json:"secret" description:"密钥"`
	Image        string `json:"image" description:"二维码"`
	//登录方式 未来动态认证时使用
	Method string `json:"method" description:"登录方式"`
}

type UserClaims struct {
	//系统用户ID
	Id uint `json:"id" description:"系统用户ID"`
	// 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Username string `json:"username" description:"用户名"`
	// 昵称，如中文名
	Nickname string `json:"nickname" description:"昵称"`
	// 系统角色
	Role  string `json:"role" description:"系统角色"`
	Nonce string `json:"nonce"`
	Email string `json:"email" description:"邮箱"`
	Phone string `json:"phone" description:"手机号码"`

	jwt.RegisteredClaims
}

// OidcCodeRequest 获取认证用户的OAuth的code
type OidcCodeRequest struct {
	//客户端ID
	ClientId string `json:"clientId" description:"客户端ID"`
	//重定向地址
	RedirectUri string `json:"redirectUri" description:"重定向地址"`
	//状态码
	State string `json:"state" description:"状态码"`
	//响应类型
	ResponseType string `json:"responseType" description:"响应类型"`
	//PKCE挑战值(camel风格)
	CodeChallenge string `json:"codeChallenge" description:"PKCE挑战值"`
	//PKCE挑战方法(camel风格)
	CodeChallengeMethod string `json:"codeChallengeMethod" description:"PKCE挑战方法,plain或S256"`
	//PKCE挑战值(标准风格)
	CodeChallengeStd string `json:"code_challenge" description:"PKCE挑战值"`
	//PKCE挑战方法(标准风格)
	CodeChallengeMethodStd string `json:"code_challenge_method" description:"PKCE挑战方法,plain或S256"`
}

func (md *OidcCodeRequest) GetCodeChallenge() string {
	if len(md.CodeChallenge) > 0 {
		return md.CodeChallenge
	}
	return md.CodeChallengeStd
}

func (md *OidcCodeRequest) GetCodeChallengeMethod() string {
	if len(md.CodeChallengeMethod) > 0 {
		return md.CodeChallengeMethod
	}
	return md.CodeChallengeMethodStd
}

type OidcCodeResponse struct {
	Code string `json:"code" description:"响应码"`
	//状态码
	State string `json:"state" description:"状态码"`
	//重定向地址
	RedirectUri string `json:"redirectUri" description:"重定向地址"`
}

type OpenIDConfiguration struct {
	Issuer                                 string   `json:"issuer" description:"发行者"`
	AuthorizationEndpoint                  string   `json:"authorization_endpoint" description:""`
	TokenEndpoint                          string   `json:"token_endpoint" description:""`
	UserinfoEndpoint                       string   `json:"userinfo_endpoint" description:""`
	JwksUri                                string   `json:"jwks_uri" description:"Token校验的"`
	IntrospectionEndpoint                  string   `json:"introspection_endpoint,omitempty" description:"Token校验"`
	RevocationEndpoint                     string   `json:"revocation_endpoint" description:"Token撤销"`
	ResponseTypesSupported                 []string `json:"response_type s_supported" description:""`
	ResponseModesSupported                 []string `json:"response_modes_supported" description:""`
	GrantTypesSupported                    []string `json:"grant_type s_supported" description:""`
	SubjectTypesSupported                  []string `json:"subject_type s_supported" description:""`
	IdTokenSigningAlgValuesSupported       []string `json:"id_token_signing_alg_values_supported" description:""`
	ScopesSupported                        []string `json:"scopes_supported" description:""`
	ClaimsSupported                        []string `json:"claims_supported" description:""`
	CodeChallengeMethodsSupported          []string `json:"code_challenge_methods_supported,omitempty" description:"支持的PKCE挑战算法"`
	RequestParameterSupported              bool     `json:"request_parameter_supported" description:""`
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported" description:""`
}

func (OpenIDConfiguration) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *OpenIDConfiguration) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal OpenIDConfiguration value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins OpenIDConfiguration) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(ins)
	return re, err
}

type ProviderExtend struct {
	WorkWeXinAgentId string `json:"workWeXinAgentId" description:""`
}

func (ProviderExtend) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *ProviderExtend) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal ProviderExtend value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins ProviderExtend) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(ins)
	return re, err
}

// LdapConfig holds configuration options for LDAP logins.
type LdapConfig struct {
	URL          string   `json:"url" yaml:"url"`
	BaseDn       string   `json:"baseDn" yaml:"baseDn"`
	UID          string   `json:"uid" yaml:"uid"`
	BindUser     string   `json:"bindUser" yaml:"bindUser"`
	BindPassword string   `json:"bindPassword" yaml:"bindPassword"`
	Filter       string   `json:"filter" yaml:"filter"`
	Attributes   []string `json:"attributes" yaml:"attributes"`
}

func (LdapConfig) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *LdapConfig) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal LdapConfig value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins LdapConfig) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(ins)
	return re, err
}

// ResponseError 错误响应
type ResponseError struct {
	//错误英文编码
	Message string `json:"message" yaml:"message" description:"错误英文编码"`
	//错误详情信息
	Detail string `json:"detail" yaml:"detail" description:"错误详情信息"`
	//支持I18N的提示信息
	Alert string `json:"alert" yaml:"alert" description:"支持I18N的提示信息"`
	//当前请求地址
	RequestURI string `json:"requestUri" description:"当前请求地址"`
}

// BatchOperationIds 需要删除的列表,根据数据库id
type BatchOperationIds struct {
	//需要删key 可以为数据库的id
	Ids []uint `json:"ids" validate:"required" description:"需要批量操作的数据库ID"`
}

// BatchOperationKeys 需要删除的列表，根据表唯一字符型字段
type BatchOperationKeys struct {
	//需要删key 可以为数据库的id
	Keys []string `json:"keys" description:"需要批量操作的数据库ID"`
}

func (receiver BatchOperationKeys) ToUint() (uints []uint) {
	for _, item := range receiver.Keys {
		key := common.StringsToUint(item)
		if key > 0 {
			uints = append(uints, key)
		}
	}
	return
}

// ExternalApp 外部应用
type ExternalApp struct {
	//应用名称
	Title string `json:"title" validate:"required" description:"应用名称"`
	//应用地址
	Url string `json:"url" validate:"required"  description:"应用地址"`
	//打开方式
	Target string `json:"target" validate:"required"  description:"打开方式"`
	//应用描述
	Desc string `json:"desc" validate:"required"  description:"应用描述"`
	//子应用
	Children []ExternalApp `json:"children" description:"子应用"`
}

// ThirdAuthMethod 支持的第三方登录方式
type ThirdAuthMethod struct {
	Oidcs           []OIDC `json:"oidcs" yaml:"oidcs" validate:"required" description:"OIDC登录"`
	MFA             bool   `json:"mfa" yaml:"mfa" validate:"required" description:"多因素认证开启"`
	FaceRecognition bool   `json:"faceRecognition" yaml:"faceRecognition" validate:"required" description:"人脸识别"`
}
type OIDC struct {
	//名称
	Name string `json:"name" validate:"required" description:"名称"`
	//认证地址
	Address string `json:"address" validate:"required" description:"授权地址"`
	//提供商类型
	Category string `json:"category" validate:"required" description:"提供商类型"`
}

// ArrayString 字符串数组
type ArrayString []string

// GormDataType gorm common data type
func (m ArrayString) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (ArrayString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "json"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (m *ArrayString) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal ArrayString value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, m)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (m ArrayString) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(m)
	return re, err
}

// ArrayUint 整数数组
type ArrayUint []uint

// GormDataType gorm common data type
func (u ArrayUint) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (ArrayUint) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "json"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (u *ArrayUint) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal ArrayUint value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, u)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (u ArrayUint) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(u)
	return re, err
}

// JsonMap json对象
type JsonMap map[string]interface{}

// GormDBDataType gorm db data type
func (JsonMap) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "json"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (u *JsonMap) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal ArrayUint value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, u)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (u JsonMap) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(u)
	return re, err
}

// ArrayFloat64 浮点数组
type ArrayFloat64 []float64

// GormDataType gorm common data type
func (u ArrayFloat64) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (ArrayFloat64) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "json"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (u *ArrayFloat64) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal ArrayUint value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, u)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (u ArrayFloat64) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(u)
	return re, err
}

func CheckFaceIdData(userFaceIdDatas []ArrayFloat64, compareFaceIdData ArrayFloat64) (ok bool) {
	for _, userFaceIdData := range userFaceIdDatas {
		var sumOfSquares float64
		for i := 0; i < len(userFaceIdData); i++ {
			diff := userFaceIdData[i] - compareFaceIdData[i]
			sumOfSquares += diff * diff
		}
		if math.Sqrt(sumOfSquares) < 0.25 {
			return true
		}
	}
	return

}

type MfaCode struct {
	//用户ID
	UserId uint `gorm:"column:user_id" json:"userId" validate:"required" description:"用户ID"`
	//验证码
	Code string `json:"code" description:"验证码"`
}
type SamlApplication struct {
	IDPMetadata string `json:"idpMetadata" description:""`
	SPMetadata  string `json:"spMetadata" description:""`
	PrivateKey  string `json:"privateKey" description:""`
	Certificate string `json:"certificate" description:""`
	MetadataURL string `json:"metadataUrl" description:""`
	SsoURL      string `json:"ssoUrl" description:""`
	LogoutURL   string `json:"logoutUrl" description:""`
	Metadata    string `json:"metadata" description:""`
}

func (app *SamlApplication) X509() (cert *x509.Certificate, private crypto.PrivateKey, err error) {
	b, _ := pem.Decode([]byte(app.Certificate))
	if b == nil {
		err = fmt.Errorf("decode application certificate fail: %s", "block is empty")
		return
	}
	cert, err = x509.ParseCertificate(b.Bytes)
	if err != nil {
		return
	}
	b, _ = pem.Decode([]byte(app.PrivateKey))
	if b == nil {
		err = fmt.Errorf("decode application private key fail: %s", "block is empty")
		return
	}
	private, err = x509.ParsePKCS8PrivateKey(b.Bytes)
	return
}

// SamlAuth saml认证
type SamlAuth struct {
	SamlRequest string `json:"SamlRequest" description:"SamlRequest"`
	RelayState  string `json:"RelayState" description:"RelayState"`
}

// UserThirdMethod 用户认证方式
type UserThirdMethod struct {
	Oidcs            []UserOIdcInfo        `json:"oidcs" validate:"required" description:"OIDC认证"`
	WebAuthns        []UserWebAuthn        `json:"webAuthns" description:"Web身份认证"`
	FaceRecognitions []UserFaceRecognition `json:"faceRecognitions" description:"人脸识别"`
	MFA              string                `json:"mfa"  validate:"required" description:"多因素认证开启"`
}

type UserOIdcInfo struct {
	//认证类型provider中的code,email,phone
	Provider string `gorm:"type:varchar(255);column:provider" json:"provider" description:"认证类型"`
	//第三方登录用户的id，邮箱，手机号
	LoginID string `gorm:"type:varchar(255);column:login_id" json:"loginId" description:"第三方认证的用户ID"`
	//第三方登录用户名
	LoginName string `gorm:"type:varchar(255);column:login_name" json:"loginName" description:"第三方认证的用户名"`
	//昵称
	Nickname string `gorm:"type:varchar(255);column:nickname" json:"nickname" description:"第三方认证的昵称"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//第三方认证的头像
	Avatar string `gorm:"type:longtext;column:avatar" json:"avatar" description:"第三方认证的头像"`
	//用户主页
	Home string `gorm:"type:varchar(500);column:home" json:"home" description:"用户主页"`
	//最近一次使用时间
	LatestUsedTime string `gorm:"type:varchar(255);column:latest_used_time" json:"latestUsedTime" description:"最近一次使用时间"`
}
type UserWebAuthn struct {
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//密钥名称
	Name string `gorm:"type:varchar(50);column:name" json:"name" validate:"required,max=50" description:"密钥名称"`
}

// UserFaceRecognition 人脸信息
type UserFaceRecognition struct {
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//录入名称
	Name string `gorm:"type:varchar(50);column:name" json:"name" validate:"required,max=50" description:"录入名称"`
}

type ApplicationRole struct {
	Name string `json:"name" description:"默认角色名称"`
	Role string `json:"role" description:"角色编码"`
	ZhCn string `json:"zhCn,omitempty" description:"中文"`
	EnUs string `json:"enUs,omitempty" description:"英文"`
}
type ApplicationRoles []ApplicationRole

// GormDataType gorm common data type
func (u ApplicationRoles) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (ApplicationRoles) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "json"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (u *ApplicationRoles) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal ArrayUint value: ", value))
	}
	err := jsoniter.Unmarshal(byteValue, u)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (u ApplicationRoles) Value() (driver.Value, error) {
	re, err := jsoniter.Marshal(u)
	return re, err
}

// CaptchaResponse 行为验证码生成数据
type CaptchaResponse struct {
	Code   string `json:"code" description:"编码"`
	Image  string `json:"image" description:"图像"`
	Prompt string `json:"prompt" description:"提示"`
	X      int    `json:"x" description:"x"`
	Y      int    `json:"y" description:"y"`
	Size   int    `json:"size" description:"size"`
	Width  int    `json:"width" description:"width"`
	Height int    `json:"height" description:"height"`
}

// CaptchaCheckResponse 行为验证码校验结果
type CaptchaCheckResponse struct {
	Success bool `json:"success" description:"校验成功"`
}

// CaptchaCheckData 行为验证码校验请求数据
type CaptchaCheckData struct {
	Code string `json:"code" description:"编码"`
	Data string `json:"data" description:"数据"`
}
