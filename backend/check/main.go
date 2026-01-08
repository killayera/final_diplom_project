package main

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "net"
    "net/smtp"
    "os"
    "path/filepath"
)

func main() {
    err := os.WriteFile("eicar.txt", []byte("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*"), 0644)
    if err != nil {
        fmt.Printf("❌ Error creating eicar.txt: %v\n", err)
        return
    }

    fmt.Println("=== 🚫 SENDING SUSPICIOUS EMAIL (FAKE DOMAIN) ===")
    sendEmail(
        "fakee@fake.com",    
        "john.doe@test.com",   
        "Urgent Invoice",      
        "./kk.exe",        
        true,                  
    )

    //fmt.Println("\n=== ✅ SENDING LEGITIMATE EMAIL ===")
    //sendEmail(
    //    "john.doe@test.com", 
    //    "damirtalipov@test.com",   
    //    "Welcome to TestMail!", 
    //    "./test.pdf",       
    //    false,                 
    //)
}

func sendEmail(from, to, subject, filePath string, fakeDKIM bool) {
    smtpAddr := "localhost:587"

    data, err := os.ReadFile(filePath)
    if err != nil {
        fmt.Printf("❌ Error reading file %s: %v\n", filePath, err)
        return
    }
    encoded := base64.StdEncoding.EncodeToString(data)
    boundary := "BOUNDARY12345"
    var msg bytes.Buffer

    
    if fakeDKIM {
       msg.WriteString("DKIM-Signature: v=1; a=rsa-sha256; d=example.com; s=fake; c=relaxed/simple; q=dns/txt; bh=somehash==; b=fakesignature==;\r\n")
    } else {
        msg.WriteString("DKIM-Signature: v=1; a=rsa-sha256; d=test.com; s=default; c=relaxed/simple; q=dns/txt; bh=validhash==; b=validsignature==;\r\n")
    }

    msg.WriteString("From: " + from + "\r\n")
    msg.WriteString("To: " + to + "\r\n")
    msg.WriteString("Subject: " + subject + "\r\n")
    msg.WriteString("MIME-Version: 1.0\r\n")
    msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
    msg.WriteString("\r\n")

    msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
    msg.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
    msg.WriteString("Please see the attached file.\r\n\r\n")

    msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
    msg.WriteString("Content-Type: application/octet-stream\r\n")
    msg.WriteString("Content-Transfer-Encoding: base64\r\n")
    msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%q\r\n\r\n", filepath.Base(filePath)))

    for i := 0; i < len(encoded); i += 76 {
        end := i + 76
        if end > len(encoded) {
            end = len(encoded)
        }
        msg.WriteString(encoded[i:end] + "\r\n")
    }
    msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

    conn, err := net.Dial("tcp", smtpAddr)
    if err != nil {
        fmt.Printf("❌ SMTP connection error: %v\n", err)
        return
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, "localtest.me")
    if err != nil {
        fmt.Printf("❌ SMTP client error: %v\n", err)
        return
    }
    defer client.Quit()

    if err := client.Mail(from); err != nil {
        fmt.Printf("❌ SMTP MAIL error: %v\n", err)
        return
    }
    if err := client.Rcpt(to); err != nil {
        fmt.Printf("❌ SMTP RCPT error: %v\n", err)
        return
    }

    w, err := client.Data()
    if err != nil {
        fmt.Printf("❌ SMTP DATA error: %v\n", err)
        return
    }
    _, err = w.Write(msg.Bytes())
    if err != nil {
        fmt.Printf("❌ SMTP write error: %v\n", err)
        return
    }
    if err := w.Close(); err != nil {
        fmt.Printf("❌ SMTP close error: %v\n", err)
        return
    }

    fmt.Printf("📤 Email from %s to %s sent successfully.\n", from, to)
}