// TypeScript mirror of the frozen Eventius links contract
// (`typespec/api4eventius-links.tsp`). Shapes MUST match the JSON the backend
// emits exactly. The QR code is rendered client-side from `url`.

/** `invitee` = bound to a single invitation; `open` = shared event link. */
export type LinkKind = 'invitee' | 'open';

/** A tokenized RSVP link (+ url for QR). */
export interface IRsvpLink {
  /** Opaque link token. */
  readonly token: string;

  /** Full shareable RSVP URL, e.g. `https://eventius.sneat.app/r/{token}`. */
  readonly url: string;

  /** Whether the link targets one invitee or is the open event link. */
  readonly kind: LinkKind;
}
