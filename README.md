# Data-Drive-System
task



# 1. Clone the repository
git clone https://github.com/yourusername/datadrive.git
cd datadrive

# 2. Install dependencies
go mod tidy

# 3. Configure email settings (required for OTP/email functionality)
# ----------------------------------------
# Go to: pkg/helper/otphelper/otphelper.go

# Locate the line:
d := gomail.NewDialer("smtp.zoho.com", 587, "YOUR_EMAIL", "YOUR_PASS")

# ➤ If you're using Zoho Mail, keep "smtp.zoho.com" as it is.
# ➤ If you're using another provider (like Gmail, Outlook, etc.), change it accordingly:
#     Gmail:    smtp.gmail.com
#     Outlook:  smtp.office365.com
#     Zoho:     smtp.zoho.com

# Replace "YOUR_EMAIL" and "YOUR_PASS" with your real email credentials.

# Save the file before proceeding.
# ----------------------------------------

# 4. Run the server
go run cmd/main.go task
