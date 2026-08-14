# ext-eventius

Public extension-definition repo for the **Eventius occasion family** — ToGethered,
RSVP.express verticals, and the Eventius/Invitus event stack — per the
[`ext-<id>` convention](https://github.com/sneat-co/sneat-specs/blob/main/standards/repo-naming.md)
and the [extension backend architecture standard](https://github.com/sneat-co/sneat-specs/blob/main/standards/extension-backend-architecture.md).

## backend/

Go module `github.com/sneat-co/ext-eventius/backend`.

**Dependency policy:** stdlib only. No dalgo, no sneat-go-core, no other extension.
This keeps the module publicly buildable in any CI without GOPRIVATE/PAT and avoids
the version treadmill the architecture standard exists to escape.

### Package `participation`

The **single vocabulary for participation-response scales** across the occasion family.
It satisfies the *"types needed by several extensions → ext-\<id\>/backend contract
module"* row of the
[extension-backend-architecture decision ladder](https://github.com/sneat-co/sneat-specs/blob/main/standards/extension-backend-architecture.md#where-shared-things-go--the-decision-ladder).

Before this module the occasion family had four diverging response scales maintained
in isolation:

| Previous location | Scale |
|---|---|
| `togethered/backend` IntentDbo.Likelihood | `no · unlikely · maybe · likely · yes` |
| `eventius/backend` RsvpStatus | `yes · no · maybe` |
| `gameboard/backend` RsvpStatus (`game_invite.go:82-88`) | `going · maybe · out` |
| `rsvp-express/practice-training` spec | `available · unavailable · maybe` |

This package is the single source of truth; those implementations should import here.

#### Types

| Type | Description |
|---|---|
| `Level` | Five-point graded scale: `no < unlikely < maybe < likely < yes`. Stored as-is in Firestore by togethered/backend. |
| `Coarse` | Three-point quick-answer scale: `yes \| no \| maybe`. Matches eventius/backend RsvpStatus and personal-events RSVP. |
| `Availability` | Alias rendering of `Coarse` for practice-training (`available=yes, unavailable=no, maybe=maybe`). Not a third scale. |
| `GameResponse` | Alias rendering of `Coarse` for gameboard/backend (`going=yes, out=no, maybe=maybe`). Not a third scale. |

#### Canonical mappings

- `Level.Coarse()` — lossy: `yes,likely→yes`; `maybe→maybe`; `unlikely,no→no`.
- `Coarse.Level()` — lossless embedding: `yes→yes`; `maybe→maybe`; `no→no`.
- `Availability.Coarse()` / `AvailabilityFromCoarse()` — alias round-trip.
- `GameResponse.Coarse()` / `GameResponseFromCoarse()` — alias round-trip.

#### Parsing

`ParseLevel(string)` and `ParseCoarse(string)` are strict wire-value parsers.
They do **not** interpret natural-language phrases ("probably" → likely);
phrase mapping is a conversational-layer concern that stays in sneat-bots.

#### NON-GOALS

This package deliberately omits:

- **Notification policy** (who gets notified at which level) — extension policy
  owned by togethered/backend; this module must never grow predicates such as
  `NotifiesFollowers()`.
- **Co-presence eligibility** — likewise extension policy; the "co-presence at
  ≥ likely" rule belongs in togethered, not in vocabulary.
- **NL phrase parsing** — see Parsing above.

#### Current and prospective consumers

| Consumer | How it uses this module |
|---|---|
| `togethered/backend` | Stores `IntentDbo.Likelihood` as `Level`; uses `Level.Coarse()` for the coarse quick-answers and `Level.AtLeast(LevelLikely)` for co-presence edge evaluation. |
| `rsvp-express` personal-events | Stores RSVP status as `Coarse`. |
| `rsvp-express` practice-training | Stores session availability as `Availability`; converts to `Coarse` for shared roster logic. |
| `eventius/backend` | Replaces local `RsvpStatus` with `Coarse`; existing stored values (`"yes"`, `"no"`, `"maybe"`) are identical. |
| `gameboard/backend` | Replaces local `RsvpStatus` (`going/maybe/out`) with `GameResponse`; backward-compatible read/write of existing records. |
| `invitus` mapping | Maps `accepted→CoarseYes`, `declined→CoarseNo`; the mapping code stays in sneat-core-modules, which imports this vocabulary. |

#### Scale changes

Adding, removing, or reordering values in `Level` or `Coarse` requires an explicit
owner decision (alexander.trakhimenok@gmail.com): every stored value and every
consumer is affected. `Availability` and `GameResponse` are derived alias renderings
and can gain new helpers without a scale-change review.

## frontend/

Nx workspace whose `ext-eventius-contract` project publishes
`@sneat/extension-eventius-contract`. Eventius implementation stays private in
`sneat-co/eventius`.

## Competios attendance facade

`backend/facade4eventius` exposes the provider-side safe attendance contract;
`typespec/api4eventius-competios-attendance.tsp` is its frozen wire source and
the published TypeScript mirror is in `frontend/`. The contract transfers no
RSVP token, URL, contact, invitee name, or payment data.

An attendance invitation is identified by the Event, Tournament, Competition,
Entry, registration and an opaque `CompetiosInviteeKey`. Competios may encode
the Entry lifecycle revision in that last value; Eventius compares it only as
an opaque string. `CompetiosAttendanceService.GetAttendanceStatus` remains for
legacy registration-only consumers and must fail closed when a registration is
not uniquely one invitee. New consumers opt into
`CompetiosAttendanceInviteeStatusService`: `GetAttendanceEventStatus` returns
only event status, while `GetAttendanceInviteeStatus` requires the complete
tuple and returns exactly one safe invitation projection.
