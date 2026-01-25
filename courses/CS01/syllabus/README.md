# Base64

## Overview

Base64 is a binary-to-text encoding scheme that represents binary data in an ASCII string format. It is commonly used in various applications such as email, file transfer, and web development.

## Usage

Base64 encoding is used to convert binary data into a format that can be safely transmitted over a network or stored in a text file. It is also used to encode binary data into a format that can be displayed in a web browser.

## Examples

### Example 1: Encoding a String

```py,playground,editable,docker-id=python3
import base64

encoded_string = base64.b64encode(b'Hello, World!')
print(encoded_string)
```

```py,playground,editable,docker-id=python3
{{#include example.py}}
```

{{#task task01}}

<task id="task01"/>
