package models

import (
	"github.com/google/uuid"
)

// StepUpMfaGrantType defines the type of grant for step-up MFA
type StepUpMfaGrantType string

const (
	// StepUpMfaGrantTypeOneTimeUse represents a one-time use grant
	StepUpMfaGrantTypeOneTimeUse StepUpMfaGrantType = "ONE_TIME_USE"
	// StepUpMfaGrantTypeTimeBased represents a time-based grant
	StepUpMfaGrantTypeTimeBased StepUpMfaGrantType = "TIME_BASED"
)

type MfaTotpType struct {
	Type string `json:"type"`
}

type MfaPhones struct {
	MfaPhoneNumberSuffix string `json:"mfa_phone_number_suffix"`
	MfaPhoneID           string `json:"mfa_phone_id"`
}

type MfaSetupType struct {
	Type         string       `json:"type"`
	PhoneNumbers *[]MfaPhones `json:"phone_numbers,omitempty"`
}

type FetchUserMfaMethodsResponse struct {
	MfaSetup MfaSetupType `json:"mfa_setup"`
}

// VerifyStepUpGrantRequest contains the parameters for verifying a step-up MFA grant
type VerifyStepUpGrantRequest struct {
	ActionType string    `json:"action_type"`
	UserID     uuid.UUID `json:"user_id"`
	Grant      string    `json:"grant"`
}

// StepUpMfaVerifyGrantResponse contains the response for step-up MFA verification
type StepUpMfaVerifyGrantResponse struct {
	Success bool `json:"success"`
}

// VerifyTotpChallengeRequest contains the parameters for verifying a TOTP challenge
type VerifyTotpChallengeRequest struct {
	ActionType      string             `json:"action_type"`
	UserID          uuid.UUID          `json:"user_id"`
	Code            string             `json:"code"`
	GrantType       StepUpMfaGrantType `json:"grant_type"`
	ValidForSeconds int                `json:"valid_for_seconds"`
}

type SendSmsMfaCodeRequest struct {
	ActionType      string             `json:"action_type"`
	UserID          uuid.UUID          `json:"user_id"`
	MfaPhoneID      uuid.UUID          `json:"mfa_phone_id"`
	GrantType       StepUpMfaGrantType `json:"grant_type"`
	ValidForSeconds int                `json:"valid_for_seconds"`
}

type VerifySmsChallengeRequest struct {
	ChallengeID string    `json:"challenge_id"`
	UserID      uuid.UUID `json:"user_id"`
	Code        string    `json:"code"`
}

// StepUpMfaVerifyTotpResponse contains the response with a step-up grant
type StepUpMfaVerifyTotpResponse struct {
	StepUpGrant string `json:"step_up_grant"`
}

type SendSmsMfaCodeResponse struct {
	ChallengeID string `json:"challenge_id"`
}

type VerifySmsChallengeResponse struct {
	StepUpGrant string `json:"step_up_grant"`
}
