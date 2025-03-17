package wxoa

type ApiWxoa interface {
	/*
	 * 参考：Js SDK签名：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
	 */
	GetJsSdkSignature(url string) (*JsSdkSignatureResult, error) // jsapi_ticket获取签名
	/*
	 * 校验微信公众号服务器
	 * https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
	 */
	VerifyServer(signature, token, timestamp, nonce string) bool
	/*
	 * 获取用户信息
	 * 参考：https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
	 */
	GetUserInfo(openid string) (*UserInfoResult, error)

	/*
	 * 接收普通消息,接收事件消息以及被动回复消息
	 * 参考：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Receiving_standard_messages.html
	 */
	HandleMessage(data []byte) ([]byte, error)

	/*
	 * 发送模板消息
	 * 参考：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
	 */
	SendTemplateMessage(toUser string, msg *TemplateMessage) error
}
