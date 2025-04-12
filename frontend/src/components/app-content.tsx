import { useState, useEffect } from "react";
import { renderIndex, renderAbout } from "@/components/app-constants";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface LibraryItem {
  library: string;
  project_team: string;
  description: string;
  versions: string[];
}

/**
 * Fetches the content for subroute.
 * @param {string} activeContent - the data field to fetch
 * @returns {React.ReactElement} - the fields to be rendered
 */
export function renderContent(activeContent: string) {
  switch (activeContent) {
    case "index":
      return renderIndex();
    case "Contribution Guide":
      return (
        <div>
          <p>This is some other text</p>
        </div>
      );
    case "About":
      return renderAbout();
    default:
      return renderLibrary(activeContent);
  }
}

/**
 * Fetches the content for library subroute.
 * @param {string} activeContent - the data field to fetch
 * @returns {React.ReactElement} - the fields to be rendered
 */
export function renderLibrary(activeContent: string) {
  const url = `http://localhost:23456/internal/links/${activeContent}`;
  const [libraryData, setLibraryData] = useState<LibraryItem>({
    library: activeContent,
    project_team: "TBD",
    description: "TBD",
    versions: ["0.1.0"],
  });

  // re-renders on url changes
  useEffect(() => {
    fetch(url)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error: status: ${response.status}`);
        }
        return response.json();
      })
      .then((text) => {
        setLibraryData(text[0]);
      })
      .catch((error) => {
        console.error("Error: ", error);
      });
  }, [url]);

  // wait until fetch is done
  if (!libraryData) return <div>Loading...</div>;

  return (
    <>
      <div>
        <h2>
          {libraryData.project_team}: {libraryData.library}
        </h2>
        <p>{libraryData.description}</p>
      </div>

      <Table className="w-2/3 p-2">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Version</TableHead>
            <TableHead className="text-right">Link</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {libraryData.versions.map((datum, index) => {
            const docUrl = `/libraries/${libraryData.library}/${datum}/`;
            return (
              <TableRow key={index}>
                <TableCell className="font-medium">{datum}</TableCell>
                <TableCell className="text-right">
                  <a href={docUrl}>docs</a>
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </>
  );
}
