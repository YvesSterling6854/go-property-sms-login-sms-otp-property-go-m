package main

import "fmt"

type MaintenanceRequest struct{ ID, Summary string }
type TenantDocument struct{ Name, Kind string }
type InspectionReminder struct{ Unit, Due string }

type SMSGateway interface {
	SendOTP(string) error
	VerifyOTP(string, string) (bool, error)
}

type LoginService struct{ SMS SMSGateway }

func (s LoginService) RequestCode(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone is required")
	}
	return s.SMS.SendOTP(phone)
}

func (s LoginService) VerifyCode(phone, code string) (string, error) {
	if phone == "" || code == "" {
		return "", fmt.Errorf("phone and code are required")
	}
	ok, err := s.SMS.VerifyOTP(phone, code)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("verification rejected")
	}
	return "tenant-session:" + phone, nil
}
