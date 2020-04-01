# Sops64

[![Build Status](https://travis-ci.org/NoUseFreak/sops64.svg?branch=master)](https://travis-ci.org/NoUseFreak/sops64)

Sops wrapper that does base64 encoding and decoding for you.

## Usage

```
sops64 --encrypt tests/plain.yml
sops64 --decrypt tests/sops.yml
```

## Install

### Official release

Download the latest [release](https://github.com/NoUseFreak/sawsh/releases).

```bash
brew install nousefreak/brew/sops64
```

or

```bash
curl -sL http://bit.ly/gh-get | PROJECT=NoUseFreak/sops64 bash
```

### Build from source

```sh
$ git clone https://github.com/NoUseFreak/sops64.git
$ cd sops64
$ make
$ make install
```

### Upgrade

To upgrade to the latest repeat the install step.
