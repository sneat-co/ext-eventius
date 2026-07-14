import { InjectionToken } from '@angular/core';
import { Observable, of } from 'rxjs';

/** A family the host can invite — id + display title only (no contact detail). */
export interface IInviteeFamily {
  readonly id: string;
  readonly title: string;
}

/** A person the host can invite — id + display name only (no contact detail). */
export interface IInviteePerson {
  readonly id: string;
  readonly name: string;
}

/**
 * Supplies the families and persons the host can pick from. Injecting this
 * (rather than calling `spaceus`/`contactus` directly) keeps the picker
 * testable and decoupled from those modules.
 */
export interface InviteeSource {
  /** Families (Sneat `family` spaces) available to invite. */
  families(spaceID: string): Observable<IInviteeFamily[]>;

  /** Persons (Sneat `contactus` contacts) available to invite. */
  persons(spaceID: string): Observable<IInviteePerson[]>;
}

export const EVENTIUS_INVITEE_SOURCE = new InjectionToken<InviteeSource>(
  'EVENTIUS_INVITEE_SOURCE',
  { providedIn: 'root', factory: () => new StubInviteeSource() },
);

/**
 * Default stub returning sample data so the picker is usable end-to-end before
 * real wiring. TODO: replace with an adapter over `spaceus` (family spaces) and
 * `contactus` (contacts) once those modules are available to the app shell.
 */
export class StubInviteeSource implements InviteeSource {
  families(): Observable<IInviteeFamily[]> {
    return of([
      { id: 'fam-sample-1', title: 'The Smith Family' },
      { id: 'fam-sample-2', title: 'The Brown Family' },
    ]);
  }

  persons(): Observable<IInviteePerson[]> {
    return of([
      { id: 'person-sample-1', name: 'Alice Smith' },
      { id: 'person-sample-2', name: 'Bob Brown' },
    ]);
  }
}
