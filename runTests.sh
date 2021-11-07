#!/bin/bash

mode=$1
mode=${mode:-development}

apiBuildFolder="./src/MicroFileServer/"

if [[ "$mode" == "development" ]]; then
  apiBuildFolder="./src/MicroFileServer/"
  echo "Run tests on development machine. Use '$apiBuildFolder' to build 'api'
  Please ensure that you stop all instances of backend
  Tests require installed node js"
elif [[ "$mode" == "ci" ]]; then
  apiBuildFolder="./deploy/"
  echo "Run tests on CI machine. Use '$apiBuildFolder' to build 'api'"
else
  echo "incorrect mode '$mode'
  usage: ./runTests.sh [development|ci]
  - development is default - use api/Dockerfile for testing
  - ci is special for CI scenario - use deploy/Dockerfile for testing"
  exit 2
fi
export API_BUILD_FOLDER=$apiBuildFolder

logGroupStart() {
  if [[ "$mode" == "ci" ]]; then
    echo "##[group]$1"
  else
    echo $1
  fi
}
logGroupEnd() {
  if [[ "$mode" == "ci" ]]; then
    echo "##[endgroup]"
  fi
}

# Check testmace tool is installed

testMaceFindResult=$(npm list -g --depth=1 | grep 'testmace/cli@1.3.1')


if [[ -z "$testMaceFindResult" ]]; then
  echo "You must install @testmace/cli@1.3.1
  npm install --global @testmace/cli@1.3.1"
  exit 2
fi

dc="docker-compose -f docker-compose.yaml -f docker-compose.override.yaml  -f docker-compose.e2e.yaml -p mfs-e2e-tests"

logGroupStart "Ensuring to e2e environment is empty..."
$dc down -v
logGroupEnd

logGroupStart "Preparing images..."
$dc pull db
$dc build api
logGroupEnd

logGroupStart "Running db and api..."
$dc up -d db api
logGroupEnd

echo "Running tests..."
./tests/e2e/Testmace/waitWebApp.sh
if [[ "$mode" == "ci" ]]; then
  testmace-cli -e e2e -o tests-out --reporter=junit ./tests/e2e/Testmace/Project
else
  testmaceOut=$(testmace-cli -e e2e ./tests/e2e/Testmace/Project)
  echo "$testmaceOut"
  mkdir -p tests-out
  echo "$testmaceOut" > "tests-out/$(date +"%Y-%m-%d-%I-%M-%p").log"
fi

logGroupStart "Show api logs..."
$dc logs api
logGroupEnd

logGroupStart "Show db logs..."
$dc logs db
logGroupEnd

logGroupStart "Clean up e2e environment..."
$dc down -v
logGroupEnd
