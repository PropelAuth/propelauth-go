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
	ActionType      string          `json:"action_type"`
	UserID          uuid.UUID       `json:"user_id"`
	Code            string          `json:"code"`
	GrantType       StepUpMfaGrantType `json:"grant_type"`
	ValidForSeconds int             `json:"valid_for_seconds"`
}

// StepUpMfaVerifyTotpResponse contains the response with a step-up grant
type StepUpMfaVerifyTotpResponse struct {
	StepUpGrant string `json:"step_up_grant"`
}