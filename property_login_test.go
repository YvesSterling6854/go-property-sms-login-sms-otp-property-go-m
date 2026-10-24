package main

import "testing"

type fakeSMS struct{ accepted bool }

func (f fakeSMS) SendOTP(string) error                   { return nil }
func (f fakeSMS) VerifyOTP(string, string) (bool, error) { return f.accepted, nil }

func TestVerifyCodeDecision(t *testing.T) {
	for _, tc := range []struct {
		name, phone, code string
		accepted          bool
		wantSession       bool
	}{
		{"valid code", "+15551234567", "246810", true, true},
		{"rejected code", "+15551234567", "000000", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := LoginService{SMS: fakeSMS{accepted: tc.accepted}}
			if tc.phone == "" || tc.code == "" {
				t.Fatal("test input must be complete")
			}
			session, err := service.VerifyCode(tc.phone, tc.code)
			if tc.wantSession && (err != nil || session == "") {
				t.Fatalf("expected session, got %q, %v", session, err)
			}
			if !tc.wantSession && err == nil {
				t.Fatal("expected rejection")
			}
		})
	}
}
