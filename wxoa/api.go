package wxoa

import (
	"github.com/hdget/lib-wechat/wxoa/types"
)

type ApiWxoa interface {
	/*
	 * 参考：Js SDK签名：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
	 */
	GetJsSdkSignature(url string) (*types.JsSdkSignatureResult, error) // jsapi_ticket获取签名
	/*
	 * 校验微信公众号服务器
	 * https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
	 */
	VerifyServer(signature, token, timestamp, nonce string) bool
	/*
	 * 获取用户信息
	 * 参考：https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	 */
	GetUserInfo(openid string) (*types.UserInfoResult, error)

	// 消息处理
	HandleMessage(data []byte) ([]byte, error)
}
