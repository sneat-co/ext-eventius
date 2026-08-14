package botcontract

import (
	"github.com/sneat-co/ext-eventius/backend/facade4eventius"
	"github.com/sneat-co/ext-eventius/backend/participation"
)

// The bot contract is a narrow, persistence-free import boundary. These aliases
// deliberately reuse the canonical Eventius facade and participation vocabulary
// rather than defining a second set of conversational event models.
type (
	EventService                             = facade4eventius.Service
	EventView                                = facade4eventius.Event
	CreateEventRequest                       = facade4eventius.CreateEventRequest
	UpdateEventRequest                       = facade4eventius.UpdateEventRequest
	CreateInviteRequest                      = facade4eventius.CreateInviteRequest
	RespondRequest                           = facade4eventius.RespondRequest
	Invite                                   = facade4eventius.Invite
	Response                                 = facade4eventius.Response
	InvitationContext                        = facade4eventius.InvitationContext
	CompetiosAttendanceService               = facade4eventius.CompetiosAttendanceService
	CompetiosAttendanceInviteeStatusService  = facade4eventius.CompetiosAttendanceInviteeStatusService
	EnsureAttendanceEventRequest             = facade4eventius.EnsureAttendanceEventRequest
	EnsureAttendanceInvitationRequest        = facade4eventius.EnsureAttendanceInvitationRequest
	EnsureAttendanceInviteeInvitationRequest = facade4eventius.EnsureAttendanceInviteeInvitationRequest
	GetAttendanceInviteeStatusRequest        = facade4eventius.GetAttendanceInviteeStatusRequest
	AttendanceStatusProjection               = facade4eventius.AttendanceStatusProjection
	CompetiosInviteeKey                      = facade4eventius.CompetiosInviteeKey
	CompetiosEntryLifecycleRevision          = facade4eventius.CompetiosEntryLifecycleRevision
	Participation                            = participation.Coarse
)

const (
	ParticipationYes   Participation = participation.CoarseYes
	ParticipationMaybe Participation = participation.CoarseMaybe
	ParticipationNo    Participation = participation.CoarseNo
)
