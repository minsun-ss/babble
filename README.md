Things to do
- convert everything to slog - done
- add a correlation id to all requests - done
- set up middleware - mostly done (prometheus left)
- set up liveness/readiness checks (if needed) - partly done
- parse the library versions - done
- add container testing - done
- finish doing tests for handlers - mostly done (missing zipfile)
- add embed for static files - done
- write a valid 404 page
- add prometheus checks as well
- figure out container within container testing (if needed - go check sdp to see how it was done there)

# Babel - What is it?

Named after the short story Library of Babel by Jorge Luis Borges, this is an attempt to store
versioned documentation for TA-maintained user-facing libraries. The front end is designed to be
minimalist and require little upkeep post deployment except for aesthetic preferences; it will update
accordingly to the files uploaded to the database.

Database schema is controlled by the repo babel-librarian via alembic in python. As it's essentially
just two tables, did not see it warranting the effort to leverage something like golembic for this.

# Documentation

Documentation by user-facing libraries are usually in the form of auto-generated HTML coming from
documentation generators such as Sphinx (python), godoc (golang), Doxygen (C++), etc. Out of scope are
documents that are tied to deployment such as REST apis, such as those coming from FastAPI, Swagger,
etc.

Documents are stored on the backend in an mariadb database as a zipped file in a
LONGBLOB binary field. As long as relative paths within the html are consistent within the zipped file,
babel should correctly serve the zipped website as a whole.

The database is intended to be shared among TA globally, although it's not mandatory, and the website
can be deployed to all regions with a specific regional specific URL if desired. There are some minor
features included in the database to allow for testing "in prod" so to speak: all docs whose "hidden"
field are set to 1 do not display or generate a menu item on the front page but can be nonetheless
transversed to as long as there are binaries stored in docs_history. This was mostly designed to
a) allow for deprecation and b) also allow for minor prod-validation as well with a small audience.
There is a test harness instead that spins up a local docker database that you can also use to serve
for testing if that is desired and the webpage served from localhost.

# Updating the front page

Content is determiend by the fixed content in /static/indexContent.html file. Feel free to change it
up as needed.

# Adding documents to babel during CI/CD

The repo babel-librarian includes the associated scripts to attach to the repo to push. Deployment
should still be done manually after a merge to master.

Further extension of this project could be to deploy a REST API to post updates/additions to the database,
but adds materially more overhead relative to benefit. If it scales outside of TA, then possibly
there is a case.

# Features

Some additional nice to haves
- healthz checks
- middleware logging
- prometheus latency logging

# Testing

Since the author hates mocks, testing is done via the testcontainers package instead. Further
integration by adapting golembic instead of alembic for migrations would make this an easier build,
but at the moment do not forsee this to need any serious refactoring.

# Notes
This project was mostly entirely written in Zed, which as an editor goes, is pretty nice, lightweight,
and works really well with golang.
