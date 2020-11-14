<!--
title: Zeebe Task
weight: 4705
-->
# Zeebe Task

**This plugin is in ALPHA stage**

This trigger allows you to listen Zeebe Task instance.

## Installation

### Flogo CLI
```bash
flogo install github.com/codelity-co/flogo-zeebeworkflow-activity
```

## Configuration

### Settings:
  | Name                | Type   | Description
  | :---                | :---   | :---
  | zeebeBrokerHost     | string | Zeebe broker host - ***REQUIRED***
  | zeebeBrokerPort     | int    | Zeebe broker port, default 26500 - ***REQUIRED***
  | bpmnProcessID       | string | BPMN process ID - ***REQUIRED***
  | serviceType         | string | BPMN Task Service Type  - ***REQUIRED***
  | command             | string | Zeebe command, Create or Cancel - ***REQUIRED***

### Handler Settings:
  None

### Output
  | Name                | Type   | Description
  | :---                | :---   | :---
  | data                | object | data object - ***REQUIRED***

### Reply:
  | Name          | Type   | Description
  | :---          | :---   | :---
  | status        | string | status text, ERROR or SUCCESS - ***REQUIRED***
  | result        | any    | trigger result

## Example

```json
{
  "id": "flogo-zeebeworkflow-trigger",
  "name": "Codelity Flogo Zeebe Task Trigger",
  "ref": "github.com/codelity-co/flogo-zeebe-trigger",
  "settings": {
    "zeebeBrokerHost": "localhost",
    "zeebeBrokerPort": 26500,
    "bpmnProcessID": "order-process",
    "serviceType": "order-service",
    "command": "Create"
  },
  "handlers": {
    "settings": {},
    "action": {
      "ref": "github.com/project-flogo/flow",
        "settings": {
          "flowURI": "res://flow:zeebe_to_some_flow"
        },
        "input":{
          "data": "=$.data"
        }
      }
    }
  }
}
```