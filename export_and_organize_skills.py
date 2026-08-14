import os
import sys
import shutil
import json
import re

REPO_DIR = "/Users/mitch.hudson/dev/skills"

IGNORE_DIRS = {
    ".git", "node_modules", "checkpoint-documents", "backups",
    ".system_generated", "cache", "downloads", "ide"
}

def parse_frontmatter(content):
    meta = {}
    if content.startswith("---"):
        parts = content.split("---", 2)
        if len(parts) >= 3:
            fm_text = parts[1]
            current_key = None
            current_val = []
            for line in fm_text.split("\n"):
                if not line.strip():
                    continue
                if ":" in line and not line.startswith(" ") and not line.startswith("\t"):
                    if current_key:
                        meta[current_key] = " ".join(current_val).strip()
                    k, v = line.split(":", 1)
                    current_key = k.strip()
                    v = v.strip().strip("\"'")
                    if v in [">-", "|", ">", "|-"]:
                        current_val = []
                    else:
                        current_val = [v] if v else []
                elif current_key:
                    current_val.append(line.strip())
            if current_key:
                meta[current_key] = " ".join(current_val).strip()
    return meta

def extract_skill_metadata(skill_dir):
    skill_md = os.path.join(skill_dir, "SKILL.md")
    if not os.path.exists(skill_md):
        return None
    try:
        with open(skill_md, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()
        fm = parse_frontmatter(content)
        name = fm.get("name", "")
        desc = fm.get("description", "")
        if not name:
            for line in content.split("\n"):
                if line.startswith("# "):
                    name = line.replace("# ", "").strip()
                    break
        if not name:
            name = os.path.basename(skill_dir)
        if not desc:
            lines = [l.strip() for l in content.split("\n") if l.strip() and not l.startswith("#") and not l.startswith("---")]
            desc = lines[0][:250] if lines else "No description provided."
        return {
            "name": name,
            "description": desc,
            "has_references": os.path.exists(os.path.join(skill_dir, "references")),
            "has_scripts": os.path.exists(os.path.join(skill_dir, "scripts")),
            "has_examples": os.path.exists(os.path.join(skill_dir, "examples")),
        }
    except Exception as e:
        return {"name": os.path.basename(skill_dir), "description": f"Error: {e}"}

def copy_skill(src_dir, dst_dir):
    os.makedirs(os.path.dirname(dst_dir), exist_ok=True)
    if os.path.exists(dst_dir):
        shutil.rmtree(dst_dir)
    
    def ignore_patterns(path, names):
        ignored = set()
        for name in names:
            if name in IGNORE_DIRS or name.endswith(('.part_00', '.part_01', '.DS_Store', '.zip', '.bak')):
                ignored.add(name)
        return ignored

    shutil.copytree(src_dir, dst_dir, symlinks=False, ignore=ignore_patterns)

def categorize_skill(name, desc):
    text = (name + " " + desc).lower()
    if any(w in text for w in ["alphafold", "protein", "genome", "genomic", "dna", "rna", "ncbi", "pdb", "chembl", "clinical", "clinvar", "dbsnp", "ols", "ccres", "ensembl", "foldseek", "gnomad", "gtex", "jaspar", "arxiv", "biorxiv", "europepmc", "openalex", "openfda", "opentargets", "pymol", "quickgo", "reactome", "string", "ucsc", "unibind", "uniprot", "science", "biology", "biomedical"]):
        return "Science, Biology & Medicine"
    elif any(w in text for w in ["bigquery", "dataform", "dbt", "spark", "lakehouse", "dataproc", "sql", "analytics", "da-", "database", "postgres", "mongodb", "mql", "elasticsearch", "splunk", "vector", "atlas"]):
        return "Databases, Data Engineering & Analytics"
    elif any(w in text for w in ["docker", "k8s", "kubernetes", "devops", "aws", "gcp", "cloud", "composer", "airflow", "dataflow", "ci/cd", "cicd", "infra", "linux", "system", "git", "evergreen"]):
        return "Cloud, DevOps & Infrastructure"
    elif any(w in text for w in ["security", "a11y", "memory-leak", "debug", "threat", "incident", "okta", "pacs", "rfid", "nfc", "wiegand", "smartcard", "credential", "auth"]):
        return "Security, Auth & Diagnostics"
    elif any(w in text for w in ["react", "vue", "frontend", "html", "css", "web", "chrome", "extension", "ui", "leafygreen", "browser", "lcp", "performance"]):
        return "Frontend & Web Development"
    elif any(w in text for w in ["agent", "llm", "gemini", "claude", "prompt", "rag", "retrieval", "antigravity", "subagent", "mcp"]):
        return "AI, Agents & Prompt Engineering"
    elif any(w in text for w in ["python", "golang", "go-", "rust", "javascript", "js-ts", "code", "programming", "software-engineering", "coding"]):
        return "Programming & Languages"
    elif any(w in text for w in ["psychology", "behavior", "influence", "clinical", "trauma", "social", "giving", "personality"]):
        return "Psychology & Behavioral Science"
    elif any(w in text for w in ["copywriting", "writing", "resume", "cv", "marketing", "executive", "comms", "content", "instruction"]):
        return "Writing, Communication & Marketing"
    elif any(w in text for w in ["firecrawl", "crawl", "scrape", "scraping", "lead"]):
        return "Web Scraping & Firecrawl"
    elif any(w in text for w in ["trading", "investing", "finance", "bank", "credit", "debt", "fsi", "regulatory", "account-intelligence", "jpmorgan", "goldman"]):
        return "Finance, Banking & Business"
    else:
        return "General & Specialized Utilities"

def main():
    print("Starting Skills Consolidation...")
    
    collected_skills = []
    
    # 1. Antigravity Built-in Skills
    ag_builtin = os.path.expanduser("~/.gemini/antigravity/builtin/skills")
    if os.path.exists(ag_builtin):
        for item in sorted(os.listdir(ag_builtin)):
            src = os.path.join(ag_builtin, item)
            if os.path.isdir(src) and os.path.exists(os.path.join(src, "SKILL.md")):
                dst = os.path.join(REPO_DIR, "antigravity", item)
                copy_skill(src, dst)
                meta = extract_skill_metadata(dst)
                if meta:
                    meta.update({
                        "platform": "Antigravity",
                        "section": "antigravity",
                        "rel_path": f"antigravity/{item}",
                        "category": categorize_skill(meta["name"], meta["description"])
                    })
                    collected_skills.append(meta)
                    print(f"[Antigravity] Copied: {meta['name']}")

    # 2. Gemini Plugins
    gemini_plugins = os.path.expanduser("~/.gemini/config/plugins")
    if os.path.exists(gemini_plugins):
        for plugin in sorted(os.listdir(gemini_plugins)):
            plugin_path = os.path.join(gemini_plugins, plugin)
            if not os.path.isdir(plugin_path) or plugin in IGNORE_DIRS:
                continue
            for root, dirs, files in os.walk(plugin_path):
                if any(x in root for x in IGNORE_DIRS):
                    continue
                if "SKILL.md" in files:
                    rel_sub = os.path.relpath(root, gemini_plugins)
                    dst = os.path.join(REPO_DIR, "gemini", rel_sub)
                    copy_skill(root, dst)
                    meta = extract_skill_metadata(dst)
                    if meta:
                        meta.update({
                            "platform": "Gemini / Plugins",
                            "section": f"gemini/{plugin}",
                            "rel_path": f"gemini/{rel_sub}",
                            "category": categorize_skill(meta["name"], meta["description"])
                        })
                        collected_skills.append(meta)
                        print(f"[Gemini Plugin: {plugin}] Copied: {meta['name']}")

    # 3. Claude Code Agents Skills (~/.agents/skills)
    agents_path = os.path.expanduser("~/.agents/skills")
    if os.path.exists(agents_path):
        for item in sorted(os.listdir(agents_path)):
            src = os.path.join(agents_path, item)
            if not os.path.isdir(src) or item in IGNORE_DIRS or item.startswith(".sync-backup"):
                continue
            if os.path.exists(os.path.join(src, "SKILL.md")):
                dst = os.path.join(REPO_DIR, "claude-code", item)
                copy_skill(src, dst)
                meta = extract_skill_metadata(dst)
                if meta:
                    meta.update({
                        "platform": "Claude Code",
                        "section": "claude-code",
                        "rel_path": f"claude-code/{item}",
                        "category": categorize_skill(meta["name"], meta["description"])
                    })
                    collected_skills.append(meta)
                    print(f"[Claude Code] Copied: {meta['name']}")
            else:
                # check nested
                for root, dirs, files in os.walk(src):
                    if "SKILL.md" in files and not any(x in root for x in IGNORE_DIRS):
                        rel_sub = os.path.relpath(root, agents_path)
                        dst = os.path.join(REPO_DIR, "claude-code", rel_sub)
                        copy_skill(root, dst)
                        meta = extract_skill_metadata(dst)
                        if meta:
                            meta.update({
                                "platform": "Claude Code",
                                "section": "claude-code",
                                "rel_path": f"claude-code/{rel_sub}",
                                "category": categorize_skill(meta["name"], meta["description"])
                            })
                            collected_skills.append(meta)
                            print(f"[Claude Code Nested] Copied: {meta['name']}")

    # 4. Claude Standalone & Workspace Skills (~/.claude)
    claude_path = os.path.expanduser("~/.claude")
    if os.path.exists(claude_path):
        for item in sorted(os.listdir(claude_path)):
            src = os.path.join(claude_path, item)
            if os.path.islink(src) or item in IGNORE_DIRS or item in ["plugins", "cache", "backups", "security", "cost-estimator", "tasks", "ide", "downloads"]:
                continue
            if os.path.isdir(src):
                if os.path.exists(os.path.join(src, "SKILL.md")):
                    dst = os.path.join(REPO_DIR, "claude", "standalone", item)
                    copy_skill(src, dst)
                    meta = extract_skill_metadata(dst)
                    if meta:
                        meta.update({
                            "platform": "Claude",
                            "section": "claude/standalone",
                            "rel_path": f"claude/standalone/{item}",
                            "category": categorize_skill(meta["name"], meta["description"])
                        })
                        collected_skills.append(meta)
                        print(f"[Claude Standalone] Copied: {meta['name']}")
                elif item == "_archived-treearch-20260621":
                    for sub in sorted(os.listdir(src)):
                        sub_src = os.path.join(src, sub)
                        if os.path.isdir(sub_src) and os.path.exists(os.path.join(sub_src, "SKILL.md")):
                            dst = os.path.join(REPO_DIR, "claude", "archived-treearch", sub)
                            copy_skill(sub_src, dst)
                            meta = extract_skill_metadata(dst)
                            if meta:
                                meta.update({
                                    "platform": "Claude (Archived)",
                                    "section": "claude/archived-treearch",
                                    "rel_path": f"claude/archived-treearch/{sub}",
                                    "category": categorize_skill(meta["name"], meta["description"])
                                })
                                collected_skills.append(meta)
                                print(f"[Claude Archived] Copied: {meta['name']}")

    # 5. Custom / Specialized Skills
    # 5a. Desktop MongoDB Skills
    dt_mongo = os.path.expanduser("~/Desktop/mongodb-skills")
    if os.path.exists(dt_mongo):
        for item in sorted(os.listdir(dt_mongo)):
            src = os.path.join(dt_mongo, item)
            if os.path.isdir(src) and os.path.exists(os.path.join(src, "SKILL.md")):
                dst = os.path.join(REPO_DIR, "custom", "mongodb-desktop", item)
                copy_skill(src, dst)
                meta = extract_skill_metadata(dst)
                if meta:
                    meta.update({
                        "platform": "Custom / MongoDB",
                        "section": "custom/mongodb-desktop",
                        "rel_path": f"custom/mongodb-desktop/{item}",
                        "category": categorize_skill(meta["name"], meta["description"])
                    })
                    collected_skills.append(meta)
                    print(f"[Custom MongoDB] Copied: {meta['name']}")

    # 5b. MQL Optimizer
    mql_path = os.path.expanduser("~/dev/mqloptimizer")
    if os.path.exists(mql_path):
        for root, dirs, files in os.walk(mql_path):
            if "SKILL.md" in files:
                skill_name = os.path.basename(root)
                dst = os.path.join(REPO_DIR, "custom", "mqloptimizer", skill_name)
                copy_skill(root, dst)
                meta = extract_skill_metadata(dst)
                if meta:
                    meta.update({
                        "platform": "Custom / MQL Optimizer",
                        "section": "custom/mqloptimizer",
                        "rel_path": f"custom/mqloptimizer/{skill_name}",
                        "category": categorize_skill(meta["name"], meta["description"])
                    })
                    collected_skills.append(meta)
                    print(f"[Custom MQL Optimizer] Copied: {meta['name']}")

    # 5c. Customer Intelligence Skills (Goldman Sachs, Okta, etc.)
    cust_path = os.path.expanduser("~/customers")
    if os.path.exists(cust_path):
        for root, dirs, files in os.walk(cust_path):
            if "SKILL.md" in files and not any(x in root for x in IGNORE_DIRS):
                skill_name = os.path.basename(root)
                dst = os.path.join(REPO_DIR, "custom", "customer-intelligence", skill_name)
                copy_skill(root, dst)
                meta = extract_skill_metadata(dst)
                if meta:
                    meta.update({
                        "platform": "Custom / Customer Intelligence",
                        "section": "custom/customer-intelligence",
                        "rel_path": f"custom/customer-intelligence/{skill_name}",
                        "category": categorize_skill(meta["name"], meta["description"])
                    })
                    collected_skills.append(meta)
                    print(f"[Custom Customer Intel] Copied: {meta['name']}")

    print(f"\nTotal unique skills copied: {len(collected_skills)}")
    
    # Save catalog.json
    with open(os.path.join(REPO_DIR, "catalog.json"), "w", encoding="utf-8") as f:
        json.dump(collected_skills, f, indent=2)
    print("Saved catalog.json")

    # Generate INDEX.md
    generate_index_md(collected_skills)
    
    # Generate README.md
    generate_readme_md(collected_skills)
    
    # Generate section READMEs
    generate_section_readmes(collected_skills)

    # Generate .gitignore
    generate_gitignore()

def generate_gitignore():
    content = """.DS_Store
*.tmp
*.log
*.lock
*.part_*
*.zip
.system_generated/
__pycache__/
*.pyc
node_modules/
"""
    with open(os.path.join(REPO_DIR, ".gitignore"), "w", encoding="utf-8") as f:
        f.write(content)
    print("Generated .gitignore")

def generate_index_md(skills):
    # Sort skills by name
    sorted_skills = sorted(skills, key=lambda x: x["name"].lower())
    
    # Group by category
    categories = {}
    for s in sorted_skills:
        cat = s.get("category", "General & Specialized Utilities")
        categories.setdefault(cat, []).append(s)

    lines = [
        "# Master Skills Index",
        "",
        f"This index provides an exhaustive directory of all **{len(skills)}** agent skills cataloged across **Google Antigravity**, **Gemini / Plugins**, **Claude Code**, **Claude Workspace**, and **Custom Tooling**.",
        "",
        "## Table of Contents",
        "- [Skills by Category](#skills-by-category)",
    ]
    for cat in sorted(categories.keys()):
        slug = cat.lower().replace(" ", "-").replace(",", "").replace("&", "")
        lines.append(f"  - [{cat} ({len(categories[cat])})](#{slug})")
    lines.append("- [Complete Alphabetical Index](#complete-alphabetical-index)")
    lines.append("")
    lines.append("---")
    lines.append("")
    lines.append("## Skills by Category")
    lines.append("")

    for cat in sorted(categories.keys()):
        slug = cat.lower().replace(" ", "-").replace(",", "").replace("&", "")
        lines.append(f"### {cat}")
        lines.append("")
        lines.append("| Skill | Platform | Path | Description |")
        lines.append("| :--- | :--- | :--- | :--- |")
        for s in categories[cat]:
            name = s["name"]
            platform = s["platform"]
            rel = s["rel_path"]
            desc = s["description"].replace("\n", " ").replace("|", "\\|")
            if len(desc) > 180:
                desc = desc[:177] + "..."
            lines.append(f"| [{name}]({rel}/SKILL.md) | `{platform}` | [`{rel}`]({rel}) | {desc} |")
        lines.append("")

    lines.append("---")
    lines.append("")
    lines.append("## Complete Alphabetical Index")
    lines.append("")
    lines.append("| Skill Name | Category | Platform | Location |")
    lines.append("| :--- | :--- | :--- | :--- |")
    for s in sorted_skills:
        name = s["name"]
        cat = s["category"]
        platform = s["platform"]
        rel = s["rel_path"]
        lines.append(f"| [{name}]({rel}/SKILL.md) | {cat} | `{platform}` | [`{rel}`]({rel}) |")
    lines.append("")

    with open(os.path.join(REPO_DIR, "INDEX.md"), "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print("Generated INDEX.md")

def generate_readme_md(skills):
    sorted_skills = sorted(skills, key=lambda x: x["name"].lower())
    
    # Counts by platform
    by_platform = {}
    for s in skills:
        p = s["platform"]
        by_platform[p] = by_platform.get(p, 0) + 1
        
    by_cat = {}
    for s in skills:
        c = s["category"]
        by_cat[c] = by_cat.get(c, 0) + 1

    lines = [
        "# Agent Skills Repository",
        "",
        "A unified, comprehensive collection of **skills**, **workflows**, **prompts**, and **domain intelligence** for modern AI agent environments: **Google Antigravity**, **Gemini**, **Claude Code**, and **Claude**.",
        "",
        f"> **Repository Status**: {len(skills)} Production Skills Cataloged & Indexed",
        "",
        "## Repository Overview",
        "",
        "```",
        "skills/",
        "├── antigravity/               # Google Antigravity core & built-in skills (AGY, SDK, rules)",
        "├── gemini/                    # Gemini ecosystem & official plugin skills (Science, Firebase, BigQuery, DevTools, etc.)",
        "├── claude-code/               # Claude Code & Agents skills library (Engineering, DevOps, RAG, Firecrawl, etc.)",
        "├── claude/                    # Claude standalone workspace skills & specialized architectural trees",
        "├── custom/                    # Dedicated domain skills (MongoDB Atlas, MQL Optimizer, Account Intel)",
        "├── catalog.json               # Machine-readable JSON metadata catalog for programmatic skill discovery",
        "├── INDEX.md                   # Full categorized and alphabetical searchable index",
        "└── README.md                  # Main entry point and usage guide",
        "```",
        "",
        "## Distribution Summary",
        "",
        "### By Platform",
        "| Platform | Total Skills | Primary Focus |",
        "| :--- | :---: | :--- |",
    ]

    platform_descs = {
        "Antigravity": "Antigravity CLI, IDE customizations, SDK orchestration, and permissioning",
        "Gemini / Plugins": "Domain plugins: Science/Bio, Data Agent Kit (BigQuery/dbt), Firebase, Chrome DevTools, Modern Web",
        "Claude Code": "Deep engineering, infrastructure, AI orchestration, Firecrawl automation, psychology & analysis",
        "Claude": "Standalone Claude workspace intelligence, security access trees, and prompt optimization",
        "Claude (Archived)": "Specialized RFID, NFC, physical security, and trading strategy reference trees",
        "Custom / MongoDB": "MongoDB Atlas connection, query optimization, stream processing, search & AI",
        "Custom / MQL Optimizer": "Deep MQL query optimization, execution plan analysis, and indexing",
        "Custom / Customer Intelligence": "Enterprise account intelligence and incident response playbooks",
    }

    for p, count in sorted(by_platform.items(), key=lambda x: x[1], reverse=True):
        desc = platform_descs.get(p, "Specialized domain workflows and agents")
        lines.append(f"| **{p}** | `{count}` | {desc} |")

    lines.extend([
        "",
        "### By Domain Category",
        "| Category | Skills Count |",
        "| :--- | :---: |",
    ])
    for c, count in sorted(by_cat.items(), key=lambda x: x[1], reverse=True):
        lines.append(f"| **{c}** | `{count}` |")

    lines.extend([
        "",
        "---",
        "",
        "## Quick Navigation",
        "",
        "- 📖 **[Full Searchable Master Index (INDEX.md)](INDEX.md)**",
        "- ⚡ **[Antigravity Skills Directory (antigravity/)](antigravity/README.md)**",
        "- 💎 **[Gemini & Plugins Directory (gemini/)](gemini/README.md)**",
        "- 🤖 **[Claude Code Skills Directory (claude-code/)](claude-code/README.md)**",
        "- 🧠 **[Claude Workspace Skills Directory (claude/)](claude/README.md)**",
        "- 🛠️ **[Custom & Specialized Skills (custom/)](custom/README.md)**",
        "",
        "---",
        "",
        "## What is an Agent Skill?",
        "",
        "An **Agent Skill** is a self-contained capability package that equips an AI agent with specialized operational procedures, domain knowledge, and execution workflows. Each skill directory contains:",
        "",
        "- `SKILL.md`: The core definition file with YAML frontmatter (`name`, `description`), activation triggers, and step-by-step procedures.",
        "- `references/` *(optional)*: In-depth technical specifications, cheat sheets, and domain reference guides.",
        "- `scripts/` *(optional)*: Automated scripts and CLI helpers executed by the agent.",
        "- `examples/` *(optional)*: Input/output examples and golden reference templates.",
        "",
        "## How to Use These Skills",
        "",
        "### 1. In Antigravity / Gemini",
        "Place skill directories in:",
        "- System-wide: `~/.gemini/antigravity/builtin/skills/<skill-name>/` or `~/.gemini/config/plugins/<plugin-name>/skills/<skill-name>/`",
        "- Project-specific: `.gemini/skills/<skill-name>/` or `.antigravity/skills/<skill-name>/`",
        "",
        "### 2. In Claude Code / Anthropic Agent System",
        "Place skill directories in:",
        "- System-wide: `~/.agents/skills/<skill-name>/` or `~/.claude/skills/<skill-name>/`",
        "- Project-specific: `.claude/skills/<skill-name>/` or `.agents/skills/<skill-name>/`",
        "",
        "### 3. Programmatic Access",
        "Use `catalog.json` to programmatically search, filter, and load skills into any custom LLM workflow or MCP server.",
        "",
        "---",
        "",
        "## License & Attribution",
        "",
        "Maintained by Mitchell Hudson. Compiled and structured for unified multi-agent development.",
        ""
    ])

    with open(os.path.join(REPO_DIR, "README.md"), "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print("Generated README.md")

def generate_section_readmes(skills):
    sections = {
        "antigravity": {
            "title": "Antigravity Built-in Skills",
            "desc": "Core skills built for Google Antigravity (AGY), including customization configuration, CLI guidance, and SDK agent orchestration."
        },
        "gemini": {
            "title": "Gemini Ecosystem & Plugin Skills",
            "desc": "Official plugins and skill suites for the Gemini ecosystem, spanning biomedical research, Google Cloud Data Agent Kit, Firebase backend development, Chrome DevTools debugging, and modern web standards."
        },
        "claude-code": {
            "title": "Claude Code & Agent Skills",
            "desc": "The extensive Claude Code agent capabilities library, covering software architecture, full-stack programming, DevOps pipelines, Firecrawl scraping, prompt engineering, behavioral psychology, and business intelligence."
        },
        "claude": {
            "title": "Claude Workspace & Architecture Skills",
            "desc": "Claude workspace capabilities, prompt deep-optimizers, telemetry pipelines, and archived physical security / RFID protocol knowledge trees."
        },
        "custom": {
            "title": "Custom & Specialized Domain Skills",
            "desc": "Tailored skills for MongoDB Atlas operations, deep MQL query optimization, and enterprise customer intelligence."
        }
    }

    for sec_key, sec_info in sections.items():
        sec_skills = [s for s in skills if s["rel_path"].startswith(sec_key + "/")]
        lines = [
            f"# {sec_info['title']}",
            "",
            sec_info["desc"],
            "",
            f"**Total Skills in Section**: `{len(sec_skills)}`",
            "",
            "## Skills Index",
            "",
            "| Skill | Relative Path | Description |",
            "| :--- | :--- | :--- |",
        ]
        for s in sorted(sec_skills, key=lambda x: x["name"].lower()):
            name = s["name"]
            rel = os.path.relpath(s["rel_path"], sec_key)
            desc = s["description"].replace("\n", " ").replace("|", "\\|")
            if len(desc) > 150:
                desc = desc[:147] + "..."
            lines.append(f"| [{name}]({rel}/SKILL.md) | [`{rel}`]({rel}) | {desc} |")
        lines.append("")
        
        target_path = os.path.join(REPO_DIR, sec_key, "README.md")
        os.makedirs(os.path.dirname(target_path), exist_ok=True)
        with open(target_path, "w", encoding="utf-8") as f:
            f.write("\n".join(lines))
        print(f"Generated {sec_key}/README.md ({len(sec_skills)} skills)")

if __name__ == "__main__":
    main()
