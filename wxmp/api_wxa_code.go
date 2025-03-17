package wxmp

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/utils/cmp"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type WxaCodeOption func(*WxmpCodeConfig)

type LimitedWxmpCode struct {
	*WxmpCodeConfig
	// 扫码进入的小程序页面路径，最大长度 128 字节，不能为空；
	// 对于小游戏，可以只传入 query 部分，来实现传参效果，如：传入 "?foo=bar"，
	// 即可在 wx.getLaunchOptionsSync 接口中的 query 参数获取到 {foo:"bar"}。
	Path string `json:"path"`
}

type UnlimitedWxmpCode struct {
	*WxmpCodeConfig
	// 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `json:"scene"`
	// 页面 page，例如 pages/index/index，根路径前不要填加 /，不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
	Page      string `json:"page"`
	CheckPath bool   `json:"check_path"`
}

// WxmpCodeConfig 微信小程序码配置
type WxmpCodeConfig struct {
	// 要打开的小程序版本。正式版为 release，体验版为 trial，开发版为 develop
	EnvVersion string `json:"env_version"`
	// 二维码的宽度，单位 px。最小 280px，最大 1280px
	Width int `json:"width"`
	// auto_color 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color"`
	// auto_color 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
	LineColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	} `json:"line_color"`
	// 是否需要透明底色，为 true 时，生成透明底色的小程序码
	IsHyaline bool `json:"is_hyaline"`
}

const (
	urlGetLimitedWxaCode   = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	urlGetUnlimitedWxaCode = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

// CreateLimitedWxaCode 创建小程序码
func (impl *wxmpImpl) CreateLimitedWxaCode(path string, width int, options ...WxaCodeOption) ([]byte, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &LimitedWxmpCode{
		Path: path,
		WxmpCodeConfig: &WxmpCodeConfig{
			EnvVersion: "release",
			Width:      width,
			AutoColor:  true,
		},
	}
	for _, opt := range options {
		opt(body.WxmpCodeConfig)
	}

	resp, err := resty.New().R().SetBody(body).Post(fmt.Sprintf(urlGetLimitedWxaCode, accessToken))
	if err != nil {
		return nil, err
	}

	// 如果不是图像数据，那就是json错误数据
	if !cmp.IsImageData(resp.Body()) {
		return nil, errors.New(convert.BytesToString(resp.Body()))
	}

	return resp.Body(), nil
}

// CreateUnLimitedWxaCode 创建小程序码
func (impl *wxmpImpl) CreateUnLimitedWxaCode(scene, page string, width int, options ...WxaCodeOption) ([]byte, error) {
	accessToken, err := impl.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// 获取post的内容
	body := &UnlimitedWxmpCode{
		Scene: scene,
		Page:  page,
		WxmpCodeConfig: &WxmpCodeConfig{
			EnvVersion: "release",
			Width:      width,
			AutoColor:  true,
		},
	}
	for _, opt := range options {
		opt(body.WxmpCodeConfig)
	}

	resp, err := resty.New().R().SetBody(body).Post(fmt.Sprintf(urlGetUnlimitedWxaCode, accessToken))
	if err != nil {
		return nil, err
	}

	// 如果不是图像数据，那就是json错误数据
	if !cmp.IsImageData(resp.Body()) {
		return nil, errors.New(convert.BytesToString(resp.Body()))
	}

	return resp.Body(), nil
}

func Trial() WxaCodeOption {
	return func(c *WxmpCodeConfig) {
		c.EnvVersion = "trial"
	}
}

func Develop() WxaCodeOption {
	return func(c *WxmpCodeConfig) {
		c.EnvVersion = "develop"
	}
}
