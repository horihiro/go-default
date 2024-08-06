# go-default

CLI tool for update default values of each input in `.vscode/tasks.json`.

## Demo
https://github.com/user-attachments/assets/e1f316e0-3817-4cca-a24c-cf57211a9a01

## Usage

```
Usage: default --target-file /PATH/TO/tasks.json --set id1=value1 id2=value2 ...

Options:
  --target-file /PATH/TO/tasks.json, -t /PATH/TO/tasks.json
                         file path of target 'tasks.json'
  --set id1=value1 id2=value2 ..., -s id1=value1 id2=value2 ...
                         pairs of id and default values to update
  --help, -h             display this help and exit
```

See [sample](./sample/.vscode/tasks.json#L17)

## Background and motivation
The value of `default` property of the each element of `inputs` in tasks.json is fixed.

Many developers want to remenber the value from user input, then some similar issues were filed and closed.

  - https://github.com/microsoft/vscode/issues/65066
  - https://github.com/microsoft/vscode/issues/78213
  - https://github.com/microsoft/vscode/issues/78422
  - https://github.com/microsoft/vscode/issues/72944

This tool might be an approach of the solution for these issues. 

## Limitation

Currently, this tool only UPDATE the value of `default` property of the element of `inputs`.  
This means this tool doesn't create the `default` property when the element doesn't have the property as following.

```json
  "version": "2.0.0",
  "tasks": [
    // ...
  ],
  "inputs": [
    {
      "id": "input1",
      "description": "...",
      "type": "promptString",
      // ...
    }
  ]
```