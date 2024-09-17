package helper

import (
	"net/smtp"
	"encoding/base64"
	"os"
	)
	

func SendEmail(userEmail, userName, hadiahORTokenORKeterangannota, typeSend, imagePath string) error {
	Auth := smtp.PlainAuth(
		"",
		"cabangbanyak@gmail.com",
		"lnbq rahl xyyg fwcy",
		"smtp.gmail.com",
	)

	var (
		subject string
		body string
		mime string
	)

	if typeSend == "AnnouncementGift" {
		subject = "Pemberitahuan: Hadiah Anda Sudah Siap untuk Diambil"
		// Isi email dalam format MIME
		body = "Dear " + userName + "\r\n\r\n" +
		"Salam sejahtera!\r\n\r\n" +
		"Kami dari Toko Bangunan Berkah Jaya ingin memberitahukan kepada Anda bahwa hadiah yang Anda tunggu telah siap untuk diambil. Kami senang memberitahu Anda bahwa Anda telah memenuhi syarat dan layak menerima hadiah Anda.\r\n\r\n" +
		"Detail Hadiah:\r\n" +
		"- Hadiah: " + hadiahORTokenORKeterangannota + "\r\n" +
		"- Cara Klaim: Hadiah Bisa langsung di Ambil di Toko\r\n\r\n" +
		"Terima kasih atas partisipasi Anda dan selamat atas penerimaan hadiah Anda. Kami berharap Anda akan menikmati pengalaman yang menyenangkan!\r\n\r\n" +
		"Salam hangat,\r\n" +
		"Toko Bangunan Berkah Jaya\r\n" +
		"Lebakwana, Kec. Kramatwatu, Kabupaten Serang, Banten"
	} else if typeSend == "ForgotPassword" {
		subject = "Ganti Password"
		body = "Dear " + userName + "\r\n" +
			"Please click http://localhost:8082/forgot/password/reset?token=" + hadiahORTokenORKeterangannota + " to reset your password."
	} else if typeSend == "NotaCancel" {
		subject = "Pemberitahuan: Nota Tidak Valid"
		mime = "MIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=BOUNDARY\r\n\r\n"

		// Read and encode the image file to base64
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			return err
		}
		encodedImage := base64.StdEncoding.EncodeToString(imageData)

		body = "--BOUNDARY\r\n" + mime +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			"Dear " + userName + ",\r\n\r\n" +
			"Kami ingin memberitahukan bahwa nota yang Anda kirim tidak valid.\r\n" +
			"Keterangan: " + hadiahORTokenORKeterangannota + "\r\n\r\n" +
			"Silakan periksa dan kirimkan ulang nota yang valid.\r\n\r\n" +
			"Terima kasih atas perhatian Anda.\r\n\r\n" +
			"Salam,\r\n" +
			"Toko Bangunan Berkah Jaya\r\n" +
			"--BOUNDARY\r\n" +
			"Content-Type: image/jpeg\r\n" +
			"Content-Transfer-Encoding: base64\r\n" +
			"Content-Disposition: attachment; filename=\"invalid_nota.jpg\"\r\n\r\n" +
			encodedImage + "\r\n" +
			"--BOUNDARY--"
	}

	msg := []byte("To: " + userEmail + "\r\n" + 
	"Subject: " + subject + "\r\n" + "\r\n" + body + "\r\n")

	// kirim mesage ke email user
	if err := smtp.SendMail(
		"smtp.gmail.com:587",
		Auth,
		"cabangbanyak@gmail.com",
		[]string{userEmail},
		msg,
	); err != nil {
		return err
	}
	
	return nil
}