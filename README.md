# App Nerds Web Application Generator

The **App Nerds Web Application Generator**, or *anwag*, is a tool to generate and scaffold a web application written in Go and and VueJS. The stack looks a little something like the following.

* Go 1.16 for building your API
* VueJS 2 for building reactive Javascript applications
* Bootstrap 4 for a great CSS foundation
* FontAwesome 5 for an icon library
* moment.js for date handling
* vue-loading-overlay to provide a spinner library. Already wired up for HTTP requests!
* vue-resource 1.5 for easy HTTP AJAX stuff
* vue-router 3 for Single Page Application routing
* vue-session for persisting data to local storage
* vuex 3 for state management

## Installation

```
go get -u github.com/app-nerds/anwag
```

## Usage

To use **anwag** simply run the following in your terminal.

```bash
anwag
```

You will receive a series of prompts. Answer the questions. Once complete you should see a screen like the one below.

![ANWAG Screenshot](anwag-screenshot.png)

