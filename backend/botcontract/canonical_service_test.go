package botcontract

import (
	"testing"

	"github.com/sneat-co/ext-eventius/backend/facade4eventius"
)

func TestCanonicalAliasesPreserveTheFacadeVocabulary(t *testing.T) {
	request := CreateEventRequest{RequestID: "request", SpaceID: "space", Title: "Night"}
	if facade4eventius.CreateEventRequest(request).Title != "Night" {
		t.Fatal("bot contract must retain the canonical create-event request")
	}
	if ParticipationYes != "yes" || ParticipationMaybe != "maybe" || ParticipationNo != "no" {
		t.Fatal("bot participation values must reuse canonical RSVP values")
	}
	inviteeKey := CompetiosInviteeKey("competios:invitee@entry-revision")
	if facade4eventius.CompetiosInviteeKey(inviteeKey) != "competios:invitee@entry-revision" {
		t.Fatal("bot contract must retain the canonical Competios invitee key")
	}
}
