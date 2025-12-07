package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type EmailData struct {
	GuestName string
}

func SendRSVPConfirmation(toEmail, guestName string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY is not set")
	}

	// tạo client Brevo API
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	client := brevo.NewAPIClient(cfg)

	// HTML template
	htmlTemplate := `
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            body { 
                font-family: Arial, sans-serif; 
                line-height: 1.6; 
                color: #333; 
                background-color: #f5f5f5;
                padding: 20px;
            }
            .container { 
                max-width: 600px; 
                margin: 0 auto; 
                background: white;
                border-radius: 10px;
                overflow: hidden;
                box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            }
            .header { 
                background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
                color: white; 
                padding: 40px 30px; 
                text-align: center;
            }
            .header h1 {
                margin: 0;
                font-size: 32px;
            }
            .content { 
                padding: 40px 30px;
            }
            .content h2 {
                color: #667eea;
                margin-top: 0;
            }
            .footer { 
                text-align: center; 
                color: #666; 
                padding: 20px;
                background: #f9fafb;
                font-size: 14px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="content">
                <h2>Xin chào {{.GuestName}}!</h2>
                <p>Cảm ơn bạn đã dành thời gian phản hồi lời mời tham dự lễ tốt nghiệp của mình.</p>
				<p>Nếu có nhu cầu, hãy nhấp vào <a href="https://calendar.app.google/uX6cR4BqkqQRan817">đây</a> để thêm sự kiện này vào ứng dụng Lịch trên điện thoại và nhận thông báo nhé!</p>
                <p>Chúc bạn thật nhiều sức khoẻ, niềm vui và có một mùa Giáng Sinh an lành!</p>
            </div>
            <div class="footer">
                <p>Trân trọng,<br><strong>Tô Hải Nhật</strong></p>
            </div>
        </div>
    </body>
    </html>
    `

	// parse template
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %v", err)
	}

	// chuẩn bị dữ liệu
	data := EmailData{
		GuestName: guestName,
	}

	// execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("template execute error: %v", err)
	}

	// lấy thông tin người gửi từ env
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderName := os.Getenv("SENDER_NAME")
	if senderEmail == "" {
		senderEmail = "noreply@example.com"
	}
	if senderName == "" {
		senderName = "Tô Hải Nhật"
	}
	// tạo email request
	sendSmtpEmail := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  senderName,
			Email: senderEmail,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: toEmail,
				Name:  guestName,
			},
		},
		Subject:     "Xác nhận tham dự - Lễ Tốt Nghiệp",
		HtmlContent: body.String(),
	}
	// send email
	_, _, err = client.TransactionalEmailsApi.SendTransacEmail(context.Background(), sendSmtpEmail)
	if err != nil {
		return fmt.Errorf("send email error: %v", err)
	}

	return nil
}
