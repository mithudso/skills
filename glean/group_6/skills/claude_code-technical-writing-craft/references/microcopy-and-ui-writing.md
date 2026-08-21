<!-- hub-reference-banner -->
> **Reference file — part of the `technical-writing-craft` hub.** Formerly the standalone `microcopy-and-ui-writing` skill.
> Sibling topics in this family are now reference files under the hubs (`writing-expert`, `technical-writing-craft`, `executive-comms`, `content-and-marketing-writing`, `career-and-formal-writing`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: microcopy-and-ui-writing
description: In-product UX writing for buttons, empty states, tooltips, form-field labels, validation messages, system messages, password-reset flows, onboarding tours, confirmation dialogs (destructive vs benign), modals, toasts, opt-in/opt-out copy, and 404/500 pages. Codifies the "verb the noun" button rule, the "what you can do here" empty-state rule, the "be specific" error rule, the "match the destructive verb" confirmation rule, sentence case vs title case decisions, NN/g's plain-language error guidelines, Mailchimp/Polaris/HIG/Material/GOV.UK conventions, Kinneret Yifrah's voice-and-tone framework, and dark-pattern avoidance (symmetry in choice, no confirmshaming, no manipulinks). TRIGGER when the user is writing or reviewing button text, an empty state, an error/validation message, a tooltip, a form label, a toast, a confirmation dialog, a 404 page, a password-reset flow, onboarding copy, modal text, an opt-in/cookie dialog, or any in-product string under ~140 characters. Also: "what should this button say", "rewrite this error message", "this empty state is bad", "the validation message is confusing", "make this toast clearer", "is this dialog destructive enough", "this CTA isn't converting". SKIP for marketing landing-page copy or conversion-optimization (use sales-and-marketing-copy), long-form prose like help docs or release notes (use writing-expert or technical-writing-craft), accessibility-driven plain-language rewrites for compliance (use plain-language), inclusive-language audits (use inclusive-language), customer-support ticket prose like apologies and holding statements (use support-ticket-writing), executive comms, and prose-craft fundamentals (use writing-expert).
---

# Microcopy and UI writing

Microcopy is every word in a product that is not body content: the buttons, the field labels, the placeholder hints, the validation messages, the empty states, the toasts, the modal titles, the 404 pages, the system status banners, the consent dialogs. Each individual piece is small. Collectively, microcopy is the voice of the product.

This skill encodes the rules from the major design systems — Mailchimp, Shopify Polaris, Apple HIG, Material Design, GOV.UK, Carbon — plus Nielsen Norman Group's error-message guidelines and Kinneret Yifrah's *Microcopy: The Complete Guide*. It is the in-product complement to support-ticket-writing (customer-facing prose) and writing-expert (long-form prose).

## When to use

- Writing or reviewing a button label, link label, or CTA.
- Drafting an empty state for a list, table, or grid.
- Writing or rewriting a form-field label, placeholder, helper text, or validation message.
- Drafting a tooltip.
- Writing a toast, banner, or system-message notification.
- Writing a confirmation dialog, especially for a destructive action.
- Writing modal titles and supporting text.
- Drafting password-reset, sign-in, or sign-up flow text.
- Writing 404, 500, or other error-page copy.
- Drafting onboarding-tour steps.
- Writing opt-in/opt-out copy, including cookie banners and consent dialogs.
- Auditing in-product strings for tone consistency.

## Skip when

- The task is marketing or landing-page copy aimed at acquisition or conversion: use sales-and-marketing-copy.
- The artifact is long-form (help center, release notes, blog): use writing-expert or technical-writing-craft.
- The audience is a customer in a support ticket: use support-ticket-writing.
- The output needs a plain-language pass for accessibility, legal, or readability-level compliance: use plain-language.
- The audit is for inclusive language: use inclusive-language.
- The work is exec or QBR communication: use executive-comms.

## Core concepts

### 1. The "verb the noun" rule for buttons

Buttons should describe the action they perform. The canonical shape is **verb + noun** (Apple HIG, GOV.UK, Polaris).

| Weak | Stronger |
|---|---|
| OK | Save changes |
| Submit | Create account |
| Yes | Delete project |
| Done | Save and close |
| Continue | Continue to payment |
| Click here | Download report |
| More | View all orders |

Three sub-rules:
- **Match the user's mental model.** If they came to delete, the button says "Delete." If they came to publish, the button says "Publish." Not "Confirm" or "OK."
- **Echo the action's verb in destructive confirmations.** If the dialog asks "Delete this project?", the destructive button says "Delete project," not "Yes" or "Confirm." This is the single highest-impact rule in destructive-action UX.
- **Cancel is the only one-word button.** "Cancel" is universal and short and always the dismissive option. Don't get cute ("Never mind," "Take me back") unless the surrounding voice is highly conversational.

Apple HIG: "Just saying 'Next' often works better than 'Let's do this!'" Clarity beats cleverness.

### 2. The "what you can do here" rule for empty states

An empty state is not a void. It is the first frame of a tutorial. The Mailchimp/Yifrah pattern:

1. **Name the absence** (one short line).
2. **Tell the user what this surface is for.**
3. **Give them the next action** (with a button when possible).

Weak:
> No items.

Stronger:
> **No saved designs yet**
> This is where your saved email designs will live. You can save any draft from the editor.
> [Create your first design]

The empty state for an unconfigured integration:
> **No integrations connected**
> Connect Slack, Jira, or GitHub to start receiving alerts in the channels your team already uses.
> [Browse integrations]

The empty state for a search with no results (notice this one is different — it is a *transient* empty, not a *first-time* empty):
> **No results for "atlasx"**
> Try a different spelling, fewer keywords, or browse all projects.

Three failure modes:
- **Apologizing for the emptiness** ("Oh no, nothing here!"). Be matter-of-fact. The user knows it's empty.
- **Selling instead of guiding.** Don't pitch a premium tier in an empty state. The user is trying to use the feature.
- **Linking out to docs without an in-product action.** A "Read the guide" link is fine *next to* a "Create your first X" button, not in place of one.

### 3. The "be specific" rule for error and validation messages

Nielsen Norman Group's canonical four-part guideline (Jakob Nielsen 2001, updated 2023):

1. **State what went wrong** in plain language.
2. **Explain why**, if useful.
3. **Tell the user how to fix it.**
4. **Don't blame the user.**

| Weak | Stronger |
|---|---|
| Invalid input. | Email must include an @ and a domain. |
| Error. | We couldn't reach the server. Check your connection and try again. |
| Password not accepted. | Password must be at least 12 characters and include a number. |
| Field required. | Please add the project name to continue. |
| Form has errors. | Two fields need attention — see them highlighted below. |
| 404 — page not found. | We couldn't find that page. It may have moved. Try the search above or [return home]. |

Polaris's two-sentence cap holds: error messages should be no longer than two sentences and never apologize hyperbolically. Mailchimp's rule: **never use exclamation points in failure messages.**

NN/g's "hostile patterns" list — what to *never* do:
- Don't say the user "failed" or "wasn't allowed."
- Don't show errors before the user has interacted (no pre-emptive red).
- Don't use ALL CAPS or red walls of text.
- Don't blame missing data the user wasn't asked for.
- Don't blow up generic system errors to the user without a "try again" path.

For real-time validation: validate *after* the user has finished a field (on blur, not on each keystroke). Re-validate after they fix it (success state, not silence). Mailchimp's password strength is the canonical positive example — graying out requirements as they're met provides reassurance, not pressure.

### 4. The match-the-verb rule for destructive confirmations

A destructive confirmation has three parts that must all align:

| Part | Pattern | Example |
|---|---|---|
| Title | "Verb this [thing]?" | "Delete project Acme?" |
| Body | What happens, what is irreversible, who is affected | "This will permanently delete all 47 issues, comments, and attachments. This cannot be undone." |
| Destructive button | The same verb, with the noun if clarifying | "Delete project" |
| Dismissive button | "Cancel" (default), never the primary visual style | "Cancel" |

Color and emphasis: the destructive button is styled as destructive (red on most systems). The dismissive button is neutral. There must be no ambiguity about which is which. Apple HIG: "If an alert button results in a destructive action, set the button's style to Destructive so that it gets appropriate formatting."

Avoid:
- "Are you sure?" as a title (vague, no information).
- "OK" / "Cancel" as button pair (no verb match).
- "Yes" / "No" as button pair (forces re-reading the title).
- Reversed visual hierarchy where the destructive button looks like the safe choice.
- Multi-step "type the name to confirm" UNLESS the action affects shared/billable resources (deleting a production database, removing all org members) — for routine deletes, this is friction theatre.

Benign vs destructive — keep them distinct. A benign confirmation is "Save changes?" with [Save changes] / [Discard]. A destructive confirmation is "Delete project?" with [Delete project] / [Cancel]. The verbs are different, the styles are different, the cognitive load is different.

### 5. Tooltips, helper text, placeholders — distinct uses, do not conflate

Three pieces of microcopy are constantly confused. They are *not* interchangeable.

| Pattern | Used when | Lives where | Example |
|---|---|---|---|
| **Label** | Always. Every input has one. | Above or beside the field. Never disappears. | "Project name" |
| **Helper text** | Persistent guidance the user needs *while* filling the field. | Below the field. Always visible. | "Letters, numbers, and dashes. 3-40 characters." |
| **Placeholder** | A formatting example. *Never the only label.* | Inside the empty field. Disappears on focus or input. | `e.g., acme-prod-cluster` |
| **Tooltip** | Edge-case info not most users need. Optional. | On hover or info-icon click. | "Project name is used in URLs and CLI invocations. It cannot be changed later." |

Anti-patterns:
- **Placeholder as label.** When the user starts typing, the label vanishes. Catastrophic for accessibility and for users who pause mid-form.
- **Tooltip as primary instruction.** If the user *needs* the info to fill the field, it belongs in helper text, not behind a hover.
- **Helper text repeating the label.** Label: "Email." Helper: "Please enter your email." Cut the helper.
- **Tooltips on mobile.** There is no hover on touch. If you ship a tooltip, ship a tap target too.

### 6. Toasts, banners, modals — match the severity to the surface

Severity dictates the surface, not the other way around.

| Surface | When | Lifespan | Examples |
|---|---|---|---|
| **Inline message** (next to the affected element) | Field-level validation, item-level status | Persistent until resolved | "This SKU is out of stock" |
| **Toast** (transient notification) | Successful actions, low-stakes info | 4-6 seconds, auto-dismiss | "Saved" "Copied" "Invitation sent" |
| **Banner** (persistent strip at top of page) | Account-wide state, time-bound info | Until user dismisses or state resolves | "Your trial ends in 3 days" |
| **Modal dialog** | Action requires confirmation, or blocks workflow until resolved | Until user resolves | Destructive confirmations, sign-in prompts |
| **System-wide alert** (full-width, top-of-app) | Service-level disruption | Until resolved | "Atlas is experiencing elevated latency. [Status page]" |

Toast copy is the shortest writing in the product. Single word or short verb-noun phrase. Past-tense for completed actions:
- "Saved" not "Your changes were saved successfully"
- "Copied" not "Link copied to clipboard"
- "Sent" not "Your message has been sent"

If a toast needs an undo or an action, add it as a button on the toast itself:
> **Project deleted.** [Undo]

That's it. Five words plus a button.

### 7. Sign-up, sign-in, password-reset — the conversion-sensitive flow

These flows have the highest abandonment rates of any in-product writing. Three rules:

**Sign-up:**
- Lead with what the user gets, not what they have to do. "Create your free account" beats "Sign up."
- Address the silent anxiety. Shopify's "You can change your store name afterwards" is the canonical example — seven words that defuse "what if I pick wrong."
- Show password requirements *before* the user types, not after they fail.
- Use a single Continue button per step. Never "Submit."

**Sign-in:**
- Label clearly: "Email" and "Password," not "Username." (Most users don't remember which.)
- The forgot-password link goes near the password field, not at the bottom of the form.
- Failed sign-in: "Email or password doesn't match" — never "Email not found" (info disclosure) and never "Wrong password" (also info disclosure).

**Password reset:**
- Acknowledge that you got the request before the email arrives: "If [email] is an account, we've sent a reset link." The "if" wording prevents email enumeration.
- The email's subject line is part of the microcopy: "Reset your [product] password" — specific product, specific verb.
- The reset page itself shows the requirements *before* the user types.

### 8. Onboarding tours and tooltip walkthroughs

Onboarding microcopy is the most-skipped writing in the product. Assume the user will skim, then act. Two principles:

1. **Show, don't tell.** A tour step that says "Click here to create a project" is weaker than a tour step that says "Create a project" with the create button already highlighted and active.
2. **Make every step skippable.** "Skip tour" is always visible; "Next" is the primary action. The tour is not a quiz.

Step copy structure:
- Title (4-7 words, describes the surface or the user goal).
- One sentence of body, max two.
- One primary action (Next, Got it, Start tour) and one escape (Skip).

> **Pin your favorite reports** (title)
> Tap the star on any report to pin it to your sidebar for fast access. (body)
> [Got it] [Skip tour] (actions)

Don't:
- Use a tour to introduce features users haven't yet earned the right to care about (premium-tier features, admin-only screens).
- Block the UI behind a tour they cannot dismiss.
- Run tours on every visit. Once and then never again, unless the user opts back in.

### 9. Consent, opt-in, opt-out — symmetry of choice

This is where microcopy meets ethics. The principle from CPPA, GDPR, and ethical-design literature: **the path to the more privacy-protective choice must not be harder than the path to the less privacy-protective choice.**

| Anti-pattern | Description | Fix |
|---|---|---|
| Confirmshaming | "No thanks, I don't want to save money" | "No thanks" — neutral phrase only |
| Pre-checked opt-in | Marketing-emails box pre-selected | Unchecked by default for non-essential consents |
| Buried "Reject all" | Big "Accept all" button, "Reject" hidden behind a menu | Same visual weight, same number of clicks |
| Manipulinks | Color-coded "Yes" in primary brand color, "No" in gray | Both options must look like real options |
| Trick wording | "Don't not opt out of unsubscribing" | Single, direct opt-in question |
| Forced action | Cookie banner that won't close until you click "Accept" | Always allow an "X" or escape |

Cookie banner copy that respects the user:
> **We use cookies to improve [product].**
> Essential cookies keep the product working. Analytics cookies help us understand how it's used. You can choose what to allow.
> [Accept all] [Reject non-essential] [Customize]

All three buttons have the same visual style. The user is not punished for clicking "Reject."

Marketing opt-in inside a sign-up:
> [ ] Email me product updates and tips. You can unsubscribe anytime.

Unchecked by default. One sentence. No emoji, no urgency.

### 10. Voice and tone — Yifrah's framework and the NN/g four dimensions

Voice is constant across the product; tone shifts by context. Mailchimp's example: voice is "fun, smart, helpful, but never silly or condescending." The tone of a 404 page is different from the tone of a billing error.

NN/g's four dimensions to calibrate tone in any single piece of microcopy:
- **Funny vs Serious** — billing errors are serious, an empty state for a personal hobby app can be playful.
- **Formal vs Casual** — enterprise tools lean formal, consumer apps lean casual.
- **Respectful vs Irreverent** — almost always respectful in product copy. Irreverence is for brand campaigns, rarely for UI.
- **Enthusiastic vs Matter-of-fact** — error states are matter-of-fact, success states can be enthusiastic but should not be exhausting.

Yifrah's voice-design exercise: write the same message three ways — as a strict bank, as a friendly assistant, as a witty co-worker. Pick the one that matches the brand. Then write three error messages and three empty states in that voice to see if it scales.

The Mailchimp humor rule: "Mailchimp has a sense of humor, so feel free to be funny when it's appropriate — but don't go out of your way to make a joke. Forced humor can be worse than none at all." Humor in microcopy is allowed in *one* place per session, max. A funny empty state plus a funny error plus a funny toast equals a product that is trying too hard.

When to drop the humor entirely:
- Anything involving money, security, or privacy.
- Anything destructive or irreversible.
- Anything during an outage, payment failure, or sign-in failure.
- Forms that are emotionally charged (cancellation, account deletion, bereavement features).

## Templates and examples

### Button label rewrites

| Context | Weak | Strong |
|---|---|---|
| Create-flow primary | Submit | Create project |
| Save-flow primary | OK | Save changes |
| Destructive | Yes / Delete | Delete project |
| Sign-up primary | Sign up | Create free account |
| Sign-in primary | Submit / Login | Sign in |
| Multi-step next | Next | Continue to payment |
| Modal dismiss | Close / No | Cancel |
| Save-and-exit | Save | Save and close |
| Confirm purchase | Confirm | Place order — $42.00 |
| Download | Click here | Download CSV |
| Add to list | + | Add to cart |

### Empty-state rewrites

**Project list, first-time:**
> **No projects yet**
> Projects organize your dashboards, alerts, and team access. Create one to get started.
> [Create project]

**Saved searches, first-time:**
> **You haven't saved any searches**
> Save a query to rerun it with one click or share the link with teammates.
> [Save current search]

**Filtered list, no results:**
> **No matches for "production-east-2"**
> Check the spelling or try a broader filter.

**Inbox, all caught up:**
> **You're all caught up**
> New messages will show up here.

(Note the last one — for an "all done" state, no CTA is correct. Don't push the user back into work.)

### Error message rewrites

| Weak | Strong |
|---|---|
| Invalid email. | Email needs an @ and a domain — for example, you@company.com. |
| Server error. | We couldn't reach the server. Check your connection or try again in a moment. |
| Password too weak. | Add at least one number and one symbol to make this password stronger. |
| Field is required. | Add a project name to continue. |
| Login failed. | Email or password doesn't match. Try again or [reset your password]. |
| Cannot delete. | This project has 12 active alerts. Pause or delete those first, then try again. |
| File too large. | This file is 22 MB. The max is 10 MB. Try compressing it or splitting it. |
| Network error 503. | We're having trouble on our end. Try again in a minute, or check [status]. |
| Unknown error. | Something went wrong on our side. Try again, or [contact support] with this ID: 7af3-2c1d. |

### Confirmation dialog rewrites

**Weak (delete project):**
> **Confirm**
> Are you sure?
> [OK] [Cancel]

**Strong:**
> **Delete project Acme?**
> This will permanently delete all 47 dashboards, alerts, and integrations in this project. This cannot be undone.
> [Delete project] [Cancel]

**Strong (high-stakes — type to confirm):**
> **Delete production cluster atlas-prod-east?**
> This will permanently delete the cluster and all backups older than 7 days. Any application still connected will lose access immediately.
>
> To confirm, type the cluster name:
> [          atlas-prod-east           ]
> [Delete cluster] (disabled until match) [Cancel]

### Toast rewrites

| Weak | Strong |
|---|---|
| Your changes have been saved successfully! | Saved |
| Link copied to clipboard. | Copied |
| Successfully sent the invitation. | Invitation sent |
| The project has been deleted. | Project deleted. [Undo] |
| Operation failed. | Couldn't save — [try again] |

### Sign-up flow

| Step | Strong copy |
|---|---|
| H1 | Create your free [product] account |
| Sub | No credit card. Cancel anytime. |
| Email label | Email |
| Password label | Password |
| Password helper | At least 12 characters with a number and a symbol. |
| Marketing opt-in | [ ] Send me product updates and tips. You can unsubscribe anytime. |
| Primary CTA | Create account |
| Already-have-account | Already have an account? [Sign in] |

### Password reset (acknowledgment screen)

> **Check your email**
> If [shown email] is an account on [product], we've sent a link to reset the password. The link expires in 30 minutes.
> Didn't get it? Check spam, or [try again with a different email].

### 404 page

> **We can't find that page**
> The link may be broken or the page may have moved.
> [Go to home] [Search]
>
> *(Optional small text)* If you were following a link from somewhere on [product], [let us know] so we can fix it.

### Onboarding tour step

> **Pin reports to your sidebar**
> Tap the star on any report to pin it. Pinned reports stay accessible from the left nav, even when you switch projects.
> [Got it] [Skip tour]

### Cookie banner

> **We use cookies to improve [product]**
> Essential cookies keep the product working. Optional cookies help us understand usage and improve the experience. You can choose what to allow.
> [Accept all] [Reject non-essential] [Customize]

## Anti-patterns

### "Click here"
The link should describe its destination. "Click here" reads to screen readers as "click here" with no context. Use "Download CSV," "View the report," "Read the API docs."

### "Submit"
Almost never the right word. The action is always more specific: Create, Save, Send, Publish, Order, Reset.

### Placeholder-as-label
Field with a placeholder "Email" and no label. When the user clicks in, the label vanishes and they lose context. Always have a visible, persistent label.

### Apologetic empty states
"Oh no, you have no projects yet!" The user knows. Skip the theater and tell them what to do.

### "Are you sure?" as the only confirmation question
Vague. Replace with "Delete project X?" — the verb-noun matches the action.

### Yes/No buttons on confirmation dialogs
Forces the user to re-read the title to know which is the dangerous answer. Use verb-noun.

### Exclamation points in error messages
Mailchimp's rule: never. An exclamation point in an error reads as either sarcasm or panic. Both are bad.

### Error messages that blame the user
"Invalid input." "Wrong password." Replace with what the user can *do*.

### Pre-checked marketing opt-ins
Dark pattern. Unchecked by default for non-essential consents.

### Tooltips containing critical information
If the user needs it to fill the field, it's not a tooltip — it's helper text.

### Forced humor everywhere
Funny empty state, funny error, funny toast, funny 404. Pick one surface to be funny, max one per session. Humor in errors that cost the user money or data is malpractice.

### Toast messages over 5 words
"Your project has been successfully created and saved." Should be "Project created."

### "Please" as a softener
"Please enter your name." "Please try again." Almost always cuttable. Polite ≠ adding "please" everywhere. Polish: "Enter your name to continue."

### Title case on body text
Title Case Is For Page Headers Only In Most Systems. Material, GOV.UK, Polaris all use sentence case for body, buttons, and labels. Title case is a brand decision; if your system is sentence case, stay there.

## Decision heuristics

**Should this button say "Submit" or something specific?**
Always specific. Submit is the word designers reach for when they have not asked what the user is *doing*. The verb the user came to perform is the button.

**Is this confirmation dialog necessary, or is it friction?**
A confirmation is justified if (a) the action is irreversible AND (b) the action is non-trivial to recover from. Saving a draft does not need a confirmation. Deleting a draft does, lightly. Deleting a project with 47 child resources does, with detail. Deleting a production cluster does, with type-to-confirm.

**Is this a tooltip or helper text?**
If most users will need it to use the field correctly, it's helper text (always visible). If only some users will want it, and the answer is non-blocking, it's a tooltip.

**Should this empty state have a button?**
Yes, unless it's an "all done" state (inbox zero, no errors, no pending tasks). For first-time states, always provide the next action.

**Is this a toast or a modal?**
Toast if the user can keep working. Modal if the next step requires their input or attention before anything else can happen.

**Is humor okay here?**
Money, security, privacy, errors, destructive actions — no. Empty states for low-stakes features, 404 pages, success states — yes, but check the brand voice and the surrounding density of jokes.

**Should I use sentence case or title case?**
Default to sentence case. Title case for proper-noun product features (e.g., "Atlas Vector Search") and for top-level page headers if the design system specifies. GOV.UK, Mailchimp, Polaris, Material — all sentence case. Apple HIG uses title case for some platform elements; consult the local system.

**How long should an error message be?**
Polaris's rule: no longer than two sentences. NN/g concurs. If you need more, you're documenting; link to docs.

**Is this opt-in copy ethical?**
Apply the symmetry test: can the user choose the more private option with the same number of clicks and the same prominence as the less private option? If not, redesign.

**Should the field validate as the user types?**
No, validate on blur (after they leave the field). The exception is positive feedback for things like password requirements, where progressive checkmarks reassure rather than nag.

**Should this toast have an undo?**
Yes if the action was destructive or affected something the user might regret. "Project deleted [Undo]" is better than "Are you sure?" before delete + irreversible after. Polish UX is moving toward undo-everywhere over confirm-everything.

## Cross-references

- **support-ticket-writing** for customer-facing prose (ticket replies, apologies, status updates). This skill is the in-product complement.
- **writing-expert** for foundational prose craft and anti-AI-isms. Layer it on top.
- **plain-language** when the audience has a defined literacy level or accessibility constraint.
- **inclusive-language** for inclusivity audit passes.
- **accessibility-ux-reviewer** for screen-reader compatibility and ARIA labels.
- **frontend-design** for the component-level conventions (where strings live, label/helper/error patterns).
- **ui-ux-pro-max** and **vanilla-js-ui-reviewer** for UI review including the words.
- **sales-and-marketing-copy** when the surface is acquisition or conversion-driven.

## References

- [Mailchimp Content Style Guide — Voice and Tone](https://styleguide.mailchimp.com/voice-and-tone/) — voice consistency, tone shifts by user state, no exclamation marks in errors.
- [Mailchimp Content Style Guide — Grammar and Mechanics](https://styleguide.mailchimp.com/grammar-and-mechanics/) — sentence case, active voice, plain English.
- [Shopify Polaris — Voice and Tone](https://polaris.shopify.com/foundations/content/voice-and-tone) — read it aloud, sound human, merchant-first.
- [Shopify Polaris — Error Messages](https://polaris-react.shopify.com/content/error-messages) — two-sentence cap, actionable not permissive, no hyperbole.
- [Apple Human Interface Guidelines — Writing](https://developer.apple.com/design/human-interface-guidelines/writing) — clarity over cleverness, verbs on buttons.
- [Apple Human Interface Guidelines — Alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) — Cancel always cancels, Destructive style for destructive actions.
- [Material Design 3 — Content Design Style Guide](https://m3.material.io/foundations/content-design/style-guide) — sentence case, second person, simple direct language.
- [Material Design 2 — Writing](https://m2.material.io/design/communication/writing.html) — voice vs tone, contextual tone shifts.
- [GOV.UK Design System — Button](https://design-system.service.gov.uk/components/button/) — sentence case, describe the action.
- [GOV.UK Style Guide](https://www.gov.uk/guidance/style-guide) — plain English, accessible writing.
- [Nielsen Norman Group — Error-Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/) — say what, why, and how to fix; do not blame the user.
- [Nielsen Norman Group — Hostile Patterns in Error Messages](https://www.nngroup.com/articles/hostile-error-messages/) — premature errors, blame, ALL CAPS.
- [Kinneret Yifrah — *Microcopy: The Complete Guide* (2nd ed.)](https://www.amazon.com/Microcopy-Complete-Guide-Kinneret-Yifrah/dp/B07N1RD7W6) — voice and tone framework, success/error/empty-state patterns.
- [Ketch — Ethical Consent UX](https://www.ketch.com/blog/posts/dark-patterns-how-to-stop) — symmetry in choice, plain language for cookie banners.
- [UX Booth — Manipulinks and Confirmshaming](https://uxbooth.com/articles/ux-dark-patterns-manipulinks-and-confirmshaming/) — naming the patterns.
- [John Saito, Dropbox — UX writing essays](https://uxdesign.cc/@jsaito) — "short always beats good," simple/straightforward/human voice principles.
- [Carbon Design System — Writing Style](https://carbondesignsystem.com/guidelines/content/writing-style/) — IBM's enterprise-focused content rules.

---

## The 5-second test — can users articulate the page purpose in 5 seconds?

**Rule.** Steve Krug (*Don't Make Me Think*, 3rd ed., 2014): show a representative user a page for 5 seconds, take it away, and ask "what is this page for? what can you do here? whose product is it?" If they cannot answer all three after 5 seconds of exposure, the page has failed the test and the copy/IA needs work.

The test isolates first-glance comprehension from any kind of deep reading. It is the microcopy equivalent of a smell test — it does not prove the page is good; it just exposes the cases where comprehension never starts.

**The three questions.**

1. *What is this?* (product, page type, surface — landing page, app dashboard, settings panel)
2. *What can I do here?* (the primary action or task)
3. *Who is it for / who made it?* (the brand or audience cue)

A page that earns all three answers in 5 seconds is *findable* in the user's mental map. A page that earns one or two is workshoppable. A page that earns zero is a redesign.

**Worked example — a SaaS dashboard.**

- Failing: a page whose header is the company logo, whose primary content is a stack of widgets labeled "Overview" "Activity" "Insights," and whose tagline is "Powerful analytics for modern teams." After 5 seconds, the user can name the brand but cannot say what *this* page is for or what they should do.
- Passing: same page, with header text "Production cluster overview — last 24 hours" and a single primary button "Investigate alert." After 5 seconds, the user can answer: *cluster health page, investigate the alert at the top, MongoDB Atlas.*

**Operational checklist.**

- [ ] The page has a single H1 that names the page or the user's current state (not the product or the company)
- [ ] The primary action is visually dominant and labeled as a verb (Save, Investigate, Connect, Continue — not "Submit" or "OK")
- [ ] The brand/product cue is present but not the loudest element on the page
- [ ] Body copy below the fold is not load-bearing for first-glance comprehension

**Tools to run the test.**

- UsabilityHub / Lyssna (formerly UsabilityHub) — purpose-built 5-second tests.
- Maze, UserTesting — both support 5-second tests as a question type.
- Lo-fi alternative: open the page for 5 seconds in front of a colleague who has not seen it; close the laptop; ask the three questions.

**When to break it.** Deep configuration screens (admin consoles, SQL editors, IDE panels) are *not* meant to be comprehensible in 5 seconds — they are meant to be efficient for trained users. The 5-second test applies to *first-encounter surfaces*: landing pages, marketing sites, onboarding flows, app dashboards, signup screens, error pages, and empty states.

**Composition with sibling rules.** Pair the 5-second test with the F-pattern / Z-pattern reading rules below — *where* the eye lands first determines *what* it can absorb in 5 seconds.

**References.**

- Krug, S. *Don't Make Me Think, Revisited: A Common Sense Approach to Web Usability*, 3rd ed., New Riders, 2014, chapter 9.
- Nielsen Norman Group — "The 5-Second Test in UX Research." https://www.nngroup.com/articles/usability-testing-101/
- Lyssna (formerly UsabilityHub) — "5-second tests." https://www.lyssna.com/guides/5-second-test/

---

## Information scent in the UI

**Rule.** This skill's parent concept (information scent — Pirolli & Card, *Information Foraging*) is covered in detail in `writing-expert`. In UI microcopy the rule applies to *labels, links, tabs, menus, and headings* — every navigation cue a user uses to decide where to go next.

**UI-specific applications.**

| Element | Weak scent | Strong scent |
|---|---|---|
| Top-nav tab | "Resources" | "Docs and API reference" |
| Settings group header | "General" | "Account, notifications, billing" |
| Menu item | "Tools" | "Export, import, bulk edit" |
| CTA button | "Learn more" | "See pricing and plans" |
| Empty-state heading | "Nothing here yet" | "Create your first cluster to get started" |

Generic labels ("More," "Other," "Settings," "Tools," "Resources") have near-zero scent — the user cannot predict the contents. They are tolerated only when the user has already learned the IA from repeated exposure. New users see them as black boxes and skip them.

**Worked example — a search results page.**

Weak-scent label: a sidebar facet labeled "Type" with values "Document, Page, Item." The user cannot predict what they will find by clicking.
Strong-scent label: the same sidebar labeled "Content type" with values "Knowledge base article, runbook, slide deck, customer story." Now the user can navigate by *prediction* rather than *trial*.

**Diagnostic.** Take a screenshot of any navigation surface (menu, tab bar, sidebar, breadcrumb trail) and show it to a target user. Cover the page contents. Ask them: "If you wanted to do X, where would you click?" If they hesitate or guess wrong, the labels have weak scent.

**Compose with this skill's `## Tour writing` section.** A guided tour is information scent applied to onboarding: each step must let the user predict what they will get from continuing.

**References.**

- Pirolli, P., & Card, S. "Information foraging." *Psychological Review* 106(4), 1999.
- Nielsen, J. "Information Scent: How Users Decide Where to Go Next." NN/g, 2003. https://www.nngroup.com/articles/information-scent/
- See also: `writing-expert` skill — *Information scent — Pirolli & Card* subsection for prose applications.

---

## F-pattern and Z-pattern reading — where the eye actually lands

**Rule.** Eye-tracking research from Jakob Nielsen / NN/g (2006, with multiple updates through 2017+) shows that web users do not read pages top-to-bottom, left-to-right; they *scan*. Two dominant patterns emerge based on content density:

- **F-pattern.** Text-heavy pages (search results, articles, long-form lists, blog posts). Users scan in an F shape: a horizontal sweep across the top, a second shorter horizontal sweep about a third of the way down, then a vertical scan down the left edge. Content on the right and bottom of dense text pages is often missed entirely.
- **Z-pattern.** Sparse pages (landing pages, hero sections, simple forms, signup screens). The eye traces a Z: top-left to top-right (logo to nav/CTA), then diagonally down to the bottom-left, then across to the bottom-right (often the primary action).

These are descriptive (what eyes do) not prescriptive (where to put things), but the implication for layout is direct: *put the highest-value content where the eye actually lands*.

**Worked example — landing page (Z-pattern target).**

Strong: logo top-left; primary nav and "Sign in" top-right; bold value proposition mid-left; product image or testimonial mid-right; primary CTA ("Start free trial") bottom-right. The Z lands the eye on the CTA.

Weak: the same content stacked center-aligned without spatial anchoring. The eye still tries to Z-scan; finds nothing at the corners; gives up.

**Worked example — article page (F-pattern target).**

Strong:
- First sentence of each paragraph carries the topic (the F's horizontal sweeps land there).
- Subheadings every 3–5 paragraphs anchor the left-edge vertical scan.
- Pull quotes and key terms bolded inline at the start of paragraphs, not the middle.
- Long horizontal text lines (>~75 characters) broken into shorter measures.

Weak: the same content as flowing prose with no subheadings, with key conclusions buried in the third sentence of each paragraph. The F-scanner misses the conclusions entirely.

**Operational rules.**

1. *Front-load.* In F-pattern surfaces, put the topic and the conclusion at the front of every paragraph and every list item. This composes with this skill's `## Headings and list items → Front-load` rule.
2. *Anchor the corners.* In Z-pattern surfaces, place a deliberate element at each Z corner (logo, secondary CTA, value prop, primary CTA). Empty corners waste the user's first scan.
3. *Right-edge content is for repeat visitors.* In F-pattern pages, content on the right edge will be missed by first-time scanners. Put it there only when the user has a reason to seek it.
4. *RTL languages flip both patterns.* Arabic, Hebrew, and other RTL UIs see a mirrored F and a mirrored Z. If you ship localized UI, plan for both.

**When to break it.** These are *averages over heatmap studies*; individual users vary, and tasks change the pattern. A user searching for a specific answer scans differently from a user exploring. Use the patterns as a *bias to plan around*, not as a contract.

**References.**

- Pernice, K. "F-Shaped Pattern of Reading on the Web: Misunderstood, But Still Relevant (Even on Mobile)." NN/g, 2017. https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/
- Nielsen, J. "F-Shaped Pattern For Reading Web Content." NN/g, 2006. https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content-discovered/
- Babich, N. "Z-Shaped Pattern For Reading Web Content." UX Planet / Adobe XD Ideas (industry summary of NN/g and supplementary eye-tracking research). https://uxplanet.org/z-shaped-pattern-for-reading-web-content-ce1135f92f1c
