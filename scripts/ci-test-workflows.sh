echo "{  
    \"local-test\": true,
    \"inputs\": {
        \"branch\": \"main\"
    }
}" >.test_filesystem/act-event.json

docker run --rm \
    -v $PWD:$PWD -w $PWD \
    -v /var/run/docker.sock:/var/run/docker.sock \
    scripts/act -e .test_filesystem/act-event.json $@