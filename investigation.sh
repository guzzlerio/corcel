#! /bin/bash

output="investigation-results"
name="$1"

if [ -z "$name" ];then
    echo "ERROR: Must supply a name for the investigation"
    exit 1
fi

investigation_output_folder="${output}/${name}"
mkdir -p ${investigation_output_folder}

for i in {1..6}
do
    echo "run ${i} starting..."
    enanos -p 1337&

    ./corcel run --workers 10 --summary --duration 10s --plan .resources/sample-plan.yml >> "${investigation_output_folder}/run_$i.log"

    kill $!
done
