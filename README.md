# Culsans

A command-line password manager using PGP encryption and Git version control.

**Culsans** is a secure, Git-backed password manager that encrypts passwords with PGP keys and maintains an encrypted password vault in a Git repository. It's rewriting of my old project [Janus](https://github.com/R2aper/Janus) that uses libgit2 + gpgmecpp so it was hard to use cross-platform.

## Features

- **PGP Encryption**: Encrypt passwords with public keys and decrypt with private keys
- **Git Integration**: Automatically commit password changes to Git
- **GPG Signing**: Optionally sign commits with your PGP key
- **Simple CLI**: Intuitive command-line interface for managing passwords

## Installation

### Build from Source

```bash
git clone https://github.com/R2aper/Culsans.git
cd Culsans
go build -o cl
```

## Usage 

```bash
Usage: cl [global-flags] <command> [args]
```

### Global Flags:

- `-pub string` Path to public PGP key file public key for add
- `-sec string` Path to public PGP key file public key for show and signing commits
- `-m string` Commit message
- `-q` Don't commit changes to git
- `-s` Sign commit
- `-h` Show this help message
- `-v` Show version

### Commands:

- `init` Initialize a new git repository
- `list` List all passwords in the vault
- `add <name>` Add a new password
- `remove <name>` Remove a password
- `show <name>` Show password content

## Example
```bash
  cl -pub ~/.pgp/pub.key add Google                 Add password with name `Google`
  cl -sec ~/.pgp/priv.key show Discord              Show password with name `Discord`      
  cl -pub ~/.pgp/pub.key -sec ~/.pgp/priv.key -s    add password with name `Github` and sign commit with private key 
```

