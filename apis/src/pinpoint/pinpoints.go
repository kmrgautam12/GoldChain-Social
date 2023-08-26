package pinpointpkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpointemail"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

func SendPinPointAccountCreationMail(username string, toAddress []*string) (bool, error) {

	emailInput := pinpointemail.SendEmailInput{
		FromEmailAddress: aws.String(fromAddress),
		Destination: &pinpointemail.Destination{
			ToAddresses: toAddress,
		},
		Content: &pinpointemail.EmailContent{
			Simple: &pinpointemail.Message{
				Subject: &pinpointemail.Content{
					Data: aws.String("Welcome to goldchain"),
				},
				Body: &pinpointemail.Body{
					Html: &pinpointemail.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(fmt.Sprintf(EmailTemplate, username)),
					},
				},
			},
		},
	}

	pinPoint := pinpointemail.New(sess)
	_, err := pinPoint.SendEmail(&emailInput)
	if err != nil {
		return false, err
	}
	return true, nil

}
