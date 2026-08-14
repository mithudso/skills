<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `react-nextjs` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: react-nextjs
title: "React & Next.js Expert"
version: "1.1.0"
updated: "2026-05-29"
description: >
  React 19 and Next.js 15+ expert reference. TRIGGER: user is building, reviewing, debugging,
  or architecting a React or Next.js application; asks about Server Components, App Router,
  server actions, React Compiler, streaming SSR, partial prerendering, state management
  (Zustand, Jotai), data fetching/caching, new hooks (useActionState, useOptimistic, use,
  useEffectEvent), Activity component, testing (Vitest, React Testing Library, Playwright),
  performance optimization, or Pages Router to App Router migration. SKIP: pure CSS/visual
  design (use frontend-design), HTML/CSS reference (use html-css), non-React frameworks
  (Vue, Svelte, Angular).
category: developer
tags:
  - react
  - react-19
  - nextjs
  - next-js-15
  - server-components
  - app-router
  - server-actions
  - react-compiler
  - streaming-ssr
  - partial-prerendering
  - zustand
  - jotai
  - state-management
  - data-fetching
  - caching
  - performance
  - testing
  - typescript
  - frontend
keywords:
  - React 19
  - Next.js 15
  - React Server Components
  - RSC
  - App Router
  - server actions
  - use server
  - use client
  - React Compiler
  - automatic memoization
  - streaming SSR
  - partial prerendering
  - PPR
  - Zustand
  - Jotai
  - TanStack Query
  - useActionState
  - useOptimistic
  - useEffectEvent
  - use API
  - Activity component
  - Suspense
  - caching
  - revalidation
  - Vitest
  - React Testing Library
  - Playwright
  - performance optimization
  - Web Vitals
whenToUse:
  - "Building or scaffolding a React or Next.js application"
  - "Choosing between server and client components"
  - "Implementing server actions or data mutations"
  - "Configuring React Compiler for automatic memoization"
  - "Designing data fetching, caching, or revalidation strategies"
  - "Choosing a state management library (Zustand, Jotai, Redux Toolkit)"
  - "Setting up streaming SSR or partial prerendering"
  - "Writing tests for React or Next.js applications"
  - "Debugging React performance issues or hydration mismatches"
  - "Reviewing React or Next.js code for anti-patterns"
  - "Migrating from Pages Router to App Router"
related_skills:
  - frontend-ui
  - testing-and-vitest-expert
context_file: SKILL.md
---

# React & Next.js Expert

Comprehensive reference for building production React 19 and Next.js 15+ applications.

## When NOT to use this skill

- **Pure CSS/visual design:** use `frontend-design`.
- **HTML/CSS element reference:** use `html-css`.
- **Accessibility audits:** use `accessibility-ux-reviewer`.
- **Non-React frameworks (Vue, Svelte, Angular):** this skill does not cover them.

---

## 1. React 19 Core Features

### 1.1 Release timeline

| Version | Date | Highlights |
|---------|------|------------|
| React 19.0 | December 2024 | Server Components stable, Actions, new hooks |
| React 19.1 | June 2025 | Bug fixes, stability improvements |
| React 19.2 | October 2025 | Activity component, useEffectEvent, React Compiler 1.0 |

### 1.2 Server Components (RSC)

Server Components render on the server and send serialized UI to the client. They never ship JavaScript to the browser.

**Use Server Components (default in App Router) when:**
- Fetching data from databases or APIs
- Accessing backend resources (file system, internal services)
- Keeping secrets server-side (API keys, tokens)
- Reducing client bundle size
- SEO-critical content

**Use Client Components (`'use client'`) when:**
- Event handlers (onClick, onChange, onSubmit)
- Browser APIs (localStorage, geolocation, IntersectionObserver)
- State and lifecycle (useState, useEffect, useRef)
- Interactive UI (modals, dropdowns, drag-and-drop)

```tsx
// ServerComponent.tsx -- NO directive needed (default)
export default async function ProductPage({ id }: { id: string }) {
  const product = await db.product.findUnique({ where: { id } });
  return (
    <article>
      <h1>{product.name}</h1>
      <ClientButton productId={id} />
    </article>
  );
}

// ClientButton.tsx
'use client';
export default function ClientButton({ productId }: { productId: string }) {
  const [added, setAdded] = useState(false);
  return <button onClick={() => setAdded(true)}>Add to Cart</button>;
}
```

**Serialization constraint:** Props from Server → Client must be JSON-safe. No functions, classes, or Dates — use ISO strings.

**Composition pattern:**
```tsx
'use client';
export function InteractiveWrapper({ children }: { children: React.ReactNode }) {
  const [open, setOpen] = useState(false);
  return <div>{open && children}</div>; // children can be Server Components
}
```

### 1.3 Actions

**Server Actions (`'use server'`):**
```tsx
// actions.ts
'use server';
import { revalidatePath } from 'next/cache';
import { redirect } from 'next/navigation';

export async function createPost(formData: FormData) {
  const title = formData.get('title') as string;
  await db.post.create({ data: { title } });
  revalidatePath('/posts');
  redirect('/posts');
}
```

**Security note:** Server Actions create POST endpoints. Always validate and authorize inputs. Never trust FormData blindly.

### 1.4 New hooks and APIs

#### useActionState

```tsx
'use client';
export default function PostForm() {
  const [state, formAction, isPending] = useActionState(createPost, null);
  return (
    <form action={formAction}>
      <input name="title" required />
      {state?.error && <p className="error">{state.error}</p>}
      <button disabled={isPending}>{isPending ? 'Creating...' : 'Create Post'}</button>
    </form>
  );
}
```

#### useOptimistic

```tsx
'use client';
export function TodoList({ todos }: { todos: Todo[] }) {
  const [optimisticTodos, addOptimisticTodo] = useOptimistic(
    todos,
    (state, newTodo: Todo) => [...state, newTodo]
  );

  async function addTodo(formData: FormData) {
    const title = formData.get('title') as string;
    addOptimisticTodo({ id: crypto.randomUUID(), title, completed: false });
    await createTodo(formData);
  }

  return (
    <>
      <ul>{optimisticTodos.map(t => <li key={t.id}>{t.title}</li>)}</ul>
      <form action={addTodo}><input name="title" /><button>Add</button></form>
    </>
  );
}
```

#### use() API

Reads a Promise or Context during render. Can be called conditionally (unlike other hooks).

```tsx
function UserProfile({ userPromise }: { userPromise: Promise<User> }) {
  const user = use(userPromise); // Suspends until resolved
  return <h1>{user.name}</h1>;
}

// Context replacement
function ThemedButton() {
  const theme = use(ThemeContext); // works inside conditionals
  return <button style={{ background: theme.primary }}>Click</button>;
}
```

#### useFormStatus

```tsx
'use client';
function SubmitButton() {
  const { pending } = useFormStatus();
  return <button disabled={pending}>{pending ? 'Saving...' : 'Save'}</button>;
}
```

#### useEffectEvent (React 19.2)

Separates event logic from Effect dependencies.

```tsx
function ChatRoom({ roomId, theme }: { roomId: string; theme: string }) {
  const onConnected = useEffectEvent(() => {
    showNotification('Connected!', theme); // always sees latest theme
  });

  useEffect(() => {
    const conn = createConnection(roomId);
    conn.on('connected', onConnected);
    conn.connect();
    return () => conn.disconnect();
  }, [roomId]); // theme NOT a dependency
}
```

**Rule:** Effect Events must NOT appear in dependency arrays.

### 1.5 Activity Component (React 19.2)

State-preserving hide/show — replaces conditional rendering for tabs and pre-loading.

```tsx
function TabContainer({ activeTab }: { activeTab: string }) {
  return (
    <>
      <Activity mode={activeTab === 'home' ? 'visible' : 'hidden'}>
        <HomePage />
      </Activity>
      <Activity mode={activeTab === 'profile' ? 'visible' : 'hidden'}>
        <ProfilePage />
      </Activity>
    </>
  );
}
```

- `visible`: renders, mounts effects, processes updates normally
- `hidden`: hides, unmounts effects, defers updates until idle

---

## 2. React Compiler (v1.0, October 2025)

Automatically memoizes components and hooks at build time. Eliminates manual `useMemo`, `useCallback`, `React.memo` in most cases.

**Performance gains at Meta:** 20–30% reduction in render time, up to 12% faster initial loads, 2.5× faster interactions.

### Installation

```bash
npm install --save-dev --save-exact babel-plugin-react-compiler@latest
```

```ts
// next.config.ts
const nextConfig = { experimental: { reactCompiler: true } };
```

```ts
// vite.config.ts
import { reactCompiler } from 'babel-plugin-react-compiler';
export default defineConfig({
  plugins: [react({ babel: { plugins: [reactCompiler] } })],
});
```

### Compatibility

- Works with React 17, 18, and 19+; React Native; existing `useMemo`/`useCallback` (become no-ops)
- Build tools: Babel, Vite, Next.js, Rsbuild (swc support experimental)
- New projects (Expo 54+, create-vite, create-next-app): enabled by default

---

## 3. Next.js 15 App Router

### 3.1 File-system routing

```
app/
  layout.tsx          # Root layout (required)
  page.tsx            # Home page /
  loading.tsx         # Auto-wrapped in Suspense
  error.tsx           # Error boundary
  not-found.tsx       # 404
  template.tsx        # Re-mounted on each navigation
  blog/[slug]/page.tsx
  (marketing)/about/page.tsx   # Route group (no URL segment)
  @modal/default.tsx           # Parallel route
  api/route.ts                 # Route handler
```

### 3.2 Layouts and templates

```tsx
// app/layout.tsx (required)
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return <html lang="en"><body>{children}</body></html>;
}
```

Nested layouts persist across navigations. Use `template.tsx` for a fresh instance on every navigation.

### 3.3 Route handlers

```tsx
// app/api/users/route.ts
export async function GET(request: NextRequest) {
  return NextResponse.json(await db.user.findMany());
}

export async function POST(request: NextRequest) {
  const body = await request.json();
  return NextResponse.json(await db.user.create({ data: body }), { status: 201 });
}
```

### 3.4 Middleware

```tsx
// middleware.ts (project root)
export function middleware(request: NextRequest) {
  if (!request.cookies.get('session')) {
    return NextResponse.redirect(new URL('/login', request.url));
  }
  return NextResponse.next();
}

export const config = { matcher: ['/dashboard/:path*', '/api/:path*'] };
```

---

## 4. Server Actions deep dive

### Validation with Zod (always required)

```tsx
'use server';
import { z } from 'zod';
import { revalidateTag } from 'next/cache';

const schema = z.object({
  title: z.string().min(1).max(200),
  content: z.string().min(1),
});

export async function createArticle(prevState: any, formData: FormData) {
  const parsed = schema.safeParse({
    title: formData.get('title'),
    content: formData.get('content'),
  });
  if (!parsed.success) return { error: parsed.error.flatten().fieldErrors };
  await db.article.create({ data: parsed.data });
  revalidateTag('articles');
  return { success: true };
}
```

### Security checklist

- Server Actions create public HTTP POST endpoints
- Always verify authentication and authorization
- Rate-limit sensitive actions
- Validate all inputs with Zod
- Never expose raw database errors to the client

---

## 5. Data fetching and caching

### 5.1 Server Component data fetching

```tsx
export default async function ProductsPage() {
  const products = await db.product.findMany({ orderBy: { createdAt: 'desc' }, take: 20 });
  return <ProductList products={products} />;
}

// Parallel fetching
export default async function Dashboard() {
  const [users, orders] = await Promise.all([getUsers(), getOrders()]);
  return <><UserStats users={users} /><OrderTable orders={orders} /></>;
}
```

### 5.2 Cache layers

| Layer | Location | Invalidation |
|-------|----------|-------------|
| Data Cache | Server | `revalidateTag`, `revalidatePath`, time-based |
| Full Route Cache | Server | `revalidatePath`, redeployment |
| Router Cache | Client | `router.refresh()`, cookie changes |

### 5.3 Fetch caching

```tsx
const data = await fetch(url);                                    // Cached (static)
const data = await fetch(url, { next: { revalidate: 3600 } });   // ISR — hourly
const data = await fetch(url, { cache: 'no-store' });             // Always fresh
const data = await fetch(url, { next: { tags: ['products'] } });  // Tagged
```

### 5.4 unstable_cache for non-fetch data

```tsx
const getCachedProducts = unstable_cache(
  async (category: string) => db.product.findMany({ where: { category } }),
  ['products-by-category'],
  { tags: ['products'], revalidate: 3600 }
);
```

**Rule:** Use `unstable_cache` for reads only. Call `revalidateTag`/`revalidatePath` only inside Server Actions or Route Handlers.

### 5.5 TanStack Query for client-side data

```tsx
'use client';
function ProductList() {
  const { data, isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => fetch('/api/products').then(r => r.json()),
    staleTime: 60_000,
  });

  const mutation = useMutation({
    mutationFn: (p) => fetch('/api/products', { method: 'POST', body: JSON.stringify(p) }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['products'] }),
  });
}
```

---

## 6. Streaming and Suspense

```tsx
export default function Page() {
  return (
    <main>
      <StaticHeader />
      <Suspense fallback={<ChartSkeleton />}><SlowChart /></Suspense>
      <Suspense fallback={<TableSkeleton />}><SlowTable /></Suspense>
    </main>
  );
}
```

**Performance rule:** the outermost Suspense boundary determines first meaningful paint. Keep the outer shell lightweight.

`loading.tsx` automatically wraps `page.tsx` in a Suspense boundary.

---

## 7. Partial Prerendering (PPR)

Hybrid rendering: static shell served from CDN edge (<50ms TTFB), dynamic regions stream from server.

```ts
// next.config.ts
const nextConfig = { experimental: { ppr: 'incremental' } };
```

```tsx
// app/product/[id]/page.tsx
export const experimental_ppr = true;

export default async function ProductPage({ params }: { params: { id: string } }) {
  const product = await getProduct(params.id); // static at build
  return (
    <article>
      <h1>{product.name}</h1>
      <Suspense fallback={<PriceSkeleton />}><DynamicPrice productId={params.id} /></Suspense>
    </article>
  );
}
```

| Aspect | Streaming SSR | PPR |
|--------|---------------|-----|
| Shell generation | Per-request | Pre-built at build time (CDN) |
| TTFB | Server-dependent | < 50ms (CDN hit) |
| Best for | Fully dynamic pages | Pages with static + dynamic mix |

---

## 8. State management

### Decision framework

```
Server data (API/DB)? → TanStack Query / Server Component fetch
Form state?           → useActionState + useFormStatus
Local to one component? → useState / useReducer
Many independent atoms? → Jotai
One coherent store?     → Zustand
Enterprise strict patterns? → Redux Toolkit
Small, rarely-changing? → React Context
```

### Zustand (~3KB)

```tsx
export const useCartStore = create<CartState>()(
  devtools(persist(
    (set, get) => ({
      items: [],
      addItem: (item) => set((state) => ({ items: [...state.items, item] })),
      totalPrice: () => get().items.reduce((sum, i) => sum + i.price, 0),
    }),
    { name: 'cart-storage' }
  ))
);

// Fine-grained selector — prevents re-render on unrelated store changes
function CartCount() {
  const count = useCartStore((state) => state.items.length);
  return <span>{count}</span>;
}
```

### Jotai (~4KB)

```tsx
export const themeAtom = atomWithStorage('theme', 'light');
export const isDarkAtom = atom((get) => get(themeAtom) === 'dark');
export const toggleThemeAtom = atom(
  (get) => get(themeAtom),
  (get, set) => set(themeAtom, get(themeAtom) === 'light' ? 'dark' : 'light')
);
```

### Comparison

| Feature | Zustand | Jotai | Redux Toolkit | React Context |
|---------|---------|-------|---------------|---------------|
| Bundle | ~3KB | ~4KB | ~15KB | 0 |
| Architecture | Store-based | Atom-based | Store + slices | Provider tree |
| Outside React | Yes | No | Yes | No |
| Best for | Global stores | Granular atoms | Enterprise | Small stable state |

---

## 9. Performance optimization

### Core Web Vitals targets

| Metric | Good | Poor |
|--------|------|------|
| LCP | < 2.5s | > 4.0s |
| INP | < 200ms | > 500ms |
| CLS | < 0.1 | > 0.25 |

### Image optimization

```tsx
<Image
  src="/hero.jpg" alt="Hero"
  width={1200} height={600}
  priority          // Preload LCP images
  placeholder="blur"
  sizes="(max-width: 768px) 100vw, 50vw"
/>
```

### Font optimization

```tsx
const inter = Inter({ subsets: ['latin'], display: 'swap', variable: '--font-inter' });
```

### Performance checklist

1. Use Server Components by default — minimize `'use client'` boundaries
2. Enable React Compiler — automatic memoization
3. `next/image` for all images
4. `next/font` — zero layout shift
5. Suspense boundaries for streaming
6. `loading.tsx` for route-level loading
7. Dynamic import heavy client components
8. `Promise.all()` for parallel data fetching
9. PPR for mixed static/dynamic pages

---

## 10. Testing

### Stack

| Layer | Tool | What it catches |
|-------|------|----------------|
| Unit | Vitest + Testing Library | Components, hooks, utilities, server actions |
| Integration | Vitest + MSW | Component + API interactions |
| E2E | Playwright | Full user flows, auth, forms |
| Visual | Storybook + Chromatic | Appearance, design system |

### Vitest setup

```ts
// vitest.config.ts
export default defineConfig({
  plugins: [react()],
  test: { environment: 'jsdom', globals: true, setupFiles: ['./test/setup.ts'] },
  resolve: { alias: { '@': './src' } },
});
```

### Testing components

```tsx
describe('Counter', () => {
  it('increments on click', () => {
    render(<Counter initialCount={0} />);
    fireEvent.click(screen.getByRole('button', { name: /increment/i }));
    expect(screen.getByText('Count: 1')).toBeInTheDocument();
  });
});
```

### Playwright Page Object Model

```ts
export class LoginPage {
  constructor(readonly page: Page) {}
  readonly emailInput = this.page.getByLabel('Email');
  readonly submitButton = this.page.getByRole('button', { name: 'Sign In' });
  async login(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.page.getByLabel('Password').fill(password);
    await this.submitButton.click();
  }
}
```

---

## 11. Anti-patterns

### Component

| Anti-pattern | Fix |
|---|---|
| `'use client'` on every component | Only add where interactivity is needed |
| Fetching in Client Components when server can | Move to Server Components |
| `'use client'` at tree root | Push boundaries as deep as possible |
| `useEffect` for data fetching | Server Components or TanStack Query |
| Server data in Zustand/Jotai | TanStack Query / SWR |

### Data fetching

| Anti-pattern | Fix |
|---|---|
| Sequential `await` in Server Components | `Promise.all()` |
| No cache tags on fetches | `next: { tags: ['entity'] }` |
| `revalidatePath` inside a read function | Only inside Server Actions / Route Handlers |
| `cache: 'no-store'` everywhere | Only for truly dynamic user-specific data |
| Slow components not wrapped in Suspense | Add Suspense boundaries |

### Security

| Anti-pattern | Fix |
|---|---|
| No input validation in Server Actions | Always validate with Zod |
| Raw database errors exposed to client | Return generic messages, log server-side |
| Client-side auth checks only | Verify in middleware + Server Actions |

---

## 12. Pages Router to App Router migration

### Step-by-step

1. Create `app/` alongside `pages/` — both routers work simultaneously
2. Move `_app.tsx` logic to `app/layout.tsx`
3. Move `_document.tsx` into root layout `<html>`/`<body>`
4. Migrate pages one at a time
5. Replace `getServerSideProps` with async Server Components
6. Replace `getStaticProps` with fetch + revalidate or `generateStaticParams`
7. Replace `useRouter` from `next/router` with `next/navigation`
8. Replace `next/head` with Metadata API

### Data fetching migration

```tsx
// BEFORE (Pages Router)
export async function getServerSideProps() {
  return { props: { data: await fetchData() } };
}

// AFTER (App Router)
export default async function Page() {
  const data = await fetchData();
  return <div>{data.title}</div>;
}
```

---

## 13. Project structure

```
src/
  app/                    # App Router pages and layouts
    (auth)/login/page.tsx
    (dashboard)/layout.tsx
    api/                  # Route handlers
  components/
    ui/                   # Primitive UI: Button, Input, Card
    features/             # Feature-specific components
  lib/                    # Shared utilities, db client, auth helpers
  hooks/                  # Custom React hooks
  stores/                 # Zustand stores, Jotai atoms
  actions/                # Server Actions
  types/                  # TypeScript types
```

**Naming conventions:**
- Components: PascalCase (`ProductCard.tsx`)
- Hooks: `use` prefix camelCase (`useDebounce.ts`)
- Actions: verb-noun camelCase (`createProduct.ts`)
- Stores/atoms: `cartStore.ts`, `themeAtom.ts`

---

## 14. TypeScript integration

```tsx
// Async params in Next.js 15
export default async function ProductPage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ sort?: string }>;
}) {
  const { id } = await params;
  const { sort } = await searchParams;
}

// Server Action types
type ActionState = { success: boolean; error?: string; fieldErrors?: Record<string, string[]> };

export async function createUser(prevState: ActionState, formData: FormData): Promise<ActionState> {
  // ...
}
```

---

## 15. Environment and configuration

```
# .env.local
DATABASE_URL=postgresql://...          # Server-only
NEXT_PUBLIC_API_URL=https://api.example.com  # Browser-accessible (NEXT_PUBLIC_ prefix)
```

```ts
// next.config.ts
const nextConfig: NextConfig = {
  experimental: { reactCompiler: true, ppr: 'incremental' },
  images: { remotePatterns: [{ protocol: 'https', hostname: 'cdn.example.com' }] },
};
```

---

## References

- [React Documentation](https://react.dev/)
- [React 19.2 Release](https://react.dev/blog/2025/10/01/react-19-2)
- [React Compiler 1.0](https://react.dev/blog/2025/10/07/react-compiler-1)
- [Next.js Documentation](https://nextjs.org/docs)
- [Next.js Caching Guide](https://nextjs.org/docs/app/getting-started/caching-and-revalidating)
- [Next.js Server Actions](https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations)
- [Zustand](https://github.com/pmndrs/zustand)
- [Jotai](https://github.com/pmndrs/jotai)
- [TanStack Query](https://tanstack.com/query)
- [Vitest](https://vitest.dev/)
- [Playwright](https://playwright.dev/)
- [Web Vitals](https://web.dev/vitals/)
