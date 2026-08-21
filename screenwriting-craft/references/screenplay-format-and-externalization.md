---
name: screenplay-format-and-externalization
description: >-
  Screenplay page format (slug lines, action lines, parentheticals, dialogue
  blocks, character-name conventions, transitions, Final Draft vs. Fountain,
  page/margin standards) and the externalization constraint that everything
  in a screenplay must be shown as filmable action, dialogue, or image.
  TRIGGER: format a slug line/scene heading; write action lines or
  parentheticals; V.O./O.S./CONT'D conventions; when to use a transition
  like CUT TO; Final Draft vs Fountain; how do I show what a character is
  thinking/feeling on screen; unfilmable interiority; screenwriting-specific
  show-don't-tell. SKIP: page-per-minute pacing and white-space-as-pacing ->
  references/scene-economy-and-script-types.md; beat systems/structure ->
  references/structure-and-beat-systems.md; prose show-don't-tell with no
  screen intent -> fiction-writing-craft.
origin: local
version: "1.0.0"
updated: "2026-07-07"
keywords:
  - slug line
  - scene heading
  - action lines
  - parentheticals
  - wrylies
  - Final Draft
  - Fountain markup
  - voice over
  - off screen
  - externalization
  - visual correlative
  - subtext
  - show dont tell screenwriting
tags:
  - screenwriting
  - format
  - film
  - television
  - craft
---

# Screenplay Format & the Externalization Constraint

A screenplay is a blueprint, not a finished text. It is read by a small,
time-pressed set of professionals (agents, producers, script readers,
eventually a crew) who need to picture a shot, a schedule, and a budget from
the page, and every format convention below exists to serve that reading,
not decoration. The externalization constraint is the deeper rule underneath
all of it: since the reader's job is to imagine a camera and a soundtrack, a
screenplay can only put on the page what a camera and a microphone could
capture. Format tells you *how* to write what the camera sees; externalization
tells you *what* is allowed to be there at all.

## Core Concepts

### (a) Slug lines (scene headings)

A slug line (master scene heading) opens every scene and gives production
three pieces of information: interior or exterior, location, and time of day.
Convention: `INT.` or `EXT.` (always abbreviated, with a trailing period),
followed by the location, then a hyphen and a time-of-day marker (DAY,
NIGHT, CONTINUOUS, LATER). `INT./EXT.` marks a scene spanning both, most
commonly a vehicle interior seen through a moving exterior.[^1]

Worked example: `INT. CAR - MOVING - NIGHT` tells a location manager,
gaffer, and AD everything they need before reading a single line of action.

When to reach for it: every scene needs one, and only one, at the top. Do
not stack multiple slug lines for camera moves within the same continuous
scene and location: that is what action lines and (sparingly) transitions
are for.

### (b) Action lines / description

Action lines render only what a camera could record: setting, character
appearance and movement, sound cues that matter to the story. Convention:
third person, present tense ("she runs," not "she ran"), active voice, and
short paragraphs: most craft sources converge on roughly 3-4 lines per
paragraph, with white space itself functioning as a pacing signal (a dense
block reads slow to a professional reader regardless of its content; see
`references/scene-economy-and-script-types.md` for the full pacing
treatment).[^2] The discipline is to include only details that affect tone,
plot, or reveal something new; visualize with the fewest words that still
land the image.

Worked example: "RAY (40s, exhausted) drops his keys twice before the door
gives" does more work in one line than a paragraph describing his emotional
state directly, because it is filmable and implies the feeling.

When to reach for it: whenever you are tempted to write more than a
sentence or two of pure scene-setting; cut to the detail that is either
plot-load-bearing or the single strongest image, and trust the actor and
director with the rest.

### (c) Character-name conventions

A character's name appears in ALL CAPS the first time they are introduced in
action/description if they will have dialogue, signaling a speaking role to
casting and the AD. In the dialogue block, the character's name is always
centered in ALL CAPS above their lines, every time they speak.[^3] Modifiers
in parentheses next to the name:

- **(V.O.)** (voice-over): the character is not physically present in the
  scene (narration, a phone caller, an unseen narrator).
- **(O.S.)** / **(O.C.)** (off-screen / off-camera): the character is
  present in the scene's world but not currently in frame.
- **(CONT'D)** (the same character's dialogue continues after an interruption or a page break):
  the label marks it.

When to reach for it: use V.O. and O.S. precisely, not interchangeably: the
distinction (absent from the scene entirely vs. present but off-frame)
matters to a director staging the shot.

### (d) Parentheticals ("wrylies")

Parentheticals are short actor/tone directions placed under a character's
name and before their dialogue, e.g. "(sarcastic)." The term "wrylies"
mocks novice writers who overuse "(wryly)" specifically. Craft consensus is
close to unanimous: parentheticals are the most commonly overused and
misused format element: they read as the writer directing the actor from
the page and second-guessing a delivery the dialogue and context should
already convey.[^4] Best practice: use sparingly (some pros suggest a rough
ceiling near one per page) and reserve them for the case where the line
would otherwise be misread (clarifying sarcasm that isn't evident from
context) rather than restating an emotion the dialogue already carries.

When to reach for it: only when cutting the parenthetical would change how
an actor reads the line in a way that breaks the scene.

### (e) Transitions

Transitions (`CUT TO:`, `DISSOLVE TO:`, `FADE TO:`, `SMASH CUT TO:`) are
right-margin instructions marking how one scene visually moves to the next.
Convention has shifted: writing `CUT TO:` between every scene is now
considered redundant and dated, since a cut is already implied by the next
slug line; it only adds page length and clutter.[^5] Modern guidance is to
use a transition only when it is doing narrative work: `CUT TO:` for a
deliberately jarring, abrupt juxtaposition; `DISSOLVE TO:` for a softer
passage of time or a dreamlike link between images.

When to reach for it: sparingly, as a chosen effect, not as connective
tissue between every pair of scenes.

### (f) Final Draft vs. Fountain

**Final Draft** is repeatedly described by independent trade sources as the
long-standing industry-standard screenwriting application; its `.fdx` file
format functions as a professional signifier, and being asked to
submit/collaborate in Final Draft format is the default expectation in
film/TV production workflows.[^6] (Final Draft's own marketing claims a
specific "used by 95% of Hollywood productions" figure; that exact
statistic is vendor-sourced and should be treated as a marketing claim, not
an independently verified fact; the broader "industry standard"
characterization is corroborated independently.)

**Fountain** is a free, open-source, plain-text markup language for
screenwriting (inspired by Markdown), letting a writer produce a properly
formatted screenplay from any plain-text editor once run through a renderer
such as Highland.[^7] It has been adopted as an import/export option by
several professional-adjacent apps (Scrivener, Slugline, Storyist), giving
it real but secondary standing: it supplements rather than replaces Final
Draft in mainstream studio workflows. No source found gives hard data on
studio-level Fountain adoption. Treat "present but secondary to Final
Draft at the studio level" as a qualified, not fully sourced, claim.

**Page and margin standards.** Screenplays are set in 12-point Courier (or
Courier Prime), a fixed-pitch font chosen specifically because uniform
character width makes page count a reliable-enough proxy for runtime.
Standard margins run roughly 1.5" on the left (for binding) and 1" on the
right/top/bottom for action, with dialogue further indented; a page holds
roughly 55 lines.[^8] This formatting standard is the entire mechanical
basis for the "one page equals one minute" convention. See
`references/scene-economy-and-script-types.md` for how reliable that
convention actually is in practice.

### (g) The externalization constraint: why it exists

Screenwriting sources converge on one causal explanation: film and TV are
audiovisual media experienced in real time: an audience gets only what the
camera records and the soundtrack carries. Unlike prose, which can narrate a
character's interior monologue directly for as long as it wants, a
screenplay has no narrator of thought; the "reader" of the finished work is
ultimately an audience watching actors, sets, and edited images. Robert
McKee's *Story* is the most frequently cited formal statement of this idea:
film's aesthetics are overwhelmingly visual, and dialogue itself is a
secondary tool behind pure image: "never write a line of dialogue when you
can create a visual expression."[^9] (McKee's book is repeatedly summarized
as putting a specific "80% visual / 20% auditory" ratio on this; that exact
figure traces to a single attribution across the sources found here and
should be treated as tentative even though the underlying visual-primacy
argument is broadly corroborated.)

### (h) Techniques for externalizing internal states

1. **Behavior as an emotional proxy.** Instead of stating a feeling, show a
   physical tell: "she picks at the label on her beer bottle" rather than
   "she's nervous." This converts an unfilmable internal state into a
   filmable, actor-playable beat.
2. **Subtext.** Dialogue that means something other than its literal
   content: characters talk around the real subject, lean on irony or
   metaphor, and the audience infers the true stakes from context. Subtext
   is the primary vehicle for conveying interiority *through dialogue*
   specifically, as distinct from action lines.[^10]
3. **Visual metaphor (the "visual correlative").** McKee's specific
   prescription: find a visual image that stands for an inner conflict,
   rather than having a character explain their inner state aloud.
4. **Reaction and staging.** Facial expression, blocking, and what a
   character does *not* say or do are legible signals of internal conflict
   without narrating it.
5. **Voice-over: the sanctioned exception.** Screenwriters can break pure
   externalization via V.O., but this is a deliberate, limited tool, not a
   workaround: it should add information the visuals and dialogue can't
   (an unreliable narrator, dramatic irony where the audience learns a
   truth the character hasn't yet). It fails as craft when it merely
   restates what is already visible or audible: the commonly cited bad
   example is a V.O. line like "I am sad because my father left me" layered
   over a shot that already shows the sadness.[^11]

### (i) "Show don't tell": the screenwriting-specific version

The general prose "show don't tell" principle is about narrative
engagement and style; the screenwriting-specific version is a harder
constraint rooted in the medium itself, not just stylistic advice. A
novelist who "tells" is making a stylistic misstep; a screenwriter who
"tells" (writing unfilmable interiority directly into action lines) is
often writing a physically unproducible instruction, because nothing on the
page corresponds to anything a crew could put in front of a camera. Only a
couple of sources draw this screen-vs-prose distinction this precisely, so
treat the sharp conceptual line as a qualified synthesis; the surrounding
externalization consensus itself is very strong.[^12]

### (j) Common failure modes

- Writing internal states directly into action lines: "He understands...",
  "She realizes...", "He feels sad" — these are textbook unfilmable
  constructions, since there is no way to physically stage "understanding"
  without an accompanying visible or audible behavior.
- A recognizable failure pattern: "Character turns the corner, wondering
  where in his life he'd gone wrong" is unfilmable absent voice-over; the
  best a camera can do is show a face "deep in thought," which reads
  ambiguously on screen.[^13]
- This mistake is specifically attributed to writers transferring prose
  habits into script form. Prose-to-screen crossover writers are called
  out as most prone to it, since direct narration of thought has no
  screenplay equivalent.
- Overreliance on unearned voice-over as a shortcut, restating emotion
  that is already visible rather than adding new information. This is
  "telling" via narration rather than a legitimate use of the V.O.
  exception.

## Anti-Patterns

**Directing from the page with parentheticals.** Stacking a parenthetical
under every line of dialogue to control an actor's delivery. Diagnose by
counting parentheticals per page; more than one or two per page is a signal
the dialogue itself isn't carrying the tone. Fix by cutting the
parenthetical and trusting the line, or rewriting the line so the tone is
unmistakable without a stage direction.

**Redundant transitions.** A `CUT TO:` after every scene heading. Diagnose
by checking whether removing the transition changes anything about how the
scene reads. If not, it is dead weight. Fix by deleting default
transitions and reserving them for scenes where the cut itself is a
deliberate effect.

**Unfilmable interiority in action lines.** Writing what a character
thinks, understands, or realizes as a bare narrative statement. Diagnose by
asking "what would the camera actually record here?" If the answer is
nothing, the line needs a behavioral or visual substitute. Fix using the
externalization techniques above: convert the stated feeling into a
physical tell, a piece of subtext-carrying dialogue, or a visual
correlative.

**Voice-over as a crutch.** Using V.O. to state an emotion or fact that the
visuals and dialogue already convey. Diagnose by removing the V.O. line
mentally and checking whether the scene still lands; if it does, the V.O.
is redundant. Fix by cutting the V.O. or replacing it with a use that adds
information the audience genuinely couldn't get otherwise (irony, an
unreliable narrator, a fact the character hasn't yet realized).

**Format rule-breaking without craft justification.** Deviating from
standard format (font, structural slug-line data, INT./EXT. conventions)
purely for stylistic effect. Craft sources are consistent that
*stylistic* rules (avoiding camera direction, minimizing transitions,
limiting parentheticals) are far more forgivable to break than
*structural* elements that carry real production information (scene
heading data, clear separation of description vs. dialogue) because
below-the-line crew depend on the latter for scheduling, lighting, and set builds. This
exception-tolerance is also gated by track record: an established
hyphenate writer breaking format reads as a bold choice; a spec-script
newcomer doing the same reads as a red flag.[^14]

## References

[^1]: StudioBinder, "What is a Slug Line? Scene Heading Screenplay Formatting," https://www.studiobinder.com/blog/what-is-a-slug-line-definition/ ; ScreenCraft, "Screenwriting Basics: The Keys to Writing Correct Scene Headings," https://screencraft.org/blog/screenwriting-basics-the-keys-to-writing-correct-scene-headings/ ; Story Sense, "Screenplay Format Guide: Scene Headings," https://www.storysense.com/format/headings.htm

[^2]: No Film School, "It's Time to Write Better Action Lines," https://nofilmschool.com/write-better-action-lines ; Backstage, "How to Write Action Lines in a Script," https://www.backstage.com/magazine/article/how-to-write-action-lines-in-script-75742/

[^3]: SoCreate, "How To Use Capitalization In Traditional Screenwriting," https://www.socreate.it/en/blogs/screenwriting/how-to-use-capitalization-in-traditional-screenwriting ; Talentville University, "Capitalization: Character Intros," https://www.talentville.com/snippet/capitalization--character-intros

[^4]: MovieOutline, "The Curse of Quirky Parentheticals — Using Wrylies in Your Screenplay," https://www.movieoutline.com/articles/writing-parentheticals-and-wrylies-in-your-screenplay.html ; Script Magazine, "Why Spec Scripts Fail: The 'Wrylie' (Parentheticals)," https://scriptmag.com/features/spec-scripts-fail-wrylie-parenthetical-2 ; Final Draft, "Screenwriter's Guide to Using Parentheticals in Screenplays," https://www.finaldraft.com/blog/lets-talk-about-parentheticals

[^5]: Final Draft Blog, "Transitions: To 'Cut To' or Not to Cut To," https://www.finaldraft.com/blog/transitions-cut-not-cut ; StudioBinder, "How to Write Transitions in a Script," https://www.studiobinder.com/blog/how-to-write-transitions-in-a-script/

[^6]: Final Draft (official), https://www.finaldraft.com/ ; Wikipedia, "Final Draft (software)," https://en.wikipedia.org/wiki/Final_Draft_(software)

[^7]: Fountain.io (official), https://fountain.io/ ; Wikipedia, "Fountain (markup language)," https://en.wikipedia.org/wiki/Fountain_(markup_language)

[^8]: Celtx Blog, "Screenplay Margins: The Formatting Rules You Need to Know," https://blog.celtx.com/screenplay-margins-guide/ ; ScreenWeaver Blog, "Screenplay Format Guide 2026: Margins, Fonts, and Slugs," https://www.screenweaver.ai/blog/screenplay-formatting-guide-2026

[^9]: Robert McKee, *Story: Substance, Structure, Style, and the Principles of Screenwriting* (1997); secondary coverage: StudioBinder, "Who is Robert McKee," https://www.studiobinder.com/blog/robert-mckee-screenwriting/

[^10]: StudioBinder, "What is Subtext — How to Use Subtext in Screenwriting," https://www.studiobinder.com/blog/what-is-subtext-definition/ ; The Write Practice, "Subtext Examples: 7 Simple Techniques to Supercharge Your Scenes," https://thewritepractice.com/subtext-examples/ ; Industrial Scripts, "Writing Subtext: What are the Different Forms of Subtext?" https://industrialscripts.com/writing-subtext/

[^11]: Beverly Boy Productions, "How to Write Inner Monologue in a Script," https://beverlyboy.com/filmmaking/how-to-write-inner-monologue-in-a-script/ ; Celtx Blog, "Internal Monologue vs. Voiceover: How to Write and Format Inner Dialogue," https://blog.celtx.com/how-to-write-internal-monologue-vs-voiceover/

[^12]: Script Reader Pro, "Show Don't Tell: How to Turn a Talky Script Into a Visual Masterpiece," https://www.scriptreaderpro.com/show-dont-tell-screenwriting/ ; Script Angel, "Show Don't Tell — Script Development," https://scriptangel.com/show-dont-tell/

[^13]: Scriptwrecked, "Character Thoughts – 'Unfilmables'?" https://scriptwrecked.com/2023/05/03/character-thoughts-unfilmables/ ; Arc Studio Blog, "Can You Write What Your Character Is Thinking in a Script?" https://www.arcstudiopro.com/blog/can-you-write-what-your-character-is-thinking-in-a-script

[^14]: Final Draft Blog, "Screenwriting Rules You Can and Shouldn't Break," https://www.finaldraft.com/blog/screenwriting-rules-you-can-and-shouldnt-break ; No Film School, "A List of Screenwriting Rules and How to Break Them," https://nofilmschool.com/screenwriting-rules
