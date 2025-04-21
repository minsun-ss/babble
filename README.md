Things to do
- separate the backend from the front end: DONE (but not really)
- design new deployment using front and backend - DONE (for docker)
- design deployment for K8s/ArgoCD
- add the nextjs front end and redesign endpoints - DONE
- add a few views in the alembic schema because it makes sense to - TURNS OUT NOT NEEDED
- set up the current models in alembic - NOT NEEDED
- set up the dockerfile for the nextjs endpoint: DONE
- set up redirects for /docs endpoints to route directly to golang: DONE
- set up post endpoints in golang - IN PROGRESS
- add some useful workflows for github - DONE
- add golang tests for huma endpoints
- set up authentication middleware
- add a cli to handle jwt creations - DONE
- consolidate the test suite to reduce failure on container spinups (and we only need to use one db per test and reset between tests) - DONE
- redesign database to accommodate users and projects modification - DONE

# endpoints
/ (23456) - Babel's original front page
/ (3000) - Babel's new front page
/healthz - health and metric endpoints
/docs/ - Babel passing through to golang's endpoing serving files
/internal/ - Babel's communication layer between JS layer and Golang layer. Public
/api/v1/ - Babel's client facing api. Requires authentication.

# api endpoints
All api endpoints require a bearer token for authentication. One token per team.
/api/docs - api docs
/api/openapi - openapi yaml
(GET)   /api/v1/libraries - list all libraries and versions available for the team. Query filtering can be on libraryName, e.g., ?library=traderpythonlib,whatever
(POST)  /api/v1/libraries - add a new library version
(GET)   /api/v1/libraryName - get specific details about a libraryName
(PATCH) /api/v1/libraryName - update specific fields for a library in general
(DELETE) /api/v1/libraryName - delete the entirety of a specific library
(GET)   /api/v1/libraryName/version - get specific details about a library version
(PATCH) /api/v1/libraryName/version - update specific fields of a specific library version
(DELETE)  /api/v1/libraryName/version - remove specific library version


# Running this POS
That's right, I went down the rabbit hole that is NiceGui, then FastHTML + MonsterUI, and then static nextjs inside golang handlers before I said egh, this setup is too difficult to maintain, so they are now all completely separate applications.

# Layouts
Because I'm a lazy mf and I can't stand frontend, I use shadcn with NextJS. All data fields used for rendering are passed through to the golang backend.

# Alembic schema
Python. Let's be real, if I could make everything python, I would because it is not cognitive overload. I've also come around on some aspect of orms being useful, that is to say, try to make all your queries be some select * from table, even if you have to generate views for this.

# Logging
Through the power of opensearch and paranoia: structlog + opensearch.

# Note to self on uv
- if you want the .venv to be in the parent folder where the pyproject is in the child folder then:
  - uv .venv --project ./foldername/ (this will create a venv in parent where name is .venv)
  - activate the venv
  - you may need to force the VIRTUAL_ENV variable to be null, that is to say `export VIRTUAL_ENV=` because uv will not respect your current venv if there is a mismatch and use VIRTUAL_ENV instead.
  - uv pip install -e ./foldername/

# Some further notes to self on nginx / docker compose / reverse proxies
- path/to/blah is not the same as path/to/blah/ and you should validate that via curl like, all the time.
- your nginx conf should use the internal ports. your app should use the external ports. this will be a source of much frustration when trying to figure out your reverse proxy configuration unless, I guess, you actually do the reasonable thing and make your internal and external ports for your containers to be exactly the same.
