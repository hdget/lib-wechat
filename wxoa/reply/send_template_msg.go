package reply

import (
	"base/g"
	"github.com/hdget/common/intf"
)

func SendTemplateMessage(redisProvider intf.RedisProvider, toUserOpenId, templateId, mpPath string, contents map[string]string) error {
	msg, err := NewTemplateSendMessage(redisProvider, &SendMessageTemplateArgument{
		AppId:        g.Config.App.Wxoa.AppId,
		AppSecret:    g.Config.App.Wxoa.AppSecret,
		ToUserOpenId: toUserOpenId,
		TemplateId:   templateId,
		MiniProgram: &SendMessageTemplateMiniProgram{
			AppId:    g.Config.App.Wxmp.AppId,
			PagePath: mpPath,
		},
	})
	if err != nil {
		return err
	}

	err = msg.Send(contents)
	if err != nil {
		return err
	}
	return nil
}
