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

# Babel

Named after the short story Library of Babel by Jorge Luis Borges, this is an attempt to store
versioned documentation for TA-maintained user-facing libraries. The front end is designed to be
minimalist and require little upkeep post deployment except for aesthetic preferences; it will update
accordingly to the files uploaded to the database.

# Documentation

Documentation by user-facing libraries are usually in the form of auto-generated HTML coming from
documentation generators such as Sphinx (python), godoc (golang), Doxygen (C++), etc. Out of scope are
documents that are tied to deployment such as REST apis, such as those coming from FastAPI, Swagger,
etc.

Documents are stored on the backend in an mariadb database as a zipped file in the
LONGBLOB binary field. As long as relative paths within the html are consistent within the page, babel
should correctly serve the zipped website as a whole.

# Adding documents to babel during CI/CD

The repo babel-librarian includes the associated scripts to attach to the repo to push. Deployment
should still be done manually after a merge to master.

Further extension of the build could be to deploy a REST API to post updates/additions to the database,
but do materially add more overhead.

# Testing

Since the author hates mocks, testing is done via the testcontainers package instead. Further
integration by adapting golembic instead of alembic for migrations would make this an easier build,
but at the moment do not forsee this to need any serious refactoring.
