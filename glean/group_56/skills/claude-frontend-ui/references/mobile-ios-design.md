<!-- hub-reference-banner -->
> **Reference file — part of the `frontend-ui` hub.** Formerly the standalone `mobile-ios-design` skill.
> Sibling topics in this family are now reference files under the hubs (`frontend-ui`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: mobile-ios-design
version: "1.1.0"
updated: "2026-05-29"
description: >
  iOS Human Interface Guidelines and SwiftUI patterns for building native iOS and iPadOS apps.
  TRIGGER: user is designing an iOS/iPadOS interface, implementing SwiftUI views or navigation,
  working with SF Symbols, Dynamic Type, Dark Mode, system materials, or ensuring HIG compliance.
  SKIP: Android/Material Design (use ui-ux-pro-max), Compose Multiplatform KMP (use
  compose-multiplatform-patterns), React Native cross-platform (use ui-ux-pro-max or
  frontend-design), web UI (use frontend-design).
whenToUse:
  - "Designing iOS app interface following Apple HIG"
  - "Building SwiftUI views and layouts"
  - "iOS navigation patterns — NavigationStack, TabView, sheets"
  - "Adaptive layouts for iPhone and iPad"
  - "SF Symbols usage and rendering modes"
  - "Supporting Dynamic Type and Dark Mode"
  - "iOS-specific gestures and interactions"
  - "Safe area insets and system chrome avoidance"
  - "SwiftUI @Observable, SwiftData, or async/await patterns"
  - "NavigationSplitView for iPad multitasking"
tags:
  - ios
  - swiftui
  - apple-hig
  - ipad
  - sf-symbols
  - dynamic-type
  - dark-mode
  - navigation
  - accessibility
related_skills:
  - compose-multiplatform-patterns
  - ui-ux-pro-max
  - accessibility-ux-reviewer
---

# iOS Mobile Design

iOS Human Interface Guidelines (HIG) and SwiftUI patterns for building polished, native iOS
and iPadOS applications.

## When NOT to use this skill

- **Android or cross-platform Material Design:** use `ui-ux-pro-max`.
- **Kotlin Multiplatform with Compose for iOS:** use `compose-multiplatform-patterns`.
- **React Native:** use `ui-ux-pro-max` or `frontend-design`.
- **Web UI:** use `frontend-design` or `html-css`.

## HIG principles

| Principle | What it means in practice |
|-----------|--------------------------|
| **Clarity** | Legible text, precise icons, subtle adornments; content is primary |
| **Deference** | UI helps users understand content without competing with it |
| **Depth** | Visual layers and motion convey hierarchy and enable navigation |

**Platform targets:**
- **iOS** — touch-first, compact displays, portrait-primary
- **iPadOS** — larger canvas, multitasking (split view, slide over), pointer support
- **visionOS** — spatial computing, eye/hand input (out of scope for most patterns here)

## Layout system

### Stack-based layouts

```swift
// Vertical stack with alignment
VStack(alignment: .leading, spacing: 12) {
    Text("Title")
        .font(.headline)
    Text("Subtitle")
        .font(.subheadline)
        .foregroundStyle(.secondary)
}

// Horizontal stack with flexible spacing
HStack {
    Image(systemName: "star.fill")
    Text("Featured")
    Spacer()
    Text("View All")
        .foregroundStyle(.tint)
}
```

### Grid layouts

```swift
// Adaptive grid — fills available width
LazyVGrid(columns: [GridItem(.adaptive(minimum: 150, maximum: 200))], spacing: 16) {
    ForEach(items) { item in ItemCard(item: item) }
}

// Fixed 3-column grid
LazyVGrid(columns: Array(repeating: GridItem(.flexible()), count: 3), spacing: 12) {
    ForEach(items) { item in ItemThumbnail(item: item) }
}
```

**Performance:** use `LazyVStack`/`LazyHStack` for long scrolling lists to avoid rendering
off-screen views.

## Navigation patterns

### NavigationStack (iOS 16+)

Type-safe, value-driven navigation. Use `NavigationPath` for programmatic control.

```swift
struct ContentView: View {
    @State private var path = NavigationPath()

    var body: some View {
        NavigationStack(path: $path) {
            List(items) { item in
                NavigationLink(value: item) {
                    ItemRow(item: item)
                }
            }
            .navigationTitle("Items")
            .navigationDestination(for: Item.self) { item in
                ItemDetailView(item: item)
            }
        }
    }
}
```

### NavigationSplitView (iPad and Mac)

Use for two- or three-column layouts on larger screens.

```swift
NavigationSplitView {
    List(items, selection: $selectedItem) { item in
        Text(item.title).tag(item)
    }
} detail: {
    if let item = selectedItem {
        ItemDetailView(item: item)
    } else {
        Text("Select an item").foregroundStyle(.secondary)
    }
}
```

### TabView (iOS 18+)

```swift
TabView(selection: $selectedTab) {
    Tab("Home", systemImage: "house", value: 0) { HomeView() }
    Tab("Search", systemImage: "magnifyingglass", value: 1) { SearchView() }
    Tab("Profile", systemImage: "person", value: 2) { ProfileView() }
}
```

### Sheets and modals

```swift
.sheet(isPresented: $showingDetail) {
    DetailSheet()
        .presentationDetents([.medium, .large])
        .presentationDragIndicator(.visible)
}
```

## System integration

### SF Symbols

```swift
// Basic symbol
Image(systemName: "heart.fill").foregroundStyle(.red)

// Multicolor rendering
Image(systemName: "cloud.sun.fill").symbolRenderingMode(.multicolor)

// Variable value (iOS 16+) — signal strength, volume, etc.
Image(systemName: "speaker.wave.3.fill", variableValue: volume)

// Symbol effect (iOS 17+)
Image(systemName: "bell.fill")
    .symbolEffect(.bounce, value: notificationCount)
```

### Dynamic Type

Always use semantic font styles — they automatically scale with user accessibility settings.

```swift
Text("Headline").font(.headline)
Text("Body text").font(.body)
Text("Caption").font(.caption)

// Custom font that respects Dynamic Type
Text("Custom").font(.custom("Avenir", size: 17, relativeTo: .body))
```

**Rule:** never use fixed `font(.system(size: 17))` for body text — it won't scale with
Dynamic Type.

### @Observable (iOS 17+)

Prefer `@Observable` over `ObservableObject` for new code — simpler syntax, better
compiler optimization.

```swift
@Observable
class ItemStore {
    var items: [Item] = []
    var isLoading = false

    func load() async {
        isLoading = true
        items = await fetchItems()
        isLoading = false
    }
}

struct ItemListView: View {
    @State private var store = ItemStore()

    var body: some View {
        List(store.items) { item in ItemRow(item: item) }
            .task { await store.load() }
    }
}
```

## Visual design

### Semantic colors and materials

Use semantic colors that automatically adapt to light/dark mode and accessibility settings.

```swift
// Semantic foreground colors
Text("Primary").foregroundStyle(.primary)
Text("Secondary").foregroundStyle(.secondary)
Text("Tertiary").foregroundStyle(.tertiary)

// System materials for blur effects
Rectangle().fill(.ultraThinMaterial)

// Background with shape
Text("Label")
    .padding()
    .background(.regularMaterial, in: RoundedRectangle(cornerRadius: 12))
```

**Rule:** never hardcode `Color(hex: "#1a1a1a")` for text or backgrounds — use semantic
tokens or asset catalog colors with light/dark variants.

### Shadows and depth

```swift
// Card shadow — subtle, directional
RoundedRectangle(cornerRadius: 16)
    .fill(.background)
    .shadow(color: .black.opacity(0.08), radius: 8, y: 4)

// Layered depth
.shadow(radius: 2, y: 1)
.shadow(radius: 8, y: 4)
```

### Safe areas

```swift
// Avoid content under notch/home indicator
.safeAreaInset(edge: .bottom) {
    BottomBar()
}

// Ignore safe area for full-bleed background only
ZStack {
    Color.blue.ignoresSafeArea()
    VStack { /* content respects safe area */ }
}
```

## Quick start component

```swift
import SwiftUI

struct FeatureCard: View {
    let title: String
    let description: String
    let systemImage: String

    var body: some View {
        HStack(spacing: 16) {
            Image(systemName: systemImage)
                .font(.title2)
                .foregroundStyle(.tint)
                .frame(width: 44, height: 44)
                .background(.tint.opacity(0.12), in: Circle())

            VStack(alignment: .leading, spacing: 4) {
                Text(title).font(.headline)
                Text(description)
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
                    .lineLimit(2)
            }

            Spacer()

            Image(systemName: "chevron.right")
                .foregroundStyle(.tertiary)
                .imageScale(.small)
        }
        .padding()
        .background(.background, in: RoundedRectangle(cornerRadius: 12))
        .shadow(color: .black.opacity(0.05), radius: 4, y: 2)
    }
}
```

## Best practices

| Practice | Rule |
|----------|------|
| Semantic colors | Always use `.primary`, `.secondary`, `.background`, `.tint` for automatic light/dark adaptation |
| SF Symbols | Prefer system symbols — consistent weight, automatic accessibility, Dynamic Type scaling |
| Dynamic Type | Use semantic fonts (`.body`, `.headline`, `.caption`) — never fixed sizes for readable text |
| Accessibility | Add `.accessibilityLabel()` and `.accessibilityHint()` on interactive elements; use `.accessibilityElement(children: .combine)` for complex views |
| Safe areas | Respect `safeAreaInset`; never hardcode padding at screen edges |
| State restoration | Use `@SceneStorage` to persist user state across app launches |
| iPad multitasking | Design for split view and slide over; use `NavigationSplitView` on larger screens |
| Async operations | Use `.task {}` modifier for view-lifecycle-bound async work; cancels automatically on disappear |

## Common issues and fixes

| Issue | Fix |
|-------|-----|
| Layout breaks on text size increase | Use `.fixedSize()` sparingly; prefer flexible layouts and `.lineLimit(nil)` |
| Performance lag in scrolling lists | Use `LazyVStack`/`LazyHStack`; add `id:` to `ForEach` |
| Navigation back gesture broken | Ensure `NavigationLink` values conform to `Hashable` |
| Dark mode color mismatch | Never hardcode colors; use semantic colors or asset catalog with appearance variants |
| VoiceOver reads wrong order | Set `.accessibilitySortPriority()` or restructure view hierarchy |
| Memory leak in async closures | Use `[weak self]` or `@MainActor` isolated properties to avoid strong captures |
| Content hidden behind tab bar | Use `.safeAreaInset(edge: .bottom)` or let `List`/`ScrollView` handle insets automatically |
