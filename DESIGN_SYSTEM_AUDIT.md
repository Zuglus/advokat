# Design System Audit — advokat

**Дата:** 2026-03-20
**Область:** Frontend (CSS / HTML templates / JS)
**Файлы:** `css/style.css`, `templates/`, `js/`

---

## Summary

| Metric | Value |
|---|---|
| Components reviewed | 16 |
| Issues found | 9 |
| Score | **74 / 100** |

The project has a well-structured design token system in `:root`, consistent component naming with BEM-like conventions, and a clean Go template component architecture. The main weaknesses are hardcoded color/shadow values that bypass tokens, two nearly identical components that should be merged, and missing interactive states.

---

## 1. Design Tokens

### Defined Tokens (`:root`)

| Category | Tokens | Names |
|---|---|---|
| Color — brand | 6 | `--navy`, `--navy-dark`, `--navy-light`, `--gold`, `--gold-dark`, `--gold-light` |
| Color — neutral | 6 | `--cream`, `--white`, `--text`, `--text-secondary`, `--text-light`, `--border`, `--border-light` |
| Shadow | 3 | `--shadow-sm`, `--shadow-md`, `--shadow-lg` |
| Radius | 2 | `--radius` (8px), `--radius-lg` (12px) |
| Typography | 2 | `--font-heading`, `--font-body` |
| **Spacing** | **0** | **— not tokenized** |
| **Font size** | **0** | **— not tokenized** |
| **Transition** | **0** | **— not tokenized** |

### Hardcoded Values Found

| Category | Instances | Examples |
|---|---|---|
| `#fff` / `#ffffff` | 12 | `.btn-primary`, `.tab-link:hover`, `.featured-card`, burger lines |
| `rgba(255,255,255, …)` | 14 | 0.8, 0.7, 0.4, 0.1, 0.08 — white overlays at 6+ different opacities |
| `rgba(201,168,76, …)` | 3 | gold-based highlight/shadow values |
| `rgba(0,0,0,0.15)` | 1 | `.header` box-shadow (not using `--shadow-*`) |
| `#fef9ed`, `#fff7e0` | 2 | alert/pricing-note gradient (not tokenized) |
| Hardcoded shadows | 3 | `.btn-cta`, `.btn-cta:hover`, `.header` — don't use `--shadow-*` |

**Total hardcoded color/shadow instances: ~35**

### Recommendations

- Add `--white-overlay-*` tokens for the 6 recurring white-with-opacity values
- Add `--alert-bg-start` / `--alert-bg-end` tokens for the gradient backgrounds
- Tokenize the `.btn-cta` shadow as `--shadow-gold`
- Move `.header` box-shadow to use `--shadow-md` or a new `--shadow-header`
- Consider adding spacing tokens (`--space-xs: 4px`, `--space-sm: 8px`, `--space-md: 16px`, `--space-lg: 24px`, `--space-xl: 32px`, `--space-2xl: 40px`) — currently ~15 unique spacing values are used inline

---

## 2. Component Inventory

| Component | CSS class | Template | States | Variants | Score |
|---|---|---|---|---|---|
| **Top Bar** | `.top-bar` | `base.html → topbar` | default | — | 7/10 |
| **Header** | `.header` | `base.html → header` | default, sticky | — | 8/10 |
| **Nav Tabs** | `.nav-tabs`, `.tab-link` | `base.html → header` | default, hover, active | — | 8/10 |
| **Burger Menu** | `.burger-btn`, `.mobile-menu` | `base.html → header` | default, open | — | 6/10 |
| **Breadcrumbs** | `.breadcrumbs` | `components.html → heading` | default | — | 8/10 |
| **Page Heading** | `.page-heading`, `.page-subheading` | `components.html → heading` | default | — | 8/10 |
| **Text Block** | `.text-paragraph` | `components.html → renderText` | — | subtitle, body2, overline | 7/10 |
| **Button** | `.btn`, `.btn-primary` | `components.html → renderButtons` | default, hover | primary, with-icon | 6/10 |
| **Button CTA** | `.btn-cta` | `index.html`, `components.html` | default, hover | — | 5/10 |
| **Button Gold** | `.btn-gold` | `index.html`, `components.html` | default, hover | — | 5/10 |
| **Price Table** | `.price-table` | `components.html → renderTable` | default, row-hover | — | 8/10 |
| **Listing** | `.listing` | `components.html → renderList` | — | divider, subheader | 7/10 |
| **Alert Block** | `.alert-block` | `components.html → renderAlert` | default | — | 6/10 |
| **Pricing Note** | `.pricing-note` | `components.html → renderPricingNote` | default | — | 5/10 |
| **Tabs** | `.tabs-container` | `components.html → renderTabs` | default, active | — | 7/10 |
| **Featured Card** | `.featured-card` | `index.html` | default | — | 7/10 |
| **Footer** | `.footer` | `base.html → footer` | default | — | 8/10 |
| **Legal Details** | `.legal-details` | `contacts.html` | default, open | — | 7/10 |

---

## 3. Naming Consistency

| Issue | Components | Recommendation |
|---|---|---|
| **3 button variants with different naming** | `.btn-primary`, `.btn-cta`, `.btn-gold` | Unify under `.btn` + modifier: `.btn--primary`, `.btn--cta`, `.btn--gold` |
| **Alert and Pricing Note are visually identical** | `.alert-block`, `.pricing-note` | Same gradient background + gold left border. Merge into one component with variants |
| **Mixed class naming convention** | `.tab-link` vs `.tabs-btn`, `.mobile-nav-link` vs `.header-contact-link` | All use flat naming (good), but some use hyphens inconsistently: `tab-link` ≠ `tabs-btn` |
| **`page-heading` vs `section-title`** | `.page-heading` (h1), `.section-title` (h2) | Naming doesn't convey hierarchy — consider `.heading-1`, `.heading-2` or keep but document |
| **Inline styles in templates** | `renderText`, `renderAlert` | `text-align` and `text-decoration` are set via `style=""` — should be CSS classes |

---

## 4. Missing States & Accessibility

| Component | Missing | Impact |
|---|---|---|
| **All buttons** | No `:focus-visible` style | Keyboard users can't see focused button |
| **All buttons** | No `:active` / pressed state | No tactile feedback on press |
| **All buttons** | No `disabled` state | No way to indicate unavailable actions |
| **Tab buttons** | No `:focus-visible` | Keyboard tab navigation invisible |
| **Burger menu** | No `aria-expanded` toggling | Screen readers don't know menu state |
| **Mobile menu** | No focus trap when open | Tab can escape to hidden content |
| **Legal details** | Uses native `<details>` (good) | — works out of the box |
| **Nav active state** | Uses `.active` class | Good — but no `aria-current="page"` |

---

## 5. Duplicate / Overlapping Components

### Alert Block ↔ Pricing Note

These two are visually identical:

```css
/* .alert-block */
background: linear-gradient(135deg, #fef9ed 0%, #fff7e0 100%);
border-left: 4px solid var(--gold);
padding: 20px 24px;
border-radius: 0 var(--radius) var(--radius) 0;

/* .pricing-note */
background: linear-gradient(135deg, #fef9ed 0%, #fff7e0 100%);
border-left: 4px solid var(--gold);
border-radius: 0 var(--radius) var(--radius) 0;
padding: 24px;  /* only difference: 4px more padding */
```

**Recommendation:** Merge into a single `.callout` component. Use `.callout--with-actions` for the pricing note variant.

### btn-cta ↔ btn-gold

Both are gold-background buttons with similar hover behavior. `btn-cta` is larger (14px 32px, 1.1rem) while `btn-gold` is standard size (12px 24px, 0.9375rem). They represent the same concept at two sizes.

**Recommendation:** Merge into `.btn--gold` with a `.btn--lg` size modifier.

---

## 6. Responsive Design

| Breakpoint | Purpose | Coverage |
|---|---|---|
| `min-width: 768px` | Show header contacts | ✅ |
| `max-width: 960px` | Collapse footer grid | ✅ |
| `max-width: 767px` | Mobile layout: burger, single-column | ✅ |

**Gap:** No breakpoint between 768–960px for medium tablets. The nav tabs could overflow on narrow desktops. Consider hiding the nav at `960px` and using the burger menu earlier.

---

## 7. JavaScript

| File | Purpose | Size | Quality |
|---|---|---|---|
| `app.js` | Burger menu toggle | 10 lines | Clean, minimal |
| `tabs.js` | Tab switching | 17 lines | Clean, conditionally loaded |

No issues. The JS is minimal, vanilla, and loaded with `defer`. `tabs.js` is only loaded on pages that need it (`{{if .UseTabs}}`).

---

## 8. Template Architecture

The template system is well-organized:

- **`base.html`** — layout shell with `topbar`, `header`, `footer` partials
- **`components.html`** — reusable render macros (`renderText`, `renderButtons`, `renderTable`, `renderList`, `renderAlert`, `renderTabs`, `renderPricingNote`, `heading`)
- **`pages/*.html`** — each page defines `{{define "content"}}` and composes from macros

**Strength:** Good separation. Pages are thin and declarative, composing from shared components.

**Weakness:** `index.html` duplicates patterns (hero CTA, featured card) that aren't in `components.html`. If other pages need a featured card or hero CTA, there's no shared template.

---

## Priority Actions

1. **Add `:focus-visible` styles to all interactive elements** — accessibility blocker, affects keyboard users and screen reader users.

2. **Tokenize remaining hardcoded colors** — 35 instances of raw hex/rgba bypass the design token system. Start with `#fff` → `var(--white)` and the white-overlay pattern.

3. **Merge `.alert-block` and `.pricing-note`** into one `.callout` component — reduces duplication and ensures visual consistency.

4. **Unify button variants** under `.btn` with modifiers (`.btn--primary`, `.btn--gold`, `.btn--lg`) — currently 3 separate button implementations with inconsistent APIs.

5. **Add `aria-expanded` to burger button** and `aria-current="page"` to active nav links — low effort, meaningful accessibility improvement.

6. **Extract inline `style=""` attributes in templates** into CSS utility classes (`.text-center`, `.text-right`, `.text-underline`).

7. **Extract `featured-card` and `hero-cta` from `index.html`** into reusable templates in `components.html`.

8. **Add spacing tokens** — at least 5 values (`4px`, `8px`, `16px`, `24px`, `32px`, `40px`) to replace the ~15 unique spacing values used across the CSS.

9. **Consider a breakpoint around 960px** for the nav tabs overflow issue on narrow desktop/wide tablet screens.
