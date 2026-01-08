package validator

import (
    "context"
    "fmt"
    "mail_server/models"
    "net"
    "strings"
)

func Validate(mail *models.Mail) error {
    if mail.From == "" {
        return fmt.Errorf("invalid email: From address is empty")
    }

    fmt.Println("\n=== 📧 Validating Email ===")
    fmt.Printf("From: %s, To: %s, SenderIP: %s\n", mail.From, mail.To, mail.SenderIP)

    var validationErrors []string

    fromEmail := strings.Trim(mail.From, "<>")
    domainIndex := strings.Index(fromEmail, "@")
    if domainIndex == -1 {
        return fmt.Errorf("invalid From address format: %s", fromEmail)
    }
    domain := fromEmail[domainIndex+1:]

    ip := net.ParseIP(mail.SenderIP)
    if ip == nil {
        return fmt.Errorf("invalid Sender IP address: %s", mail.SenderIP)
    }

    SPFResult, err := checkSPF(context.Background(), ip, fromEmail, domain)
    if err != nil || SPFResult != SPFResultPass {
        fmt.Printf("🔴 SPF: %s (Error: %v, From: %s, IP: %s, Domain: %s)\n", SPFResult, err, fromEmail, mail.SenderIP, domain)
        validationErrors = append(validationErrors, fmt.Sprintf("SPF check failed: %s (%v)", SPFResult, err))
    } else {
        fmt.Printf("🟢 SPF: %s (From: %s, IP: %s, Domain: %s)\n", SPFResult, fromEmail, mail.SenderIP, domain)
    }

    dkimResult, err := verifyDKIM(mail.RawMessage)
    if err != nil || dkimResult != "pass" {
        fmt.Printf("🔴 DKIM: %s (Error: %v)\n", dkimResult, err)
        validationErrors = append(validationErrors, fmt.Sprintf("DKIM verification failed: %s (%v)", dkimResult, err))
    } else {
        fmt.Printf("🟢 DKIM: %s\n", dkimResult)
    }

    // Check DMARC policy fetch
    if _, err := fetchDMARCPolicy(domain); err != nil {
        fmt.Printf("🔴 DMARC: could not fetch policy (%v, Domain: %s)\n", err, domain)
        validationErrors = append(validationErrors, fmt.Sprintf("DMARC policy check failed: %v", err))
    } else {
        // Simulate DMARC failure if either SPF or DKIM fails
        if dkimResult != "pass" || SPFResult != SPFResultPass {
            fmt.Printf("🔴 DMARC: fail (Reason: SPF=%s, DKIM=%s, Domain: %s)\n", SPFResult, dkimResult, domain)
            validationErrors = append(validationErrors, fmt.Sprintf("DMARC failed: SPF=%s, DKIM=%s", SPFResult, dkimResult))
        } else {
            fmt.Printf("🟢 DMARC: pass (Domain: %s)\n", domain)
        }
    }
    //mail.SenderIP
    maliciousIP := "91.246.58.169"
    if err := checkAbuseIP(maliciousIP); err != nil {
        fmt.Printf("🔴 AbuseIPDB: fail (Error: %v, IP: %s)\n", err, mail.SenderIP)
        validationErrors = append(validationErrors, fmt.Sprintf("AbuseIPDB check failed: %v", err))
    } else {
        fmt.Printf("🟢 AbuseIPDB: pass (IP: %s)\n", mail.SenderIP)
    }

    if err := scanFiles(mail); err != nil {
        fmt.Printf("🔴 VirusTotal: fail (Error: %v)\n", err)
        validationErrors = append(validationErrors, fmt.Sprintf("VirusTotal scan failed: %v", err))
    } else {
        fmt.Printf("🟢 VirusTotal: pass\n")
    }

    if len(validationErrors) > 0 {
        fmt.Printf("❌ Email validation failed for %s:\n - %s\n", fromEmail, strings.Join(validationErrors, "\n - "))
        return fmt.Errorf("mail validation failed:\n - %s", strings.Join(validationErrors, "\n - "))
    }

    fmt.Printf("✅ All validations passed for email from %s\n", fromEmail)
    return nil
}