package mail

import (
	"sun-panel/global"
	"sun-panel/models"

	"gopkg.in/gomail.v2"
)

type EmailInfo struct {
	Username string // 账号
	Password string // 密码
	Host     string // 服务器地址
	Port     int    // 端口 默认465
}

type Emailer struct {
	EmailInfo EmailInfo
	Dialer    *gomail.Dialer
}

func NewEmailer(emailInfo EmailInfo) *Emailer {
	dialer := gomail.NewDialer(emailInfo.Host, emailInfo.Port, emailInfo.Username, emailInfo.Password)
	return &Emailer{
		Dialer:    dialer,
		EmailInfo: emailInfo,
	}
}

func (e *Emailer) Send(mailTo []string, send_name, title, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.EmailInfo.Username)
	//这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)  //发送给多个用户
	m.SetHeader("Subject", title) //设置邮件主题
	m.SetBody("text/html", body)  //设置邮件正文

	return e.Dialer.DialAndSend(m)
}

// 发送邮件
func (m *Emailer) SendMail(mailTo string, title, content string) error {
	fromUrl := "127.0.0.1"
	appName := global.Lang.Get("common.app_name")
	mUser := models.User{Mail: mailTo}
	userInfo := mUser.GetUserInfoByMail()

	nickName := ""
	if userInfo != nil && userInfo.Name != "" {
		nickName = " " + userInfo.Name
	}

	body := `<meta charset="utf-8">
<table width="600px"  style="max-width: 600px;" align="center">
    <tr style="width: 600px;background-color: rgb(28, 197, 249);">
        <td align="left" style="width: 600px;padding: 22px 18px 11px;display: inline-block;">
            <div style="font-weight: 900;font-size: 18px;">
                <p>Hi` + nickName + `:</p>
            </div>
        </td>
        <td style="width: 100%;display: inline-block;border-top: 4px dashed rgb(255, 255, 255);"> </td>
        <td style="width: 600px;padding: 18px;display: inline-block;">
			<div align="left" style="color: rgb(57, 57, 57); line-height: 1.6; font-size: 14px; margin: 0px;font-weight: bolder;">
					` + content + `
			</div>
        </td>
        <td style="width: 600px;padding: 18px;display: inline-block;">
            <div align="rignt">
                <div style="font-size: 14px; margin: 0px;text-align: right;font-size: 14px; font-weight: bolder;">
                    -- ` + global.Lang.Get("mail.from") + `[<a href="` + fromUrl + `" style="color: #575757;">` + appName + `</a>]</div>
            </div>
        </td>
    </tr>
</table>`
	return SendMail(m, []string{mailTo}, appName, title, body)
}

// 发送链接邮件
func (m *Emailer) SendMailOfLink(mailTo, title, content, btn_name, url string) error {
	content = content + getLabelHtmlBtn(btn_name, url)
	return m.SendMail(mailTo, title, content)
}

// 发送注册邮件
func (m *Emailer) SendMailOfRegister(mailTo, key string) error {
	fromUrl := "127.0.0.1"
	appName := global.Lang.Get("common.app_name")
	title := global.Lang.GetWithFields("mail.register_title", map[string]string{
		"AppName": appName,
	})
	content := global.Lang.GetWithFields("mail.register_content", map[string]string{
		"AppName": appName,
	})
	return m.SendMailOfLink(mailTo, title, content, global.Lang.Get("mail.register_click_btn"), fromUrl+"/profile/auth.html#/linkRegister?code="+key)
}

// 发送验证码邮件
func (m *Emailer) SendMailOfVCode(mailTo, title, content, vCode string) error {
	// appName := global.Lang.Get("common.app_name")
	content = content + getLabelHtmlVcode(vCode)
	return m.SendMail(mailTo, title, content)
}

// 发送邮件
//
//	@param emailer
//	@param mailTo
//	@param send_name
//	@param title
//	@param body
//	@return error
func SendMail(emailer *Emailer, mailTo []string, send_name, title, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	if emailer.EmailInfo.Port == 0 {
		emailer.EmailInfo.Port = 465
	}
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(emailer.EmailInfo.Username, send_name))
	//这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)  //发送给多个用户
	m.SetHeader("Subject", title) //设置邮件主题
	m.SetBody("text/html", body)  //设置邮件正文

	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 跳过证书验证,不推荐
	err := emailer.Dialer.DialAndSend(m)
	return err
}

func getLabelHtmlBtn(btn_name string, href string) string {
	return `<div><a style="color: #fff;background-color: #2e2e2e;display: inline-block;padding: 10px 30px;border-radius: 5px;" href="` + href + `">` + btn_name + `</a></div>`
}

func getLabelHtmlVcode(vcode string) string {
	return `<p><div style="color: #fff;background-color: #2e2e2e;display: inline-block;padding: 10px 30px;border-radius: 5px;">` + vcode + `</div></p>`
}
