#!/bin/bash
echo "{\"user\":\"foo\",\"id\":\"goo\",\"description\":\"test image\",\"main-image\":\"`base64 $1`\",\"label-image\":\"`base64 $2`\"}" > $3
