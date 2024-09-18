#!/bin/sh

pid=$(ps x | awk '/swagger.*8085/ && !/awk/ {print $1}')
kill $pid
