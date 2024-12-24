# gotta start somewhere

Things to do
- convert everything to slog - done
- add a correlation id to all requests - done
- set up middleware
- parse the library versions
- add container testing
- finish doing tests for handlers
- add embed for static files
- write a valid 404 page

# Babel

Named after the short story Library of Babel by Jorge Luis Borges, this is an attempt to store
versioned documentation for TA-maintained user-facing libraries. The front end is designed to be
minimalist and require little upkeep beyond deployment except for aesthetic preferences.

# Documentation

Documentation by user-facing libraries are usually in the form of auto-generated HTML coming from
documentation generators such as Sphinx (python), godoc (golang), Doxygen (C++), etc. Out of scope are
documents that are tied to deployment such as REST apis, such as those coming from FastAPI, Swagger,
etc.

# Testing

Since the author hates mocks, testing is done via the testcontainers package instead. Further
integration by adapting golembic instead of alembic for migrations would make this an easier build,
but at the moment do not forsee this to need any serious refactoring.
