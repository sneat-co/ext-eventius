import { describe, expect, expectTypeOf, it } from 'vitest';

import {
  COMPETIOS_ATTENDANCE_COMMAND_CONFLICT_CODE,
  COMPETIOS_ATTENDANCE_ID_MAX_UTF8_BYTES,
  COMPETIOS_ATTENDANCE_REASON_MAX_UTF8_BYTES,
  type ICompetiosAttendanceCommandBinding,
  type ICancelCompetiosAttendanceEventCommand,
  type ICompetiosAttendanceStatus,
  type IEnsureCompetiosAttendanceInviteeInvitationRequest,
  type IGetCompetiosAttendanceInviteeStatusRequest,
  type IRevokeCompetiosAttendanceInvitationCommand,
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

  it('publishes bounded ID/reason and global command-conflict vocabulary', () => {
    const status: ICompetiosAttendanceStatus = {
      competiosEventKey: 'event-1',
      attendanceEventID: 'eventius-event-1',
      eventState: 'active',
    };
    const binding: ICompetiosAttendanceCommandBinding = {
      servicePrincipalID: 'competios-production',
      requestID: 'request-1',
      operation: 'cancel_attendance_event',
      payloadFingerprint: '0'.repeat(64),
      projection: status,
    };

    expect(COMPETIOS_ATTENDANCE_ID_MAX_UTF8_BYTES).toBe(128);
    expect(COMPETIOS_ATTENDANCE_REASON_MAX_UTF8_BYTES).toBe(512);
    expect(COMPETIOS_ATTENDANCE_COMMAND_CONFLICT_CODE).toBe('command_conflict');
    expect(binding.payloadFingerprint).toHaveLength(64);
  });
});
