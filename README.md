# Local .env copy CLI

![GitHub Tag](https://img.shields.io/github/v/tag/y3owk1n/cpenv)
![NPM Downloads](https://img.shields.io/npm/dm/cpenv)
![GitHub License](https://img.shields.io/github/license/y3owk1n/cpenv)

EnvCopy CLI is a powerful command-line tool that simplifies the process of copying environment files for different projects. With just a few commands, you can effortlessly manage and replicate environment configurations across your projects.

This is useful when it comes to working within git worktrees, and you need the same .env(s) across multiple worktrees. Also sometimes when you need to run commands like `git reset --hard; and git clean -dfx`, you can always get your .env file back easily.

## Simple Project Demo

<https://github.com/y3owk1n/cpenv/assets/62775956/a61b0944-1507-4291-a81f-2fd4d198a572>

## Monorepo Project Demo

<https://github.com/y3owk1n/cpenv/assets/62775956/d36b733d-a222-46ff-befe-bac8b6fd73ea>

<!--toc:start-->

- [Prerequisites](#prerequisites)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Options](#options)
  - [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Contributions](#contributions)
- [License](#license)
<!--toc:end-->

## Prerequisites

Before using EnvCopy CLI, ensure you have the following prerequisites:

1. Running `cpenv` for the first time will prompt you to setup your `env-files` folder.

- You can set it to any folder you like, but the default is `~/.env-files`

2. Organize your projects within your chosen directory. Each project should have its own subdirectory.

```plaintext
  ~/.env-files
  ├── single-env-project
  │   ├── .env
  ├── multi-env-project
  │   ├── .env
  │   ├── apps
  │   │   ├── web
  │   │   │   ├── .env
  │   │   ├── api
  │   │   │   ├── .env
  └── other-projects...
```

## Features

- **Automatic Project Setup:** Automatically setup the `env-files` folder if it doesn't exist through simple prompts.

- **Interactive Project Selection:** Easily choose the project for which you want to copy environment files using a user-friendly interactive prompt or specify it directly through command-line options.

- **Automatic Overwrite:** Use the --auto-yes option to automatically overwrite existing files without prompting for confirmation.

- **Global Overwrite Option:** Opt for a global overwrite to replace all existing files in the current project with a single confirmation.

- **Backup Env(s) To Vault:** Back up your project env files to vault and ignore `*.template` and `*.example`.

## Getting Started

### Installation

Install via brew:

```bash
brew tap y3owk1n/tap
brew install y3owk1n/tap/cpenv
```

### Usage

Go to your project directory and run the following command in your terminal:

```bash
cpenv // without any subcommands will go into interactive mode
```

This will launch the interactive mode, guiding you through project selection, file copying and backups.

By supplying subcommands, you can skip the action flow

```bash
cpenv copy [options] // start copy interactive flow
cpenv backup // start backup interactive flow
```

### Options

#### For root

- -V, --version: Output the version number
- -h, --help: Display help for command

#### For `cpenv copy`

- -p, --project <project>: Specify the project for which you want to copy environment files.
- -y, --auto-yes: Automatically overwrite files without prompting for confirmation.

#### For `cpenv backup`

- No options for now

### Examples

Interactive Mode

```bash
cpenv // choose action flow
cpenv copy [options] // start copy flow
cpenv backup // start backup flow
```

Specify Project

```bash
cpenv copy -p single-env-project
```

Specify Project and Enable Auto-Overwrite

```bash
cpenv copy -p multi-env-project --auto-yes
```

## Troubleshooting

If you encounter any issues or errors, please refer to the ~~troubleshooting section in the wiki~~ (Not ready yet).

## Contributions

Feel free to contribute by opening issues, suggesting enhancements, or submitting pull requests. We value your feedback and ideas to enhance the capabilities of `cpenv` further!

## License

This project is licensed under the MIT License. Feel free to use, modify, and distribute it as you see fit.

Happy coding!
