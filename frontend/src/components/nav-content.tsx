import { useState } from "react";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
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
            Deployment is Golang + Mariadb + HTMX for the simplicity of this
            combo relative to its features.
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
      const hellodata = [
        {
          title: activeContent,
          text: "this is more content but no clue what it should be",
          links: ["a", "b"],
        },
      ];
      return (
        <>
          <div>
            <h2>{hellodata[0].title}</h2>
            <p>{hellodata[0].text}</p>
          </div>

          <Table className="w-3/4 p-2">
            <TableHeader>
              <TableRow>
                <TableHead className="w-[100px]">Version</TableHead>
                <TableHead className="text-right">Link</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {hellodata[0].links.map((datum, index) => (
                <TableRow key={index}>
                  <TableCell className="font-medium">{datum}</TableCell>
                  <TableCell className="text-right">{datum}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </>
      );
  }
}
