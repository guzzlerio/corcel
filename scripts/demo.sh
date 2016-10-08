#! /bin/bash
enanos --random-size --max-sleep 3s --max-size 1MB -p 1337 &
./corcel run --progress bar --summary --duration 10s --plan .resources/sample-plan.yml
kill $!
