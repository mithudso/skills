# Customer Context Architect v2.3.0 Changelog

## Summary
Enhanced document discovery to identify and surface documents created by all MongoDB personnel on an account (TAMs, SAs, AEs, CSMs, case engineers), solving the problem of "I know Russell (SA) has great docs for this customer but I forget where he put them."

## Changes Made

### 1. New Section: MongoDB Team Discovery (Unit 2)
**Location:** After alias resolution, before main source retrieval  
**Purpose:** Systematically discover documents owned/created by MongoDB team members

**Process:**
1. **Extract team from SFDC:** Pull all account team assignments:
   - `Account.TAM__c`, `Solutions_Architect__c`, `Account_Executive__c`, `Customer_Success_Manager__c`
   - All `AccountTeamMember` records
   - `Account.Owner` (fallback)
   - Case owners and engineers from all cases on the account

2. **Per-person document discovery:** For each team member:
   - Normalize name to search variants (`First Last`, `First.Last`, `FLast`, `First_Last`)
   - Use Google Drive API/search to find documents where:
     - Owner = team member email (from SFDC `User.Email`)
     - Creator = team member name
     - Sharing permissions include team member
   - **Scope:** Search all accessible shared drives:
     - `Shared drives/TS Premium Services - TAM & NTSE/` (all subfolders)
     - `Shared drives/Solutions Architecture/` (if SA on account)
     - `Shared drives/Customer Success/` (if CSM on account)
   - **Filter:** Include only if:
     - Filename/path contains customer alias, OR
     - Modified in last 180 days AND content mentions customer alias (first 10 pages)
     - Document is shared (not owner-private)

3. **Deduplication:** Cross-reference against canonical customer folder; mark duplicates but preserve owner attribution

4. **Output:** New `## MongoDB team documents` table:
   ```markdown
   | Owner | Role | Document | Last Modified | Path |
   |-------|------|----------|---------------|------|
   | Russell Easton | Solutions Architect | GS Architecture Deep Dive.pdf | 2026-06-15 | [link] |
   ```

**Fallback:** If Drive API unavailable, fall back to manual folder scan + Glean search by owner

### 2. Updated Sources List
**File:** `customer-context-architect/SKILL.md`, lines 184-186

**Changes:**
- **Source 2 (SFDC):** Now explicitly lists team assignment fields extracted for document discovery
- **Source 2a (NEW):** MongoDB team documents source with full description

### 3. Updated Precedence Rules
**File:** `customer-context-architect/SKILL.md`, line 209

**New precedence tier:**
- **Rank 3:** MongoDB team documents (SA architecture reviews, TAM runbooks)
- When conflict with canonical folder: prefer **most recent** unless canonical marked "final"/"approved"

### 4. Version Bump
- **v2.2.0 → v2.3.0**
- Updated: 2026-07-02

## Problem Solved

**Before:** TAMs would forget where SAs like Russell Easton stored valuable customer documents (architecture deep-dives, migration runbooks, technical reviews). These weren't in the canonical customer folder, so they were invisible to the context architect.

**After:** The system now:
1. Identifies all MongoDB team members from SFDC
2. Searches for documents they own/created across all shared drives
3. Filters to customer-relevant documents
4. Surfaces them in a dedicated `## MongoDB team documents` table with owner attribution

**Example:** For Goldman Sachs, if Russell Easton (SA) has 5 architecture docs scattered across `Shared drives/Solutions Architecture/` that mention Goldman or GS, they'll now appear in the context file with attribution, even if they're not in the canonical GS engagement folder.

## Implementation Notes

### Google Drive Search Methods
The skill supports multiple implementation paths:

**Option 1: Google Drive API (preferred)**
```python
from googleapiclient.discovery import build

service = build('drive', 'v3', credentials=creds)
query = f"'{user_email}' in owners and (name contains '{customer_alias}' or fullText contains '{customer_alias}')"
results = service.files().list(
    q=query,
    corpora='drive',
    driveId=shared_drive_id,
    fields='files(id, name, owners, modifiedTime, webViewLink)'
).execute()
```

**Option 2: Drive CLI tools**
```bash
# Using drive-cli or rclone
drive list --query "owner:'russell.easton@mongodb.com' AND name:Goldman"
```

**Option 3: Manual grep (fallback)**
```bash
find "/Users/mitch.hudson/Library/CloudStorage/GoogleDrive-mitch.hudson@mongodb.com/Shared drives/" \
  -type f -user "Russell Easton" -mtime -180 | grep -i "goldman\|GS"
```

### SFDC Fields to Extract
Ensure your SFDC query includes:
```soql
SELECT 
  TAM__c, 
  Solutions_Architect__c,
  Account_Executive__c,
  Customer_Success_Manager__c,
  (SELECT UserId, TeamMemberRole FROM AccountTeamMembers)
FROM Account
WHERE Name = 'Goldman Sachs'
```

## Testing Checklist

- [ ] Run on Goldman Sachs account (Russell Easton SA test case)
- [ ] Verify SFDC team extraction includes all roles
- [ ] Confirm Drive search finds SA-owned documents
- [ ] Check deduplication against canonical folder
- [ ] Verify fallback when Drive API unavailable
- [ ] Test with account that has no SA (should gracefully skip SA-specific drives)
- [ ] Confirm 180-day filter excludes stale documents
- [ ] Verify output table format matches spec

## Related Skills
- `tam-operations` — Calls this skill for account reviews
- `solve-case` — May call this for case-specific context
- `atlas-diagnostics-expert` — Can use discovered docs for troubleshooting

## Future Enhancements (Not in v2.3.0)
- **Glean owner filter:** Add Glean search by document owner
- **Slack thread attribution:** Identify which team member started high-value threads
- **Meeting note ownership:** Tag meeting notes by primary MongoDB attendee
- **Staleness alerts:** Flag team documents >180 days old that should be refreshed
- **Auto-suggest canonical promotion:** Recommend high-value team docs for canonical folder
