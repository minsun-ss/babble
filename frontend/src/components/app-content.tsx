import { useState } from "react";
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
      return (
        <>
          <blockquote>
            <i>
              "I repeat: In order for a book to exist, it is sufficient that it
              be possible. Only the impossible is excluded."
            </i>
            <footer>Jorge Luis Borges</footer>
          </blockquote>
          <h2>Introduction</h2>
          <p>
            Named after the short story "The Library of Babel" by Jorge Luis
            Borges, this attempts to consolidate all the documentation of
            TA-managed user-facing libraries into one place, including (if
            available) relevant changelogs and versioning.
          </p>
          <p>
            Deployment was in Golang + Mariadb + HTMX for the simplicity of this
            combo relative to its features before the author decided that she
            was going to commit to Javascript and fell headlong into NextJS
            (advance warning: don't do this). Now instead of one error, the
            author suffers random CORS errors intead all day, every day and went
            down the highway to hell that are reverse proxy configurations.
          </p>
          <h2>Contact</h2>
          <p>TA Global: TA.Global@flowtraders.com</p>
        </>
      );
    case "Contribution Guide":
      return (
        <div>
          <p>This is some other text</p>
        </div>
      );
    case "About":
      return (
        <div>
          <p>This is a lot of text to talk about stuff</p>
        </div>
      );
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
