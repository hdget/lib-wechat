package wxoa

type ApiWxoa interface {
	/*
	 * 参考：Js SDK签名：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
	 */
	GetJsSdkSignature(url string) (*JsSdkSignature, error) // jsapi_ticket获取签名
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

	// 接收消息处理
	HandleRecv(data []byte) ([]byte, error)
}

type MessageHandler interface {
	Handle() ([]byte, error)             // 处理
	GetDefaultCallback() MessageCallback // 缺省的处理方法
}

type MessageCallback interface {
	Execute(msg any) ([]byte, error)
}
