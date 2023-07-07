# Habitify CLI
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0)
[![Twitter Follow](https://img.shields.io/twitter/follow/samy_osmium?style=for-the-badge)](https://twitter.com/intent/follow?screen_name=samy_osmium)
[![GitHub issues](https://img.shields.io/github/issues/samyosm/habitify-cli?style=for-the-badge)](https://github.com/samyosm/habitify-cli/issues)

Habitify CLI is a text user interface that allows you to access and manage your habits in the [Habitify](https://www.habitify.me/) habit tracker.
Use it to see the progress of your habits, change their status, and much more.

## Installation

### Go
```bash
go install github.com/samyosm/habitify-cli@latest
```

### Manual Build
```bash
git clone https://github.com/samyosm/habitify-cli.git
cd habitify-cli
go build
```

## Usage
It make use of the [Habitify API](https://docs.habitify.me/) so you will first need to obtain your api key.

### 1. Obtaining API Key 
From Mobile Apps(iOS, Android)
1. Open Settings
2. Open API Credential
3. Copy API Key or send to your desktop by tapping Send Via...

From Web App (unavailable, yet)
1. Open [API Credential](https://app.habitify.me/preference/api-credential) by navigating to Profile & Settings > API Credential
2. Copy API Key

### 2. Using API Key
Option 1: Making it an environement variable e.g.
```bash
HABITIFY_API_KEY=xxx habitify
```

Option 2: Put it in a configuration file at `$HOME/.config/habitify/config.yml`
```yml
api-key: xxx
```

Option 3: Use initialize command
```bash
habitify init xxx
habitify
```

## License
[GPL-3.0](./LICENSE)
