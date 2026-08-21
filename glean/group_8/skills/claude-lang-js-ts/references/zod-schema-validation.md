<!-- hub-reference-banner -->
> **Reference file — part of the `lang-js-ts` hub.** Formerly the standalone `zod-schema-validation` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: zod-schema-validation
version: "1.2.0"
updated: "2026-05-29"
category: developer
tags: [zod, validation, schema, typescript, runtime-validation, type-safety, z.infer, safeParse, transform, pipe, form-validation, api-validation, zod-error, branded-types, discriminatedUnion]
description: "Zod runtime schema validation expert — core API, schema composition, transforms/pipes, error handling, TypeScript inference (z.infer), API validation middleware, form integration, branded types, recursive schemas, and performance. TRIGGER: user is implementing or reviewing runtime validation in TypeScript/JavaScript; writing Zod schemas for API request/response bodies, form inputs, environment variables, config files, or database rows; asking about z.infer, safeParse, discriminatedUnion, branded types, transforms, or ZodError handling; migrating from Zod 3 to Zod 4; integrating Zod with react-hook-form, tRPC, or Next.js Server Actions; needing to enforce TypeScript nominal typing at runtime. SKIP: user needs JSON Schema validation without TypeScript (use AJV directly); user is asking about Joi, Yup, or Valibot; pure TypeScript compile-time type checking with no runtime validation needed."
related_skills: [programming-languages, software-engineering-patterns, react-nextjs, testing-and-vitest-expert]
---

# Zod Schema Validation

## Overview

Zod is a TypeScript-first schema declaration and validation library. It eliminates the gap between runtime validation and static types by letting you define a schema once and derive both a runtime validator and a TypeScript type from the same source. Zod has zero dependencies, works in Node.js, Deno, Bun, and every modern browser, and is the de-facto standard for runtime validation in the TypeScript ecosystem.

**When to reach for Zod**: API request/response validation, form input validation, environment variable parsing, configuration file validation, database row validation, message contract enforcement, and anywhere you need a runtime guarantee that data matches a TypeScript type.

**Version note**: Zod 4 (released mid-2025) introduced significant performance gains (14x faster string parsing, 7x faster array parsing, 6.5x faster object parsing), a 57% smaller core, built-in JSON Schema conversion, `z.interface()`, `z.stringbool()`, file validation, metadata support, and Zod Mini for bundle-sensitive projects. Patterns below cover both Zod 3 and Zod 4; differences are called out where relevant.

**Python analog**: The equivalent runtime-validation library in Python is **Pydantic v2** (Rust `pydantic-core`). The concepts map closely — `z.object` ≈ `BaseModel`, `.parse()` ≈ `model_validate()`, `.safeParse()` ≈ `model_validate` + `try/except ValidationError`, `z.discriminatedUnion` ≈ `Field(discriminator=...)`, `.transform()` ≈ `field_validator`, env parsing ≈ `pydantic-settings` `BaseSettings`, and JSON Schema export ≈ `model_json_schema()`. The same anti-patterns apply (recreate-schema-per-request, validate-then-pass-dicts). For Python, load `references/pydantic-v2.md` in the `programming-languages` hub.

---

## Core API

### Primitives

```ts
import { z } from "zod";

z.string();  z.number();  z.bigint();  z.boolean();
z.date();    z.symbol();  z.undefined(); z.null();
z.void();    z.any();     z.unknown();   z.never();

z.function().args(z.string(), z.number()).returns(z.boolean());
z.promise(z.string());

// File validation (Zod 4)
z.file().maxSize(5_000_000).mimeType(["image/png", "image/jpeg"]);
```

### Strings

```ts
const Email   = z.string().email();
const Url     = z.string().url();
const Uuid    = z.string().uuid();
const Trimmed = z.string().trim().min(1, "Required");

z.string().min(3).max(255);
z.string().length(6);
z.string().regex(/^[A-Z]{2}\d{4}$/, "Must be 2 uppercase letters + 4 digits");
z.string().toLowerCase();
z.string().toUpperCase();
z.string().normalize(); // Zod 4
```

### Numbers

```ts
z.number().int();
z.number().positive();
z.number().nonnegative();
z.number().min(0).max(100);
z.number().multipleOf(5);
z.number().finite();
z.number().safe();  // within Number.MAX_SAFE_INTEGER

z.int(); // Zod 4 shorthand for z.number().int()
```

### Objects

```ts
const User = z.object({
  id:    z.string().uuid(),
  name:  z.string().min(1),
  email: z.string().email(),
  age:   z.number().int().nonnegative().optional(),
});

type User = z.infer<typeof User>;
// { id: string; name: string; email: string; age?: number | undefined }

User.strict();      // throws on unknown keys
User.strip();       // silently removes unknown keys (default)
User.passthrough(); // preserves unknown keys
```

### Arrays, Tuples, Enums

```ts
const Tags = z.array(z.string()).min(1).max(10);
const Pair = z.tuple([z.string(), z.number()]);
const Rest = z.tuple([z.string()]).rest(z.number()); // [string, ...number[]]

const Role = z.enum(["admin", "editor", "viewer"]);
type Role = z.infer<typeof Role>; // "admin" | "editor" | "viewer"

enum Direction { Up, Down, Left, Right }
const DirectionSchema = z.nativeEnum(Direction);
```

### Unions and Discriminated Unions

```ts
// Simple union — tries each branch in order (slower)
const StringOrNumber = z.union([z.string(), z.number()]);

// Discriminated union — fast lookup on discriminator key (prefer this)
const ApiResult = z.discriminatedUnion("status", [
  z.object({ status: z.literal("success"), data: z.string() }),
  z.object({ status: z.literal("error"),   message: z.string() }),
]);

// Zod 4: supports nested discriminated unions
const NestedResult = z.discriminatedUnion("status", [
  z.object({ status: z.literal("ok"), payload: z.unknown() }),
  z.discriminatedUnion("code", [
    z.object({ status: z.literal("fail"), code: z.literal(400), detail: z.string() }),
    z.object({ status: z.literal("fail"), code: z.literal(500), detail: z.string() }),
  ]),
]);
```

### Records, Maps, Sets, Literals

```ts
z.literal("admin");
z.record(z.string(), z.number());  // Record<string, number>
z.map(z.string(), z.object({ v: z.number() }));
z.set(z.string()).min(1).max(5);
```

### `z.interface()` — optional property control (Zod 4)

```ts
const UserProfile = z.interface({
  name: z.string(),
  bio:  z.optional(z.string()),  // bio?: string (truly optional, not string | undefined)
  age:  z.optional(z.number()),
});
```

### Nullables and Optionals

```ts
z.string().optional();  // string | undefined
z.string().nullable();  // string | null
z.string().nullish();   // string | null | undefined
```

---

## Schema Composition

```ts
const BaseUser = z.object({ id: z.string(), name: z.string() });

// .extend() — add or override fields
const AdminUser = BaseUser.extend({ role: z.literal("admin"), permissions: z.array(z.string()) });

// .merge() — combine two object schemas
const WithTimestamps = z.object({ createdAt: z.date(), updatedAt: z.date() });
const TimestampedUser = BaseUser.merge(WithTimestamps);

// .pick() / .omit() — select or exclude fields
const UserName = User.pick({ name: true, email: true });
const NoAge    = User.omit({ age: true });

// .partial() / .required() — toggle optionality
const UpdateUser = User.partial();            // all fields optional
const PatchUser  = User.partial({ age: true }); // only age optional

// .deepPartial() — recursive partial
const Config = z.object({ db: z.object({ host: z.string(), port: z.number() }) });
const PatchConfig = Config.deepPartial();
```

---

## Transforms and Pipes

```ts
// .transform() — map output to a new type
const StringToNumber = z.string().transform((val) => Number(val));
type Out = z.output<typeof StringToNumber>; // number

// .pipe() — transform then validate
const PositiveFromString = z.string()
  .transform((val) => Number(val))
  .pipe(z.number().positive());

// z.preprocess() — transform input before validation
const CoercedInt = z.preprocess(
  (val) => typeof val === "string" ? parseInt(val, 10) : val,
  z.number().int()
);

// z.coerce — built-in type coercion
const Port = z.coerce.number().int().min(1).max(65535);
Port.parse("3000"); // 3000

// z.stringbool() — env-style boolean coercion (Zod 4)
// "true"|"1"|"yes"|"on" → true; "false"|"0"|"no"|"off" → false
const FeatureFlag = z.stringbool();
```

### Codec pattern (Zod 4) — bidirectional transform

```ts
const DateCodec = z.codec(
  z.iso.datetime(),                     // parse: string -> Date
  (date: Date) => date.toISOString(),   // serialize: Date -> string
);
DateCodec.parse("2025-06-15T00:00:00Z"); // Date object
DateCodec.serialize(new Date());          // ISO string
```

---

## Error Handling

```ts
// .parse() throws ZodError on failure
try {
  Name.parse("");
} catch (err) {
  if (err instanceof z.ZodError) console.error(err.issues);
}

// .safeParse() never throws — returns discriminated result (prefer in production)
const result = Name.safeParse("");
if (!result.success) {
  console.error(result.error.issues);
  // [{ code: "too_small", minimum: 1, type: "string", message: "...", path: [] }]
} else {
  console.log(result.data);
}

// Formatting errors
const flat      = result.error.flatten();   // Zod 3: fieldErrors + formErrors
const formatted = result.error.format();    // Zod 3: nested object with _errors
const tree      = z.treeifyError(result.error); // Zod 4 (flatten/format deprecated)

// Custom messages
z.string().min(1, "Name is required");
z.string().email("Please enter a valid email");
// Zod 4: single `error` param
z.string({ error: (issue) => `Got ${typeof issue.input}, expected string` });

// Global error map
// Zod 3: z.setErrorMap(...)
// Zod 4: z.config({ customError: (issue) => ... })
```

---

## TypeScript Integration

```ts
// z.infer — extract the output type
type User = z.infer<typeof UserSchema>;

// z.input vs z.output — distinguish pre- and post-transform types
const DateFromString = z.string().transform((s) => new Date(s));
type In  = z.input<typeof DateFromString>;  // string
type Out = z.output<typeof DateFromString>; // Date (z.infer is an alias for z.output)

// Branded types — nominal typing at zero runtime cost
const UserId  = z.string().uuid().brand<"UserId">();
const OrderId = z.string().uuid().brand<"OrderId">();
type UserId   = z.infer<typeof UserId>;   // string & { __brand: "UserId" }

function getUser(id: UserId) { /* ... */ }
// getUser("raw-string");             // TS error — not branded
// getUser(OrderId.parse("..."));     // TS error — wrong brand
getUser(UserId.parse("550e8400-...")); // OK

// Generic schema functions
function validate<T extends z.ZodTypeAny>(schema: T, data: unknown): z.infer<T> | null {
  const result = schema.safeParse(data);
  return result.success ? result.data : null;
}
```

---

## API Validation Middleware

### Express middleware

```ts
function validateRequest(schemas: { body?: ZodSchema; params?: ZodSchema; query?: ZodSchema }) {
  return (req: Request, res: Response, next: NextFunction) => {
    for (const [key, schema] of Object.entries(schemas)) {
      const result = schema.safeParse(req[key as keyof typeof schemas]);
      if (!result.success) {
        return res.status(400).json({
          error: "Validation failed",
          field: key,
          issues: result.error.flatten().fieldErrors,
        });
      }
      (req as any)[key] = result.data; // replace with validated + stripped data
    }
    next();
  };
}

const CreateUserBody = z.object({
  name:  z.string().min(1).max(100),
  email: z.string().email(),
  age:   z.number().int().min(0).max(150).optional(),
});

app.post("/users", validateRequest({ body: CreateUserBody }), (req, res) => {
  const user = req.body as z.infer<typeof CreateUserBody>;
  res.json({ created: user });
});
```

### Next.js Server Action

```ts
export async function updateProfile(formData: FormData) {
  "use server";
  const result = UpdateProfileSchema.safeParse(Object.fromEntries(formData));
  if (!result.success) return { errors: result.error.flatten().fieldErrors };
  await db.profiles.update(result.data);
  return { success: true };
}
```

### tRPC

```ts
const appRouter = t.router({
  createUser: t.procedure
    .input(z.object({ name: z.string().min(1), email: z.string().email() }))
    .mutation(async ({ input }) => db.users.create(input)),
});
```

### Environment variable validation

```ts
const EnvSchema = z.object({
  NODE_ENV:     z.enum(["development", "staging", "production"]),
  PORT:         z.coerce.number().int().min(1).max(65535).default(3000),
  DATABASE_URL: z.string().url(),
  REDIS_URL:    z.string().url().optional(),
  LOG_LEVEL:    z.enum(["debug", "info", "warn", "error"]).default("info"),
  ENABLE_CACHE: z.stringbool().default("true"),  // Zod 4
});

export const env = EnvSchema.parse(process.env);
```

---

## Form Validation

### React Hook Form + Zod

```tsx
const SignupSchema = z.object({
  email:    z.string().email("Invalid email"),
  password: z.string()
    .min(8, "At least 8 characters")
    .regex(/[A-Z]/, "Must include an uppercase letter")
    .regex(/\d/, "Must include a digit"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords must match",
  path: ["confirmPassword"],
});

const { register, handleSubmit, formState: { errors } } = useForm<z.infer<typeof SignupSchema>>({
  resolver: zodResolver(SignupSchema),
});
```

### Shared schema (client + server)

```ts
// schemas/user.ts — shared between frontend and backend
export const CreateUserSchema = z.object({
  name:  z.string().min(1).max(100),
  email: z.string().email(),
  role:  z.enum(["admin", "editor", "viewer"]).default("viewer"),
});
export type CreateUserInput  = z.input<typeof CreateUserSchema>;   // form input type
export type CreateUserOutput = z.output<typeof CreateUserSchema>;  // validated output type
```

---

## Advanced Patterns

### Recursive schemas with z.lazy()

```ts
interface TreeNode { value: string; children: TreeNode[]; }

const TreeNodeSchema: z.ZodType<TreeNode> = z.object({
  value:    z.string(),
  children: z.lazy(() => z.array(TreeNodeSchema)),
});
```

### .check() — lightweight validation (Zod 4)

```ts
const PositiveInt = z.number().check((val) => val > 0, "Must be positive");
```

### .refine() and .superRefine()

```ts
const FutureDate = z.date().refine((d) => d > new Date(), "Date must be in the future");

// Cross-field validation
const DateRange = z.object({ start: z.date(), end: z.date() })
  .superRefine((data, ctx) => {
    if (data.end <= data.start) {
      ctx.addIssue({ code: z.ZodIssueCode.custom, message: "End must be after start", path: ["end"] });
    }
  });
```

### Async validation

```ts
const UniqueEmail = z.string().email().refine(
  async (email) => !(await db.users.findByEmail(email)),
  { message: "Email is already registered" },
);
const result = await UniqueEmail.safeParseAsync("user@example.com");
```

### Default values, catch, metadata

```ts
z.string().default("untitled");  // uses "untitled" when input is undefined
z.number().catch(0);              // uses 0 when validation fails

// Schema metadata (Zod 4)
const Price = z.number().positive().meta({ label: "Price", description: "In USD", example: 29.99 });

// JSON Schema conversion (Zod 4)
const jsonSchema = z.toJSONSchema(UserSchema);
```

---

## Performance

| Library | Relative speed | Notes |
|---------|---------------|-------|
| AJV | Fastest | JSON Schema compiled; best raw throughput |
| VineJS | Very fast | Purpose-built for Node.js |
| Zod 4 | Fast | 7x faster than Zod 3 for objects |
| Zod 3 | Moderate | Adequate for most apps |
| Joi | Moderate | Rich features; heavier |
| Yup | Slower | Form-validation focus |

**Performance tips:**

1. Define schemas at module scope — never inside request handlers.
2. Use `z.discriminatedUnion()` instead of `z.union()` when a discriminator key exists.
3. Validate at the boundary; trust typed data downstream.
4. Use `.strip()` (the default) unless `.passthrough()` or `.strict()` is required.
5. Use `.safeParse()` in production — avoids throw/catch overhead on invalid input.
6. For API hot paths (>10k req/s): consider AJV for the validation step.

**Bundle size:** Zod core is ~13 kB gzipped. Use Zod Mini (`zod/mini`) in Zod 4 for tree-shakeable functional API in client-heavy apps:

```ts
import * as z from "zod/mini";
const UserSchema = z.object({ name: z.pipe(z.string(), z.minLength(1)) });
```

---

## Anti-Patterns

### 1. Recreating schemas per request

```ts
// BAD — allocates a new schema object on every call
app.post("/users", (req, res) => {
  const schema = z.object({ name: z.string() }); // never do this inside a handler
  schema.parse(req.body);
});

// GOOD — define once at module scope
const CreateUserBody = z.object({ name: z.string() });
app.post("/users", (req, res) => { CreateUserBody.parse(req.body); });
```

### 2. Using .parse() without try/catch on untrusted input

```ts
// BAD — unhandled ZodError crashes the process
const data = schema.parse(untrustedInput);

// GOOD — use safeParse for untrusted data at boundaries
const result = schema.safeParse(untrustedInput);
if (!result.success) return res.status(400).json({ errors: result.error.flatten().fieldErrors });
```

### 3. Using .refine() for normalization instead of .transform()

```ts
// BAD — refine should validate, not mutate
z.string().refine((val) => { val = val.trim(); return val.length > 0; });

// GOOD — transform for normalization, then validate
z.string().trim().min(1);
```

### 4. Using z.infer when z.input is needed for forms

```ts
// WRONG: z.infer gives the OUTPUT type after transforms
// CORRECT: use z.input<> for pre-transform types (what forms and APIs send)
const DateSchema = z.string().transform((s) => new Date(s));
type FormValues = z.input<typeof DateSchema>;  // string
type Parsed     = z.output<typeof DateSchema>; // Date
```

### 5. Overusing .passthrough()

```ts
// BAD — unknown keys bypass validation
const LooseUser = z.object({ name: z.string() }).passthrough();

// GOOD — strip (default) or strict
const SafeUser = z.object({ name: z.string() }); // .strip() is default
```

---

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| `Expected string, received undefined` | Missing field or wrong path | Check the input object has the expected key; use `.optional()` if field can be absent |
| TypeScript "not assignable" error | Schema has transforms; using `z.infer` where `z.input` is needed | Use `z.input<typeof Schema>` for pre-transform types |
| `.refine()` inside `discriminatedUnion` fails | ZodEffects wrapping breaks discriminator lookup | Move `.refine()`/`.superRefine()` to after the `discriminatedUnion()` call |
| `safeParseAsync is not a function` | Using `.safeParse()` with an async refinement | Switch to `.safeParseAsync()` |
| Branded type not assignable | Passing raw unbranded data | Parse through the branded schema first |
| `z.coerce.number()` returns NaN | Input is non-numeric string | Add `.pipe(z.number().finite())` or use `z.preprocess` with explicit `parseInt` |
| `flatten()` / `format()` deprecated (Zod 4) | API changed in Zod 4 | Use `z.treeifyError(error)` instead |
| Bundle too large for client app | Full Zod import | Use Zod Mini (`zod/mini`) in Zod 4 |
| Infinite type / slow compilation with `z.lazy()` | Missing explicit type annotation | Add `const Schema: z.ZodType<YourType> = z.object(...)` |
| `.merge()` drops refinements | `.merge()` only combines shape, not effects | Use `.extend()` instead, or re-apply refinements after merge |

---

## References

- Official docs: https://zod.dev
- Zod 4 release notes: https://zod.dev/v4
- Zod 4 migration guide: https://zod.dev/v4/changelog
- GitHub: https://github.com/colinhacks/zod
- @hookform/resolvers (zodResolver): https://github.com/react-hook-form/resolvers
- tRPC + Zod: https://trpc.io/docs
- Benchmark data: https://github.com/vinejs/vine/blob/develop/benchmarks.md
