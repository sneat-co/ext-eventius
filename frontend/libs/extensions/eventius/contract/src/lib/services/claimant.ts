import { InjectionToken } from '@angular/core';
import { ISlotClaimant } from '../models/bring-along';

/**
 * The identity used when the current user claims or releases a bring-along
 * slot. Injecting this (rather than hard-coding) keeps claim/release testable
 * and decoupled from RSVP identity.
 *
 * TODO: replace the stub default with the real responder identity once RSVP
 * lands (Tasks 3/4) — claim/release should use the link-gated guest token or
 * the signed-in family/person ref instead of this placeholder.
 */
export const EVENTIUS_CURRENT_CLAIMANT = new InjectionToken<ISlotClaimant>(
  'EVENTIUS_CURRENT_CLAIMANT',
  {
    providedIn: 'root',
    factory: () => ({ ref: 'stub-claimant', label: 'Smith Family' }),
  },
);
