#! /bin/bash
enanos -p 1337 &
corcel run --progress bar --summary --duration 1s --plan .resources/sample-plan.yml
kill $!
