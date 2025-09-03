# go-markdown-to-notion
Markdown to Notion converter.
Append markdown content to a specified Notion block.

<p> 
  <img src="./images/image1.png" width="300">
  <img src="./images/image2.png" width="300">
</p>

# Usage
```shell
# help
go-markdown-to-notion help

# delete existing block children
go-markdown-to-notion delete-all-blocks --notion-page-or-block-id xxxxx

# convert and upload markdown file to Notion
go-markdown-to-notion upload --notion-block-id xxxxx --source-md-filepath sample.md
```

# License
The MIT License

Copyright Shohei Koyama / sioncojp

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
