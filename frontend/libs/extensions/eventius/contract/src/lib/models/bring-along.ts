// TypeScript mirror of the frozen Eventius bring-along contract
// (`typespec/api4eventius-bringalong.tsp`). Shapes MUST match the JSON the
// backend emits/accepts exactly. The host defines slots; responders claim an
// open slot and may release it back to open, preventing duplicate contributions.

/**
 * Who claimed a slot. Stored as an opaque reference plus a display label so the
 * slot can render "Drinks — Smith Family" without resolving the reference.
 */
export interface ISlotClaimant {
  /** Opaque claimant reference (e.g. a family space id, contact id, or guest token). */
  readonly ref: string;

  /** Display label shown next to the slot, e.g. "Smith Family". */
  readonly label: string;
}

/** A bring-along item slot on an event. `claimedBy` absent => the slot is open. */
export interface IBringAlongSlot {
  /** Slot id, unique within its event. */
  readonly id: string;

  /** What's needed, e.g. "Drinks". */
  readonly label: string;

  /** Optional free-text quantity/notes, e.g. "2 bottles". */
  readonly note?: string;

  /** Present when claimed; absent means the slot is open. */
  readonly claimedBy?: ISlotClaimant;

  /** Audit: when the slot was defined (RFC3339). */
  readonly createdAt: string;
}

/** POST body to define a slot. Host-only. */
export interface IDefineSlotRequest {
  label: string;
  note?: string;
}

/** POST body to claim a slot. */
export interface IClaimSlotRequest {
  ref: string;
  label: string;
}

/** POST body to release a slot — only the claiming `ref` may release. */
export interface IReleaseSlotRequest {
  ref: string;
}
