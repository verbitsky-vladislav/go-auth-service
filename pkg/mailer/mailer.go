package mailer

const (
	SMTPServer = "smtp.gmail.com"
)

type mailer struct {
	User     string
	Password string
}

func NewMailer(Username, Password string) mailer {
	return mailer{Username, Password}
}
