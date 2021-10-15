# ant

CLI application that manipulates **ant cli specification document**. 
<br/>The **ant CLI specification document** is a file that describes a cli application interface.

The ant CLI has four (4) essential capabilities which are:
- Lint an ant CLI specification document represented inm **json** or **yaml** format
- Generate the CLI project from an ant CLI specification document
- Generate the integration test project from an ant CLI specification document
- Generate an HTML document that describes an CLI from the ant CLI specification document

PS: The current version of the CLI only supports the first capability.

## Installation

```sh
yarn add command-line-application
# or
npm i --save command-line-application
```

## Usage
The current version of ant CLI defines two (2) commands which are:
- export - Exports an ant CLI object into a specific file
- lint - Verifies if a specific ant CLI document complies with the ant CLI document schema

### Lint 
The lint command consumes a file **json** or **yaml** in order to search for any deviation of ant CLI document schema.
```sh
    ant lint [path-to-file]
```
The argument **path-to-file** specifies the file which will be consumed. In case the argument isn't specified, the CLI assumes the working directory **index.json** as default.

### Export
The export command exports an object into a file as shown bellow:

```sh
    ant export [object-type] [path-to-file]
```

The current version can only export schema as **object-type** , as shown bellow:
```sh
    ant export schema [path-to-file]
```
The argument **path-to-file** specifies where the schema is exported to. In case the argument isn't specified, the CLI assumes schema.json.
