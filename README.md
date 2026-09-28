<p align="center">
  <img src="https://www.freebsd.org/images/beastie-right.svg" width="96" alt="The BSD Daemon" />
</p>

## FreeBSD support

[![Go](https://img.shields.io/badge/go-1.25%2B-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![FreeBSD](https://img.shields.io/badge/FreeBSD-15.1%20tested-AB2B3B?style=flat-square&logo=freebsd&logoColor=white)](https://www.freebsd.org)
[![License](https://img.shields.io/badge/license-MIT-lightgrey?style=flat-square)](https://github.com/zsnero/wails-freebsd/blob/port/freebsd/LICENSE)
[![branch](https://img.shields.io/badge/branch-port%2Ffreebsd-8B6CEF?style=flat-square)](https://github.com/zsnero/wails-freebsd/tree/port/freebsd)
[![Wails](https://img.shields.io/badge/fork%20of-wailsapp%2Fwails-v3.0.0--beta.26-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/wailsapp/wails)

This fork tracks [wailsapp/wails](https://github.com/wailsapp/wails) and adds
FreeBSD as a target for the Wails v3 desktop framework. The work lives on the
`port/freebsd` branch; `master` is kept identical to upstream so it stays easy
to rebase against.

Wails builds on GTK3 and WebKitGTK 4.1, both of which FreeBSD packages, so there
is no emulation layer involved. The app window, webview, dialogs, menus,
clipboard, keyboard handling and single-instance locking all work. GTK3 is the
backend that gets selected, not GTK4. `wails3 doctor` knows about `pkg` and will
report missing packages.

Known gaps:

- The system tray is a stub. There is no StatusNotifierItem implementation on
  FreeBSD, so the API is present for apps that guard their tray code, but every
  call returns `errSystemTrayUnsupported` instead of doing nothing quietly.
- Global shortcuts use the X11 path. There is no FreeBSD equivalent of the Linux
  portal backend, so registration is unreliable outside a running session.

Notifications go through the freedesktop backend, which FreeBSD supports.

Yaria, a desktop video and audio downloader, runs on this port. See
[yaria.live](https://yaria.live).

Issues and pull requests still go to [wailsapp/wails](https://github.com/wailsapp/wails/issues).

### Building and running on FreeBSD

Install the dependencies first. `webkit2-gtk_41` is the GTK3 and libsoup3
flavour; it is the one that installs `webkit2gtk-4.1.pc`, which is what the cgo
directives link against. Note the underscore, the port is `www/webkit2-gtk`.

```sh
sudo pkg install gtk3 webkit2-gtk_41 pkgconf www/npm
```

Base system has clang, so no compiler package is needed. Then:

```sh
git clone -b port/freebsd https://github.com/zsnero/wails-freebsd
cd wails-freebsd/v3
go build ./...
go test ./...
```

The wails3 CLI builds and runs on the port as well:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 doctor
```

`wails3 doctor` reports which of the above are missing and prints the `pkg`
command to fix them. It knows the difference between a package name and the
pkg-config module name, which is the usual thing to get wrong here.

To try a real application, the examples in `v3/examples` build and run:

```sh
cd v3/examples/dev
wails3 dev
```

<p align="center" style="text-align: center">
  <img src="./assets/images/logo-universal.png" width="55%"><br/>
</p>

<p align="center">
  Build desktop applications using Go & Web Technologies.
  <br/>
  <br/>
  <a href="https://github.com/wailsapp/wails/blob/master/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/wailsapp/wails"/>
  </a>
  <a href="https://goreportcard.com/report/github.com/wailsapp/wails">
    <img src="https://goreportcard.com/badge/github.com/wailsapp/wails" />
  </a>
  <a href="https://pkg.go.dev/github.com/wailsapp/wails">
    <img src="https://pkg.go.dev/badge/github.com/wailsapp/wails.svg" alt="Go Reference"/>
  </a>
  <a href="https://github.com/wailsapp/wails/issues">
    <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" />
  </a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_shield" alt="FOSSA Status">
    <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=shield" />
  </a>
  <a href="https://github.com/avelino/awesome-go" rel="nofollow">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="Awesome" />
  </a>
  <a href="https://discord.gg/BrRSWTaxVK">
    <img alt="Discord" src="https://img.shields.io/discord/1042734330029547630?logo=discord"/>
  </a>
  <br/>
  <a href="https://github.com/wailsapp/wails/actions/workflows/build-and-test.yml" rel="nofollow">
    <img src="https://img.shields.io/github/actions/workflow/status/wailsapp/wails/build-and-test.yml?branch=master&logo=Github" alt="Build" />
  </a>
  <a href="https://github.com/wailsapp/wails/tags" rel="nofollow">
    <img alt="GitHub tag (latest SemVer pre-release)" src="https://img.shields.io/github/v/tag/wailsapp/wails?include_prereleases&label=version"/>
  </a>
</p>

<div align="center">
<strong>
<samp>

[English](README.md) · [简体中文](README.zh-Hans.md) · [日本語](README.ja.md) ·
[한국어](README.ko.md) · [Español](README.es.md) · [Português](README.pt-br.md) ·
[Русский](README.ru.md) · [Francais](README.fr.md) · [Uzbek](README.uz.md) · [Deutsch](README.de.md) ·
[Türkçe](README.tr.md) · [Bahasa Indonesia](README.id.md)

</samp>
</strong>
</div>

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Introduction](#introduction)
- [Features](#features)
  - [Roadmap](#roadmap)
- [Getting Started](#getting-started)
- [Sponsors](#sponsors)
- [FAQ](#faq)
- [Stargazers over time](#stargazers-over-time)
- [Contributors](#contributors)
- [License](#license)
- [Inspiration](#inspiration)

## Introduction

The traditional method of providing web interfaces to Go programs is via a built-in web server. Wails offers a different
approach: it provides the ability to wrap both Go code and a web frontend into a single binary. Tools are provided to
make this easy for you by handling project creation, compilation and bundling. All you have to do is get creative!

## Features

- Use standard Go for the backend
- Use any frontend technology you are already familiar with to build your UI
- Quickly create rich frontends for your Go programs using pre-built templates
- Easily call Go methods from Javascript
- Auto-generated Typescript definitions for your Go structs and methods
- Native Dialogs & Menus
- Native Dark / Light mode support
- Supports modern translucency and "frosted window" effects
- Unified eventing system between Go and Javascript
- Powerful cli tool to quickly generate and build your projects
- Multiplatform
- Uses native rendering engines - _no embedded browser_!

### Roadmap

The project roadmap may be found [here](https://github.com/wailsapp/wails/discussions/1484). New functionality and public behaviour changes use a [WEP (Wails Enhancement Proposal)](v3/wep/README.md), submitted as a draft pull request rather than a feature-request issue.

## Getting Started

Wails has two active versions:

| Version | Status | Install | Docs |
|---|---|---|---|
| v2 | Stable | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` | [wails.io](https://wails.io/) |
| v3 | Beta | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` | [v3.wails.io](https://v3.wails.io/) |

Full installation instructions are available for [v2](https://wails.io/docs/gettingstarted/installation) and [v3](https://v3.wails.io).

## Sponsors

This project is supported by these kind people / companies:
<img src="website/static/img/sponsors.svg" style="width:100%;max-width:800px;"/>

## Powered By

[![JetBrains logo.](https://resources.jetbrains.com/storage/products/company/brand/logos/jetbrains.svg)](https://jb.gg/OpenSource)

## FAQ

- Is this an alternative to Electron?

  Depends on your requirements. It's designed to make it easy for Go programmers to make lightweight desktop
  applications or add a frontend to their existing applications. Wails does offer native elements such as menus
  and dialogs, so it could be considered a lightweight electron alternative.

- Who is this project aimed at?

  Go programmers who want to bundle an HTML/JS/CSS frontend with their applications, without resorting to creating a
  server and opening a browser to view it.

- What's with the name?

  When I saw WebView, I thought "What I really want is tooling around building a WebView app, a bit like Rails is to
  Ruby". So initially it was a play on words (Webview on Rails). It just so happened to also be a homophone of the
  English name for the [Country](https://en.wikipedia.org/wiki/Wales) I am from. So it stuck.

## Stargazers over time

<a href="https://github.com/wailsapp/wails/stargazers">
  <img alt="Wails star history chart" src="website/static/img/star-history.svg" width="800" />
</a>

## Contributors

The contributors list is getting too big for the readme! All the amazing people who have contributed to this
project have their own page [here](https://wails.io/credits#contributors).

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwailsapp%2Fwails.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwailsapp%2Fwails?ref=badge_large)

## Inspiration

This project was mainly coded to the following albums:

- [Manic Street Preachers - Resistance Is Futile](https://open.spotify.com/album/1R2rsEUqXjIvAbzM0yHrxA)
- [Manic Street Preachers - This Is My Truth, Tell Me Yours](https://open.spotify.com/album/4VzCL9kjhgGQeKCiojK1YN)
- [The Midnight - Endless Summer](https://open.spotify.com/album/4Krg8zvprquh7TVn9OxZn8)
- [Gary Newman - Savage (Songs from a Broken World)](https://open.spotify.com/album/3kMfsD07Q32HRWKRrpcexr)
- [Steve Vai - Passion & Warfare](https://open.spotify.com/album/0oL0OhrE2rYVns4IGj8h2m)
- [Ben Howard - Every Kingdom](https://open.spotify.com/album/1nJsbWm3Yy2DW1KIc1OKle)
- [Ben Howard - Noonday Dream](https://open.spotify.com/album/6astw05cTiXEc2OvyByaPs)
- [Adwaith - Melyn](https://open.spotify.com/album/2vBE40Rp60tl7rNqIZjaXM)
- [Gwidaith Hen Fran - Cedors Hen Wrach](https://open.spotify.com/album/3v2hrfNGINPLuDP0YDTOjm)
- [Metallica - Metallica](https://open.spotify.com/album/2Kh43m04B1UkVcpcRa1Zug)
- [Bloc Party - Silent Alarm](https://open.spotify.com/album/6SsIdN05HQg2GwYLfXuzLB)
- [Maxthor - Another World](https://open.spotify.com/album/3tklE2Fgw1hCIUstIwPBJF)
- [Alun Tan Lan - Y Distawrwydd](https://open.spotify.com/album/0c32OywcLpdJCWWMC6vB8v)
