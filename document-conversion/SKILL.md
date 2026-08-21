---
name: document-conversion
description: >-
  Uses pandoc and poppler (pdftotext) to convert PDFs, Word documents, and EPUBs into LLM-readable Markdown.
---

# Document Conversion Instructions

When the user asks you to read, analyze, summarize, or extract information from a local binary document (.pdf, .docx, .epub), you MUST NOT try to read the file directly. Instead, you must first convert it to Markdown or plain text.

## Reading PDFs
Use `pdftotext` (from the poppler package) to extract text from PDFs.
```bash
# Extract text to a temporary file
pdftotext "/path/to/document.pdf" "/path/to/temp.txt"

# Read the temporary file
cat "/path/to/temp.txt"
```

## Reading Word Docs (.docx) and EPUBs
Use `pandoc` to convert complex documents into Markdown, which is the optimal format for you to read.
```bash
# Convert docx to Markdown
pandoc "/path/to/document.docx" -o "/path/to/temp.md" --extract-media=/path/to/media_folder

# Convert EPUB to Markdown
pandoc "/path/to/book.epub" -t markdown -o "/path/to/book.md"
```

After converting, read the resulting `.md` or `.txt` file to answer the user's request.
