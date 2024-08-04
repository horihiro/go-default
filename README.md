# tasksjson-updater

CLI tool for update default values of each input in `.vscode/tasks.json`.

## Demo
https://github.com/user-attachments/assets/e1f316e0-3817-4cca-a24c-cf57211a9a01

## Usage

```
Usage: tasksjson-updater --target-file TARGET-FILE --default-value DEFAULT-VALUE

Options:
  --target-file TARGET-FILE, -t TARGET-FILE
                         file path of target 'tasks.json' 
  --default-value DEFAULT-VALUE, -v DEFAULT-VALUE
                         id and default values to update, the format of 'DEFAULT-VALUE' is '${ID}=${DEFAULT_VALUE}'
  --help, -h             display this help and exit
```

See [sample](./sample/.vscode/tasks.json#L17)

## Background
The value of `default` property of the each element of `inputs` in tasks.json is fixed.

Many developers want to remenber the value from user input, then some similar issues were filed and closed without 

## Limitation

Currently, this tool only update the value of `default` property of the element of `inputs`.  
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