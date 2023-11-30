# wait-action
Github action allowing to wait until after some jobs complete

## Usage examples

```yaml
name: test-after-image-build
on:
  ...
env:
  GITHUB_TOKEN: ${{ secrets.TOKEN }}
jobs:
  e2e-tests:
    runs-on: ubuntu-22.04
    timeout-minutes: 10
    steps:
      - name: Get SHA
        run: |
          if [ ${{ github.event.pull_request.head.sha }} != "" ]; then
            echo "GIT_SHA=${{ github.event.pull_request.head.sha }}" >> "$GITHUB_ENV"
          else
            echo "GIT_SHA=${{ github.sha }}" >> "$GITHUB_ENV"
          fi
      ...
      - name: Wait for images appear in remote container registry
        uses: yandzee/wait-action@main
        env:
          GITHUB_TOKEN: ${{ env.GITHUB_TOKEN }}
        with:
          # Required param
          head-sha: ${{ env.GIT_SHA }}
          poll-delay: 10s
          # Comma separated paths to workflow files
          workflows: .github/workflows/build-images.yaml

```

## Supported kinds of events to wait

- Workflows: it waits until all of them have finished or at least one of them failed
