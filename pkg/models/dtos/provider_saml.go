package dtos

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"
	"time"

	"github.com/efucloud/common"
	"github.com/go-playground/validator/v10"
)

// ProviderSamlDetailList SAML提供商列表响应
type ProviderSamlDetailList struct {
	//当前页数据
	Data []*ProviderSamlDetail `json:"data"`
	//数据库满足条件的数据总数
	Total int64 `json:"total,omitempty" validate:"required"`
}

// ProviderSamlDetail SAML提供商详情
type ProviderSamlDetail struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" description:"记录ID"`
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt" description:"创建时间"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt,omitempty" description:"更新时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"提供商名称"`
	//提供商编码
	Category string `gorm:"type:varchar(255)" json:"category" validate:"required" description:"提供商编码"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//IdP EntityID
	EntityId string `gorm:"type:varchar(500);column:entity_id" json:"entityId" validate:"required" description:"IdP EntityID"`
	//SAML SSO地址
	SsoURL string `gorm:"type:varchar(1000);column:sso_url" json:"ssoUrl" validate:"required" description:"SAML SSO地址"`
	//回调地址(AssertionConsumerService URL)
	AcsURL string `gorm:"type:varchar(1000);column:acs_url" json:"acsUrl" validate:"required" description:"回调地址"`
	//SAML证书
	Certificate string `gorm:"type:longtext;column:certificate" json:"certificate" validate:"required" description:"SAML证书"`
	//元数据地址
	MetadataURL string `gorm:"type:varchar(1000);column:metadata_url" json:"metadataUrl" description:"元数据地址"`
	//元数据内容
	Metadata string `gorm:"type:longtext;column:metadata" json:"metadata" description:"元数据内容"`
	//登录ID字段映射
	LoginIDAttr string `gorm:"type:varchar(255);column:login_id_attr" json:"loginIdAttr" description:"登录ID字段映射"`
	//登录名字段映射
	LoginNameAttr string `gorm:"type:varchar(255);column:login_name_attr" json:"loginNameAttr" description:"登录名字段映射"`
	//邮箱字段映射
	EmailAttr string `gorm:"type:varchar(255);column:email_attr" json:"emailAttr" description:"邮箱字段映射"`
	//手机号字段映射
	PhoneAttr string `gorm:"type:varchar(255);column:phone_attr" json:"phoneAttr" description:"手机号字段映射"`
	//昵称字段映射
	NicknameAttr string `gorm:"type:varchar(255);column:nickname_attr" json:"nicknameAttr" description:"昵称字段映射"`
	//头像字段映射
	AvatarAttr string `gorm:"type:varchar(255);column:avatar_attr" json:"avatarAttr" description:"头像字段映射"`
}

// ProviderSamlCreate SAML提供商创建
type ProviderSamlCreate struct {
	//创建时间
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-" description:"创建时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"提供商名称"`
	//提供商编码
	Category string `gorm:"type:varchar(255)" json:"category" validate:"required" description:"提供商编码"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//IdP EntityID
	EntityId string `gorm:"type:varchar(500);column:entity_id" json:"entityId" validate:"required" description:"IdP EntityID"`
	//SAML SSO地址
	SsoURL string `gorm:"type:varchar(1000);column:sso_url" json:"ssoUrl" validate:"required" description:"SAML SSO地址"`
	//回调地址(AssertionConsumerService URL)
	AcsURL string `gorm:"type:varchar(1000);column:acs_url" json:"acsUrl" validate:"required" description:"回调地址"`
	//SAML证书
	Certificate string `gorm:"type:longtext;column:certificate" json:"certificate" validate:"required" description:"SAML证书"`
	//元数据地址
	MetadataURL string `gorm:"type:varchar(1000);column:metadata_url" json:"metadataUrl" description:"元数据地址"`
	//元数据内容
	Metadata string `gorm:"type:longtext;column:metadata" json:"metadata" description:"元数据内容"`
	//登录ID字段映射
	LoginIDAttr string `gorm:"type:varchar(255);column:login_id_attr" json:"loginIdAttr" description:"登录ID字段映射"`
	//登录名字段映射
	LoginNameAttr string `gorm:"type:varchar(255);column:login_name_attr" json:"loginNameAttr" description:"登录名字段映射"`
	//邮箱字段映射
	EmailAttr string `gorm:"type:varchar(255);column:email_attr" json:"emailAttr" description:"邮箱字段映射"`
	//手机号字段映射
	PhoneAttr string `gorm:"type:varchar(255);column:phone_attr" json:"phoneAttr" description:"手机号字段映射"`
	//昵称字段映射
	NicknameAttr string `gorm:"type:varchar(255);column:nickname_attr" json:"nicknameAttr" description:"昵称字段映射"`
	//头像字段映射
	AvatarAttr string `gorm:"type:varchar(255);column:avatar_attr" json:"avatarAttr" description:"头像字段映射"`
}

func (ins *ProviderSamlCreate) Default(ctx context.Context) {
	ins.CreatedAt = time.Now()
	applyProviderSamlDefaults(&ins.EntityId, &ins.SsoURL, &ins.Certificate, &ins.LoginIDAttr, &ins.LoginNameAttr, &ins.EmailAttr, &ins.PhoneAttr, &ins.NicknameAttr, &ins.AvatarAttr, ins.Metadata)
}
func (ins *ProviderSamlCreate) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(ins)
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

// ProviderSamlUpdate SAML提供商修改
type ProviderSamlUpdate struct {
	//主键
	ID uint `gorm:"primarykey;column:id" json:"id" validate:"required" description:"记录ID"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//提供商名称
	Name string `gorm:"type:varchar(255)" json:"name" validate:"required" description:"提供商名称"`
	//提供商编码
	Category string `gorm:"type:varchar(255)" json:"category" validate:"required" description:"提供商编码"`
	//是否有效
	Enable bool `gorm:"column:enable;default:true" json:"enable" description:"是否有效"`
	//IdP EntityID
	EntityId string `gorm:"type:varchar(500);column:entity_id" json:"entityId" validate:"required" description:"IdP EntityID"`
	//SAML SSO地址
	SsoURL string `gorm:"type:varchar(1000);column:sso_url" json:"ssoUrl" validate:"required" description:"SAML SSO地址"`
	//回调地址(AssertionConsumerService URL)
	AcsURL string `gorm:"type:varchar(1000);column:acs_url" json:"acsUrl" validate:"required" description:"回调地址"`
	//SAML证书
	Certificate string `gorm:"type:longtext;column:certificate" json:"certificate" validate:"required" description:"SAML证书"`
	//元数据地址
	MetadataURL string `gorm:"type:varchar(1000);column:metadata_url" json:"metadataUrl" description:"元数据地址"`
	//元数据内容
	Metadata string `gorm:"type:longtext;column:metadata" json:"metadata" description:"元数据内容"`
	//登录ID字段映射
	LoginIDAttr string `gorm:"type:varchar(255);column:login_id_attr" json:"loginIdAttr" description:"登录ID字段映射"`
	//登录名字段映射
	LoginNameAttr string `gorm:"type:varchar(255);column:login_name_attr" json:"loginNameAttr" description:"登录名字段映射"`
	//邮箱字段映射
	EmailAttr string `gorm:"type:varchar(255);column:email_attr" json:"emailAttr" description:"邮箱字段映射"`
	//手机号字段映射
	PhoneAttr string `gorm:"type:varchar(255);column:phone_attr" json:"phoneAttr" description:"手机号字段映射"`
	//昵称字段映射
	NicknameAttr string `gorm:"type:varchar(255);column:nickname_attr" json:"nicknameAttr" description:"昵称字段映射"`
	//头像字段映射
	AvatarAttr string `gorm:"type:varchar(255);column:avatar_attr" json:"avatarAttr" description:"头像字段映射"`
}

func (ins *ProviderSamlUpdate) Default(ctx context.Context) {
	ins.UpdatedAt = time.Now()
	applyProviderSamlDefaults(&ins.EntityId, &ins.SsoURL, &ins.Certificate, &ins.LoginIDAttr, &ins.LoginNameAttr, &ins.EmailAttr, &ins.PhoneAttr, &ins.NicknameAttr, &ins.AvatarAttr, ins.Metadata)
}
func (ins *ProviderSamlUpdate) Validate(ctx context.Context) (err error) {
	validate := validator.New()
	lang := common.GetLangFromCtx(ctx, "")
	validate.RegisterTagNameFunc(common.TagNameI18N(lang))
	trans := common.LoadValidateTranslator(lang, validate)
	err = validate.Struct(ins)
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

// ProviderSamlStatus 认证提供商状态
type ProviderSamlStatus struct {
	//主键
	Ids []uint `json:"ids" validate:"required" description:"主键"`
	//更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"-" description:"更新时间"`
	//是否有效
	Enable bool `json:"enable" description:"是否有效"`
}

func (t *ProviderSamlStatus) Default(ctx context.Context) {
	t.UpdatedAt = time.Now()
}
func (t *ProviderSamlStatus) Validate(ctx context.Context) (err error) {
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

func applyProviderSamlDefaults(entityID, ssoURL, certificate, loginIDAttr, loginNameAttr, emailAttr, phoneAttr, nicknameAttr, avatarAttr *string, metadata string) {
	metaEntityID, metaSSOURL, metaCert := parseSamlMetadata(metadata)
	if len(strings.TrimSpace(*entityID)) == 0 {
		*entityID = metaEntityID
	}
	if len(strings.TrimSpace(*ssoURL)) == 0 {
		*ssoURL = metaSSOURL
	}
	if len(strings.TrimSpace(*certificate)) == 0 {
		*certificate = metaCert
	}
	*certificate = normalizeSamlCertificate(*certificate)
	if len(strings.TrimSpace(*loginIDAttr)) == 0 {
		*loginIDAttr = "NameID"
	}
	if len(strings.TrimSpace(*loginNameAttr)) == 0 {
		*loginNameAttr = "name"
	}
	if len(strings.TrimSpace(*emailAttr)) == 0 {
		*emailAttr = "email"
	}
	if len(strings.TrimSpace(*phoneAttr)) == 0 {
		*phoneAttr = "phone_number"
	}
	if len(strings.TrimSpace(*nicknameAttr)) == 0 {
		*nicknameAttr = "nickname"
	}
	if len(strings.TrimSpace(*avatarAttr)) == 0 {
		*avatarAttr = "picture"
	}
}

type samlMetadataEntitiesDescriptor struct {
	EntityDescriptors []samlMetadataEntityDescriptor `xml:"EntityDescriptor"`
}

type samlMetadataEntityDescriptor struct {
	EntityID          string                         `xml:"entityID,attr"`
	IDPSSODescriptors []samlMetadataIDPSSODescriptor `xml:"IDPSSODescriptor"`
}

type samlMetadataIDPSSODescriptor struct {
	SingleSignOnServices []samlMetadataService       `xml:"SingleSignOnService"`
	KeyDescriptors       []samlMetadataKeyDescriptor `xml:"KeyDescriptor"`
}

type samlMetadataService struct {
	Binding  string `xml:"Binding,attr"`
	Location string `xml:"Location,attr"`
}

type samlMetadataKeyDescriptor struct {
	Use     string              `xml:"use,attr"`
	KeyInfo samlMetadataKeyInfo `xml:"KeyInfo"`
}

type samlMetadataKeyInfo struct {
	X509Data samlMetadataX509Data `xml:"X509Data"`
}

type samlMetadataX509Data struct {
	Certificates []string `xml:"X509Certificate"`
}

func parseSamlMetadata(raw string) (entityID, ssoURL, certificate string) {
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		return
	}
	var entity samlMetadataEntityDescriptor
	if err := xml.Unmarshal([]byte(raw), &entity); err == nil && (len(entity.EntityID) > 0 || len(entity.IDPSSODescriptors) > 0) {
		return parseSamlEntityDescriptor(entity)
	}
	var entities samlMetadataEntitiesDescriptor
	if err := xml.Unmarshal([]byte(raw), &entities); err != nil {
		return
	}
	for _, item := range entities.EntityDescriptors {
		entityID, ssoURL, certificate = parseSamlEntityDescriptor(item)
		if len(ssoURL) > 0 || len(certificate) > 0 {
			return
		}
	}
	return
}

func parseSamlEntityDescriptor(entity samlMetadataEntityDescriptor) (entityID, ssoURL, certificate string) {
	entityID = strings.TrimSpace(entity.EntityID)
	for _, descriptor := range entity.IDPSSODescriptors {
		if len(ssoURL) == 0 {
			ssoURL = selectSamlSSOServiceURL(descriptor.SingleSignOnServices)
		}
		if len(certificate) == 0 {
			certificate = selectSamlCertificate(descriptor.KeyDescriptors)
		}
	}
	certificate = normalizeSamlCertificate(certificate)
	return
}

func selectSamlSSOServiceURL(services []samlMetadataService) string {
	var first string
	for _, service := range services {
		location := strings.TrimSpace(service.Location)
		if len(location) == 0 {
			continue
		}
		if len(first) == 0 {
			first = location
		}
		binding := strings.TrimSpace(service.Binding)
		if strings.HasSuffix(binding, "HTTP-Redirect") {
			return location
		}
	}
	return first
}

func selectSamlCertificate(keys []samlMetadataKeyDescriptor) string {
	var first string
	for _, key := range keys {
		for _, cert := range key.KeyInfo.X509Data.Certificates {
			cert = strings.TrimSpace(cert)
			if len(cert) == 0 {
				continue
			}
			if len(first) == 0 {
				first = cert
			}
			if strings.EqualFold(strings.TrimSpace(key.Use), "signing") {
				return cert
			}
		}
	}
	return first
}

func normalizeSamlCertificate(cert string) string {
	cert = strings.TrimSpace(cert)
	if len(cert) == 0 {
		return cert
	}
	if strings.Contains(cert, "BEGIN CERTIFICATE") {
		return cert
	}
	clean := strings.Map(func(r rune) rune {
		switch r {
		case '\n', '\r', '\t', ' ':
			return -1
		default:
			return r
		}
	}, cert)
	if len(clean) == 0 {
		return cert
	}
	var b strings.Builder
	b.WriteString("-----BEGIN CERTIFICATE-----\n")
	for len(clean) > 64 {
		b.WriteString(clean[:64])
		b.WriteString("\n")
		clean = clean[64:]
	}
	b.WriteString(clean)
	b.WriteString("\n-----END CERTIFICATE-----")
	return b.String()
}
