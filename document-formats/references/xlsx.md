<!-- hub-reference-banner -->
> **Reference file — part of the `document-formats` hub.** Formerly the standalone `xlsx` skill.
> Sibling topics in this family are now reference files under the hubs (`document-formats`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: xlsx
version: "1.1.0"
updated: "2026-05-29"
category: developer
tags: [xlsx, excel, spreadsheet, openpyxl, pandas, formulas, financial-models, data-export, libreoffice]
description: "Excel/spreadsheet expert — create, read, edit, and fix .xlsx files using openpyxl and pandas in Python. TRIGGER: user wants to open, read, edit, or fix an existing .xlsx/.xlsm/.csv/.tsv file (adding columns, formulas, formatting, charts, cleaning messy data); create a spreadsheet from scratch or from other data; convert between tabular formats; clean or restructure messy tabular data into a proper spreadsheet; any task where a spreadsheet file is the primary deliverable. TRIGGER especially when the user references a spreadsheet by name or path — even casually — and wants something done with it. SKIP: primary deliverable is a Word document, HTML report, standalone Python script, database pipeline, or Google Sheets API integration, even if tabular data is involved; CSV-only processing with no formatting requirements (use csv skill instead)."
related_skills: [csv, docx, pdf, pptx]
license: Proprietary. LICENSE.txt has complete terms
---

# XLSX Creation, Editing, and Analysis

## Output requirements (always apply)

### All Excel files

- **Professional font**: Use a consistent, professional font (Arial, Times New Roman) unless the user specifies otherwise.
- **Zero formula errors**: Deliver every Excel model with zero formula errors (`#REF!`, `#DIV/0!`, `#VALUE!`, `#N/A`, `#NAME?`).
- **Preserve existing templates**: When modifying files, study and exactly match existing format, style, and conventions. Existing template conventions override these guidelines.

### Financial models

**Color coding** (industry-standard; skip when user or template specifies otherwise):

| Color | Hex / RGB | Use for |
|-------|-----------|---------|
| Blue text | `0,0,255` | Hardcoded inputs; numbers users change for scenarios |
| Black text | `0,0,0` | All formulas and calculations |
| Green text | `0,128,0` | Links pulling from other worksheets in the same workbook |
| Red text | `255,0,0` | External links to other files |
| Yellow background | `255,255,0` | Key assumptions needing attention or cells to update |

**Number formatting**:
- Years: text strings (`"2024"`, not `"2,024"`)
- Currency: `$#,##0`; always specify units in headers (`"Revenue ($mm)"`)
- Zeros: use `$#,##0;($#,##0);-` so all zeros display as `-`
- Percentages: `0.0%` (one decimal)
- Multiples: `0.0x`
- Negative numbers: parentheses `(123)`, not minus `-123`

**Documentation for hardcoded values**: Add a comment or adjacent cell note in this format:
`Source: [System/Document], [Date], [Specific Reference], [URL if applicable]`

Examples:
- `Source: Company 10-K, FY2024, Page 45, Revenue Note, [SEC EDGAR URL]`
- `Source: Bloomberg Terminal, 8/15/2025, AAPL US Equity`

## Critical rule: use formulas, not hardcoded values

Always write Excel formulas instead of calculating in Python and hardcoding results. This keeps the spreadsheet dynamic.

```python
# WRONG — hardcodes calculated values
sheet['B10'] = df['Sales'].sum()       # Hardcodes 5000
sheet['C5'] = (new - old) / old        # Hardcodes 0.15

# CORRECT — let Excel calculate
sheet['B10'] = '=SUM(B2:B9)'
sheet['C5'] = '=(C4-C2)/C2'
sheet['D20'] = '=AVERAGE(D2:D19)'
```

This applies to all calculations: totals, percentages, ratios, differences.

## LibreOffice requirement

Formula recalculation requires LibreOffice. Assume it is installed. The `scripts/recalc.py` script auto-configures LibreOffice on first run, including in sandboxed environments (handled via `scripts/office/soffice.py`).

## Common workflow

1. **Choose tool**: pandas for data analysis and bulk operations; openpyxl for formulas, formatting, and Excel-specific features.
2. **Create or load**: create new workbook or load existing file.
3. **Modify**: add/edit data, formulas, formatting.
4. **Save**: write to file.
5. **Recalculate formulas (mandatory when formulas are present)**:
   ```bash
   python scripts/recalc.py output.xlsx
   ```
6. **Verify and fix errors**: the script returns JSON; if `status` is `errors_found`, fix the identified errors and recalculate again.

## Tool usage

### Reading and analyzing data (pandas)

```python
import pandas as pd

df = pd.read_excel('file.xlsx')                        # First sheet
all_sheets = pd.read_excel('file.xlsx', sheet_name=None)  # All sheets as dict
```

### Creating new Excel files (openpyxl)

```python
from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill, Alignment

wb = Workbook()
sheet = wb.active
sheet['A1'] = 'Hello'
sheet['B2'] = '=SUM(A1:A10)'
sheet['A1'].font = Font(bold=True, color='FF0000')
sheet['A1'].fill = PatternFill('solid', start_color='FFFF00')
sheet.column_dimensions['A'].width = 20
wb.save('output.xlsx')
```

### Editing existing files (openpyxl)

```python
from openpyxl import load_workbook

wb = load_workbook('existing.xlsx')
sheet = wb['SheetName']   # or wb.active
sheet['A1'] = 'New Value'
sheet.insert_rows(2)
sheet.delete_cols(3)
wb.save('modified.xlsx')
```

**Warning**: `data_only=True` reads calculated values but permanently loses formulas if the file is then saved. Use it only for read-only analysis.

## Recalculation

```bash
python scripts/recalc.py output.xlsx [timeout_seconds]
```

Returns JSON:
```json
{
  "status": "success",
  "total_errors": 0,
  "total_formulas": 42,
  "error_summary": {
    "#REF!": { "count": 2, "locations": ["Sheet1!B5", "Sheet1!C10"] }
  }
}
```

If `status` is `errors_found`, fix the locations in `error_summary` and recalculate.

## Formula verification checklist

Before returning any Excel file with formulas:

- [ ] Test 2–3 sample cell references — verify they pull the correct values
- [ ] Column mapping confirmed (column 64 = BL, not BK; remember zero-indexing vs openpyxl 1-indexing)
- [ ] Row offset confirmed (DataFrame row 5 = Excel row 6 due to 1-indexing)
- [ ] NaN handling in place (`pd.notna()` checks before referencing null cells)
- [ ] Division-by-zero guards on all `/` formulas (`#DIV/0!`)
- [ ] Cross-sheet references use correct format (`Sheet1!A1`)
- [ ] No `#REF!` from deleted rows/columns
- [ ] Recalc script confirms zero errors

## Code style

Write minimal, concise Python — no unnecessary comments, no verbose variable names, no redundant print statements. For the Excel file itself, add cell comments for complex formulas and document data sources for hardcoded values.
