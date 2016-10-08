#! /bin/bash
enanos --random-size --max-sleep 3s --max-size 50KB -p 1337 &
./corcel run --progress bar --workers 10 --summary --duration 10s --plan .resources/sample-plan.yml
kill $!
