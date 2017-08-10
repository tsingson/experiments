package system_notify

import (
	"strings"
	"github.com/go-gomail/gomail"
	"html/template"
	"bytes"
	"os_adv_engine/system_cfg"
	"strconv"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type CsNotify interface {
	SendNotify(string, string) bool
}

var notifyList []CsNotify

func init() {
	smtp_s_str, _ := system_cfg.System_cfg.GetValue("email_notify", "smtp_server")
	smtp_p_str, _ := system_cfg.System_cfg.GetValue("email_notify", "smtp_port")
	sender_str, _ := system_cfg.System_cfg.GetValue("email_notify", "sender")
	passwd_str, _ := system_cfg.System_cfg.GetValue("email_notify", "passwd")

	receivers := []string{}
	receiversStr, _ := system_cfg.System_cfg.GetValue("email_notify", "receivers")
	for _, receiverStr := range strings.Split(receiversStr, ";") {
		receivers = append(receivers, strings.TrimSpace(receiverStr))
	}

	smtp_p_int, _ := strconv.Atoi(smtp_p_str)

	en := &EmailNotify{
		smtp_s: smtp_s_str,
		smtp_p: smtp_p_int,
		fromer: sender_str,
		toers: receivers,
		ccers: []string{},
		e_user: strings.Split(sender_str, "@")[0],
		e_passwd: passwd_str,
	}
	notifyList = append(notifyList, en)

	ln := &LanxinNotify{
		app_id:"",
		app_secret:"",
		host:"",
		tousers:"",
	}
	notifyList = append(notifyList, ln)
}

func Notify(title string, content string) {
	for _,value := range notifyList {
		value.SendNotify(title, content)
	}
}

type (
	EmailNotify struct {
		smtp_s string
		smtp_p int
		fromer string
		toers  []string
		ccers  []string
		e_user string
		e_passwd string
	}

	LanxinNotify struct {
		app_id string
		app_secret string
		host string
		tousers string
	}
)

func (ln * LanxinNotify)SendNotify(title string, content string) bool {

	return true
}

func (en *EmailNotify)SendNotify(title string, content string) bool {
	msg := gomail.NewMessage()
	msg.SetHeader("From", en.fromer)
	msg.SetHeader("To", en.toers...)
	msg.SetHeader("Subject", title)

	msg.SetBody("text/html", en.renderNotify(content))

	mailer := gomail.NewDialer(en.smtp_s, en.smtp_p, en.e_user, en.e_passwd)
	if err := mailer.DialAndSend(msg); err != nil {
		panic(err)
	}
	return true

}

func (en *EmailNotify) renderNotify(content string) string {
	tplStr := `<html>
<body>
 {{.}}
</table>
</body>
</html>`

	outBuf := &bytes.Buffer{}
	tpl := template.New("email notify template")
	tpl, _ = tpl.Parse(tplStr)
	tpl.Execute(outBuf, content)

	return outBuf.String()
}

func httpGet(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return fmt.Sprintf(string(body))
}
