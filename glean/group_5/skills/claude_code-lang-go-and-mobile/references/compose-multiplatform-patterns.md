<!-- hub-reference-banner -->
> **Reference file — part of the `lang-go-and-mobile` hub.** Formerly the standalone `compose-multiplatform-patterns` skill.
> Sibling topics in this family are now reference files under the hubs (`lang-python`, `lang-js-ts`, `lang-go-and-mobile`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: compose-multiplatform-patterns
version: "1.1.0"
updated: "2026-05-29"
description: >
  Compose Multiplatform and Jetpack Compose patterns for KMP projects — state management,
  navigation, composable design, platform-specific UI, and performance optimization.
  TRIGGER: user is building Compose UI (Jetpack Compose or Compose Multiplatform), managing
  ViewModel state with StateFlow, implementing Navigation in Android/KMP, designing reusable
  composables, or optimizing recomposition. SKIP: pure Swift/iOS UI (use mobile-ios-design),
  React/web UI (use frontend-design or react-nextjs), non-Compose Android views.
origin: ECC
whenToUse:
  - "Building Compose UI or Jetpack Compose screens"
  - "State management with ViewModel and StateFlow in Compose"
  - "Type-safe navigation in Compose Navigation 2.8+"
  - "Designing reusable composables with slot-based APIs"
  - "Optimizing recomposition and Compose performance"
  - "KMP platform-specific UI with expect/actual"
  - "Material 3 theming with dynamic color"
  - "LazyColumn performance with stable keys"
tags:
  - kotlin
  - compose
  - jetpack-compose
  - compose-multiplatform
  - kmp
  - android
  - state-management
  - navigation
  - material3
related_skills:
  - mobile-ios-design
  - typescript-expert
---

# Compose Multiplatform Patterns

Patterns for building shared UI with Compose Multiplatform and Jetpack Compose across Android,
iOS, desktop, and web. Covers state management, navigation, composable design, and performance.

## When NOT to use this skill

- **Pure iOS native UI:** use `mobile-ios-design` for SwiftUI and HIG patterns.
- **React/web UI:** use `frontend-design` or `react-nextjs`.
- **Non-Compose Android Views (XML layouts):** this skill does not cover legacy View system.

## State management

### ViewModel + single state object

Use one data class per screen. Expose as `StateFlow`; collect with `collectAsStateWithLifecycle`.

```kotlin
data class ItemListState(
    val items: List<Item> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val searchQuery: String = ""
)

class ItemListViewModel(
    private val getItems: GetItemsUseCase
) : ViewModel() {
    private val _state = MutableStateFlow(ItemListState())
    val state: StateFlow<ItemListState> = _state.asStateFlow()

    fun onSearch(query: String) {
        _state.update { it.copy(searchQuery = query) }
        loadItems(query)
    }

    private fun loadItems(query: String) {
        viewModelScope.launch {
            _state.update { it.copy(isLoading = true) }
            getItems(query).fold(
                onSuccess = { items -> _state.update { it.copy(items = items, isLoading = false) } },
                onFailure = { e -> _state.update { it.copy(error = e.message, isLoading = false) } }
            )
        }
    }
}
```

### Collecting state in Compose

Split screen composables into a stateful container (collects state) and a stateless content
composable (receives data as params). The stateless version is easy to preview and test.

```kotlin
@Composable
fun ItemListScreen(viewModel: ItemListViewModel = koinViewModel()) {
    val state by viewModel.state.collectAsStateWithLifecycle()
    ItemListContent(state = state, onSearch = viewModel::onSearch)
}

@Composable
private fun ItemListContent(
    state: ItemListState,
    onSearch: (String) -> Unit
) {
    // Stateless — easy to preview and test
}
```

### Event receiver pattern

For complex screens, use a sealed interface for events rather than multiple callback lambdas.

```kotlin
sealed interface ItemListEvent {
    data class Search(val query: String) : ItemListEvent
    data class Delete(val itemId: String) : ItemListEvent
    data object Refresh : ItemListEvent
}

// In ViewModel
fun onEvent(event: ItemListEvent) {
    when (event) {
        is ItemListEvent.Search -> onSearch(event.query)
        is ItemListEvent.Delete -> deleteItem(event.itemId)
        is ItemListEvent.Refresh -> loadItems(_state.value.searchQuery)
    }
}

// In composable — single lambda instead of many
ItemListContent(state = state, onEvent = viewModel::onEvent)
```

**Anti-pattern:** `mutableStateOf` in ViewModel. Prefer `MutableStateFlow` with
`collectAsStateWithLifecycle` — it respects lifecycle and is safe across configuration changes.

## Navigation

### Type-safe navigation (Compose Navigation 2.8+)

Define routes as `@Serializable` objects. Compile-time type safety, no string route IDs.

```kotlin
@Serializable data object HomeRoute
@Serializable data class DetailRoute(val id: String)
@Serializable data object SettingsRoute

@Composable
fun AppNavHost(navController: NavHostController = rememberNavController()) {
    NavHost(navController, startDestination = HomeRoute) {
        composable<HomeRoute> {
            HomeScreen(onNavigateToDetail = { id -> navController.navigate(DetailRoute(id)) })
        }
        composable<DetailRoute> { backStackEntry ->
            val route = backStackEntry.toRoute<DetailRoute>()
            DetailScreen(id = route.id)
        }
        composable<SettingsRoute> { SettingsScreen() }
    }
}
```

### Dialogs and bottom sheets

Use `dialog()` and overlay patterns instead of imperative show/hide state.

```kotlin
NavHost(navController, startDestination = HomeRoute) {
    composable<HomeRoute> { /* ... */ }
    dialog<ConfirmDeleteRoute> { backStackEntry ->
        val route = backStackEntry.toRoute<ConfirmDeleteRoute>()
        ConfirmDeleteDialog(
            itemId = route.itemId,
            onConfirm = { navController.popBackStack() },
            onDismiss = { navController.popBackStack() }
        )
    }
}
```

**Anti-pattern:** passing `NavController` deep into composables. Pass lambda callbacks instead.

## Composable design

### Slot-based APIs

Use slot parameters for flexible, composable APIs.

```kotlin
@Composable
fun AppCard(
    modifier: Modifier = Modifier,
    header: @Composable () -> Unit = {},
    content: @Composable ColumnScope.() -> Unit,
    actions: @Composable RowScope.() -> Unit = {}
) {
    Card(modifier = modifier) {
        Column {
            header()
            Column(content = content)
            Row(horizontalArrangement = Arrangement.End, content = actions)
        }
    }
}
```

### Modifier order

Apply modifiers in this order: layout → shape → drawing → interaction.

```kotlin
Text(
    text = "Hello",
    modifier = Modifier
        .padding(16.dp)                       // 1. Layout
        .clip(RoundedCornerShape(8.dp))       // 2. Shape
        .background(Color.White)              // 3. Drawing
        .clickable { }                        // 4. Interaction
)
```

## KMP platform-specific UI

### expect/actual for platform composables

```kotlin
// commonMain
@Composable
expect fun PlatformStatusBar(darkIcons: Boolean)

// androidMain
@Composable
actual fun PlatformStatusBar(darkIcons: Boolean) {
    val systemUiController = rememberSystemUiController()
    SideEffect { systemUiController.setStatusBarColor(Color.Transparent, darkIcons) }
}

// iosMain
@Composable
actual fun PlatformStatusBar(darkIcons: Boolean) {
    // iOS handles this via UIKit interop or Info.plist
}
```

## Performance

### Stable types for skippable recomposition

Mark classes `@Stable` or `@Immutable` when all properties are stable. The Compose compiler
skips recomposition of composables whose stable inputs haven't changed.

```kotlin
@Immutable
data class ItemUiModel(
    val id: String,
    val title: String,
    val description: String,
    val progress: Float
)
```

### LazyColumn with stable keys

```kotlin
LazyColumn {
    items(
        items = items,
        key = { it.id }  // Stable keys enable item reuse and animations
    ) { item ->
        ItemRow(item = item)
    }
}
```

### derivedStateOf for deferred reads

Use `derivedStateOf` to avoid recomposing the parent when only a derived value changes.

```kotlin
val listState = rememberLazyListState()
val showScrollToTop by remember {
    derivedStateOf { listState.firstVisibleItemIndex > 5 }
}
```

### Avoid allocations during recomposition

```kotlin
// BAD — new lambda and filtered list every recomposition
items.filter { it.isActive }.forEach { ActiveItem(it, onClick = { handle(it) }) }

// GOOD — stable key per item, filtered list memoized
val activeItems = remember(items) { items.filter { it.isActive } }
activeItems.forEach { item ->
    key(item.id) {
        ActiveItem(item, onClick = { handle(item) })
    }
}
```

## Material 3 theming

### Dynamic color scheme (Android 12+)

```kotlin
@Composable
fun AppTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    dynamicColor: Boolean = true,
    content: @Composable () -> Unit
) {
    val colorScheme = when {
        dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S -> {
            if (darkTheme) dynamicDarkColorScheme(LocalContext.current)
            else dynamicLightColorScheme(LocalContext.current)
        }
        darkTheme -> darkColorScheme()
        else -> lightColorScheme()
    }

    MaterialTheme(colorScheme = colorScheme, content = content)
}
```

## Anti-patterns

| Anti-pattern | Preferred alternative |
|---|---|
| `mutableStateOf` in ViewModel | `MutableStateFlow` + `collectAsStateWithLifecycle` |
| Passing `NavController` deep into composables | Pass lambda callbacks |
| Heavy computation in `@Composable` functions | Move to ViewModel or `remember {}` |
| `LaunchedEffect(Unit)` for ViewModel initialization | Initialize in ViewModel `init {}` |
| Creating new object instances in composable params | Hoist to `remember {}` or ViewModel |

## Related skills

- `mobile-ios-design` — SwiftUI and iOS HIG patterns for KMP iOS targets
