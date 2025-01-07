package svc

import (
	"fmt"
	"testing"
)

func TestEmailService_parseEmailBody(t *testing.T) {
	type args struct {
		body string
		data map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				body: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>Approval Email</title>\n    <style>\n        body {\n            font-family: Arial, sans-serif;\n            background-color: #f4f4f4;\n            margin: 0;\n            padding: 0;\n        }\n        .container {\n            width: 100%;\n            max-width: 600px;\n            margin: 0 auto;\n            background-color: #ffffff;\n            padding: 20px;\n            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\n        }\n        .header {\n            background-color: #4CAF50;\n            color: #ffffff;\n            padding: 10px 0;\n            text-align: center;\n        }\n        .content {\n            padding: 20px;\n        }\n        .content h1 {\n            color: #333333;\n        }\n        .content p {\n            color: #666666;\n            line-height: 1.6;\n        }\n        .footer {\n            text-align: center;\n            padding: 10px 0;\n            color: #999999;\n            font-size: 12px;\n        }\n        .button {\n            display: inline-block;\n            padding: 10px 20px;\n            margin: 20px 0;\n            background-color: #4CAF50;\n            color: #ffffff;\n            text-decoration: none;\n            border-radius: 5px;\n        }\n    </style>\n</head>\n<body>\n<div class=\"container\">\n    <div class=\"header\">\n        <h1>Approval Request</h1>\n    </div>\n    <div class=\"content\">\n        <h1>Hello {{Recipient Name}},</h1>\n        <p>We are pleased to inform you that your request for {{Request Details}} has been approved. Please find the details below:</p>\n        <p><strong>Request ID:</strong> {{id}}</p>\n        <p><strong>Approval Date:</strong> {{created_date}}</p>\n        <p>If you have any questions or need further assistance, please do not hesitate to contact us.</p>\n        <a href=\"{{Approval Link}}\" class=\"button\">View Details</a>\n    </div>\n    <div class=\"footer\">\n        <p>&copy; 2023 Your Company. All rights reserved.</p>\n    </div>\n</div>\n</body>\n</html>",
				data: map[string]interface{}{
					"id": "dwadwa",
				},
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EmailService{}
			got := e.parseEmailBody(tt.args.body, tt.args.data)
			fmt.Println("result ", got)
		})
	}
}
