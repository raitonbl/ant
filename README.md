# ant

CLI application that manipulates **ant cli specification document**. 
<br/>The **ant CLI specification document** is a file that describes a cli application interface.

The ant CLI has four (4) essential capabilities which are:
- Lint an ant CLI specification document represented inm **json** or **yaml** format
- Generate the CLI project from an ant CLI specification document
- Generate the integration test project from an ant CLI specification document
- Generate an HTML document that describes an CLI from the ant CLI specification document

PS: The current version of the CLI only supports the 1st capability.

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


## document example
ant CLI document that describe the CLI tool in yaml format:
```yaml
name: ant
version: 1.0.0
description: manipulates ant cli specification document
commands:
  - name: lint
    description: >-
      allows the validation of an CLI specification file
    parameters:
      - in: arguments
        name: path-to-file
        description: the CLI specification file URI
        schema:
            type: string
    exit:
      - code: 0
        message: Document is valid
      - code: 1
        message: unexpected problem occurred
      - code: 2
        message: Document isn't valid
  - name: export
    description: exports an ant object
    parameters:
      - in: arguments
        index: 0
        name: object-type
        description: the CLI specification file URI
        schema:
          type: string
          enum:
            - schema
      - in: arguments
        index: 1
        name: path-to-file
        description: path to the file to export to
        schema:
          type: string
    exit:
      - code: 0
        message: Successfully exported
      - code: 1
        message: unexpected problem occurred
```
ant CLI document that describe the CLI tool in json format:
```json
{
  "name": "ant",
  "version": "1.0.0",
  "description": "manipulates ant cli specification document",
  "commands": [
    {
      "name": "lint",
      "description": "allows the validation of an CLI specification file",
      "parameters": [
        {
          "in": "arguments",
          "name": "path-to-file",
          "description": "the CLI specification file URI",
          "schema": {
            "type": "string"
          }
        }
      ],
      "exit": [
        {
          "code": 0,
          "message": "Document is valid"
        },
        {
          "code": 1,
          "message": "unexpected problem occurred"
        },
        {
          "code": 2,
          "message": "Document isn't valid"
        }
      ]
    },
    {
      "name": "export",
      "description": "exports an ant object",
      "parameters": [
        {
          "in": "arguments",
          "index": 0,
          "name": "object-type",
          "description": "the CLI specification file URI",
          "schema": {
            "type": "string",
            "enum": [
              "schema"
            ]
          }
        },
        {
          "in": "arguments",
          "index": 1,
          "name": "path-to-file",
          "description": "path to the file to export to",
          "schema": {
            "type": "string"
          }
        }
      ],
      "exit": [
        {
          "code": 0,
          "message": "Successfully exported"
        },
        {
          "code": 1,
          "message": "unexpected problem occurred"
        }
      ]
    }
  ]
}
```