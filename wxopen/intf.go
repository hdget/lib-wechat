package wxopen

type ApiWxopen interface {
	/*
	 * 网站应用快速扫码登录
	 * 参考：https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	 */
	Login(code string) (string, string, error) // 网站应用扫码快速登录
}
