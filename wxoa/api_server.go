package wxoa

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

// VerifyServer 公众号接入时校验服务器
// 参考： https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
func (impl *wxoaImpl) VerifyServer(signature, token, timestamp, nonce string) bool {
	si := []string{token, timestamp, nonce}
	sort.Strings(si)              //字典序排序
	str := strings.Join(si, "")   //组合字符串
	s := sha1.New()               //返回一个新的使用SHA1校验的hash.Hash接口
	_, _ = io.WriteString(s, str) //WriteString函数将字符串数组str中的内容写入到s中
	calculatedSignature := fmt.Sprintf("%x", s.Sum(nil))
	return signature == calculatedSignature
}
