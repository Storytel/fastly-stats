# Fastly Stackdriver Exporter

A simple utility to pull Fastly stats from their Realtime API and report it to stackdriver.

## Getting started

You must create a [Fastly API Key][fastly-api-key]. Read-only access is enough.

Create all metric descriptors in Stackdriver. Examples are using the `gcloud-config` container as
described in the [README of google/cloud-sdk][google-cloud-sdk]. *This step is optional* but recommended
to get good metric units and descriptions in Stackdriver.

```
docker run --rm -it --volumes-from gcloud-config -e GOOGLE_APPLICATION_CREDENTIALS=/root/.config/gcloud/legacy_credentials/<your-email-here>/adc.json storytel/fastly-stackdriver-exporter -project <GCP-project> -rebuild-metric-descriptors
```

Start the metric collector and reporter. This will run indefinitely and report metrics to Stackdriver

```
docker run --rm -it --volumes-from gcloud-config -e GOOGLE_APPLICATION_CREDENTIALS=/root/.config/gcloud/legacy_credentials/<your-email-here>/adc.json -e FASTLY_API_KEY=<fastly-api-key> -e FASTLY_SERVICE=<fastly-service> storytel/fastly-stackdriver-exporter -project <GCP-project>
```

[google-cloud-sdk]: https://hub.docker.com/r/google/cloud-sdk/
[fastly-api-key]: https://docs.fastly.com/en/guides/using-api-tokens

## Release

The release process is manual (fow now).

Create a git tag with the next version number and push it to origin.
```sh
git tag vX.Y.Z
git push origin v1.5.0
```

Build the docker image and push it to Docker Hub.
```sh
docker build -t storytel/fastly-stats:v1.5.0 .
docker push storytel/fastly-stats:v1.5.0
```
