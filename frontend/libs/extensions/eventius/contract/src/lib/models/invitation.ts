// TypeScript mirror of the frozen Eventius invitations contract
// (`typespec/api4eventius-invitations.tsp`). Shapes MUST match the JSON the
// backend emits/accepts exactly. An invitation references the invitee by id
// only — Sneat (`spaceus`/`contactus`) stays the source of truth for names.

/** Whether an invitation targets a whole family (space) or a single person. */
export type InviteeType = 'family' | 'person';

/** An invitation of a family or person to an event. */
export interface IInvitation {
  /** Invitation id, unique within its event. */
  readonly id: string;

  /** Discriminates which reference field is populated. */
  readonly inviteeType: InviteeType;

  /** Set when `inviteeType === 'family'`: the invited family's space id. */
  readonly familySpaceID?: string;

  /** Set when `inviteeType === 'person'`: the invited contact's id. */
  readonly contactID?: string;

  /** Audit: when the invitation was added (RFC3339). */
  readonly createdAt: string;
}

/**
 * POST body to add one invitee. Exactly one of `familySpaceID` / `contactID`
 * is set, matching `inviteeType`. No name or contact detail is sent — only id.
 */
export interface IAddInviteeRequest {
  inviteeType: InviteeType;
  familySpaceID?: string;
  contactID?: string;
}
