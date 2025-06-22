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
      <p>Versioned library documentation.</p>
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
        the simplicity of its features and this was pretty great. Then nextjs
        came next and it's still a bit buggy. :(
      </p>
      <h2>Contact</h2>
      <p>TA Global: TA.Global@flowtraders.com</p>
    </>
  );
}
