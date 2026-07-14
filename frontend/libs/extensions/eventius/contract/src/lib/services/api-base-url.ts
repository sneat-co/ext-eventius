import { InjectionToken } from '@angular/core';

/**
 * Base URL the Eventius API is served from (without a trailing slash), e.g.
 * `https://api.sneat.app`. Defaults to '' so requests are relative to the
 * app origin (`/api4eventius/...`). The app shell can override it.
 */
export const EVENTIUS_API_BASE_URL = new InjectionToken<string>(
  'EVENTIUS_API_BASE_URL',
  { providedIn: 'root', factory: () => '' },
);
