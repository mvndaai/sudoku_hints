## Get the config so api keys are not in the re
# gh auth login

API_URL=https://api.github.com/repos/mvndaai/secrets/contents/sudoku_hints.jsonc
gh api $API_URL -H "Accept: application/vnd.github.raw" > config.jsonc

# make it so this file never gets uploaded to github
git update-index --assume-unchanged config.jsonc