# Patchachu

Go package for fetching, locally storing, and querying google patch status info

## Introducing Patchachu

<https://github.com/exokernel/patchachu>

This will be the main library that can be used in a command line tool or a service.

Patchachu uses GCP APIs to build a local database of patch deployment information of the GCP projects you point it at.

It includes a CLI tool with it called Patcha but the patchachu library itself is intended to be generally useful for your own services as well.

The first thing Patcha does is build a local database. Once the db is built you can ask it to tell you various information such as "is host h covered by a patch deployment and if so which?" or "show me all the host covered by patch deployment d" or "show me all the hosts not covered by a patch deployment" etc. Patchu can also generate JSON and CSV reports on the patch status of all the hosts in your project!

All you need to do is download a release of Patcha or install it via homebrew, apt, yum etc

Then set up a config to point it at your project and start querying it for patch info! Patcha also accepts configuration on the command line!
