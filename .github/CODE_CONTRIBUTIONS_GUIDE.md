# Development Guide

- [Development Guide](#development-guide)
  - [Welcome](#welcome)
  - [Preparing Your Local Operating System](#preparing-your-local-operating-system)
    - [Setting Up macOS](#setting-up-macos)
    - [Setting Up Windows](#setting-up-windows)
  - [Installing Required Software](#installing-required-software)
    - [Installing on macOS](#installing-on-macos)
    - [Installing on Linux](#installing-on-linux)
  - [Installing Docker](#installing-docker)
  - [GitHub Workflow](#github-workflow)

## Welcome

This document is the canonical source of truth for building and contributing to the [Go SDK][project].

Please submit an [issue] on GitHub if you:

- Notice a requirement that this doc does not capture.
- Find a different doc that specifies requirements (the doc should instead link here).

## Preparing Your Local Operating System

Where needed, each piece of required software will have separate instructions for Linux, Windows, or macOS.

### Setting Up macOS

Parts of this project assume you are using GNU command line tools; you will need to install those tools on your system. [Follow these directions to install the tools](https://ryanparman.com/posts/2019/using-gnu-command-line-tools-in-macos-instead-of-freebsd-tools/).

In particular, this command installs the necessary packages:

```bash
brew install coreutils ed findutils gawk gnu-sed gnu-tar grep make jq
```

You will want to include this block or something similar at the end of your `.bashrc` or shell init script:

```bash
GNUBINS="$(find `brew --prefix`/opt -type d -follow -name gnubin -print)"

for bindir in ${GNUBINS[@]}
do
  export PATH=$bindir:$PATH
done

export PATH
```

This ensures that the GNU tools are found first in your path. Note that shell init scripts work a little differently for macOS. [This article can help you figure out what changes to make.](https://scriptingosx.com/2017/04/about-bash_profile-and-bashrc-on-macos/)

### Setting Up Windows

If you are running Windows, you will need to use one of two methods to set up your machine. To figure out which method is the best choice, you will first need to determine which version of Windows you are running. To do this, press Windows logo key + R, type winver, and click OK. You may also enter the ver command at the Windows Command Prompt.

- If you're using Windows 10, Version 2004, Build 19041 or higher, you can use Windows Subsystem for Linux (WSL) to perform various tasks. [Follow these instructions to install WSL2](https://docs.microsoft.com/en-us/windows/wsl/install-win10).
- If you're using an earlier version of Windows, then create a Linux virtual machine with at least 8GB of memory and 60GB of disk space.

Once you have finished setting up your WSL2 installation or Linux VM, follow the instructions below to configure your system for building and developing code.

## Installing Required Software

After setting up your operating system, you will be required to install software dependencies required to run examples, perform static checks, linters, execute tests, etc.

### Installing on macOS

Some of the build tools were installed when you prepared your system with the GNU command line tools earlier. However, you will also need to install the [Command Line Tools for Xcode](https://developer.apple.com/library/archive/technotes/tn2339/_index.html).

### Installing on Linux

All Linux distributions have the GNU tools available. The most popular distributions and commands used to install these tools are below.

- Debian/Ubuntu

  ```bash
  sudo apt update
  sudo apt install build-essential
  ```

- Fedora/RHEL/CentOS

  ```bash
  sudo yum update
  sudo yum groupinstall "Development Tools"
  ```

- OpenSUSE

  ```bash
  sudo zypper update
  sudo zypper install -t pattern devel_C_C++
  ```

- Arch

  ```bash
  sudo pacman -Sy base-devel
  ```

### Installing Go

The Go SDK is written in [Go](http://golang.org). If you need to setup a Go development environment, please follow the instructions in the [Go Getting Started guide](https://golang.org/doc/install).

Confirm that your `GOPATH` and `GOBIN` environment variables are correctly set as detailed in [How to Write Go Code](https://golang.org/doc/code.html) before proceeding.

### Installing Docker

Some aspects of development require Docker. To install Docker in your development environment, [follow the instructions from the Docker website](https://docs.docker.com/get-docker/).

**Note:** If you are running macOS, make sure that `/usr/local/bin` is in your `PATH`.

### Project Specific Software

Once you have the basics, you can download and install any project specific dependencies by navigating to the root your fork and running:

```bash
make ensure-deps
```

If you have not forked and `git clone`'ed your fork, please review the next section.

## GitHub Workflow

To check out code to work on, please refer to [this guide][github_workflow].

> Attribution: This was in part barrowed from this [document](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md) but tailored for our use case.

[project]: https://github.com/deepgram/deepgram-go-sdk
[issue]: https://github.com/deepgram/deepgram-go-sdk/issues
[github_workflow]: https://github.com/deepgram/deepgram-go-sdk/.github/GITHUB_WORKFLOW.md
