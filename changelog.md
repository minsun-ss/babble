# Changelog

All notable changes to this project will be documented in this file.

## [2.5.0] - 2025-04-21

### Bug Fixes

- Bugfix on arrays returning "" on nonexistent versions
- Fix the failing tests
- Updated gitignore to only ignore certain binaries

### Features

- Setting up human api and redirecting paths to /api endpoint
- Setup generic API response wrapper
- Updated readme for the endpoints desired
- Update the viper check for the private key
- Setup jwt, cli, rearranged backend to accomodate
- Adjusted the database api to create the users and api setup
- Now with the correct modelling for constraints
- Building out the admin cli
- Updated the test suite to streamline the spinup of containers
- Updated alembic user table to contain roles
- Set up tests for creating a user
- Integration test on creation of a user
- Added delete user and its accompanying tests
- Added project to cli
- Made passing tests for create and delete project endpoint
- Added project access grant to the cli
- Set up retrieve api key endpoint

### Miscellaneous Tasks

- Version bump
- Updated readme

## [2.4.0] - 2025-04-12

### Bug Fixes

- Updated tests
- Golang tests
- Indentation
- Path
- Combined download & build
- Permissions on the tests
- The working directory means i don't need cd i think
- Removed the build part
- Fixing the liveness test check
- Rename/consolidate some handlers

### Features

- Trying to set up tests
- Moving some build files into the appropriate directories
- Adjusted the alembic models to account for new team accounts
- Setting up the external api to be used by third parties
- Adding bcrypt to allow hash checks
- Cleanup a little in the golang library
- Setting up for the /docs endpoint for the rest api

### Miscellaneous Tasks

- Updated version
- Updated readme

## [2.3.0] - 2025-04-11

### Features

- Set up app constants
- Updated the app content to reflect the constants
- Figured out how to get the app main to generate the correct links
- Fixing order on rendering some cleanup on library items

### Miscellaneous Tasks

- Version bump, changelog

## [2.2.0] - 2025-04-01

### Bug Fixes

- Updated table but idk for
- Typescript loves types, naturally

### Features

- Fixed the precommit configuration for ruff and golang and ran changes
- Updated the yaml hook check
- Slowly going through the front end app
- Remove superfluous images and tsx
- Updated templates
- New content tsx
- Added alembic schema for project key and project name
- Set up two endpoints for my lovely nextjs app

### Miscellaneous Tasks

- Version bump

## [2.1.5] - 2025-03-30

### Bug Fixes

- Remove extraneous files
- Update the check versions to only work on branch updates

### Features

- Change name, update version
- Removed extraneous dist
- Set up the model and upgraded accordingly

### Miscellaneous Tasks

- Updated gitignore
- Update changelog

## [2.1.3] - 2025-03-30

### Bug Fixes

- On merge

### Miscellaneous Tasks

- Version bump

## [2.1.2] - 2025-03-30

### Bug Fixes

- Fix tag

### Miscellaneous Tasks

- Updated version
- Changelog

## [2.1.1] - 2025-03-30

### Bug Fixes

- Check the script for origin master
- Try one more time on the head

### Features

- Add version check to github action
- Updated version check
- Updated the check versions action
- Added a token to github action to tag
- Updated versions with better logging
- Fixed the check version build action

### Miscellaneous Tasks

- Updated readme todo
- Formatting, changelog

## [2.1.0] - 2025-03-30

### Features

- Set up the nextjs app dockerfile and accompanying configs
- Adjusted package names so that they're all consistent with each other
- Updated makefile test containers
- Updated backend dockerfile
- Set up docker compose with direct
- Updated webserver to expose to 23456
- Fixed the port mapping for the backend
- Updated github actions

### Miscellaneous Tasks

- Updated README
- Update changelog

## [2.0.0] - 2025-03-30

### Bug Fixes

- Forgot to attach the pvc to the /tmp
- Bump image up one to change logging levels
- Deployment?
- Removed main files

### Features

- Add readme
- First commit
- Starting moving out the models and handlers elsewhere
- Set up database connections
- My page is gonna be amazing
- Renamed everything to babel
- Removed extraneous zippygo go file
- Added handlers for zipfiles
- Moved out to separate sections
- Set up the more info section
- Set up gorm
- Updated makefile
- Wow the database works
- Cleaning up connections
- Finished generating the templates page
- Updated
- Fixed the path on the test files
- Make cleanup its own handler
- Added a simple cleanup handler
- Updated the page to look better and swapped to slog
- Removed the use of gorm.go
- Removed the use of gorm.go
- Fixed css and maybe the js
- Fix logging
- Moved around structures, set up middleware
- Add health and logging handlers
- Updated go mod
- Updated configuration to also hold db
- Added healthz check
- Renamed menu.go to index.go
- Added some more metrics on middleware logging
- Updated handlers logging
- Still not sure where to put these functions
- Updated all the items in the main plus go mod
- Added documentation for config go
- Updated env go with panics
- Updated handlers
- Updated configuration
- Wrote a cool testing library here
- Updated zipfiles handler
- Added config.go
- Renamed some files
- Setting up file cleanup handlers
- Streamlined configuration file, updated tests and libraries to reflect this refactor
- Fixed the configuration file
- Figured out how to do embeds
- Added a liveness check test
- Set up test harness as separate item
- Updated readme
- Added go.sum for version control
- Test env
- Now i'm endlessly yak shaving, should work on other things
- Added a fixed home button
- Some adjustments to font colors
- Now made the content of the front page dynamic
- Streamlined the embeds
- More yak shaving
- Prune stale files
- Prune stale files
- Pruning stale files
- Remove stale files
- Remove stale directories
- Remove stale files
- More golike?
- Stale folder
- Stale folders
- Remove stale folders
- Stale folders
- Stale folders
- Consolidated static to assets folder
- Set up paths
- Moved static files out
- Precommit fixes
- Precommits
- Updated makefile
- Added prometheus endpoint
- Adding the ability to serve on 80 and 443, added version
- Set up the image for docker push
- Updated dockerfile and commands
- Added deployment folder
- Adjusted the deploys folder
- Deployment fix to yaml
- Changed the path to the image
- Changed deployment yaml
- Fixed the secrets handling of the env
- Fml don't compile images between mac and linux for k8s
- Updating logging for 1.0.4
- Set up the node js static html/css
- Removed 443 since i'm using https redirect in the k8s deployment
- Updated the makefile shortcuts and the dockerfile
- Moved entirety of golang to backend folder
- Added frontend
- Removed stray dockerfile in parent
- By your powers combined, an app
- Also added this
- Set up git cliff
- Updated and added changelog
- Workflows
- Workflows setup
- Set up models
- Updated schemas
- Updated gitignore
- Updated title to the frontend

<!-- generated by git-cliff -->
