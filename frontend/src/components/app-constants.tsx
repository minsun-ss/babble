/**
 * Fetches the content for index route ("/").
 * @returns {React.ReactElement} - the fields to be rendered
 */
export function renderIndex(): React.ReactElement {
  return (
    <>
      <blockquote>
        <i>
          "I repeat: In order for a book to exist, it is sufficient that it be
          possible. Only the impossible is excluded."
        </i>
        <footer>Jorge Luis Borges</footer>
      </blockquote>
      <h2>Introduction</h2>
      <p>
        Named after the short story "The Library of Babel" by Jorge Luis Borges,
        this attempts to consolidate all the documentation of TA-managed
        user-facing libraries into one place, including (if available) relevant
        changelogs and versioning.
      </p>
    </>
  );
}

/**
 * Fetches the content for about route ("/about").
 * @returns {React.ReactElement} - the fields to be rendered
 */
export function renderAbout(): React.ReactElement {
  return (
    <>
      <p>
        The first iteration of this page was done in golang + mariadb + htmx for
        the simplicity of its features and this was pretty great. Then the
        author had a fit of madness and decided that Golang + Mariadb + NextJS
        was it (spoiler alert: don't do this). Now instead of some nice fixable
        generic golang error, now it's all CORS and reverse proxy configuration
        madness just for some generic shadcn layout, plus learning React the
        hard way (the golang part remains the same, at least). The long and
        short of it is: it is really really difficult to get away from
        Javascript if you need a pretty front end and although the path to all
        frameworks end in Javascript anyway, really try to stay away as long as
        you can if you don't need a super beautiful or fancy frontend.
      </p>
      <h2>Contact</h2>
      <p>TA Global: TA.Global@flowtraders.com</p>
    </>
  );
}
