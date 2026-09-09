# caraxes

<img width="2172" height="724" alt="image" src="https://github.com/user-attachments/assets/99f2d5b4-3b1a-4047-a1c2-17576923655a" />

**Caraxes** is a simple Git-like version control system implemented in Go. It supports basic Git functionality such as initializing a repository, staging files, storing objects, creating trees, and committing changes. Caraxes is designed as a learning project to understand the inner workings of Git and content-addressable storage.

`caraxes init` greets you with the glowing, drop-shadowed ASCII banner shown above, rendered straight in the terminal.

---

## Features

- Initialize a Caraxes repository (`caraxes init`)
- Stage files into the index (`caraxes add <file>...`)
- Inspect the working tree against the index (`caraxes status`)
- Commit staged changes (`caraxes commit -m <message>`)
- Walk commit history from the current branch (`caraxes log`)
- Supports parent commits and updates the `master` ref automatically
- Content integrity ensured through SHA-1 hashes
- Tree entries sorted by file name for consistent hashing
- Author identity via environment variables (`XGIT_AUTHOR_NAME`, `XGIT_AUTHOR_EMAIL`)
- Colorized status output (Charm's `lipgloss`) and a fire-gradient startup banner

---

## Installation

1. Ensure Go is installed (1.27+). If not, install it from [go.dev](https://go.dev/dl/).
2. Clone the repository:

```bash
git clone https://github.com/MiteshBari23/Caraxes
cd Caraxes/caraxes
```

3. Build the project:

```bash
go build -o caraxes .
```

4. Add the binary to your PATH (optional):

```bash
export PATH="$PATH:$(pwd)"
```

---

## Usage

All commands are executed inside a Caraxes repository (created via `caraxes init`).

### Initialize Repository

```bash
caraxes init [path]
```

Creates a `.caraxes` folder with the following structure:

```
.caraxes/
├─ objects/
│  ├─ info/
│  └─ pack/         # Stores blob, tree, and commit objects
├─ refs/
│  ├─ heads/
│  │  └─ master      # Stores the latest commit hash for the master branch
│  └─ tags/
└─ HEAD              # Points to the current branch (refs/heads/master)
```

Running it a second time reinitializes the existing repository instead of erroring out.

---

### Add

```bash
caraxes add <file>...
```

* Hashes each file into a **blob object** (`blob <size>\0` header + content, SHA-1 addressed).
* Stores the compressed blob under `.caraxes/objects/<hash_prefix>/<hash_suffix>`.
* Records the file path → blob hash mapping in `.caraxes/index`.

Example:

```bash
caraxes add file1.txt file2.txt
```

---

### Status

```bash
caraxes status
```

* Walks the working tree (skipping `.caraxes`) and hashes every file.
* Compares each hash against what's recorded in the index.
* Prints **modified** files (tracked, but content changed) in red and **untracked** files in orange.

Example output:

```
Changes not staged for commit:
  (use "caraxes add <file>..." to update what will be committed)
        modified:   file1.txt

Untracked files:
  (use "caraxes add <file>..." to include in what will be committed)
        file2.txt
```

---

### Commit

```bash
caraxes commit -m "<message>"
```

* Builds a **tree object** from everything currently in the index.
* Creates a **commit object** pointing at that tree, with an optional parent (the current `master` HEAD, if any).
* Records author and committer information from environment variables.
* Updates `refs/heads/master` to the new commit.

Example:

```bash
export XGIT_AUTHOR_NAME="Alice"
export XGIT_AUTHOR_EMAIL="alice@example.com"

caraxes commit -m "Initial commit"
```

Output:

```
committed as 189d8004443db9b25142154c16540dbe218a3dfa
```

---

### Log

```bash
caraxes log
```

* Starts at `refs/heads/master` and walks parent commits back to the root.
* Prints each commit's hash, author, and message.

Example output:

```
commit 189d8004443db9b25142154c16540dbe218a3dfa
Author: Alice <alice@example.com> 1760211794 +0530

    Initial commit
```

---

## Environment Variables

* `XGIT_AUTHOR_NAME` – author name used when committing (required for `commit`)
* `XGIT_AUTHOR_EMAIL` – author email used when committing (required for `commit`)

Example:

```bash
export XGIT_AUTHOR_NAME="Alice"
export XGIT_AUTHOR_EMAIL="alice@example.com"
```

---

## Internal Structure

Caraxes objects mimic Git's content-addressable structure:

```
Object types:
- Blob:   raw file content (with "blob <size>\0" header)
- Tree:   directory structure (file/directory names, hashes)
- Commit: tree hash, parent, author/committer, and commit message
```

### Object → Tree → Commit Diagram

```
[file1.txt]       [file2.txt]         [src/]
    │                 │                  │
    └──> Blob        └──> Blob         └──> Tree
               ┌───────────────────────────────┐
               │              Tree               │
               │   (directory structure hash)   │
               └───────────────────────────────┘
                           │
                           └──> Commit
                                 (tree hash + parent + author + message)
```

* **Blob**: content of a single file
* **Tree**: organizes blobs and subtrees, stores SHA-1 hashes
* **Commit**: points to a tree, optionally references a parent commit

---

## Notes

* Caraxes uses **SHA-1 hashes** for content addressing.
* The index (`.caraxes/index`) is a flat JSON map of path → blob hash.
* Trees store only object hashes, not actual file content.
* Tree entries are **sorted by filename** to ensure consistent hashing.
* Only the `master` branch is supported; `HEAD` points to it.
* `caraxes init` also prints a glowing, drop-shadowed CARAXES banner (`internal/banner`) — set `NO_COLOR` to disable it.

---

## Example Workflow

```bash
# Initialize repository
caraxes init

# Stage files
caraxes add file1.txt file2.txt

# Check what's staged / changed
caraxes status

# Commit staged changes
export XGIT_AUTHOR_NAME="Alice"
export XGIT_AUTHOR_EMAIL="alice@example.com"
caraxes commit -m "Initial commit"

# View commit history
caraxes log
```
