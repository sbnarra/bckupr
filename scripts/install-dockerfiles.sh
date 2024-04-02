#!bash

tools=""
for dockerfile in $(find scripts -name "*.Dockerfile"); do
    tool=$(echo $dockerfile | cut -d'.' -f1)
    tools="$tool $tools"
    docker build . -t $tool -f $dockerfile
done

echo
echo "images [ $tools] installed"