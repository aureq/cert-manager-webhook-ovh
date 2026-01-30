# Commit Tool

## Description

This tool automates the process of creating a Git commit.

## Instructions to follow

Write a commit message for the currently staged file(s) going through each file one by one.
ALWAYS Use the command `git diff --cached --name-status` to determine which files are staged for the commit.
NEVER add any files that are not staged for commit.
IGNORE any previous commit messages.
Prefix commit messages with an emoji according to the type of change.
ALWAYS add a space between the emoji and the rest of the commit message.
ALWAYS create commit message based on importance of the change.
Types of change rules (Priority based from most important to least important):
- major new features start with 🎉
- small improvements start with 🌿
- bugs, fixes and corrections start with 🐛
- minor changes, tweaks and renaming variables, tidying code, making code more consistent start with 🌱
- deleted files messages start with 🔥
- dependency upgrades start with ⏩ (`go.mod`, `go.sum`, etc)
- specifications stored in the `specs/` directory start with 🧠.
- documentation and comments start with 📝 (`docs/*`, `README.md` files, comments in code, etc)
- tooling related messages start with ⚙️ (`Dockerfile`, `Makefile`, CI/CD, `.claude/commands/*` etc)
If a change occurs in `.gitignore` file the commit message start with 🙈.
ALL commit messages should be written in the imperative mood, e.g. "Add feature" not "Added feature" or "Adds feature".
If a change related to a command, make sure to use quotes around the command name, e.g. 'core hasync ...' or 'firewall alias ...'.
For large commits, try to summarise the changes. Do NOT add a message saying "Generated with Claude Code".
