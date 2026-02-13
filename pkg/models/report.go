package models

import "github.com/google/uuid"

type ReportPagination struct {
	PageSize   *int
	PageNumber *int
}

// org report types

type OrgReportRecord struct {
	Id              string         `json:"id"`
	ReportId        string         `json:"report_id"`
	OrgId           uuid.UUID      `json:"org_id"`
	Name            string         `json:"name"`
	NumUsers        int            `json:"num_users"`
	OrgCreatedAt    int64          `json:"org_created_at"`
	ExtraProperties map[string]any `json:"extra_properties"`
}

type OrgReport struct {
	OrgReports     []OrgReportRecord `json:"org_reports"`
	CurrentPage    int               `json:"current_page"`
	TotalCount     int               `json:"total_count"`
	PageSize       int               `json:"page_size"`
	HasMoreResults bool              `json:"has_more_results"`
	ReportTime     int               `json:"report_time"`
}

// user report types

type UserOrgMembershipForReport struct {
	DisplayName string    `json:"display_name"`
	OrgId       uuid.UUID `json:"org_id"`
	UserRole    string    `json:"user_role"`
}

type UserReportRecord struct {
	Id              string                       `json:"id"`
	ReportId        string                       `json:"report_id"`
	UserId          uuid.UUID                    `json:"user_id"`
	Email           string                       `json:"email"`
	UserCreatedAt   int64                        `json:"user_created_at"`
	LastActiveAt    int64                        `json:"last_active_at"`
	Username        *string                      `json:"username"`
	FirstName       *string                      `json:"first_name"`
	LastName        *string                      `json:"last_name"`
	OrgData         []UserOrgMembershipForReport `json:"org_data"`
	ExtraProperties map[string]any               `json:"extra_properties"`
}

type UserReport struct {
	UserReports    []UserReportRecord `json:"user_reports"`
	CurrentPage    int                `json:"current_page"`
	TotalCount     int                `json:"total_count"`
	PageSize       int                `json:"page_size"`
	HasMoreResults bool               `json:"has_more_results"`
	ReportTime     int64              `json:"report_time"`
}
