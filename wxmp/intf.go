package wxmp

type ApiWxmp interface {
	/*
	 * 通过code换取openid和unionid
	 *
	 * 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	 */
	Login(code string) (string, string, error) // 获取OpenId和UnionId

	/*
	 * 解密用户信息
	 *
	 * 参考：https://developers.weixin.qq.com/miniprogram/dev/api/open-api/user-info/wx.getUserInfo.html
	 */
	DecryptUserInfo(encryptedData, iv string) (*UserInfo, error) // 解密有用户信息

	/*
	 * Deprecated: 已弃用，请使用GetUserPhoneNumber
	 */
	DecryptMobileInfo(encryptedData, iv string) (*MobileInfo, error) // 解密手机号码信息
	/*
	 * 创建永久有效的小程序码，可接受path参数较长，生成个数受限
	 * 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/createQRCode.html
	 */
	CreateLimitedWxaCode(path string, width int, options ...WxaCodeOption) ([]byte, error)
	/*
	 * 生成小程序码，可接受页面参数较短，生成个数不受限
	 * 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getUnlimitedQRCode.html
	 */
	CreateUnLimitedWxaCode(scene, page string, width int, options ...WxaCodeOption) ([]byte, error)
	/*
	 * 快速手机号验证
	 * 参考: https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-info/phone-number/getPhoneNumber.html
	 */
	GetUserPhoneNumber(code string) (*MobileInfo, error)
}
