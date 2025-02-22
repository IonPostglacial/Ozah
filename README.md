# What is this software?

Hazo is a Taxonomy software intending to be simple and ergonomic to use.
The plan is to have a Progressive Web App that works offline that can be used when internet connection is not great, for instand during field work, as well as a server side application that works even with JavaScript disabled (with gracefully degraded ergonomy) that can be used with outdated machines or machines with small computing power.

This repository currently contains the Server Side part of Hazo, the PWA will be merged in this repository once this Server Side app is mature enough.

# How it is developed?

The server side of Hazo is developed using Go, with very few dependencies to avoid bitrot. The dependencies we have should be easy to replace like the SQLite 3 driver, the crypto library and SQLC which is used to generate SQL Query boilerplate.

To generate the queries, go to the src/db folder and use the following command:

```shell
sqlc generate
```

UI Components are Go HTML templates and use Web Components and HTMX for progressive enhancements. The JavaScript and CSS code for the components are currently bundled in the executable and served as a single JS and a single CSS file.
The Air tool can be used to automatically reload the server on code change, by default the app will be server on localhost:8080.
In debug mode the web pages should refresh automatically when the server is updated.