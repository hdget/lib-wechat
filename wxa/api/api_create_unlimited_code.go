package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/hdget/utils/cmp"
	"github.com/pkg/errors"
)

type getUnlimitedWxaCodeRequest struct {
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
	// 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `json:"scene"`
	// 页面 page，例如 pages/index/index，根路径前不要填加 /，不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
	Page      string `json:"page"`
	CheckPath bool   `json:"check_path"`
}

const (
	// 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getUnlimitedQRCode.html
	urlGetUnlimitedWxaCode = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

// ApiCreateUnLimitedWxaCode 创建小程序码
func (impl apiImpl) CreateUnLimitedWxaCode(accessToken string, scene, page string, width int) ([]byte, error) {
	// 获取post的内容
	req := &getUnlimitedWxaCodeRequest{
		Scene:      scene,
		Page:       page,
		EnvVersion: "release",
		Width:      width,
		AutoColor:  true,
	}

	url := fmt.Sprintf(urlGetUnlimitedWxaCode, accessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	// 如果不是图像数据，那就是json错误数据
	// 如果不是图像数据，那就是json错误数据
	if !cmp.IsImageData(resp.Body()) {
		var result api.Result
		err = json.Unmarshal(resp.Body(), &result)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(result.ErrMsg)
	}

	return resp.Body(), nil
}
