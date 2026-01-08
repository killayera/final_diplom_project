package validator

import (
	"fmt"
	"net"
	"strings"
)

func fetchDMARCPolicy(domain string) (string, error) {
	records, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return "", err
	}
	for _, record := range records {
		if strings.Contains(record, "v=DMARC1") {
			if strings.Contains(record, "p=reject") {
				return "reject", fmt.Errorf("DKIM verification failed")
			}
		}
	}
	return "pass", nil
}
