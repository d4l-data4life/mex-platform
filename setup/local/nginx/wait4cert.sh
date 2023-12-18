#!/bin/sh

echo "waiting for intermediate certificate files to appear"

while [ ! -f /pairgen/inter.pem ]; do
    echo "/pairgen/inter.pem not yet there"
    sleep 1
done

while [ ! -f /pairgen/inter.key.pem ]; do
    echo "/pairgen/inter.key.pem not yet there"
    sleep 1
done

echo "both files present, continuing with startup"
