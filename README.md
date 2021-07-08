# Squat

<p align="center">
  <img src="https://github.com/Shanduur/squat/blob/main/other/squat.png?raw=true"/>
</p>

[![license](https://img.shields.io/github/license/shanduur/squat?)](LICENSE)
![go version](https://img.shields.io/github/go-mod/go-version/shanduur/squat?)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanduur/squat?)](https://goreportcard.com/report/github.com/shanduur/squat)

Squat is an application that provides simple SQL data generation functionality. 

It generates synthetic SQL data based on the table definition, that is gathered from the DBMS. 
Squat supports IBM Informix and PostgreSQL, with planned support for all major databases, including MySQL, CockroachDB and MariaDB.

# Requirements

The Informix provider uses [alexbrainman/odbc](github.com/alexbrainman/odbc) package. This means, that for Linux and other \*NIX operating systems you have to install *unixODBC* application. Additionally, for the compilation, the development version of that is needed (e.g. *unixodbc-dev* on Debian). In Windows it calls directly to the *odbc.dll* - that also needs to be installed. Additionally, you have to provide your own Informix CSDK, that includes client driver for ODBC.

## Docker

For Docker image, as for now, you need to manually log in to the container, and install Informix CSDK, as it's non-free software distributed by IBM.

# Configuration

The app is configured using multiple small configuration files and environment variables. Each file is used for the database connection provider package. Unfortunately, now you can connect only to single database of given type.

Each config file is read from config location provided to the app through environmental variable.

- `CONFIG_LOCATION` - sets the directory containing the configuration files for providers.
- `DATA_LOCATION`- sets the directory containing the *data.json* and *data.gob* files.


For examples of the configuration files, look at the [bin/config](./bin/config) folder in the root of the repository.

Additionally you can provide your own *data.gob* file. The best way of doing this, is to parse the *JSON* file with *gob-generator*, which source code can be found in the tools directory.

In Docker all env variables are set by default, and all exemplary files are loaded into the specific dirs. You should use *Docker Compose* to create volumes and edit the configs from the host machine. 

# Usage

Usage of the application is quite simple. If you are running it bare-metal, just run the app, with the proper configuration files and environmental variables. Then head to [localhost:8080](localhost:8080) (or change the port correspondingly to the port you have passed with `--port` argument), and start using the web user interface.
