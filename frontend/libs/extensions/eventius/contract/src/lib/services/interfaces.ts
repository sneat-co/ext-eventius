import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  IBringAlongSlot,
  IClaimSlotRequest,
  IDefineSlotRequest,
  IReleaseSlotRequest,
} from '../models/bring-along';
import {
  ICreateEventResponse,
  ICreateEventRequest,
  IEvent,
  IUpdateEventRequest,
} from '../models/event';
import { IEventiusEventListItem } from '../models/eventius-event';
import { IAddInviteeRequest, IInvitation } from '../models/invitation';
import { IRsvp, IRsvpContext, ISubmitRsvpRequest } from '../models/rsvp';
import { IRsvpLink } from '../models/rsvp-link';

// Runtime-light service contracts the eventius pages/space-menu depend on. Each
// interface mirrors the public surface of the concrete service in the internal
// lib; the implementation is bound to the matching token below via
// provideEventiusInternal(). Contract must never import from internal.

/** Bespoke HTTP client for the Eventius event API (read/edit). */
export interface IEventService {
  listEvents(spaceID: string): Observable<IEvent[]>;
  getEvent(spaceID: string, eventID: string): Observable<IEvent>;
  updateEvent(
    spaceID: string,
    eventID: string,
    request: IUpdateEventRequest,
  ): Observable<IEvent>;
  cancelEvent(spaceID: string, eventID: string): Observable<IEvent>;
}

export const EVENT_SERVICE = new InjectionToken<IEventService>('EventService');

/** Firestore-direct access to a space's eventius events + event creation. */
export interface IEventiusEventService {
  watchEvents(spaceID: string): Observable<IEventiusEventListItem[]>;
  createEvent(
    spaceID: string,
    request: ICreateEventRequest,
  ): Observable<ICreateEventResponse>;
}

export const EVENTIUS_EVENT_SERVICE = new InjectionToken<IEventiusEventService>(
  'EventiusEventService',
);

/** Token-gated, public RSVP API client. */
export interface IRsvpService {
  resolveToken(token: string): Observable<IRsvpContext>;
  submitRsvp(token: string, request: ISubmitRsvpRequest): Observable<IRsvp>;
  updateRsvp(
    token: string,
    rsvpID: string,
    request: ISubmitRsvpRequest,
  ): Observable<IRsvp>;
  listRsvps(spaceID: string, eventID: string): Observable<IRsvp[]>;
}

export const RSVP_SERVICE = new InjectionToken<IRsvpService>('RsvpService');

/** Eventius invitations API client. */
export interface IInvitationService {
  addInvitee(
    spaceID: string,
    eventID: string,
    request: IAddInviteeRequest,
  ): Observable<IInvitation>;
  listInvitations(
    spaceID: string,
    eventID: string,
  ): Observable<IInvitation[]>;
}

export const INVITATION_SERVICE = new InjectionToken<IInvitationService>(
  'InvitationService',
);

/** Eventius links API client (host-only). */
export interface ILinkService {
  issueInviteeLink(
    spaceID: string,
    eventID: string,
    invitationID: string,
  ): Observable<IRsvpLink>;
  issueOpenLink(spaceID: string, eventID: string): Observable<IRsvpLink>;
}

export const LINK_SERVICE = new InjectionToken<ILinkService>('LinkService');

/** Eventius bring-along (slots) API client. */
export interface IBringAlongService {
  defineSlot(
    spaceID: string,
    eventID: string,
    request: IDefineSlotRequest,
  ): Observable<IBringAlongSlot>;
  listSlots(spaceID: string, eventID: string): Observable<IBringAlongSlot[]>;
  claimSlot(
    spaceID: string,
    eventID: string,
    slotID: string,
    request: IClaimSlotRequest,
  ): Observable<IBringAlongSlot>;
  releaseSlot(
    spaceID: string,
    eventID: string,
    slotID: string,
    request: IReleaseSlotRequest,
  ): Observable<IBringAlongSlot>;
}

export const BRING_ALONG_SERVICE = new InjectionToken<IBringAlongService>(
  'BringAlongService',
);
