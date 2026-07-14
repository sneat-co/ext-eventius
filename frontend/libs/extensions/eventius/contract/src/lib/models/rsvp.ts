// TypeScript mirror of the frozen Eventius RSVP contract
// (`typespec/api4eventius-rsvp.tsp`) plus the public resolve context from
// (`typespec/api4eventius-links.tsp`). Shapes MUST match the JSON the backend
// emits/accepts exactly. Submitting an RSVP is pure event-attendance and NEVER
// joins the host's space.

import { LinkKind } from './rsvp-link';

/** Responder's attendance answer. */
export type RsvpStatus = 'yes' | 'no' | 'maybe';

/** A submitted RSVP. */
export interface IRsvp {
  /** RSVP id, unique within its event. */
  readonly id: string;

  /** Invitation this RSVP answers (per-invitee link); empty for an open link. */
  readonly invitationID?: string;

  /** Attribution for an open-link response: the responder's self-identified name. */
  readonly selfIdentifiedName?: string;

  /** Optional self-identified family name (open link). */
  readonly selfIdentifiedFamilyName?: string;

  readonly status: RsvpStatus;

  /** Adults attending. Forced to 0 by the backend when `status == 'no'`. */
  readonly adults: number;

  /** Children attending. Forced to 0 by the backend when `status == 'no'`. */
  readonly children: number;

  /** Optional dietary/allergy notes. */
  readonly dietary?: string;

  /** Optional free-text comment. */
  readonly comment?: string;

  /** True when submitted by a signed-in Sneat account. */
  readonly viaAccount: boolean;

  /** Audit: submission time (RFC3339). */
  readonly submittedAt: string;
}

/** POST body to submit an RSVP. Headcounts are zeroed by the backend when `status == 'no'`. */
export interface ISubmitRsvpRequest {
  status: RsvpStatus;
  adults: number;
  children: number;
  dietary?: string;
  comment?: string;

  /**
   * REQUIRED for an open-link token (the responder self-identifies); omitted for
   * a per-invitee token (attribution comes from the invitation). Never a
   * guest-list selection — only this free-text name. (AC: attribution-by-link-type)
   */
  selfIdentifiedName?: string;

  /** Optional self-identified family name (open link). */
  selfIdentifiedFamilyName?: string;
}

/**
 * Build the submit body from raw form values. When `status == 'no'` the
 * headcounts are forced to 0 (the backend zeroes them anyway), and blank
 * dietary/comment/self-identified names are dropped. For an open-link response
 * the trimmed `selfIdentifiedName` is carried through (the caller enforces it is
 * non-empty). Pure so it can be unit-tested without rendering.
 * (AC: no-rsvp-zeroes-headcount, rsvp-full-fields-captured, attribution-by-link-type)
 */
export function buildRsvpRequest(form: {
  status: RsvpStatus;
  adults: number;
  children: number;
  dietary: string;
  comment: string;
  selfIdentifiedName?: string;
  selfIdentifiedFamilyName?: string;
}): ISubmitRsvpRequest {
  const attending = form.status !== 'no';
  const dietary = form.dietary.trim();
  const comment = form.comment.trim();
  const name = form.selfIdentifiedName?.trim();
  const familyName = form.selfIdentifiedFamilyName?.trim();
  return {
    status: form.status,
    adults: attending ? form.adults : 0,
    children: attending ? form.children : 0,
    ...(dietary ? { dietary } : {}),
    ...(comment ? { comment } : {}),
    ...(name ? { selfIdentifiedName: name } : {}),
    ...(familyName ? { selfIdentifiedFamilyName: familyName } : {}),
  };
}

/**
 * Public context returned when a token is resolved — only what the RSVP form
 * needs. NEVER includes the guest list or any other invitee.
 */
export interface IRsvpContext {
  /** Display fields for the event. */
  readonly eventTitle: string;
  readonly eventStart: string;
  readonly eventLocation: string;

  /** Whether the event still accepts RSVPs (false once cancelled). */
  readonly open: boolean;

  /** `invitee` → pre-attributed to one invitee; `open` → responder self-identifies. */
  readonly kind: LinkKind;

  /** Set only when `kind == 'invitee'`: the invitation this link is bound to. */
  readonly invitationID?: string;
}
