// TypeScript mirror of the frozen Eventius wire contract (`typespec/api4eventius.tsp`).
// These shapes MUST match the JSON the backend emits/accepts exactly.

/** Lifecycle status of an event. */
export type EventStatus = 'active' | 'cancelled';

/**
 * An Eventius event. Persisted by the backend at
 * `/spaces/{spaceID}/ext/eventius/events/{eventID}`.
 */
export interface IEvent {
  /** Event id, unique within its host space. */
  readonly id: string;

  /** Host space id this event belongs to. */
  readonly spaceID: string;

  /** Display name, e.g. "Sophie's 10th Birthday". */
  readonly title: string;

  /** Event start date/time as an RFC3339 / ISO string. */
  readonly start: string;

  /** Free-text location. */
  readonly location: string;

  /** Optional free-text description. */
  readonly description?: string;

  /** Lifecycle status; `active` on creation. */
  readonly status: EventStatus;

  /** Audit: creation timestamp (RFC3339). */
  readonly createdAt: string;

  /** Audit: last-modification timestamp (RFC3339). */
  readonly updatedAt: string;
}

/** POST body to create an event. */
export interface ICreateEventRequest {
  /** Display name (required, non-empty). */
  title: string;

  /** Event start date/time (RFC3339). */
  start: string;

  /** Free-text location (required, non-empty). */
  location: string;

  /** Optional free-text description. */
  description?: string;

  /**
   * Event length in minutes. Optional — when omitted the backend's calendarius
   * facade applies its default (60m).
   */
  durationMinutes?: number;
}

/**
 * Response of the create-event endpoint. The event itself is read
 * Firestore-direct (happening + overlay); the create response only identifies
 * the new happening so the UI can navigate to it.
 */
export interface ICreateEventResponse {
  /** New event/happening id. */
  readonly id: string;
}

/**
 * PUT body to edit an event. All fields optional; only present fields are
 * updated. Does not change `spaceID` or `status` (use cancel for status).
 */
export interface IUpdateEventRequest {
  title?: string;
  start?: string;
  location?: string;
  description?: string;
}

/** Standard error body returned by the backend with a non-2xx status. */
export interface IApiError {
  /** Machine-readable error code, e.g. `unauthorized`, `not_found`. */
  code: string;

  /** Human-readable message. */
  message: string;
}
