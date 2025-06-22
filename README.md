# Babble
Internal website to host versioned documentation across multiple teams.
- front end: nextjs
- backend: golang + mariadb db
- schema management: python

Includes:
- api to manage admin related items (e.g., creation/deletion of jwt)
- api to manage updating docs by teams
- k8s manifests for deployment

# TODO
- design deployment manifests for K8s/ArgoCD
- design api endpoints for updating libraries

# endpoints
/ (23456) - Babble's original front page
/ (3000) - Babble's new front page
/healthz - health and metric endpoints
/docs/ - Babble passing through to golang's endpoing serving files
/internal/ - Babble's communication layer between JS layer and Golang layer. Public
/api/v1/ - Babble's client facing api. Requires authentication.

# api endpoints
All api endpoints require a bearer token for authentication. One token per team.
/api/docs - api docs
(GET)   /api/v1/libraries - list all libraries and versions available for the team. Query filtering can be on libraryName, e.g., ?library=mylibrary
(POST)  /api/v1/libraries - add a new library version
(GET)   /api/v1/libraryName - get specific details about a libraryName
(PATCH) /api/v1/libraryName - update specific fields for a library in general
(DELETE) /api/v1/libraryName - delete the entirety of a specific library
(GET)   /api/v1/libraryName/version - get specific details about a library version
(PATCH) /api/v1/libraryName/version - update specific fields of a specific library version
(DELETE)  /api/v1/libraryName/version - remove specific library version

# Note to self on uv
- if you want the .venv to be in the parent folder where the pyproject is in the child folder then:
  - uv .venv --project ./foldername/ (this will create a venv in parent where name is .venv)
  - activate the venv
  - you may need to force the VIRTUAL_ENV variable to be null, that is to say `export VIRTUAL_ENV=` because uv will not respect your current venv if there is a mismatch and use VIRTUAL_ENV instead.
  - uv pip install -e ./foldername/

# Some further notes to self on nginx / docker compose / reverse proxies
- path/to/blah is not the same as path/to/blah/ and you should validate that via curl like, all the time.
- your nginx conf should use the internal ports. your app should use the external ports. this will be a source of much frustration when trying to figure out your reverse proxy configuration unless, I guess, you actually do the reasonable thing and make your internal and external ports for your containers to be exactly the same.
