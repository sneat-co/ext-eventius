import { IIdAndBriefAndDbo } from '@sneat/core';
import {
  IHappeningBrief,
  IHappeningDbo,
  ISingleHappeningDbo,
} from '@sneat/extension-calendarius-contract';

// An eventius event IS a Calendarius happening of kind `'event'`. There is no
// separate overlay document: the happening carries the `kind` marker (so events
// can be queried with `kind == 'event'`) and any eventius-only extras live under
// `happening.ext.eventius`. Display fields (title, start, location, status) are
// the happening's own fields.

/** Eventius-only extras embedded under `happening.ext.eventius`. */
export interface IEventiusEventBrief {
  /** listus list id holding the bring-along / give-away items, when attached. */
  readonly bringAlongListID?: string;
}

type Happening = IIdAndBriefAndDbo<IHappeningBrief, IHappeningDbo>;

/** The Calendarius HappeningKind value that marks a happening as an eventius event. */
export const EVENT_HAPPENING_KIND = 'event';

/** A row for the eventius events list: a happening that is an eventius event. */
export interface IEventiusEventListItem {
  /** happeningID. */
  readonly id: string;
  readonly title: string;
  /** Display start: ISO from the single happening, or `date[ time]` from its slot. */
  readonly start?: string;
  readonly happening: Happening;
}

/**
 * Pure projection: the events list is the set of happenings whose kind is
 * `'event'`, projected to list rows. Kept pure so it is unit-testable without
 * Firestore. The service already narrows the query to `kind == 'event'`; this
 * filters again defensively so the projection is correct for any input.
 */
export function selectEventListItems(
  happenings: readonly Happening[],
): IEventiusEventListItem[] {
  return happenings
    .filter((h) => (h.dbo?.kind ?? h.brief?.kind) === EVENT_HAPPENING_KIND)
    .map((h) => ({
      id: h.id,
      title: h.dbo?.title ?? h.brief?.title ?? '',
      start: eventStart(h.dbo),
      happening: h,
    }));
}

function eventStart(dbo?: IHappeningDbo | null): string | undefined {
  if (!dbo) {
    return undefined;
  }
  const dtStarts = (dbo as ISingleHappeningDbo).dtStarts;
  if (dtStarts) {
    return new Date(dtStarts).toISOString();
  }
  const slot = Object.values(dbo.slots ?? {})[0];
  if (slot?.start?.date) {
    return slot.start.time
      ? `${slot.start.date} ${slot.start.time}`
      : slot.start.date;
  }
  return undefined;
}
