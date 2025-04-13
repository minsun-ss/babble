package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ListLibrariesOperation is the docstring details for the get endpoint
func ListLibrariesOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-libraries",
		Method:      http.MethodGet,
		Path:        "/libraries/",
		Summary:     "List all library documents",
		Description: "List all library documents available for the specific auth key. Filtering can be done via library.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func PostLibrariesOperation() huma.Operation {
	return huma.Operation{
		OperationID: "post-libraries",
		Method:      http.MethodPost,
		Path:        "/libraries/",
		Summary:     "Add new library",
		Description: "Add documentation for a new library version. If the library does not yet exist, information about the library is also added.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func GetLibraryOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-library",
		Method:      http.MethodGet,
		Path:        "/library/{libraryName}",
		Summary:     "Retrieve library details",
		Description: "Retrieve details about a specific library.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func GetLibraryVersionOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-library-version",
		Method:      http.MethodGet,
		Path:        "/library/{libraryName}/{libraryVersion}",
		Summary:     "Retrieve library version details",
		Description: "Retrieve details about a specific library version",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func PatchLibraryOperation() huma.Operation {
	return huma.Operation{
		OperationID: "patch-library",
		Method:      http.MethodPatch,
		Path:        "/library/{libraryName}",
		Summary:     "Update library details",
		Description: "Update specific details about a library.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func PatchLibraryVersionOperation() huma.Operation {
	return huma.Operation{
		OperationID: "patch-library-version",
		Method:      http.MethodPatch,
		Path:        "/library/{libraryName}/{libraryVersion}",
		Summary:     "Update library version details",
		Description: "Update specific details about a specific library version.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func DeleteLibraryOperation() huma.Operation {
	return huma.Operation{
		OperationID: "delete-library",
		Method:      http.MethodDelete,
		Path:        "/library/{libraryName}",
		Summary:     "Delete library series",
		Description: "Delete all details about a specific library. This will also effect a purge of all versions of this specific library.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func DeleteLibraryVersionOperation() huma.Operation {
	return huma.Operation{
		OperationID: "delete-library-version",
		Method:      http.MethodDelete,
		Path:        "/library/{libraryName}/{libraryVersion}",
		Summary:     "Delete specific library version",
		Description: "Delete a specific library version.",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}
