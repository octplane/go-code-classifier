#!/bin/sh

find '../src' -maxdepth 1 -type d -not -name "src" | xargs -n1 ln -s
# find '../src' -maxdepth 1 -type d -not -name "src" | xargs basename >> .gitignore

