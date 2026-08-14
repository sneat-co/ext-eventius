import { describe, expect, expectTypeOf, it } from 'vitest';

import type {
  ICancelCompetiosAttendanceEventCommand,
  ICompetiosAttendanceStatus,
  IEnsureCompetiosAttendanceInviteeInvitationRequest,
  IGetCompetiosAttendanceInviteeStatusRequest,
  IRevokeCompetiosAttendanceInvitationCommand,
} from './competios-attendance';

describe('Competios attendance invitee contract', () => {
  it('requires the complete opaque invitee lifecycle tuple on an exact ensure command', () => {
    const request: IEnsureCompetiosAttendanceInviteeInvitationRequest = {
      requestID: 'request-1',
      attendanceEventID: 'eventius-event-1',
      competiosEventKey: 'event-1',
      competiosRegistrationKey: 'registration-1',
      competiosTournamentKey: 'tournament-1',
      competiosCompetitionKey: 'competition-1',
      competiosEntryKey: 'entry-1',
      competiosInviteeKey: 'invitee-1',
      competiosEntryLifecycleRevision: 'entry-revision-2',
      responder: { kind: 'account', accountID: 'account-1' },
    };

    expect(Object.keys(request)).toEqual([
      'requestID',
      'attendanceEventID',
      'competiosEventKey',
      'competiosRegistrationKey',
      'competiosTournamentKey',
      'competiosCompetitionKey',
      'competiosEntryKey',
      'competiosInviteeKey',
      'competiosEntryLifecycleRevision',
      'responder',
    ]);
  });

  it('keeps the exact lookup and safe projection tuple aligned', () => {
    const lookup: IGetCompetiosAttendanceInviteeStatusRequest = {
      competiosEventKey: 'event-1',
      competiosTournamentKey: 'tournament-1',
      competiosCompetitionKey: 'competition-1',
      competiosEntryKey: 'entry-1',
      competiosRegistrationKey: 'registration-1',
      competiosInviteeKey: 'invitee-1@entry-revision-2',
      competiosEntryLifecycleRevision: 'entry-revision-2',
    };
    const status: ICompetiosAttendanceStatus = {
      ...lookup,
      attendanceEventID: 'eventius-event-1',
      attendanceInvitationID: 'eventius-invitation-1',
      eventState: 'active',
      invitationState: 'active',
    };

    expect(status.competiosInviteeKey).toBe(lookup.competiosInviteeKey);
    expect(status.competiosEntryLifecycleRevision).toBe(
      lookup.competiosEntryLifecycleRevision,
    );
    expectTypeOf<
      Extract<
        keyof ICompetiosAttendanceStatus,
        'token' | 'rsvpToken' | 'contact' | 'contactID' | 'payment' | 'paymentID'
      >
    >().toEqualTypeOf<never>();
  });

  it('mirrors exact revoke and cancel commands with request IDs and audit reasons', () => {
    const revoke: IRevokeCompetiosAttendanceInvitationCommand = {
      requestID: 'revoke-1',
      attendanceEventID: 'eventius-event-1',
      attendanceInvitationID: 'eventius-invitation-1',
      competiosEventKey: 'event-1',
      competiosTournamentKey: 'tournament-1',
      competiosCompetitionKey: 'competition-1',
      competiosEntryKey: 'entry-1',
      competiosRegistrationKey: 'registration-1',
      competiosInviteeKey: 'invitee-1',
      competiosEntryLifecycleRevision: 'entry-revision-2',
      reason: 'entry withdrawn',
    };
    const cancel: ICancelCompetiosAttendanceEventCommand = {
      requestID: 'cancel-1',
      attendanceEventID: revoke.attendanceEventID,
      competiosEventKey: revoke.competiosEventKey,
      reason: 'competition event cancelled',
    };

    expect(revoke.requestID).toBeTruthy();
    expect(revoke.reason).toBeTruthy();
    expect(cancel.requestID).toBeTruthy();
    expect(cancel.reason).toBeTruthy();
  });
});
