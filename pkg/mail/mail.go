package mail

import(
	"net/smtp"
)

func SendEmail(to []string, subject string, body string) error{
	auth := smtp.PlainAuth(
		"",
		"snetverify@gmail.com",
		"wcdu aylg kkqx bnqc",
		"smtp.gmail.com",
	)
	message := "Subject: " + subject + "\r\n" +
               "Content-Type: text/html; charset=\"utf-8\"\r\n" +
               "\r\n" + 
               body

    return smtp.SendMail("smtp.gmail.com:587", auth, "snetverify@gmail.com", to, []byte(message))
}