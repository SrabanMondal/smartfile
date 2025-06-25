# 🔍 SmartFile CLI

A fast, minimal yet powerful file system toolkit written in Go.  
Organize, clean, analyze, flatten, and search your file system with surgical precision.

> 🧠 Designed for humans. Works like a charm on cluttered directories.

---

## 📦 Features

✅ File Organizer  
✅ File Cleaner (empty file remover)  
✅ Archive Old Files  
✅ Tree + Storage Analyzer  
✅ Flatten Nested Directories  
✅ Smart Search (by name, size, date, type)  
✅ Smart Summary  
✅ Keyword Matching in Text/Code Files

---

## 🚀 Installation

### Option 1: Download Prebuilt Binary

```bash
curl -L https://github.com/SrabanMondal/smartfile/releases/latest/download/smartfile.exe -o smartfile
chmod +x smartfile
```

### Option 2: Build from Source

This project is built using go version 1.24.3

```bash
git clone https://github.com/SrabanMondal/smartfile
cd smartfile
go build -o smartfile main.go
./smartfile
```

## 🧩 Commands & Usage

You can use --help flag to get all information from cli as well

```bash
smartfile --help
smartfile [command] --help
```

### Brief description of commands

#### 🔄 organize

Organize files by extension or by modified year-month.

```bash
smartfile organize --type=ext        # Group into folders like PDF, DOCX, JPG
smartfile organize --type=date       # Group into 2025/06 etc.
smartfile organize --type=ext --depth=5     # Look into subfolders too upto depth you want. Default is scanning only top files. Use cdepth=-1 for scanning till bottom
```

#### 🧹 clean

Delete all empty files (size 0B) from current folder or recursively.

```bash
smartfile clean                    # Removes empty folders from current directory
```

#### 📦 archive

Archive old files into archive/ or zip if flag is set.

```bash
smartfile archive --months=6       # Archive files older than 6 months
smartfile archive --month=12 --zip            # Zip them into archive.zip as well
```

#### 🌲 analyze (tree + storage)

Tree view with folder summaries (folder count, file count, extensions, size).

```bash
smartfile analyze                  # Basic summary
smartfile analyze --detailed       # Includes file names as well
smartfile analyze --max      # Show largest size file in each folder
```

#### 📁 flatten

Bring files from deep structure into one flat folder.

```bash
smartfile flatten --level=2        # Flatten upto 2 levels
smartfile flatten --move           # Move instead of copy
smartfile flatten --output=flat    # Output folder name. Default is flattened
smartfile flatten --unique=false   # Overwrite if name clashes
smartfile flatten --here           # Flatten directly into current dir
```

#### 🔍 search (SmartSearch)

Search files by extension, size, name, date, or even content (in .txt, .md, .csv, .go etc.)

```bash
smartfile search --help # for detailed descriptions
smartfile search --ext=".txt,.md" --contains="Diabetes"
smartfile search --name="report"
smartfile search --min="1MB" --max="100MB"
smartfile search --after="2024-01-01" --before="2024-12-31"
smartfile search --sort=size --asc
smartfile search --limit=10
```

#### 📊 summary

Summarize file types, size usage, modified time histograms.

```bash
smartfile summary                  # Overview of files
smartfile summary --ext=".pdf,.txt"     # Limit to PDF and txt files
smartfile summary --within-days=30        # Focus on recent 30 days
```

#### 🔒 Supported Extensions for --contains

You can use contains flag to search keywords inside files as well but will only work if extension filter is provided too, for preventing too many scans in deep level directories

```bash
smartfile search --contains flatten --ext .go  #searches flatten keyboard in all .go files
```

Only works safely on:
.txt, .md, .csv, .log
.go, .py, .js, .html
(Not for binary files like .pdf, .docx)

## Future Enhancements

Adding semantic file search and directory index creation for fast retrievel.
