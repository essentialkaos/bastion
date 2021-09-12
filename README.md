<p align="center"><a href="#readme"><img src="https://gh.kaos.st/bastion.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/bastion/ci"><img src="https://kaos.sh/w/bastion/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/bastion"><img src="https://kaos.sh/r/bastion.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/bastion"><img src="https://kaos.sh/b/a35d2d0c-7416-4287-84bb-a6919d894271.svg" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/w/bastion/codeql"><img src="https://kaos.sh/w/bastion/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

**Bastion** is utility for temporary access limitation to server.

### Installation

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install bastion
```

### Usage

* Configure your Bastion instance through configuration file (`/etc/bastion.knf`)
* Start Bastion daemon by command `sudo service bastion start` (even if you use CentOS 7)
* After start daemon return unique URL for enabling bastion mode
* Send any request (`GET`/`POST`/`HEAD`/etc...) to generated URL

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/bastion/ci.svg?branch=master)](https://kaos.sh/w/bastion/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/bastion/ci.svg?branch=master)](https://kaos.sh/w/bastion/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
