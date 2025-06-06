#!/bin/bash
set -e

_si="si"
_go="golang"
starter="${_si}${_go}"

replace() {
  old="$1"
  new="$2"

  if [ -z "$old" ] || [ -z "$new" ]; then
    echo "Usage: $0 OLD_STRING NEW_STRING"
    exit 1
  fi

  # Detect macOS or GNU sed
  if sed --version >/dev/null 2>&1; then
    # GNU sed
    sed_cmd=(sed -i)
  else
    # BSD sed (macOS)
    sed_cmd=(sed -i '')
  fi

  # Replace in files containing the pattern
  git grep -l "$old" | while read -r file; do
    "${sed_cmd[@]}" "s/${old//\//\\/}/${new//\//\\/}/g" "$file"
  done
}

remote() {
  URL=$(git remote get-url origin)
  if [[ "$URL" == *"peruri-dev/sigolang"* ]]; then
    git remote remove origin
    git remote add sigolang "$URL"
  fi
}

read -p "Enter your project name [$starter]: " projectname
echo "Using projectname: $projectname"

if [[ "$projectname" == "" ]]; then
  echo "Skip"
  exit 1
fi

git co -b dev
replace $starter $projectname
remote
