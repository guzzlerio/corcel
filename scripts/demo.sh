#! /bin/bash
enanos --random-size --max-sleep 3s --max-size 50KB -p 1337 &
./corcel run --workers 50 --summary --duration 10s --plan .resources/sample-plan.yml
kill $!
