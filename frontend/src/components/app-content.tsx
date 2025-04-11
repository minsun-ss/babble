import { useState } from "react";
import { renderIndex, renderAbout } from "@/components/app-constants";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

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
  const data = [
    {
      library: activeContent,
      project_team: "myeh",
      description: "this is more content but no clue what it should be",
      versions: ["a", "b"],
    },
  ];

  const url = "http://localhost:23456/api/links/" + activeContent;
  fetch(url)
    .then((response) => response.text())
    .then((text) => console.log(text));

  return (
    <>
      <div>
        <h2>
          {data[0].project_team}: {data[0].library}
        </h2>
        <p>{data[0].description}</p>
      </div>

      <Table className="w-3/4 p-2">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Version</TableHead>
            <TableHead className="text-right">Link</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data[0].versions.map((datum, index) => {
            const docUrl = "/docs/" + data[0].library + "/" + datum + "/";
            return (
              <TableRow key={index}>
                <TableCell className="font-medium">{datum}</TableCell>
                <TableCell className="text-right">
                  <a href={docUrl}>{datum}</a>
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </>
  );
}
