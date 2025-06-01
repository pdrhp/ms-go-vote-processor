package entity

import "fmt"

type VoteStatus string

const (
    VoteStatusReceived VoteStatus = "RECEIVED"
    VoteStatusSent     VoteStatus = "SENT"

    VoteStatusProcessing VoteStatus = "PROCESSING"
    VoteStatusProcessed  VoteStatus = "PROCESSED"
    VoteStatusFailed     VoteStatus = "FAILED"
)

func (s VoteStatus) IsValid() bool {
    switch s {
    case VoteStatusReceived, VoteStatusSent, VoteStatusProcessing, VoteStatusProcessed, VoteStatusFailed:
        return true
    }
    return false
}

func (s VoteStatus) Validate() error {
    if !s.IsValid() {
        return fmt.Errorf("status inválido: %s", s)
    }
    return nil
}

func VoteStatusFromString(s string) (VoteStatus, error) {
    status := VoteStatus(s)
    if err := status.Validate(); err != nil {
        return "", err
    }
    return status, nil
}

func MustVoteStatusFromString(s string) VoteStatus {
    status, err := VoteStatusFromString(s)
    if err != nil {
        panic(err)
    }
    return status
}

func (s VoteStatus) CanTransitionTo(newStatus VoteStatus) bool {
    switch s {
    case VoteStatusReceived:
        return newStatus == VoteStatusProcessing
    case VoteStatusSent:
        return newStatus == VoteStatusProcessing
    case VoteStatusProcessing:
        return newStatus == VoteStatusProcessed || newStatus == VoteStatusFailed
    case VoteStatusProcessed:
        return false
    case VoteStatusFailed:
        return newStatus == VoteStatusProcessing
    }
    return false
}