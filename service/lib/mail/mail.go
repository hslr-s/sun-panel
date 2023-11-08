package mail

import (
	"sun-panel/global"
)

// 发送注册验证码
//
//	@param emailer
//	@param mailTo 收件人
//	@param vcode 验证码
//	@return error
func SendRegisterEmail(emailer *Emailer, mailTo, vcode string) error {
	appName := global.Lang.Get("common.app_name")
	title := global.Lang.GetWithFields("mail.register_vcode_title", map[string]string{
		"AppName": appName,
	})
	content := global.Lang.GetWithFields("mail.register_vcode_content", map[string]string{
		"AppName": appName,
		"Minute":  "10",
	})
	err := emailer.SendMailOfVCode(mailTo, title, content, vcode)
	if err != nil {
		global.Logger.Errorf("failed to send email to %s, err:%+v\n", mailTo, err)
	}
	return err
}

// 发送重置密码验证码
//
//	@param emailer
//	@param mailTo
//	@param vcode
//	@return error
func SendResetPasswordVCode(emailer *Emailer, mailTo, vcode string) error {
	title := global.Lang.Get("mail.reset_password_password_title")
	content := global.Lang.Get("mail.reset_password_password_content")
	err := emailer.SendMailOfVCode(mailTo, title, content, vcode)
	if err != nil {
		global.Logger.Errorf("failed to send email to %s, err:%+v\n", mailTo, err)
	}
	return err
}

// // 事件提醒
// //
// //	@param emailer
// //	@param mailTo
// //	@param eventReminder
// //	@return error
// func SendEventReminder(emailer *Emailer, mailTo string, eventReminder models.EventReminder) error {
// 	startTime := eventReminder.Event.StartTime.Time.Format(cmn.TimeFormatMode4)
// 	endTime := eventReminder.Event.EndTime.Time.Format(cmn.TimeFormatMode4)
// 	title := global.Lang.GetWithFields("mail.reminder_title", map[string]string{
// 		"Title": eventReminder.Event.Title,
// 		"Time":  startTime,
// 	})
// 	contentParam := map[string]string{
// 		"ItemTitle": eventReminder.Event.Item.Title,
// 		"Time":      startTime,
// 	}
// 	content := "<p>" + global.Lang.GetWithFields("mail.reminder_content", contentParam) + "</p>"
// 	content += "<p>" + global.Lang.Get("mail.reminder_event_title") + " : " + eventReminder.Event.Title + "</p>"
// 	content += "<p>" + global.Lang.Get("common.start_time") + " : " + startTime + "</p>"
// 	content += "<p>" + global.Lang.Get("common.end_time") + " : " + endTime + "</p>"

// 	err := emailer.SendMail(mailTo, title, content)
// 	if err != nil {
// 		global.Logger.Errorf("failed to send email to %s, err:%+v\n", mailTo, err)
// 	}
// 	return err
// }
