package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

type EmailData struct {
	GuestName string
}

func SendRSVPConfirmation(toEmail, guestName string) error {
	// Load SMTP config
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	fromEmail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	// Create message
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Xác nhận tham dự - Lễ Tốt Nghiệp")

	// HTML template đơn giản
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
                <p>Cảm ơn bạn đã dành thời gian phản hồi lời mời tham dự lễ tốt nghiệp của tôi.</p>
				<p>Nếu có nhu cầu, hãy nhấp vào <a href="https://calendar.app.google/ns4vqiWnxy2PiYpdA">đây</a> để thêm sự kiện này vào ứng dụng Lịch (hoặc Google Calendar) và nhận thông báo nhé!</p>
                <p>Chúc bạn thật nhiều sức khoẻ và niềm vui!</p>
            </div>
            <div class="footer">
                <p>Trân trọng,<br><strong>Tô Hải Nhật</strong></p>
            </div>
        </div>
    </body>
    </html>
    `

	// Parse template
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %v", err)
	}

	// Prepare data - chỉ cần họ tên
	data := EmailData{
		GuestName: guestName,
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("template execute error: %v", err)
	}

	// Set body as HTML
	m.SetBody("text/html", body.String())

	// Create dialer
	d := gomail.NewDialer(smtpHost, smtpPort, fromEmail, password)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email error: %v", err)
	}

	return nil
}
