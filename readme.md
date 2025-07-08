# syncsh

A command-line tool for synchronizing shell sessions across multiple machines in a network using WireGuard VPN tunnels.

## Overview

syncsh enables real-time synchronization of shell command history between multiple machines, allowing developers to maintain consistent command history across their development environments. The tool establishes secure WireGuard VPN connections between machines and monitors shell history files for changes, propagating updates across the network.

## Features

- **Secure Network Tunneling**: Uses WireGuard VPN for encrypted peer-to-peer connections
- **Real-time History Synchronization**: Monitors shell history files and synchronizes changes across connected machines
- **Cross-platform Support**: Written in Go for compatibility across different operating systems
- **Configurable Interface**: Customizable WireGuard interface names and history file paths
- **Automatic Diff Detection**: Intelligent diffing system to track and merge history changes
- **SQLite Storage**: Local database for configuration and state management

## Architecture

The application consists of several key components:

- **Command Interface**: Cobra-based CLI with `init` and `connect` commands
- **Configuration Management**: YAML-based configuration with SQLite backend
- **Network Layer**: WireGuard tunnel management and peer connectivity
- **File Monitoring**: Real-time file system watcher for history file changes
- **History Diffing**: Semantic diff algorithm for merging shell history changes
- **Secret Management**: Secure handling of WireGuard private/public key pairs

## Installation

```bash
go install github.com/TheRealSibasishBehera/syncsh@latest
```

## Usage

### Initialize a Machine

Set up syncsh on a machine to prepare it for synchronization:

```bash
syncsh init [flags]
```

**Flags:**
- `--history-path`: Custom path to shell history file (default: auto-detect)
- `--interface`: WireGuard interface name (default: "syncsh0")

### Connect to Remote Machine

Connect to a remote syncsh-enabled machine:

```bash
syncsh connect
```

## Configuration

syncsh uses a YAML configuration file with the following structure:

```yaml
sql_path: syncsh.db
history: /path/to/shell/history
interface: syncsh0
```

## Technical Details

### Network Protocol

- **VPN Technology**: WireGuard for secure tunneling
- **Default Port**: 51820
- **Keepalive Interval**: 25 seconds
- **MTU**: Configurable (default: WireGuard default)

### Security

- **Key Generation**: Automatic WireGuard key pair generation
- **Endpoint Discovery**: Dynamic endpoint resolution
- **Encrypted Communication**: All traffic encrypted via WireGuard

### File Monitoring

- **Watch System**: Uses fsnotify for efficient file system monitoring
- **Diff Algorithm**: Semantic diff matching for intelligent history merging
- **Conflict Resolution**: Automatic handling of concurrent history changes

## Dependencies

- **WireGuard**: `golang.zx2c4.com/wireguard` for VPN functionality
- **Cobra**: `github.com/spf13/cobra` for CLI interface
- **fsnotify**: `github.com/fsnotify/fsnotify` for file system monitoring
- **go-diff**: `github.com/sergi/go-diff` for history diffing
- **go-yaml**: `github.com/goccy/go-yaml` for configuration parsing

## Development Status

This project is currently under active development. Core functionality is being implemented across the following areas:

- Network tunnel establishment and management
- Configuration system and persistence
- File monitoring and synchronization logic
- Command-line interface and user experience

## License

Licensed under the Apache License, Version 2.0. See LICENSE file for details.

## Contributing

This project is maintained by Sibasish Behera. Contributions, issues, and feature requests are welcome.