<p align="center"><a href="#readme"><img src="https://gh.kaos.st/bastion.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://github.com/essentialkaos/bastion/actions"><img src="https://github.com/essentialkaos/bastion/workflows/CI/badge.svg" alt="GitHub Actions Status" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/bastion"><img src="https://goreportcard.com/badge/github.com/essentialkaos/bastion" alt="GoReportCard" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-bastion-master"><img alt="codebeat badge" src="https://codebeat.co/badges/a35d2d0c-7416-4287-84bb-a6919d894271" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

**Bastion** is utility for temporary access limitation to server.

### Installation

#### From ESSENTIAL KAOS Public repo for RHEL6/CentOS6

```bash
[sudo] yum install -y https://yum.kaos.st/kaos-repo-latest.el6.noarch.rpm
[sudo] yum install bastion
```

#### From ESSENTIAL KAOS Public repo for RHEL7/CentOS7

```bash
[sudo] yum install -y https://yum.kaos.st/kaos-repo-latest.el7.noarch.rpm
[sudo] yum install bastion
```

### Usage

* Configure your Bastion instance through configuration file (`/etc/bastion.knf`)
* Start Bastion daemon by command `sudo service bastion start` (even if you use CentOS 7)
* After start daemon return unique URL for enabling bastion mode
* Send any request (`GET`/`POST`/`HEAD`/etc...) to generated URL

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://github.com/essentialkaos/bastion/workflows/CI/badge.svg?branch=master)](https://github.com/essentialkaos/bastion/actions) |
| `develop` | [![CI](https://github.com/essentialkaos/bastion/workflows/CI/badge.svg?branch=develop)](https://github.com/essentialkaos/bastion/actions) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
