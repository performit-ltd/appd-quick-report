<a name="readme-top"></a>



<!-- PROJECT LOGO -->
<br />
<div align="center">

  <h3 align="center">AppD Quick Report Builder</h3>

  <p align="center">
    Automated generation of AppD statistics reports.
    <br />
    <br />
    ·
    <a href="https://github.com/sivanovie/appd-quick-report/issues">Report Bug / Request Feature</a>
    ·
  </p>
</div>

<!-- Overview -->
## Overview

### Capabilities

* Generate Excel .xlsx report file for a given Controller instance.
* Include APM application statistics for every application - Number of Errors, Number of Calls and number of health rules (by status i.e. active/inactive).
* Use a config file to customise the report outlook.

<!-- Usage -->
## Usage

### Download

* Download one of the available releases from release directory.
* MacOS: 
* Linux (Fedora/RHEL/CentOS): 
* Download configuration file conf.yaml.

### Configure

* Edit conf.yaml
* Read the comments for every flag, it is self-explainable

### Run

* This is a standard OS executable file, so run as any other executable: ./appd-stats
* Program expects conf.yaml to be present in same dir, where the executable is.

### Troubleshoot

* The program generates a log called appd-stats.log.
* Check log for errors and issues.
