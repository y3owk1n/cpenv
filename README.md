# Local .env copy CLI

![GitHub Tag](https://img.shields.io/github/v/tag/y3owk1n/cpenv)
![NPM Downloads](https://img.shields.io/npm/dm/cpenv)
![GitHub License](https://img.shields.io/github/license/y3owk1n/cpenv)

EnvCopy CLI is a powerful command-line tool that simplifies the process of copying environment files for different projects. With just a few commands, you can effortlessly manage and replicate environment configurations across your projects.

This is useful when it comes to working within git worktrees, and you need the same .env(s) across multiple worktrees. Also sometimes when you need to run commands like `git reset --hard; and git clean -dfx`, you can always get your .env file back easily.

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

1. **Node.js:** Make sure you have Node.js installed on your machine.

2. Running `cpenv` for the first time will prompt you to setup your `env-files` folder.

- You can set it to any folder you like, but the default is `~/.env-files`

3. Organize your projects within your chosen directory. Each project should have its own subdirectory.

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

## Getting Started

### Installation

Install cpenv globally:

```bash
npm install -g cpenv
```

or locally:

```bash
npm install --save-dev cpenv
yarn add --dev cpenv
pnpm add -D cpenv
bun add -d cpenv
```

### Usage

If you install globally, go to your project directory and run the following command in your terminal:

```bash
cpenv
```

If you install locally within a project, you can run the following command in your terminal:

```bash
pnpm cpenv # whatever package manager you are using
```

This will launch the interactive mode, guiding you through project selection and file copying.

### Options

- -p, --project <project>: Specify the project for which you want to copy environment files.
- -y, --auto-yes: Automatically overwrite files without prompting for confirmation.

### Examples

Interactive Mode

```bash
cpenv # pnpm cpenv if running locally
```

Specify Project

```bash
cpenv -p single-env-project # pnpm cpenv if running locally
```

Specify Project and Enable Auto-Overwrite

```bash
cpenv -p multi-env-project --auto-yes # pnpm cpenv if running locally
```

## Troubleshooting

If you encounter any issues or errors, please refer to the ~~troubleshooting section in the wiki~~ (Not ready yet).

## Contributions

Feel free to contribute by opening issues, suggesting enhancements, or submitting pull requests. We value your feedback and ideas to enhance the capabilities of `cpenv` further!

## License

This plugin is licensed under the MIT License. Feel free to use, modify, and distribute it as you see fit.

Happy coding!
